package handler

import (
   "fmt"
   "log"
   "net/http"

   "strings"

   "github.com/urfave/cli/v2"
)

var token string

// Webserver - creates a webserver
func Webserver(c *cli.Context){
   port := c.String("port")
   address := c.String("address")
   bind := address + ":" + port

   token = c.String("token")

   if token == "" {
      log.Fatal("No Signalfx Token provided")
   }

   http.HandleFunc("/", handleRoot)
   http.HandleFunc("/dashboard/", handler)
   http.HandleFunc("/detector/", handler)
   http.HandleFunc("/api/metrics", handleMetrics)

   fmt.Println("Starting server on " + bind)

   if err := http.ListenAndServe(bind, nil); err != nil {
      log.Fatal("Cannot bind to " + bind, err)
   }
}

// handleRoot - Handle root path
func handleRoot(w http.ResponseWriter, r *http.Request) {
   if r.Header.Get("User-Agent") == "ELB-HealthChecker/2.0" {
      return
   }
   fmt.Fprintf(w, "Check out signalfx2terroform repository README to know how to use this")
}

// handler - handler for detector and dashboard resources
func handler(w http.ResponseWriter, r *http.Request) {
   // TODO: Improve logging
   log.Printf("New request from <%s> and User-agent <%s> and URL: <%s>", r.Header.Get("X-Forwarded-For"), r.Header.Get("User-Agent"), r.URL.String())

   out, err := importResource(r.URL.String())

   if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintln(w, err)
   }

   fmt.Fprintln(w, out)
}

// importResource - returns the id for the resource to import
// returns string with output and error
func importResource(url string) (string, error) {
   split := strings.Split(url, "/")

   switch i := split[1]; i {
      case "dashboard":
         return string(dashboardProcessor(strings.Split(split[2], "?")[0], token)), nil
      case "detector":
         return string(detectorProcessor(strings.Split(split[3], "?")[0], token)), nil
      default:
         return "", fmt.Errorf("Cannot import %s", i)
   }
}

// handleMetrics - print out string "up 1"
// TODO: implement prometheus exporter for metrics
func handleMetrics(w http.ResponseWriter, r *http.Request) {
   fmt.Fprintf(w, "up 1")
   return
}
