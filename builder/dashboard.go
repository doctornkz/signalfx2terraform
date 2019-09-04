package builder

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/signalfx/signalfx-go"
	"github.com/signalfx/signalfx-go/dashboard"
	"github.com/zclconf/go-cty/cty"
)

// CreateDashboard - function for generating dashboard
func CreateDashboard(f *hclwrite.File, dashboard *dashboard.Dashboard, client *signalfx.Client) *hclwrite.Body {
	rootBody := f.Body()
	dashBlock := rootBody.AppendNewBlock("resource", []string{"signalfx_dashboard", dashboard.Id})
	dashBody := dashBlock.Body()
	dashBody.SetAttributeValue("dashboard_group", cty.StringVal(dashboard.GroupId))
	dashBody.SetAttributeValue("name", cty.StringVal(fmt.Sprintf("test-%s", dashboard.Name))) // TODO: Hardcode to prevent self-destroy
	dashBody.SetAttributeValue("description", cty.StringVal(dashboard.Description))
	dashBody.SetAttributeValue("charts_resolution", densityProc(dashboard.ChartDensity))
	// Complex `Time` logic here.
	// Terraform provider has different description,
	// SignalFX API has different fields,
	// Variables from documentation doesn't implemented.
	// Meh
	if dashboard.Filters.Time != nil { // If fields don't exist - SignalFX will use system vars.
		start := string(dashboard.Filters.Time.Start)
		end := string(dashboard.Filters.Time.End)
		if end == "Now" { // If `Now` - scale relative.
			dashBody.SetAttributeValue("time_range", cty.StringVal(start))
		} else { // If not - absolute
			dashBody.SetAttributeValue("start_time", cty.NumberIntVal(strToInt(start)))
			dashBody.SetAttributeValue("end_time", cty.NumberIntVal(strToInt(end)))
		}
	}

	dashBody.AppendNewline()

	// Filter section processing
	for _, filter := range dashboard.Filters.Sources {
		filterBlock := dashBody.AppendNewBlock("filter", nil)
		filterBody := filterBlock.Body()
		filterBody.SetAttributeValue("property", cty.StringVal(filter.Property))
		filterBody.SetAttributeValue("values", filterValueProc(filter))
		filterBody.SetAttributeValue("negated", cty.BoolVal(filter.NOT))
		filterBody.SetAttributeValue("apply_if_exist", cty.BoolVal(filter.ApplyIfExists))
		// negated
	}

	// Variables section processing
	for _, variable := range dashboard.Filters.Variables {
		// Not full implementation,
		// see https://github.com/signalfx/terraform-provider-signalfx/blob/master/website/docs/r/dashboard.html.markdown

		variableBlock := dashBody.AppendNewBlock("variable", nil)
		variableBody := variableBlock.Body()
		variableBody.SetAttributeValue("property", cty.StringVal(variable.Property))
		variableBody.SetAttributeValue("description", cty.StringVal(variable.Description))
		variableBody.SetAttributeValue("alias", cty.StringVal(variable.Alias))
		variableBody.SetAttributeValue("values", variableValueProc(variable))
		variableBody.SetAttributeValue("value_required", cty.BoolVal(variable.Required))
		variableBody.SetAttributeValue("replace_only", cty.BoolVal(variable.ReplaceOnly))
		variableBody.SetAttributeValue("apply_if_exist", cty.BoolVal(variable.ApplyIfExists))

	}

	// Charts position processing
	for _, chart := range dashboard.Charts {
		// Receive data about chart from API
		// TODO: Need to implement init() section and Client class
		chartHelper, err := client.GetChart(chart.ChartId)
		if err != nil {
			log.Printf("Chart error: %v", err)
			log.Fatal("Can't get chart")
		}

		chartID := fmt.Sprintf("\"${%s.%s.id}\"", t.types[chartHelper.Options.Type], labelProc(chartHelper.Id))

		chartPosBlock := dashBody.AppendNewBlock("chart", nil)
		chartPosBody := chartPosBlock.Body()

		chartPosBody.SetAttributeTraversal("chart_id", hcl.Traversal{hcl.TraverseRoot{Name: chartID}})
		chartPosBody.SetAttributeValue("column", cty.NumberIntVal(int64(chart.Column)))
		chartPosBody.SetAttributeValue("row", cty.NumberIntVal(int64(chart.Row)))
		chartPosBody.SetAttributeValue("width", cty.NumberIntVal(int64(chart.Width)))
		chartPosBody.SetAttributeValue("height", cty.NumberIntVal(int64(chart.Height)))
	}
	dashBody.AppendNewline()
	return dashBody
}
