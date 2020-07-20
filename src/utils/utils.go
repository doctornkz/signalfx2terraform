package utils

import (
	"fmt"
	"log"
	"strconv"
	"encoding/json"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/signalfx/signalfx-go/chart"
	"github.com/signalfx/signalfx-go/dashboard"
	"github.com/signalfx/signalfx-go/detector"
	"github.com/signalfx/signalfx-go/notification"

	"github.com/signalfx/signalfx-go"
	"github.com/zclconf/go-cty/cty"
)

// LabelProc ...
func LabelProc(name string) string {
	return fmt.Sprintf("sfx-%s", name)

}

// ProgramTextProc ...
func ProgramTextProc(programText string) string {
	return fmt.Sprintf("<<EOF\n%s\nEOF", programText)
}

// MaxDelayProc ...
func MaxDelayProc(chart *chart.Chart) int64 {
	if chart.Options != nil || chart.Options.ProgramOptions != nil {
		return 0
	}
	return int64(*chart.Options.ProgramOptions.MaxDelay / 1000) // Convert to sec
}

// MaxDelayDetectorProc ...
func MaxDelayDetectorProc(detector *detector.Detector) int64 {
	if detector.MaxDelay == nil {
		return 0
	}
	return int64(*detector.MaxDelay / 1000) // Convert to sec
}

// RefreshIntervalProc ...
func RefreshIntervalProc(chart *chart.Chart) int64 {
	if chart.Options.RefreshInterval == nil {
		return 0
	}
	return int64(*chart.Options.RefreshInterval / 1000) // Convert to sec
}

// MinResolutionProc ...
func MinResolutionProc(chart *chart.Chart) int64 {
	if chart.Options != nil || chart.Options.ProgramOptions != nil {
		return 0
	}
	return int64(*chart.Options.ProgramOptions.MinimumResolution / 1000) // Convert to sec
}

// StrToInt ...
func StrToInt(s string) int64 {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int64(num / 1000) // Convert to sec
}

// GroupByProc ...
func GroupByProc(chart *chart.Chart) cty.Value {
	groupBy := chart.Options.GroupBy
	groupByList := []cty.Value{}

	if len(groupBy) == 0 {
		return cty.ListValEmpty(cty.String)
	}
	for _, v := range groupBy {
		ctyV := cty.StringVal(v)
		groupByList = append(groupByList, ctyV)
	}
	return cty.ListVal(groupByList)
}

// DisableSamplingProc ...
func DisableSamplingProc(chart *chart.Chart) cty.Value {
	if chart.Options != nil && chart.Options.ProgramOptions != nil {
		return cty.BoolVal(chart.Options.ProgramOptions.DisableSampling)
	}
	return cty.BoolVal(false)
}

// ColorRangeProc ...
func ColorRangeProc(c *chart.Chart, cb *hclwrite.Body) {
	colorRange := c.Options.ColorRange
   cr := ColorRangeOptions{}
   cr.Color = colorRange.Color
   cr.Max = colorRange.Max
   cr.Min = colorRange.Min

   js, err := json.Marshal(cr)
   if err != nil {
      log.Fatal("Cannot marshal SecondaryVisualization{} structure", err)
   }
   // map json strings to an interface
   // convert from []byte to strings and value type
   s := make(map[string]interface{})
   json.Unmarshal(js, &s)

   setAttributeOptions(cb, s, "color_range")
}

// ColorScale2Proc - create color_scale body
func ColorScale2Proc(c *chart.Chart, cb *hclwrite.Body, color map[int32]string) {
   for _, f := range c.Options.ColorScale2 {
      sv := SecondaryVisualization{}
      sv.Gt = f.Gt
      sv.Gte = f.Gte
      sv.Lt = f.Lt
      sv.Lte = f.Lte
      sv.PaletteIndex = f.PaletteIndex

      js, err := json.Marshal(sv)
      if err != nil {
         log.Fatal("Cannot marshal SecondaryVisualization{} structure", err)
      }

      // map json strings to an interface
      // convert from []byte to strings and value type
      s := make(map[string]interface{})
      json.Unmarshal(js, &s)

      setAttributeOptions(cb, s, "color_scale")
   }
}

