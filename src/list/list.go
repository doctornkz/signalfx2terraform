package list

import (
   "github.com/doctornkz/signalfx2terraform/src/utils"
   "github.com/hashicorp/hcl2/hcl"
   "github.com/hashicorp/hcl2/hclwrite"
   "github.com/signalfx/signalfx-go/chart"
   "github.com/zclconf/go-cty/cty"
)

// Chart - function for generating line chart
func Chart(f *hclwrite.File, chart *chart.Chart) *hclwrite.Body {

   // program_text wrapper
   programText := utils.ProgramTextProc(chart.ProgramText)
   // group_by wrapper

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

   chartBody.SetAttributeTraversal("program_text", hcl.Traversal{
      hcl.TraverseRoot{
         Name: programText,
      },
   })
   chartBody.SetAttributeValue("disable_sampling", utils.DisableSamplingProc(chart))
   chartBody.SetAttributeValue("unit_prefix", cty.StringVal(chart.Options.UnitPrefix))
   if chart.Options.ColorBy == "Range" {
      // chartBody.SetAttributeValue("color_range", utils.ColorRangeProc(chart))
      utils.ColorRangeProc(chart, chartBody)
   }

   cby := chart.Options.ColorBy
   chartBody.SetAttributeValue("color_by", cty.StringVal(cby))
   if cby == "Scale" {
      utils.ColorScale2Proc(chart, chartBody)
   }

   chartBody.SetAttributeValue("max_delay", cty.NumberIntVal(utils.MaxDelayProc(chart)))
   chartBody.SetAttributeValue("refresh_interval", cty.NumberIntVal(utils.RefreshIntervalProc(chart)))

   // legend_options_fields
   utils.GetLegendOptionsBlock(chart, chartBody)

   // create viz_options
   utils.GetVizOptions(chart, chartBody)

   chartBody.SetAttributeValue("secondary_visualization", cty.StringVal(chart.Options.SecondaryVisualization))
   chartBody.AppendNewline()
   return chartBody
}
