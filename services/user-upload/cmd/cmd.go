package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	v1 "github.com/omaressameldin/lazy-panda-upload-service/core/pkg/service/v1"
	"github.com/omaressameldin/lazy-panda-upload-service/core/pkg/uploader"
	"github.com/omaressameldin/lazy-panda-upload-service/services/user-upload/server"
	"github.com/omaressameldin/lazy-panda-utils/app/pkg/firebase"
)

// Config is configuration for Server
type Config struct {
	Port   string
	Config string
	Bucket string
}

var v1API *v1.UploadServiceServer

// RunServer get the flags and starts the server
func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.Port, "port", "", "port to bind")
	flag.StringVar(&cfg.Config, "config", "", "firebase json config file")
	flag.StringVar(&cfg.Bucket, "bucket", "", "firebase storage bucket")
	flag.Parse()

	if len(cfg.Port) == 0 {
		return fmt.Errorf("invalid TCP port for server: '%s'", cfg.Port)
	}

	if len(cfg.Bucket) == 0 {
		return fmt.Errorf("invalid Bucket for firebase database: '%s'", cfg.Bucket)
	}

	_, err := os.Stat(cfg.Config)
	if os.IsNotExist(err) {
		return fmt.Errorf("File does not exist: '%s'", cfg.Config)
	}
	connector := initConnector(cfg.Config, cfg.Bucket)

	if len(cfg.Port) == 0 {
		return fmt.Errorf("invalid TCP port for server: '%s'", cfg.Port)
	}

	v1API = v1.NewUploadServiceServer(connector)

	return server.RunServer(ctx, v1API, cfg.Port)
}

// initConnector initializes database connector
func initConnector(firebaseConfig, bucket string) uploader.Uploader {
	connector, err := firebase.StartConnection(firebaseConfig, "", bucket)
	if err != nil {
		panic(err)
	}

	return connector
}

// CloseServer closes all connections such as database connection
func CloseServer() error {
	return nil
}
