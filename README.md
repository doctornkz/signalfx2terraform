# SignalFX2Terraform

## Simple converter hand-made dashboards to terraform configuration
Inspired by Ricardo Silveira(OLX) and Luis Rodriges(OLX)

## How it works:
Application based on generic libraries from Hashicorp and Signalfx (AKA SFX) team:

[hclwrite](https://godoc.org/github.com/hashicorp/hcl2/hclwrite) - Package hclwrite deals with the problem of generating HCL configuration and of making specific surgical changes to existing HCL configurations.

[signalfx-go](https://github.com/signalfx/signalfx-go) - EXPERIMENTAL Go client library and instrumentation bindings for SignalFx

Converter uses SignalFX API Token for fetching dashboard/detector data, see details here: [API Docs](https://developers.signalfx.com/basics/basics_overview.html).
After that application unpacks JSON to SFX objects with `signalfx-go` library. We are using special `hclwriter` to convert objects to HashiCorp Language. Also we have ugly pack of wrappers for Type manipulation inside code. After `hclwriter` creates the terraform config file and fills it.
## Usage:
```
git clone  git@github.com:doctornkz/signalfx2terraform.git
cd signalfx2terraform
make mod && make build
./bin/signalfx2terraform help
NAME:
   signalfx2terraform - Signalfx to Terraform converter cli

USAGE:
   signalfx2terraform [global options] command [command options] [arguments...]

COMMANDS:
   import     Import signalfx resources
   webserver  Create webserver to interact with signalfx resources
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

###### Legend
`import` - use this subcommand to import SFX resources\
`webserver` - will create a local webserver for you to see the terraform file

#### Import
This subcommand will import the SFX resource and translate it to terraform code
To see the description of it you can execute
```
./bin/signalfx2terraform import --help
NAME:
   signalfx2terraform import - Import signalfx resources

USAGE:
   signalfx2terraform import [command options] [arguments...]

OPTIONS:
   --token value, -t value      Signalfx token
   --dashboard value, -d value  Signalfx dashboard id
   --detector value, -x value   Signalfx detector id
   --help, -h                   show help (default: false)
```

You need to use correct SFX API Token and your dashboard or detector ID. These IDs presented as part SFX URL.

Like that:\
`https://<REALM>.signalfx.com/#/dashboard/DxuFENBAAJI` > `DxuFENBAAJI`\
or\
`https://REALM.signalfx.com/#/detector/v2/D-9Usa2AIAA/` > `D-9Usa2AIAA`

Run command with parameters:

- For **Dashboard**:

```
./bin/signalfx2terraform import -t <TOKEN> -d <DASHBOARD_ID>
resource "signalfx_dashboard" "***" {
  dashboard_group = "***"
  name            = "test-Apache_Dashboard_tests_graph"
  description     = ""
  time_range      = "-30m"
  .....
```
- For **Detector**:

```
./bin/signalfx2terraform -t <TOKEN> -x <DETECTOR_ID>
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

###### Legend
If you installed this cli by using `make install` then use your path instead of `./bin/signalfx2terraform`

#### Webserver
This subcommand will create local webserver in order to see the whole SFX resource translated to terraform code
To see the description of it you can execute

```
./bin/signalfx2terraform webserver --help
NAME:
   signalfx2terraform webserver - Create webserver to interact with signalfx resources

USAGE:
   signalfx2terraform webserver [command options] [arguments...]

OPTIONS:
   --port value, -p value     Webserver port to bind (default: 8080) [$PORT]
   --address value, -a value  Webserver address to use (default: localhost) [$ADDRESS]
   --token value, -t value    Signalfx token [$SIGNALFX_TOKEN]
   --help, -h                 show help (default: false)
```

You need to use correct SFX API Token. These IDs presented as part SFX URL.
Like that:\
`https://<REALM>.signalfx.com/#/dashboard/DxuFENBAAJI` > `DxuFENBAAJI`\
or\
`https://REALM.signalfx.com/#/detector/v2/D-9Usa2AIAA/` > `D-9Usa2AIAA`

By default you don't need to change the port and local address to create the webserver. To create the webserver execute
```
./bin/signalfx2terraform import -t <TOKEN>
Starting server on localhost:8080
```

Then go to your browser to `http://localhost:8080/`.
To import the SFX resources and visualize it in this webserver, grab your SFX resource and change the it like this:

- from: `https://signalfx.com/#/dashboard/`
- to:   `http://localhost:8080/#/dashboard/`

### You should know:
 - Work in progress, now covered only 70% of documented functionality
 - Dashboard name renamed to `test-<Dashboard Name>`. It's hardcoded to prevent destroying original dashboard. The same with detectors.
 - There can be some bugs in colors, not imported properly. SFX only recently standartized it. This will be fixed also soon.

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
