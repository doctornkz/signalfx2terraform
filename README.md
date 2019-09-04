# SignalFX2Terraform 

#### Simple converter hand-made dashboards to terraform configuration
Inspired by Ricardo Silveira(OLX) and Luis Rodriges(OLX)

#### How it works:
Application based on generic libraries from Hashicorp and Signalfx team:

[hclwrite](https://godoc.org/github.com/hashicorp/hcl2/hclwrite) - Package hclwrite deals with the problem of generating HCL configuration and of making specific surgical changes to existing HCL configurations.

[signalfx-go](https://github.com/signalfx/signalfx-go) - EXPERIMENTAL Go client library and instrumentation bindings for SignalFx

Converter uses SignalFX API Token for fetching dashboard/detector data, see details here: [API Docs](https://developers.signalfx.com/basics/basics_overview.html).
After that application unpacks JSON to Signalfx objects with `signalfx-go` library. We are using special `hclwriter` to convert objects to HashiCorp Language. Also we have ugly pack of wrappers for Type manipulation inside code. After `hclwriter` creates the terraform config file and fills it.      
#### Usage:
```
git clone  git@github.com:doctornkz/signalfx2terraform.git
cd signalfx2terraform
go build
./signalfx2terraform --help
Usage of ./signalfx2terraform:
  -dashboard string
        Dashboard ID, without URL
  -detector string
        Detector ID, without URL
  -token string
        SignalFX Token
```
You need to use correct API Token and your dashboard or detector ID. These IDs presented as part SignalFX's URL. 
Like that:\
`https://<REALM>.signalfx.com/#/dashboard/DxuFENBAAJI` > `DxuFENBAAJI`\
or\
`https://REALM.signalfx.com/#/detector/v2/D-9Usa2AIAA/` > `D-9Usa2AIAA`

Run command with parameters:
- For Dashboard
```
./signalfx2terraform --token=<TOKEN> --dashboard=<DASHBOARD_ID>
resource "signalfx_dashboard" "***" {
  dashboard_group = "***"
  name            = "test-Apache_Dashboard_tests_graph"
  description     = ""
  time_range      = "-30m"
  .....
```
- For Detector
```
./signalfx2terraform --token=<TOKEN> --detector=<DETECTOR_ID>
resource "signalfx_detector" "sfx-***" {
  name         = "test-Mysql-server: Memory-Utilization"
  description  = ""
  max_delay    = 30
  program_text = <<EOF
D = data('memory.utilization', filter=filter('application_name', 'mysql')).publish(label='D', enable=False)
detect(when(D > 90, lasting='10m')).publish('Memory-Utilization')
EOF
  rule {
    description   = "The value of memory.utilization is above 90 for 10m."
    severity      = "Critical"
    .....
```
You can reorder STDOUT to file for editing, checking or passing with terraform.

### You should know:
 - Work in progress, now covered only 40% of documented functionality, but people never use even 60%
 - Dashboard name renamed to `test-<Dashboard Name>`. It's hardcoded to prevent destroying original dashboard. The same with detectors.
 - Converter knows nothing about guidlines, vars, tags, etc. It gives you only Raw dashboard ot detector configuration.
 - Detectors v1 (handmade) badly JSON-structured, after exporting, please, check configuration (and fix if it needed). 

### TODO: 
 - Cover ~90% functionality
 - Tests 
 - Simplify code

### Help:
 - MR to repo
 - Bugs to issues

### Thanks to:
 - OLX Group, for motivation and opportunity to build exporter
 - SignalFX, for fresh monitoring and alerting system for SRE's
 - HashiCorp, for awesome resource management tool - terraform

### Special thanks to:
 - Cory, https://github.com/cory-signalfx, for fixing bugs in libraries.
