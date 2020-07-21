package timeseries

import (
	"github.com/doctornkz/signalfx2terraform/src/utils"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/signalfx/signalfx-go/chart"
	"github.com/zclconf/go-cty/cty"
)

// Chart - function for generating time series chart
func Chart(f *hclwrite.File, chart *chart.Chart) *hclwrite.Body {

	// program_text wrapper
	programText := utils.ProgramTextProc(chart.ProgramText)
	// group_by wrapper
	groupBy := chart.Options.GroupBy
	if groupBy == nil {
		groupBy = []string{}
	}
	// tags wrapper
	tags := chart.Tags
	if tags == nil {
		tags = []string{}
	}
	// wrapper around label
	label := utils.LabelProc(chart.Id)

	rootBody := f.Body()
	chartBlock := rootBody.AppendNewBlock("resource", []string{utils.Type[chart.Options.Type], label})
	chartBody := chartBlock.Body()
	chartBody.SetAttributeValue("name", cty.StringVal(chart.Name))
	chartBody.SetAttributeValue("description", cty.StringVal(chart.Description))
	chartBody.SetAttributeValue("plot_type", cty.StringVal(chart.Options.DefaultPlotType))
	chartBody.SetAttributeValue("stacked", cty.BoolVal(chart.Options.Stacked))
	chartBody.SetAttributeValue("axes_include_zero", cty.BoolVal(chart.Options.IncludeZero))

	// Histograms processing
	if chart.Options.HistogramChartOptions != nil {
		histogramOptionsBlock := chartBody.AppendNewBlock("histogram_options", nil)
		histogramOptionsBody := histogramOptionsBlock.Body()
		histogramOptionsBody.SetAttributeValue("color_theme", cty.StringVal(utils.OptionsColor[*chart.Options.HistogramChartOptions.ColorThemeIndex]))
	}

   // legend_options_fields
   utils.GetLegendOptionsBlock(chart, chartBody)

   // create viz_options
	utils.GetVizOptions(chart, chartBody)


	// create event_options
	utils.EventProc(chart, chartBody)
	// chartBody.SetAttributeValue("event_options", utils.EventProc(chart))

	chartBody.SetAttributeTraversal("program_text", hcl.Traversal{
		hcl.TraverseRoot{
			Name: programText,
		},
	})

	chartBody.SetAttributeValue("disable_sampling", utils.DisableSamplingProc(chart))
	chartBody.SetAttributeValue("minimum_resolution", cty.NumberIntVal(utils.MinResolutionProc(chart)))
	chartBody.SetAttributeValue("unit_prefix", cty.StringVal(chart.Options.UnitPrefix))
	chartBody.SetAttributeValue("max_delay", cty.NumberIntVal(utils.MaxDelayProc(chart)))
	chartBody.SetAttributeValue("color_by", cty.StringVal(chart.Options.ColorBy))
	chartBody.SetAttributeValue("on_chart_legend_dimension", utils.OnChartLegendProc(chart))
	// Time range processing
	if chart.Options.Time != nil { // Checking Time structure against nil

		if chart.Options.Time.Type == "relative" {
			/* Convert to sec, we have different format here:
			See details here: https://github.com/terraform-providers/terraform-provider-signalfx/issues/55
			*/
			chartBody.SetAttributeValue("time_range", cty.NumberIntVal(*chart.Options.Time.Range/1000))
		}

		if chart.Options.Time.Type == "absolute" {
			chartBody.SetAttributeValue("start_time", cty.NumberIntVal(*chart.Options.Time.Start/1000))
			chartBody.SetAttributeValue("end_time", cty.NumberIntVal(*chart.Options.Time.End/1000))
		}
	}
	return chartBody
}
