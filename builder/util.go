package builder

import (
	"fmt"
	"strconv"

	"github.com/signalfx/signalfx-go/chart"
	"github.com/signalfx/signalfx-go/dashboard"
	"github.com/signalfx/signalfx-go/detector"
	"github.com/signalfx/signalfx-go/notification"

	"github.com/zclconf/go-cty/cty"
)

func labelProc(name string) string {
	return fmt.Sprintf("sfx-%s", name)

}

func programTextProc(programText string) string {
	return fmt.Sprintf("<<EOF\n%s\nEOF", programText)
}

//

func maxDelayProc(chart *chart.Chart) int64 {
	if chart.Options != nil || chart.Options.ProgramOptions != nil {
		return 0
	}
	return int64(*chart.Options.ProgramOptions.MaxDelay / 1000) // Convert to sec
}

func maxDelayDetectorProc(detector *detector.Detector) int64 {

	return int64(*detector.MaxDelay / 1000) // Convert to sec
}

func refreshIntervalProc(chart *chart.Chart) int64 {
	if chart.Options.RefreshInterval == nil {
		return 0
	}
	return int64(*chart.Options.RefreshInterval / 1000) // Convert to sec
}

func minResolutionProc(chart *chart.Chart) int64 {
	if chart.Options != nil || chart.Options.ProgramOptions != nil {
		return 0
	}
	return int64(*chart.Options.ProgramOptions.MinimumResolution / 1000) // Convert to sec
}

func strToInt(s string) int64 {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int64(num / 1000) // Convert to sec
}

func groupByProc(chart *chart.Chart) cty.Value {
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

func disableSamplingProc(chart *chart.Chart) cty.Value {
	if chart.Options != nil && chart.Options.ProgramOptions != nil {
		return cty.BoolVal(chart.Options.ProgramOptions.DisableSampling)
	}
	return cty.BoolVal(false)
}

func colorRangeProc(chart *chart.Chart) cty.Value {
	attr := map[string]cty.Value{}
	colorRange := chart.Options.ColorRange
	if colorRange != nil {
		if colorRange.Color == "" {
			attr["color"] = cty.StringVal("#05ce00")
		} else {
			attr["color"] = cty.StringVal(colorRange.Color)
		}

		attr["min_value"] = cty.NumberFloatVal(colorRange.Min)
		attr["max_value"] = cty.NumberFloatVal(colorRange.Max)
	} else {
		attr["color"] = cty.StringVal("#05ce00")
	}
	return cty.ObjectVal(attr)
}

func colorScale2Proc(chart *chart.Chart) cty.Value {
	colorScale := chart.Options.ColorScale2
	var colorScaleList []cty.Value
	if len(colorScale) == 0 {
		return cty.ListValEmpty(cty.String)
	}

	if chart.Options.ColorBy != "Scale" {
		return cty.ListValEmpty(cty.String)
	}

	for _, v := range colorScale {

		gt := cty.NullVal(cty.Number)
		gte := cty.NullVal(cty.Number)
		lt := cty.NullVal(cty.Number)
		lte := cty.NullVal(cty.Number)

		attr := map[string]cty.Value{}

		isEmpty := true

		if v.Gt != nil {
			gt = cty.NumberFloatVal(float64(*v.Gt))
			attr["gt"] = gt
			isEmpty = false
		}
		if v.Gte != nil {
			gte = cty.NumberFloatVal(float64(*v.Gte))
			attr["gte"] = gte
			isEmpty = false
		}
		if v.Lt != nil {
			lt = cty.NumberFloatVal(float64(*v.Lt))
			attr["lt"] = lt
			isEmpty = false
		}
		if v.Lte != nil {
			lte = cty.NumberFloatVal(float64(*v.Lte))
			attr["lte"] = lte
			isEmpty = false
		}
		if isEmpty == false {
			attr["color"] = cty.StringVal(c.color[*v.PaletteIndex])
		}

		ctyV := cty.ObjectVal(attr)

		colorScaleList = append(colorScaleList, ctyV)

	}
	return cty.TupleVal(colorScaleList)
}

func filterValueProc(filter *dashboard.ChartsSingleFilter) cty.Value {
	var valueList []cty.Value
	for _, v := range filter.Value {
		valueList = append(valueList, cty.StringVal(v))
	}
	return cty.ListVal(valueList)
}

func variableValueProc(filter *dashboard.ChartsWebUiFilter) cty.Value {
	var valueList []cty.Value
	for _, v := range filter.Value {
		valueList = append(valueList, cty.StringVal(v))
	}
	return cty.ListVal(valueList)
}

func densityProc(density *dashboard.DashboardChartDensity) cty.Value {

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

func severityProc(rule detector.Rule) cty.Value {

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

/*
Adopted from:
https://github.com/terraform-providers/terraform-provider-signalfx/blob/master/signalfx/notifications.go
Thanks, Cory
*/
func notificationProc(rule detector.Rule) cty.Value {
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

func notificationProcV1(notifications []map[string]string) cty.Value {
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

func onChartLegendProc(chart *chart.Chart) cty.Value {

	if chart.Options.OnChartLegendOptions != nil {
		return cty.StringVal(chart.Options.OnChartLegendOptions.DimensionInLegend)
	}
	return cty.StringVal("")
}

func legendShowProc(chart *chart.Chart) cty.Value {
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

func vizProc(chart *chart.Chart) cty.Value {
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
			attr["color"] = cty.StringVal(c.color[*v.PaletteIndex])
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

func eventProc(chart *chart.Chart) cty.Value {
	var valueList []cty.Value
	if len(chart.Options.PublishLabelOptions) == 0 && len(chart.Options.EventPublishLabelOptions) == 0 {
		return cty.ListValEmpty(cty.String)
	}

	attr := map[string]cty.Value{}
	for _, v := range chart.Options.EventPublishLabelOptions {

		attr["display_name"] = cty.StringVal(v.DisplayName)
		attr["label"] = cty.StringVal(v.Label)
		if v.PaletteIndex != nil {
			// due issue https://github.com/signalfx/signalfx-go/issues/26
			// need cast int32 against int
			paletteIntexInt32 := int32(*v.PaletteIndex)
			attr["color"] = cty.StringVal(c.color[paletteIntexInt32])
		}

		ctyV := cty.ObjectVal(attr)
		valueList = append(valueList, ctyV)
	}
	return cty.TupleVal(valueList)

}

func shortVizProc(chart *chart.Chart) cty.Value {
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
			attr["color"] = cty.StringVal(c.color[*v.PaletteIndex])
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
