package main

import (
	"ctes/cmd"
	"os"

	cli "github.com/jawher/mow.cli"
)

func main() {

	args := make(map[string]*string)
	app := cli.App("ctes", "Simple CLI tool for printing CloudTrail Logs or Pushing them to Elasticsearch")
	args["bucket"] = app.StringArg("BUCKET", "", `Cloudtrail S3 Bucket`)
	args["prefix"] = app.StringOpt("p prefix", "AWSLogs", `S3 Object Prefix, useful for narrowing Cloudtrail searches/results`)
	args["url"] = app.StringOpt("u url", "", `Elasticsearch URL, if not specified, results are printed to stdout`)
	args["region"] = app.StringOpt("r region", "eu-west-1", `AWS Region`)

	app.Spec = "BUCKET [OPTIONS]"

	app.Cmd.Action = func() {
		r := cmd.Request{
			Bucket: *args["bucket"],
			Prefix: *args["prefix"],
			URL:    *args["url"],
			Region: *args["region"],
		}

		r.Log()

	}

	app.Run(os.Args)
}
