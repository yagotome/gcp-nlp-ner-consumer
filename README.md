# GCP NLP NER Consumer

This is an application to read a CSV of texts to have entities extracted and request Google Cloud Natural Language to retrieve Named Entities and store them in a TSV format (as format expected by Stanford CoreNLP)


## Pre-requisites

You'll need Go v1.11 to run.


## Running

To run project, first set environment variable GOOGLE_APPLICATION_CREDENTIALS to credentials path:

```
$ export GOOGLE_APPLICATION_CREDENTIALS=<your-credentials-path>
```

Then, just run the main: (if source is in you `GOPATH`, set `GO111MODULE` environment variable first)

```
$ go run cmd/createtsv/main.go
```