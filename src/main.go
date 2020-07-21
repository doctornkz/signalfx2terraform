package main

//https://godoc.org/github.com/hashicorp/hcl2/hclwrite#Block.BuildTokens

import (
   "log"
   "os"

   "github.com/urfave/cli/v2"
	"github.com/doctornkz/signalfx2terraform/src/handler"
)

var (
	signalfxAPIToken string
	dashboardID      string
	detectorID       string
	version          string
)


func main() {
   app := &cli.App{
      Usage: "Signalfx to Terraform converter cli",
      EnableBashCompletion: true,
      Commands: []*cli.Command{
         {
            Name: "import",
            Usage: "Import signalfx resources",
            Flags: []cli.Flag{
               &cli.StringFlag{
                 Name: "token",
                 Aliases: []string{"t"},
                 Usage: "Signalfx token",
                 Required: true,
               },
               &cli.StringFlag{
                  Name: "dashboard",
                  Aliases: []string{"d"},
                  Usage: "Signalfx dashboard id",
               },
               &cli.StringFlag{
                  Name: "detector",
                  Usage: "Signalfx detector id",
                  Aliases: []string{"x"},
               },
            },
            Action: func(c *cli.Context) error {
               handler.Import(c)
               return nil
            },
         },
         {
            Name: "webserver",
            Usage: "Create webserver to interact with signalfx resources",
            Flags: []cli.Flag{
               &cli.IntFlag{
                 Name: "port",
                 Aliases: []string{"p"},
                 Usage: "Webserver port to bind",
                 Value: 8080,
                 DefaultText: "8080",
               },
               &cli.StringFlag{
                  Name: "address",
                  Aliases: []string{"a"},
                  Usage: "Webserver address to use",
                  Value: "localhost",
                  DefaultText: "localhost",
               },
            },
            Action: func(c *cli.Context) error {
               handler.Webserver()
               return nil
            },
         },
      },
   }

   err := app.Run(os.Args)
   if err != nil {
      log.Fatal(err)
   }
}