// FilterValueProc ...
func FilterValueProc(filter *dashboard.ChartsSingleFilter) cty.Value {
	var valueList []cty.Value
	for _, v := range filter.Value {
		valueList = append(valueList, cty.StringVal(v))
	}
	return cty.ListVal(valueList)
}

// VariableValueProc ...
func VariableValueProc(filter *dashboard.ChartsWebUiFilter) cty.Value {
	var valueList []cty.Value
	for _, v := range filter.Value {
		valueList = append(valueList, cty.StringVal(v))
	}
	return cty.ListVal(valueList)
}

// DensityProc ...
func DensityProc(density *dashboard.DashboardChartDensity) cty.Value {

	switch *density {
	case dashboard.DEFAULT:
		return cty.StringVal("default")

	case dashboard.HIGH:
		return cty.StringVal("high")

	case dashboard.HIGHEST:
		return cty.StringVal("highest")

	case dashboard.LOW:
		return cty.StringVal("low")
	}

	return cty.StringVal("default")
}

// SeverityProc ...
func SeverityProc(rule detector.Rule) cty.Value {

	switch rule.Severity {
	case detector.CRITICAL:
		return cty.StringVal("Critical")

	case detector.WARNING:
		return cty.StringVal("Warning")

	case detector.MAJOR:
		return cty.StringVal("Major")

	case detector.MINOR:
		return cty.StringVal("Minor")
	}

	return cty.StringVal("Info")
}

// NotificationProc ...
/*
Adopted from:
https://github.com/terraform-providers/terraform-provider-signalfx/blob/master/signalfx/notifications.go
Thanks, Cory
*/
func NotificationProc(rule detector.Rule) cty.Value {
	notifications := rule.Notifications
	var notificationList []cty.Value
	if len(notifications) == 0 {
		return cty.ListValEmpty(cty.String)
	}

	for _, n := range notifications {
		route := ""
		nt := n.Type
		switch nt {
		case "BigPanda":
			bp := n.Value.(*notification.BigPandaNotification)
			route = fmt.Sprintf("%s,%s", nt, bp.CredentialId)
		case "Email":
			em := n.Value.(*notification.EmailNotification)
			route = fmt.Sprintf("%s,%s", nt, em.Email)
		case "Office365":
			off := n.Value.(*notification.Office365Notification)
			route = fmt.Sprintf("%s,%s", nt, off.CredentialId)
		case "Opsgenie":
			og := n.Value.(*notification.OpsgenieNotification)
			route = fmt.Sprintf("%s,%s,%s,%s,%s", nt, og.CredentialId, og.ResponderName, og.ResponderId, og.ResponderType)
		case "PagerDuty":
			pd := n.Value.(*notification.PagerDutyNotification)
			route = fmt.Sprintf("%s,%s", nt, pd.CredentialId)
		case "ServiceNow":
			sn := n.Value.(*notification.ServiceNowNotification)
			route = fmt.Sprintf("%s,%s", nt, sn.CredentialId)
		case "Slack":
			sl := n.Value.(*notification.SlackNotification)
			route = fmt.Sprintf("%s,%s,%s", nt, sl.CredentialId, sl.Channel)
		case "Team":
			t := n.Value.(*notification.TeamNotification)
			route = fmt.Sprintf("%s,%s", nt, t.Team)
		case "TeamEmail":
			te := n.Value.(*notification.TeamEmailNotification)
			route = fmt.Sprintf("%s,%s", nt, te.Team)
		case "VictorOps":
			vo := n.Value.(*notification.VictorOpsNotification)
			route = fmt.Sprintf("%s,%s,%s", nt, vo.CredentialId, vo.RoutingKey)
		case "Webhook":
			wh := n.Value.(*notification.WebhookNotification)
			route = fmt.Sprintf("%s,%s,%s,%s", nt, wh.CredentialId, wh.Secret, wh.Url)
		case "XMatters":
			xm := n.Value.(*notification.XMattersNotification)
			route = fmt.Sprintf("%s,%s", nt, xm.CredentialId)
		}

		ctyV := cty.StringVal(route)

		notificationList = append(notificationList, ctyV)
	}
	return cty.TupleVal(notificationList)
}

