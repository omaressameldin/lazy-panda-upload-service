package cmd

import (
	"context"
	"flag"
	"fmt"

	v1 "github.com/omaressameldin/lazy-panda-upload-service/pkg/service/v1"

	"github.com/omaressameldin/lazy-panda-upload-service/server"
)

// Config is configuration for Server
type Config struct {
	Port           string
	FirebaseConfig string
	Collection     string
}

var v1API *v1.UploadServiceServer

// RunServer get the flags and starts the server
func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.Port, "port", "", "port to bind")

	flag.Parse()

	if len(cfg.Port) == 0 {
		return fmt.Errorf("invalid TCP port for server: '%s'", cfg.Port)
	}

	v1API = v1.NewUploadServiceServer()

	return server.RunServer(ctx, v1API, cfg.Port)
}

// CloseServer closes all connections such as database connection
func CloseServer() error {
	return nil
}
