package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	elastic "gopkg.in/olivere/elastic.v5"
	elogrus "gopkg.in/sohlich/elogrus.v2"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Request - Type used for handling requests to S3, Elasticsearch and Logging
type Request struct {
	URL    string
	Bucket string
	Index  string
	Prefix string
	Region string
}

func s3log(d []*map[string]string, log *logrus.Logger) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	for _, item := range d {
		o := make(map[string]interface{})

		for k, v := range *item {
			o[k] = v
		}

		log.WithFields(o).Info("Cloudtrail Log")
	}
}

// addHook - Adds the Elasticsearch hook/endpoint to the logger
func (r *Request) addHook(log *logrus.Logger) {
	c, err := elastic.NewClient(elastic.SetURL(r.URL))
	if err != nil {
		log.Panic(err)
	}
	hostname, _ := os.Hostname()
	hook, err := elogrus.NewElasticHook(c, hostname, logrus.DebugLevel, "cloudtrail")
	if err != nil {
		log.Panic(err)
	}

	log.Hooks.Add(hook)
}

// flatten - Recursively traverses a given map[string]interface{} and creates a flat version of the same, i.e No sub-branches
func (r *Request) flatten(m map[string]interface{}, o map[string]string, key string) {
	// takes a map[string]interface{} from arbitrary json output and flattens to a map
	for k, v := range m {
		if val, ok := v.(map[string]interface{}); ok {
			r.flatten(val, o, k+".")
		} else {
			// fmt.Printf("%s%s --> %s"+"\n", key, k, v)
			mapKey := key + k
			if _, ok := v.(string); ok {
				o[mapKey] = v.(string)
			}
		}

	}
}

// getResp - Reads the Body of an http response and gives me a string
func (r *Request) getResp(b io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(b)
	return buf.String()
}

// s3List - populates a channel with objects listed from S3
func (r *Request) s3List(c chan string) {
	todo := make(chan []*s3.Object)
	done := make(chan bool)

	svc := s3.New(session.New(), &aws.Config{Region: aws.String(r.Region)})
	if svc == nil {
		fmt.Println("s3list: missing s3 client")
	}

	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(r.Bucket), // Required
		Prefix: aws.String(r.Prefix),
	}

	// Feed Jobs
	go func() {
		//TODO: Populate Job Channel
		for {
			resp, err := svc.ListObjectsV2(params)
			if err != nil {
				fmt.Println(err.Error())
			}

			// Populate todos
			todo <- resp.Contents

			// Continue
			if *resp.IsTruncated {
				params.ContinuationToken = resp.NextContinuationToken
				continue
			}

			break

		}
		close(todo)
	}()

	// Execute Jobs
	go func() {
		for {
			j, ok := <-todo
			if ok {
				//TODO: Do something with Items from Job Queue
				for _, obj := range j {
					if strings.Contains(*obj.Key, ".gz") && !strings.Contains(*obj.Key, "Digest") {
						// fmt.Println("add")
						c <- *obj.Key // Pass Key to channel (sort of like yield in python)

						// log := fmt.Sprintf("INFO: Found %d keys [%s]", *f, r.Prefix)
					}
				}

			} else {
				done <- true
				return
			}
		}
	}()

	// wait
	<-done
}

// getRecords - Creates an array of individual cloudtrail records from a group of same.
func (r *Request) getRecords(key string) ([]*map[string]string, string) {
	svc := s3.New(session.New(), &aws.Config{Region: aws.String("eu-west-1")})
	params := &s3.GetObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(key),
	}

	// Create a map to store results
	result := make([]*map[string]string, 0)

	resp, err := svc.GetObject(params)

	if err != nil {
		return result, err.Error()
	}

	// Unbuffer response body
	str := r.getResp(resp.Body)

	// Create byte string array
	b := []byte(str)

	// Create empty interface for arbitrary data
	var f interface{}

	if err := json.Unmarshal(b, &f); err != nil {
		fmt.Println(err.Error())
	}

	// Type assertion
	m := f.(map[string]interface{})

	// Use case statement to identify underlying type in arbitrary interface{}
	for _, v := range m {
		switch vv := v.(type) {
		case []interface{}:
			// iterate interface array
			for _, u := range vv {
				j := u.(map[string]interface{})
				o := make(map[string]string) // For storing result maps
				r.flatten(j, o, "")          // calling my flatten function
				result = append(result, &o)
			}
		}
	}
	return result, ""
}

//Log - Prints CloudTrail logs
func (r *Request) Log() {
	jobs := make(chan string)
	done := make(chan bool)

	//Define Logger
	log := logrus.New()
	if r.URL != "" {
		r.addHook(log)
	}

	// Feed Jobs
	go func() {
		r.s3List(jobs)
		close(jobs)
	}()

	// Execute Jobs
	go func() {
		for {
			k, ok := <-jobs
			if ok {
				o, err := r.getRecords(k)

				// fmt.Println(o)
				s3log(o, log)
				if err != "" {
					// fmt.Printf("ERROR: %s\n", err)
					continue
				}

			} else {

				done <- true
				return
			}
		}
	}()

	// wait
	<-done
}