// NotificationProcV1 ...
func NotificationProcV1(notifications []map[string]string) cty.Value {
	var notifocationList []cty.Value
	if len(notifications) == 0 {
		return cty.ListValEmpty(cty.String)
	}
	for _, item := range notifications {
		route := "" // Create new string with notification routing
		switch item["type"] {
		case "email":
			route = fmt.Sprintf("Email,%s", item["email"])
		case "pagerduty":
			route = fmt.Sprintf("PagerDuty,%s", item["credentialId"])
		case "bigpanda":
			route = fmt.Sprintf("BigPanda,%s", item["credentialId"])
		case "office365":
			route = fmt.Sprintf("Office365,%s", item["credentialId"])
		case "servicenow":
			route = fmt.Sprintf("ServiceNow,%s", item["credentialId"])
		case "xmatters":
			route = fmt.Sprintf("XMatters,%s", item["credentialId"])
		case "slack":
			route = fmt.Sprintf("Slack,%s,%s", item["credentialId"], item["channel"])
		case "webhook":
			route = fmt.Sprintf("WebHook%s,%s", item["secret"], item["url"])
		case "team":
			route = fmt.Sprintf("Team,%s", item["team"])
		case "teamemail":
			route = fmt.Sprintf("TeamEmail,%s", item["team"])
		case "opsgenie":
			route = fmt.Sprintf("OpsGenie,%s,%s,%s,%s,%s",
				item["credentialId"],
				item["credentialName"],
				item["responderName"],
				item["responderId"],
				item["responderType"])
		case "victorops":
			route = fmt.Sprintf("VictorOps,%s,%s", item["credentialId"], item["routingKey"])
		}
		ctyV := cty.StringVal(route)
		notifocationList = append(notifocationList, ctyV)
	}
	return cty.TupleVal(notifocationList)
}

// OnChartLegendProc ...
func OnChartLegendProc(chart *chart.Chart) cty.Value {

	if chart.Options.OnChartLegendOptions != nil {
		dimensionInLegend := chart.Options.OnChartLegendOptions.DimensionInLegend

		/* Two specific cases:
		https://www.terraform.io/docs/providers/signalfx/r/time_chart.html
			* `property` The name of the property to display.
			Note the special values of `plot_label` (corresponding with the API's `sf_metric`)
			which shows the label of the time series `publish()` and
			`metric` (corresponding with the API's `sf_originatingMetric`)
			that shows the name of the metric for the time series being displayed
		*/

		if dimensionInLegend == "sf_metric" {
			dimensionInLegend = "plot_label"
		}
		if dimensionInLegend == "sf_originatingMetric" {
			dimensionInLegend = "metric"
		}

		return cty.StringVal(dimensionInLegend)
	}
	return cty.StringVal("")
}

// LegendShowProc ...
func LegendShowProc(chart *chart.Chart) cty.Value {
	var valueList []cty.Value
	if len(chart.Options.LegendOptions.Fields) == 0 {
		return cty.ListValEmpty(cty.String)
	}
	for _, v := range chart.Options.LegendOptions.Fields {
		attr := map[string]cty.Value{}
		attr["property"] = cty.StringVal(v.Property)
		attr["enabled"] = cty.BoolVal(v.Enabled)
		ctyV := cty.ObjectVal(attr)
		valueList = append(valueList, ctyV)

	}
	return cty.TupleVal(valueList)
}

