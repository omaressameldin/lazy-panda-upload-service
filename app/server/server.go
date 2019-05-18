package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	v1 "github.com/omaressameldin/lazy-panda-upload-service/pkg/api/v1"

	"google.golang.org/grpc"
)

// RunServer runs service to publish User service
func RunServer(ctx context.Context, v1API v1.UploadServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	v1.RegisterUploadServiceServer(server, v1API)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down upload server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	log.Println("starting upload server...")
	return server.Serve(listen)
}
