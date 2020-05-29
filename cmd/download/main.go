package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-openapi/swag"
	"github.com/pottava/gcs-utils/log"
	cli "gopkg.in/alecthomas/kingpin.v2"
)

// for compile flags
var (
	ver    = "dev"
	commit string
	date   string
)

// Config is set of configurations
type Config struct { // nolint
	AppVersion string
	BucketName *string
	Object     *string
	OutputFile *string
	Timeout    *int64
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			log.Errors.Fatal(err)
		}
	}()

	app := cli.New("gcs-download", "downloads a file from Google Cloud Storage")
	if len(commit) > 0 && len(date) > 0 {
		app.Version(fmt.Sprintf("%s-%s (built at %s)", ver, commit, date))
	} else {
		app.Version(ver)
	}
	conf := &Config{}
	conf.AppVersion = ver
	conf.BucketName = app.Flag("bucket", "GCS bucket name.").
		Short('b').Envar("BUCKET_NAME").Required().String()
	conf.Object = app.Flag("object", "GCS object name.").
		Short('o').Envar("OBJECT_NAME").Required().String()
	conf.OutputFile = app.Flag("output", "Output file name.").
		Envar("OUTPUT_FILE").String()
	conf.Timeout = app.Flag("timeout", "Timeout seconds for the download.").
		Short('t').Envar("TIMEOUT").Default("30").Int64()

	cli.MustParse(app.Parse(os.Args[1:]))

	// Cancel
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(swag.Int64Value(conf.Timeout))*time.Second)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		os.Exit(1)
	}()

	// Get an object
	client, cErr := storage.NewClient(ctx)
	if cErr != nil {
		log.Errors.Fatal(cErr)
	}
	bucket := client.Bucket(swag.StringValue(conf.BucketName))
	object := bucket.Object(swag.StringValue(conf.Object))
	reader, rErr := object.NewReader(ctx)
	if rErr != nil {
		log.Errors.Fatal(rErr)
	}
	defer reader.Close()

	// Write the result
	if conf.OutputFile == nil || len(swag.StringValue(conf.OutputFile)) == 0 {
		io.Copy(os.Stdout, reader) // nolint
	} else {
		file, fErr := os.Create(swag.StringValue(conf.OutputFile))
		if fErr != nil {
			log.Errors.Fatal(fErr)
		}
		writer := bufio.NewWriter(file)
		io.Copy(writer, reader) // nolint

		if err := writer.Flush(); err != nil {
			log.Errors.Fatal(err)
		}
	}
}
