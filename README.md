# Ctes

Ctes is a simple compiled CLI tool for printing Cloudtrail logs from S3 or pushing them to an Elasticsearch endpoint.


## Objectives

The project was born of a very specific need to print and grep Cloudtrail logs, as well as push specific bits or logs to an Elasticsearch endpoint in a repeatable fashion.


## Requirements

- AWS Credentials/Config with sufficient privileges for accessing the CloudTrail S3 Bucket.

- An Elasticsearch 5.0 endpoint [_Optional_]


## Usage

```
$ ctes

Usage: ctes [OPTIONS] BUCKET

Simple CLI tool for printing CloudTrail Logs or Pushing them to Elasticsearch

Arguments:
  BUCKET=""    Cloudtrail S3 Bucket

Options:
  -p, --prefix="AWSLogs"   S3 Object Prefix, useful for narrowing Cloudtrail searches/results
  -u, --url=""             Elasticsearch URL, if not specified, results are printed to stdout
```


## Install

__If you've got Golang setup up, then:__

    go get -v github.com/daidokoro/ctes


__If not, here are some precompiled links:__

- [MacOS Darwin amd64]()
- [Linux x86]()
- [Linux x64]()

__Untested Builds!__

- [Windows x86]()
- [Windows x64]()
