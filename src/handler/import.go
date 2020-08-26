package handler

import (
   "fmt"
   "log"
   "github.com/hashicorp/hcl2/hclwrite"
   "github.com/signalfx/signalfx-go"
   "github.com/urfave/cli/v2"

   "github.com/doctornkz/signalfx2terraform/src/utils"
   "github.com/doctornkz/signalfx2terraform/src/detectors"
   "github.com/doctornkz/signalfx2terraform/src/timeseries"
   "github.com/doctornkz/signalfx2terraform/src/list"
   "github.com/doctornkz/signalfx2terraform/src/heatmap"
   "github.com/doctornkz/signalfx2terraform/src/singlevalue"
   "github.com/doctornkz/signalfx2terraform/src/text"
)

const (
   // APIURL : default entrypoint for customer's requests
   APIURL = "https://api.eu0.signalfx.com"
)

// Import - import signalfx resource
func Import(c *cli.Context){
   token := c.String("token")

   if c.IsSet("dashboard") {
      if dId := c.String("dashboard"); dId != "" {
         fmt.Printf("%s",dashboardProcessor(dId, token))
      } else {
         log.Fatal("Dashboard Id not specified")
      }
   }

   if c.IsSet("detector") {
      if dId := c.String("detector"); dId != "" {
         fmt.Printf("%s",detectorProcessor(dId, token))
      } else {
         log.Fatal("Detector Id not specified")
      }
   }
}

// dashboardProcessor - process dashboard import
func dashboardProcessor(d string, t string) []byte {

   client, err := signalfx.NewClient(t, signalfx.APIUrl(APIURL))

   if err != nil {
      log.Fatal("Something wrong with API client")
   }
   dashboard, err := client.GetDashboard(d)

   if err != nil {
      log.Printf("Dashboard error: %v", err)
      log.Fatal("Can't fetch dashboard")
   }
   charts := dashboard.Charts

   f := hclwrite.NewEmptyFile()

   utils.CreateDashboard(f, dashboard, client)

   for _, v := range charts {
      chart, err := client.GetChart(v.ChartId)

      if err != nil {
         log.Printf("Chart error: %v", err)
         log.Fatal("Can't get chart")
      }

      switch types := chart.Options.Type; types {
      case "SingleValue":
         singlevalue.Chart(f, chart)
      case "Heatmap":
         heatmap.Chart(f, chart)
      case "TimeSeriesChart":
         timeseries.Chart(f, chart)
      case "List":
         list.Chart(f, chart)

      case "Text":
         text.Chart(f, chart)
      default:
         continue
      }

   }
   return f.Bytes()
}

// detectorProcessor - process detector import
func detectorProcessor(d string, t string) []byte {
   client, err := signalfx.NewClient(t, signalfx.APIUrl(APIURL))

   if err != nil {
      log.Fatal("Something wrong with API client")
   }

   detector, err := client.GetDetector(d)

   if err != nil {
      // Failover for V1 detector
      // TODO: Implement reliable check
      //log.Printf("Detector error: %v", err)
      //log.Println("Can't fetch detector with V2 API, trying failover method...")
      f := hclwrite.NewEmptyFile()

      detectors.CreateDetectorV1(f, APIURL, d, t)

      return f.Bytes()
   } else {
      f := hclwrite.NewEmptyFile()
      detectors.CreateDetector(f, detector)

      return f.Bytes()
   }
}
