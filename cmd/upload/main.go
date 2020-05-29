package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
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
	InputFile  *string
	Object     *string
	Timeout    *int64
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			log.Errors.Fatal(err)
		}
	}()

	app := cli.New("gcs-upload", "uploads a file from Google Cloud Storage")
	if len(commit) > 0 && len(date) > 0 {
		app.Version(fmt.Sprintf("%s-%s (built at %s)", ver, commit, date))
	} else {
		app.Version(ver)
	}
	conf := &Config{}
	conf.AppVersion = ver
	conf.BucketName = app.Flag("bucket", "GCS bucket name.").
		Short('b').Envar("BUCKET_NAME").Required().String()
	conf.InputFile = app.Flag("input", "Input file name.").
		Short('i').Envar("INPUT_FILE").Required().String()
	conf.Object = app.Flag("object", "GCS object name.").
		Short('o').Envar("OBJECT_NAME").String()
	conf.Timeout = app.Flag("timeout", "Timeout seconds for the upload.").
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

	// Open a file
	file, fErr := os.Open(swag.StringValue(conf.InputFile))
	if fErr != nil {
		log.Errors.Fatal(fErr)
	}
	defer file.Close()

	// Create an object
	if conf.Object == nil || len(swag.StringValue(conf.Object)) == 0 {
		_, name := filepath.Split(swag.StringValue(conf.InputFile))
		conf.Object = swag.String(name)
	}
	client, cErr := storage.NewClient(ctx)
	if cErr != nil {
		log.Errors.Fatal(cErr)
	}
	bucket := client.Bucket(swag.StringValue(conf.BucketName))
	writer := bucket.Object(swag.StringValue(conf.Object)).NewWriter(ctx)
	defer writer.Close()

	if _, err := io.Copy(writer, file); err != nil {
		log.Errors.Fatal(err)
	}
}
