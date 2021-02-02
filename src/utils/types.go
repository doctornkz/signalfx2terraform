package utils

// https://github.com/terraform-providers/terraform-provider-signalfx/blob/master/signalfx/resource_signalfx_time_chart.go

var Color = map[int32]string{
   0: "gray",
   1: "blue",
   2: "azure",
   3: "navy",
   4: "brown",
   5: "orange",
   6: "yellow",
   7: "iris",
   8: "magenta",
   9: "pink",
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

// SecondaryVisualization - color_scale struct
type SecondaryVisualization struct {
   Gt            *float32          `json:"gt,omitempty"`
   Gte           *float32          `json:"gte,omitempty"`
   Lt            *float32          `json:"lt,omitempty"`
   Lte           *float32          `json:"lte,omitempty"`
   PaletteIndex  *int32            `json:"color"`
}

// PublishLabelOptions - viz_options structure
type PublishLabelOptions struct {
   DisplayName    string `json:"display_name,omitempty"`
   Label          string `json:"label"`
   PaletteIndex   *int32 `json:"color,omitempty"`
   PlotType       string `json:"plot_type,omitempty"`
   ValuePrefix    string `json:"value_prefix,omitempty"`
   ValueSuffix    string `json:"value_suffix,omitempty"`
   ValueUnit      string `json:"value_unit,omitempty"`
   YAxis          int32  `json:"axis,omitempty"`
}

//EventPublishLabelOptions - event_options structure
type EventPublishLabelOptions struct {
   DisplayName  string              `json:"display_name,omitempty"`
   Label        string              `json:"label,omitempty"`
   PaletteIndex *int32              `json:"color,omitempty"`
   Color        *map[int32]string   `json:"color_structure"`
}

// OptionsColor - viz_options, event_options, color_scale colors
var OptionsColor = map[int32]string{
   0: "gray",
   1: "blue",
   2: "azure",
   3: "navy",
   4: "brown",
   5: "orange",
   6: "yellow",
   7: "iris",
   8: "magenta",
   9: "pink",
   10: "purple",
   11: "violet",
   12: "lilac",
   13: "emerald",
   14: "green",
   15: "aquamarine",
}

var ColorScalePallete = map[int32]string {
   0: "gray",
   1: "blue",
   2: "light_blue",
   3: "navy",
   4: "dark_orange",
   5: "orange",
   6: "dark_yellow",
   7: "magenta",
   8: "cerise",
   9: "pink",
   10: "violet",
   11: "purple",
   12: "gray_blue",
   13: "dark_green",
   14: "green",
   15: "aquamarine",
   16: "red",
   17: "yellow",
   18: "vivid_yellow",
   19: "light_green",
   20: "lime_green",
}

// ColorRangeOptions - heatmap chart color struct
type ColorRangeOptions struct {
   Color string  `json:"color,omitempty,default:#05ceI00"`
   Max   float64 `json:"max_value,omitempty"`
   Min   float64 `json:"min_value,omitempty"`
}

// ValueColor - color pallete for Singlevalue and Heatmap charts
var ValueColor = map[int32]string {
   0: "gray",
   1: "blue",
   2: "navy",
   3: "orange",
   4: "yellow",
   5: "magenta",
   6: "purple",
   7: "violet",
   8: "lilac",
   9: "green",
   10: "aquamarine",
}

var Type = map[string]string{
   "Heatmap":         "signalfx_heatmap_chart",
   "SingleValue":     "signalfx_single_value_chart",
   "TimeSeriesChart": "signalfx_time_chart",
   "List":            "signalfx_list_chart",
   "Text":            "signalfx_text_chart",
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

// argparse - argument parser struct
var argparse struct {
   Command  string   `docopt:"<cmd>"`
   Tries    int      `docopt:"-n"`
   Force    bool     // Gets the value of --force
}

// DetectorRule structure
type DetectorRule struct {
   Description string `json:"description,omitempty"`
   DetectLabel string `json:"detectLabel,omitempty"`
   Disabled bool `json:"disabled,omitempty"`
   ParameterizedBody string `json:"parameterizedBody,omitempty"`
   ParameterizedSubject string `json:"parameterizedSubject,omitempty"`
   RunbookUrl string   `json:"runbookUrl,omitempty"`
   Severity   string `json:"severity,omitempty"`
   Tip        string   `json:"tip,omitempty"`
}

