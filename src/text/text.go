package text

import (
	"git.naspersclassifieds.com/infrastructure/verticals/signalfx2terraform/src/utils"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/signalfx/signalfx-go/chart"
	"github.com/zclconf/go-cty/cty"
)

// Chart - function for generating text chart
func Chart(f *hclwrite.File, chart *chart.Chart) *hclwrite.Body {

	// wrapper around label
	label := utils.LabelProc(chart.Id)

	rootBody := f.Body()
	chartBlock := rootBody.AppendNewBlock("resource", []string{utils.Type[chart.Options.Type], label})
	chartBody := chartBlock.Body()
	chartBody.SetAttributeValue("name", cty.StringVal(chart.Name))
	chartBody.SetAttributeValue("description", cty.StringVal(chart.Description))
	chartBody.SetAttributeValue("markdown", cty.StringVal(chart.Options.Markdown))
	chartBody.AppendNewline()
	return chartBody
}