// VizProc ...
func VizProc(chart *chart.Chart) cty.Value {
	var valueList []cty.Value
	if len(chart.Options.PublishLabelOptions) == 0 && len(chart.Options.EventPublishLabelOptions) == 0 {
		return cty.ListValEmpty(cty.String)
	}

	/* Processing data label, like
	A = data('gitlab.deployment', filter=filter('application_name', 'carspt') and \
	filter('status', 'succeed'), rollup='sum').publish(label='A')
	*/

	for _, v := range chart.Options.PublishLabelOptions {
		attr := map[string]cty.Value{}
		attr["display_name"] = cty.StringVal(v.DisplayName)
		attr["label"] = cty.StringVal(v.Label)
		// FIXME: Ugly hack here, 0 not null!
		if v.PaletteIndex != nil {
			attr["color"] = cty.StringVal(Color[*v.PaletteIndex])
		}

		/*
		 Specifies the position of the Y-axis for the plot associated with the SignalFlow statement.
		 If `yAxis` is set to 0, the axis is on the left side; otherwise it's on the right.
		 The default is 0 (left side).<br> **Note** -- This option is only available if 'options.type`
		 is `TimeSeriesChart`, `List`, or `SingleValue`.
		*/
		if v.YAxis != 0 { // Empty field not allowed
			attr["axis"] = cty.StringVal("right")
		}

		if v.PlotType != "" { // Empty string not allowed
			attr["plot_type"] = cty.StringVal(v.PlotType)
		}
		if v.ValueUnit != "" { // Empty string not allowed
			attr["value_unit"] = cty.StringVal(v.ValueUnit)
		}
		attr["value_prefix"] = cty.StringVal(v.ValuePrefix)
		attr["value_suffix"] = cty.StringVal(v.ValueSuffix)

		ctyV := cty.ObjectVal(attr)
		valueList = append(valueList, ctyV)

	}
	return cty.TupleVal(valueList)
}

// EventProc - create event_options body
func EventProc(c *chart.Chart, cb *hclwrite.Body) {
   // assign all options available
   // if no options it wont panic, structure uses `omitempty`
   if c.Options.EventPublishLabelOptions != nil {
      for _, f := range c.Options.EventPublishLabelOptions {
         e := EventPublishLabelOptions{}
         e.DisplayName  = f.DisplayName
         e.Label        = f.Label
         e.PaletteIndex = f.PaletteIndex

         // marshal structure to json strings
         js, err := json.Marshal(e)
         if err != nil {
            log.Fatal("Cannot marshal EventPublishLabelOptions{} structure", err)
         }

         // map json strings to an interface
         // convert from []byte to strings and value type
         s := make(map[string]interface{})
         json.Unmarshal(js, &s)

         setAttributeOptions(cb, s, "event_options")
      }
   }
}

// ShortVizProc ...
func ShortVizProc(chart *chart.Chart) cty.Value {
	var valueList []cty.Value
	if len(chart.Options.PublishLabelOptions) == 0 {
		return cty.ListValEmpty(cty.String)
	}
	for _, v := range chart.Options.PublishLabelOptions {
		attr := map[string]cty.Value{}
		attr["display_name"] = cty.StringVal(v.DisplayName)
		attr["label"] = cty.StringVal(v.Label)
		// FIXME: Ugly hack here, 0 not null!
		if v.PaletteIndex != nil {
			attr["color"] = cty.StringVal(Color[*v.PaletteIndex])
		}
		if v.ValueUnit != "" { // Empty string not allowed
			attr["value_unit"] = cty.StringVal(v.ValueUnit)
		}
		attr["value_prefix"] = cty.StringVal(v.ValuePrefix)
		attr["value_suffix"] = cty.StringVal(v.ValueSuffix)

		ctyV := cty.ObjectVal(attr)
		valueList = append(valueList, ctyV)

	}
	return cty.TupleVal(valueList)
}

