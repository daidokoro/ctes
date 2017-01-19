# Ctes

Ctes is a simple compiled CLI tool for printing Cloudtrail logs from S3 or pushing them to an Elasticsearch endpoint.


## Objectives

The project was born of a very specific need to print and grep Cloudtrail logs, as well as push specific bits or logs to an Elasticsearch endpoint in a repeatable fashion.


## Requirements

- AWS Credentials/Config with sufficient privileges for accessing the CloudTrail S3 Bucket.

- An Elasticsearch 5.0 endpoint [_Optional_]


## Usage

```sh
$ ctes

Simple CLI tool for printing CloudTrail Logs or Pushing them to Elasticsearch

Usage:
  ctes [s3 bucket] [flags]

Flags:
  -p, --prefix string   Prefix of the S3 Key, useful for narrowing searches and output (default "AWSLogs")
  -r, --region string   AWS Region where S3 Bucket resides (default "eu-west-1")
  -u, --url string      Elasticsearch URL Endpoint

```

## Example

![Alt text](demo.gif?raw=true "Demo")


## Install

__If you've got Golang setup up, then:__

    go get -v github.com/daidokoro/ctes


Compiled Binaries available [here!](https://github.com/daidokoro/ctes/releases)
