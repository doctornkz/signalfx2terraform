package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/signalfx/signalfx-go/detector"

	"github.com/zclconf/go-cty/cty"
)

// CreateDetector - function for generating detector from API
func CreateDetector(f *hclwrite.File, detector *detector.Detector) *hclwrite.Body {

	// wrapper around label
	label := labelProc(detector.Id)

	rootBody := f.Body()
	detectorBlock := rootBody.AppendNewBlock("resource", []string{"signalfx_detector", label})
	detectorBody := detectorBlock.Body()
	detectorBody.SetAttributeValue("name", cty.StringVal(fmt.Sprintf("test-%s", detector.Name))) // TODO: Hardcoded to prevent self-destroy
	detectorBody.SetAttributeValue("description", cty.StringVal(detector.Description))
	detectorBody.SetAttributeValue("max_delay", cty.NumberIntVal(maxDelayDetectorProc(detector)))
	detectorBody.SetAttributeTraversal("program_text", hcl.Traversal{
		hcl.TraverseRoot{
			Name: programTextProc(detector.ProgramText),
		},
	})

	// Rules processing
	for _, rule := range detector.Rules {
		ruleBlock := detectorBody.AppendNewBlock("rule", nil)
		ruleBody := ruleBlock.Body()
		ruleBody.SetAttributeValue("description", cty.StringVal(rule.Description))
		ruleBody.SetAttributeValue("severity", severityProc(*rule))
		ruleBody.SetAttributeValue("detect_label", cty.StringVal(rule.DetectLabel))
		ruleBody.SetAttributeValue("notifications", notificationProc(*rule))

	}
	return detectorBody
}

// CreateDetectorV1 - function for generating detector from old version API.
func CreateDetectorV1(f *hclwrite.File, api string, detectorID string, token string) *hclwrite.Body {
	client := &http.Client{}
	detectorURL := fmt.Sprintf("%v/v1/detector/%v", api, detectorID)
	req, err := http.NewRequest("GET", detectorURL, nil)
	req.Header.Add("X-SF-TOKEN", token)
	detectorResponse, err := client.Do(req)

	if err != nil {
		log.Fatalf("Can't fetch data from API %v, %v", detectorURL, err)
	}
	defer detectorResponse.Body.Close()
	body, err := ioutil.ReadAll(detectorResponse.Body)
	if err != nil {
		log.Fatalf("Can't read body JSON, %v", err)
	}

	var detector DetectorV1
	err = json.Unmarshal(body, &detector)
	if err != nil {
		log.Fatalf("Can't load JSON, %v", err)
	}

	// wrapper around label
	label := labelProc(detector.Sf_id)

	rootBody := f.Body()
	detectorBlock := rootBody.AppendNewBlock("resource", []string{"signalfx_detector", label})
	detectorBody := detectorBlock.Body()
	detectorBody.SetAttributeValue("name", cty.StringVal(fmt.Sprintf("testV1-%s", detector.Sf_detector))) // TODO: Hardcoded to prevent self-destroy
	detectorBody.SetAttributeValue("description", cty.StringVal(detector.Sf_description))
	detectorBody.SetAttributeValue("max_delay", cty.NumberIntVal(detector.Sf_jobMaxDelay))
	detectorBody.SetAttributeTraversal("program_text", hcl.Traversal{
		hcl.TraverseRoot{
			Name: programTextProc(detector.Sf_programText),
		},
	})

	// Rules processing
	for _, rule := range detector.Sf_rules {
		ruleBlock := detectorBody.AppendNewBlock("rule", nil)
		ruleBody := ruleBlock.Body()
		ruleBody.SetAttributeValue("severity", cty.StringVal(rule.Severity))
		ruleBody.SetAttributeValue("detect_label", cty.StringVal(rule.DetectLabel))
		ruleBody.SetAttributeValue("description", cty.StringVal(rule.Readable))
		ruleBody.SetAttributeValue("notifications", notificationProcV1(rule.Notifications))

	}

	return detectorBody

}
