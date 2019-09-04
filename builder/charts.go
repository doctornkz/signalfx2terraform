package builder

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/signalfx/signalfx-go/chart"

	"github.com/zclconf/go-cty/cty"
)

// HeatMapChart - function for generating heatmap chart
func HeatMapChart(f *hclwrite.File, chart *chart.Chart) *hclwrite.Body {

	// program_text wrapper
	programText := programTextProc(chart.ProgramText)

	// tags wrapper
	tags := chart.Tags
	if tags == nil {
		tags = []string{}
	}

	// wrapper around label
	label := labelProc(chart.Id)

	rootBody := f.Body()
	chartBlock := rootBody.AppendNewBlock("resource", []string{t.types[chart.Options.Type], label})
	chartBody := chartBlock.Body()
	chartBody.SetAttributeValue("name", cty.StringVal(chart.Name))
	chartBody.SetAttributeValue("description", cty.StringVal(chart.Description))
	chartBody.SetAttributeTraversal("program_text", hcl.Traversal{
		hcl.TraverseRoot{
			Name: programText,
		},
	})
	chartBody.SetAttributeValue("unit_prefix", cty.StringVal(chart.Options.UnitPrefix))
	chartBody.SetAttributeValue("max_delay", cty.NumberIntVal(maxDelayProc(chart)))
	chartBody.SetAttributeValue("group_by", groupByProc(chart))
	chartBody.SetAttributeValue("color_scale", colorScale2Proc(chart))
	chartBody.SetAttributeValue("minimum_resolution", cty.NumberIntVal(minResolutionProc(chart)))
	chartBody.SetAttributeValue("disable_sampling", disableSamplingProc(chart))
	chartBody.SetAttributeValue("refresh_interval", cty.NumberIntVal(refreshIntervalProc(chart)))
	chartBody.AppendNewline()
	return chartBody
}

// SingleValueChart - function for generating single value chart
func SingleValueChart(f *hclwrite.File, chart *chart.Chart) *hclwrite.Body {

	// program_text wrapper
	programText := programTextProc(chart.ProgramText)
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
	label := labelProc(chart.Id)

	rootBody := f.Body()
	chartBlock := rootBody.AppendNewBlock("resource", []string{t.types[chart.Options.Type], label})
	chartBody := chartBlock.Body()
	chartBody.SetAttributeValue("name", cty.StringVal(chart.Name))
	chartBody.SetAttributeValue("description", cty.StringVal(chart.Description))
	chartBody.SetAttributeTraversal("program_text", hcl.Traversal{
		hcl.TraverseRoot{
			Name: programText,
		},
	})
	chartBody.SetAttributeValue("unit_prefix", cty.StringVal(chart.Options.UnitPrefix))
	chartBody.SetAttributeValue("max_delay", cty.NumberIntVal(maxDelayProc(chart)))
	chartBody.SetAttributeValue("color_by", cty.StringVal(chart.Options.ColorBy))
	chartBody.SetAttributeValue("color_scale", colorScale2Proc(chart))
	chartBody.SetAttributeValue("refresh_interval", cty.NumberIntVal(refreshIntervalProc(chart)))
	chartBody.SetAttributeValue("viz_options", shortVizProc(chart))
	chartBody.SetAttributeValue("secondary_visualization", cty.StringVal(chart.Options.SecondaryVisualization))
	chartBody.AppendNewline()
	return chartBody
}

// TimeSeriesChart - function for generating time series chart
func TimeSeriesChart(f *hclwrite.File, chart *chart.Chart) *hclwrite.Body {

	// program_text wrapper
	programText := programTextProc(chart.ProgramText)
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
	label := labelProc(chart.Id)

	rootBody := f.Body()
	chartBlock := rootBody.AppendNewBlock("resource", []string{t.types[chart.Options.Type], label})
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
		histogramOptionsBody.SetAttributeValue("color_theme", cty.StringVal(vizcolor.color[*chart.Options.HistogramChartOptions.ColorThemeIndex]))
	}
	if chart.Options.LegendOptions != nil {
		for _, field := range chart.Options.LegendOptions.Fields {
			legendOptionsBlock := chartBody.AppendNewBlock("legend_options_fields", nil)
			legendOptionsBody := legendOptionsBlock.Body()
			legendOptionsBody.SetAttributeValue("property", cty.StringVal(field.Property))
			legendOptionsBody.SetAttributeValue("enabled", cty.BoolVal(field.Enabled))
		}
	}
	chartBody.SetAttributeValue("viz_options", vizProc(chart))
	chartBody.SetAttributeValue("event_options", eventProc(chart))
	chartBody.SetAttributeTraversal("program_text", hcl.Traversal{
		hcl.TraverseRoot{
			Name: programText,
		},
	})

	chartBody.SetAttributeValue("disable_sampling", disableSamplingProc(chart))
	chartBody.SetAttributeValue("minimum_resolution", cty.NumberIntVal(minResolutionProc(chart)))
	chartBody.SetAttributeValue("unit_prefix", cty.StringVal(chart.Options.UnitPrefix))
	chartBody.SetAttributeValue("max_delay", cty.NumberIntVal(maxDelayProc(chart)))
	chartBody.SetAttributeValue("color_by", cty.StringVal(chart.Options.ColorBy))
	chartBody.SetAttributeValue("on_chart_legend_dimension", onChartLegendProc(chart))
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

// ListChart - function for generating line chart
func ListChart(f *hclwrite.File, chart *chart.Chart) *hclwrite.Body {

	// program_text wrapper
	programText := programTextProc(chart.ProgramText)

	// tags wrapper
	tags := chart.Tags
	if tags == nil {
		tags = []string{}
	}
	// wrapper around label
	label := labelProc(chart.Id)

	rootBody := f.Body()
	chartBlock := rootBody.AppendNewBlock("resource", []string{t.types[chart.Options.Type], label})
	chartBody := chartBlock.Body()
	chartBody.SetAttributeValue("name", cty.StringVal(chart.Name))
	chartBody.SetAttributeValue("description", cty.StringVal(chart.Description))

	chartBody.SetAttributeTraversal("program_text", hcl.Traversal{
		hcl.TraverseRoot{
			Name: programText,
		},
	})
	chartBody.SetAttributeValue("disable_sampling", disableSamplingProc(chart))
	chartBody.SetAttributeValue("unit_prefix", cty.StringVal(chart.Options.UnitPrefix))

	chartBody.SetAttributeValue("max_delay", cty.NumberIntVal(maxDelayProc(chart)))
	chartBody.SetAttributeValue("refresh_interval", cty.NumberIntVal(refreshIntervalProc(chart)))
	chartBody.SetAttributeValue("viz_options", shortVizProc(chart))
	chartBody.AppendNewline()
	return chartBody
}

// TextChart - function for generating text chart
func TextChart(f *hclwrite.File, chart *chart.Chart) *hclwrite.Body {

	// wrapper around label
	label := labelProc(chart.Id)

	rootBody := f.Body()
	chartBlock := rootBody.AppendNewBlock("resource", []string{t.types[chart.Options.Type], label})
	chartBody := chartBlock.Body()
	chartBody.SetAttributeValue("name", cty.StringVal(chart.Name))
	chartBody.SetAttributeValue("description", cty.StringVal(chart.Description))
	chartBody.SetAttributeValue("markdown", cty.StringVal(chart.Options.Markdown))
	chartBody.AppendNewline()
	return chartBody
}