// CreateDashboard - function for generating dashboard
func CreateDashboard(f *hclwrite.File, dashboard *dashboard.Dashboard, client *signalfx.Client) *hclwrite.Body {
	rootBody := f.Body()
	dashBlock := rootBody.AppendNewBlock("resource", []string{"signalfx_dashboard", dashboard.Id})
	dashBody := dashBlock.Body()
	dashBody.SetAttributeValue("dashboard_group", cty.StringVal(dashboard.GroupId))
	dashBody.SetAttributeValue("name", cty.StringVal(fmt.Sprintf("test-%s", dashboard.Name))) // TODO: Hardcode to prevent self-destroy
	dashBody.SetAttributeValue("description", cty.StringVal(dashboard.Description))
	dashBody.SetAttributeValue("charts_resolution", DensityProc(dashboard.ChartDensity))
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
			dashBody.SetAttributeValue("start_time", cty.NumberIntVal(StrToInt(start)))
			dashBody.SetAttributeValue("end_time", cty.NumberIntVal(StrToInt(end)))
		}
	}

	dashBody.AppendNewline()

	// Filter section processing
	for _, filter := range dashboard.Filters.Sources {
		filterBlock := dashBody.AppendNewBlock("filter", nil)
		filterBody := filterBlock.Body()
		filterBody.SetAttributeValue("property", cty.StringVal(filter.Property))
		filterBody.SetAttributeValue("values", FilterValueProc(filter))
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
		variableBody.SetAttributeValue("value_required", cty.BoolVal(variable.Required))
		if variable.Value != nil {
			variableBody.SetAttributeValue("values", VariableValueProc(variable))
		}
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

		chartID := fmt.Sprintf("%s.%s.id", Type[chartHelper.Options.Type], LabelProc(chartHelper.Id))

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

// GetVizOptions - create viz_options body
func GetVizOptions(c *chart.Chart, cb *hclwrite.Body){

   // assign all options available
   // if no options it wont panic, structure uses `omitempty`
   if c.Options.LegendOptions != nil {
      for _, f := range c.Options.PublishLabelOptions {
         p := PublishLabelOptions{}
         p.Label = f.Label
         p.DisplayName = f.DisplayName
         p.PaletteIndex = f.PaletteIndex
         p.YAxis = f.YAxis
         p.PlotType = f.PlotType
         p.ValueUnit = f.ValueUnit
         p.ValuePrefix = f.ValuePrefix
         p.ValueSuffix = f.ValueSuffix

         // marshal structure to json strings
         js, err := json.Marshal(p)
         if err != nil {
            log.Fatal("Cannot marshal PublishLabelOptions{} structure", err)
         }

         // map json strings to an interface
         // convert from []byte to strings and value type
         s := make(map[string]interface{})
         json.Unmarshal(js, &s)

         setAttributeOptions(cb, s, "viz_options")
      }
   }
}

// setAttributeOptions - fill Chart Body with attributes
func setAttributeOptions(cb *hclwrite.Body, s map[string]interface{}, name string) {
   b := cb.AppendNewBlock(name, nil)
   bc := b.Body()

   // range over interface and create the viz_options
   // change the type of cty accordingly
   for k, v := range s {
      switch v.(type){
         case float64:
            if k == "color"{
               c := OptionsColor[int32(v.(float64))]
               bc.SetAttributeValue(k, cty.StringVal(c))
            } else {
               bc.SetAttributeValue(k, cty.NumberFloatVal(v.(float64)))
            }
         case string:
            bc.SetAttributeValue(k, cty.StringVal(v.(string)))
      }
   }
}

// GetLegendOptionsBlock - create legend_options_fields body
func GetLegendOptionsBlock(c *chart.Chart, cb *hclwrite.Body){
   legend := &hclwrite.Block{}
   if c.Options.LegendOptions != nil {
      for _, field := range c.Options.LegendOptions.Fields {
         legend = cb.AppendNewBlock("legend_options_fields", nil)
         legendOptionsBody := legend.Body()
         legendOptionsBody.SetAttributeValue("property", cty.StringVal(field.Property))
         legendOptionsBody.SetAttributeValue("enabled", cty.BoolVal(field.Enabled))
      }
   }
}
