package builder

// https://github.com/terraform-providers/terraform-provider-signalfx/blob/master/signalfx/resource_signalfx_time_chart.go

var c = struct {
	color map[int32]string
}{
	color: map[int32]string{
		0:  "gray",
		1:  "blue",
		2:  "azure",
		3:  "navy",
		4:  "brown",
		5:  "orange",
		6:  "yellow",
		7:  "magenta",
		8:  "purple",
		9:  "pink",
		10: "violet",
		11: "lilac",
		12: "iris",
		13: "emerald",
		14: "green",
		15: "aquamarine",
		16: "red",
		17: "yellow",
		18: "yellow",
		19: "green",
		20: "green",
		21: "gray",
	},
}

var vizcolor = struct {
	color map[int32]string
}{
	color: map[int32]string{
		0:  "gray",
		1:  "blue",
		2:  "azure",
		3:  "navy",
		4:  "brown",
		5:  "orange",
		6:  "yellow",
		7:  "pink",
		8:  "magenta",
		9:  "iris",
		10: "purple",
		11: "violet",
		12: "lilac",
		13: "emerald",
		14: "green",
		15: "aquamarine",
		16: "red",
		17: "gold",
		18: "greenyellow",
		19: "chartreuse",
		20: "jade",
	},
}

/*
var ChartColorsSlice = []chartColor{
	{"gray", "#999999"},
	{"blue", "#0077c2"},
	{"light_blue", "#00b9ff"},
	{"navy", "#6CA2B7"},
	{"dark_orange", "#b04600"},
	{"orange", "#f47e00"},
	{"dark_yellow", "#e5b312"},
	{"magenta", "#bd468d"},
	{"cerise", "#e9008a"},
	{"pink", "#ff8dd1"},
	{"violet", "#876ff3"},
	{"purple", "#a747ff"},
	{"gray_blue", "#ab99bc"},
	{"dark_green", "#007c1d"},
	{"green", "#05ce00"},
	{"aquamarine", "#0dba8f"},
	{"red", "#ea1849"},
	{"yellow", "#ea1849"},
	{"vivid_yellow", "#ea1849"},
	{"light_green", "#acef7f"},
	{"lime_green", "#6bd37e"},
}

*/

var t = struct {
	types map[string]string
}{
	types: map[string]string{
		"Heatmap":         "signalfx_heatmap_chart",
		"SingleValue":     "signalfx_single_value_chart",
		"TimeSeriesChart": "signalfx_time_chart",
		"List":            "signalfx_list_chart",
		"Text":            "signalfx_text_chart",
	},
}

type testRulesV1 struct {
	Above               string
	Duration            string
	Invalid             bool
	IsCustomizedMessage bool
	JobResolution       string
	Name                string
	Notifications       []map[string]interface{}
	Parameterized       string
	PercentOfDuration   int64
	Readable            string
	SeverityLevel       string
	ShowThreshold       bool
	TargetPlot          string
	ThresholdMode       string
	TriggerMode         string
	UniqueKey           int64
}

type RulesV1 struct {
	DetectLabel         string
	Disabled            bool
	IsCustomizedMessage bool
	Notifications       []map[string]string
	Parameterized       string
	Readable            string
	Severity            string
}

type DetectorV1 struct {
	Sf_description                     string
	Sf_createdOnMs                     int
	Sf_creator                         string
	Sf_currentJobIds                   []string
	Sf_detector                        string
	Sf_id                              string
	Sf_jobIdsHistory                   []map[string]interface{}
	Sf_jobLabelResolutions             []string
	Sf_jobMaxDelay                     int64
	Sf_labelResolutions                map[string]int64
	Sf_memberOf                        []string
	Sf_organizationID                  string
	Sf_overMTSLimit                    bool
	Sf_packageSpecifications           string
	Sf_programText                     string
	Sf_programs                        []string
	Sf_rules                           []RulesV1
	Sf_signalflowVersion               int64
	Sf_sourceSelectorEquivalentFilters [][]string
	Sf_sourceSelectors                 []string
	Sf_timezone                        string
	Sf_type                            string
	sf_uiModel                         interface{}
	Sf_updatedBy                       string
	Sf_updatedOnMs                     int64
}
