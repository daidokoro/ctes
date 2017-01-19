package main

import "github.com/spf13/cobra"

// Global variables for flags
var region string
var prefix string
var url string

func init() {
	// Initialise cobra stuff
	rootCmd.PersistentFlags().StringVarP(&prefix, "prefix", "p", "AWSLogs", "Prefix of the S3 Key, useful for narrowing searches and output")
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "Elasticsearch URL Endpoint")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "eu-west-1", "AWS Region where S3 Bucket resides")
}

// root command (calls all other commands)
var rootCmd = &cobra.Command{
	Use:   "ctes [s3 bucket]",
	Short: "Simple CLI tool for printing CloudTrail Logs or Pushing them to Elasticsearch",
	Run: func(cmd *cobra.Command, args []string) {

		// If no arguments are given
		if len(args) < 1 {
			cmd.Help()
			return
		}

		bucket := args[0]

		r := &Request{
			Bucket: bucket,
			Prefix: prefix,
			URL:    url,
			Region: region,
		}

		// Run Log Job
		r.Log()

	},
}
