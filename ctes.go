package main

import (
	"buscando/s3trail"
	"os"

	cli "github.com/jawher/mow.cli"
)

func main() {

	// r := s3trail.Request{
	// 	Bucket: "daidokoro-log",
	// 	Prefix: "AWSLogs",
	// 	Region: "eu-west-1",
	// 	URL:    "http://localhost:9200",
	// }
	//
	// r.Log()

	// Defining CLI options in main function
	args := make(map[string]*string)
	app := cli.App("translate", "Simple app for translating text via Google Translate API")
	args["prefix"] = app.StringOpt("p prefix", "AWSLogs", `S3 Object Prefix, useful for narrowing Cloudtrail searchess`)
	args["bucket"] = app.StringOpt("b bucket", "", `Cloudtrail S3 Bucket`)
	args["url"] = app.StringOpt("u url", "", `Elasticsearch URL`)

	app.Cmd.Action = func() {
		r := s3trail.Request{
			Bucket: *args["bucket"],
			Prefix: *args["prefix"],
			Region: "eu-west-1",
			URL:    *args["url"],
		}

		r.Log()

	}

	app.Run(os.Args)
}
