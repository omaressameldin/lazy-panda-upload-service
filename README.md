# lazy-panda-upload-service
Lazy Panda is a project aimed on managing employee time for consulting companies.

## What this is
- A monorepo of services to upload files for **Lazy Panda** services using Go, [GRPC](https://github.com/grpc/grpc-go)
- the services use [GRPC streams](https://ops.tips/blog/sending-files-via-grpc/) to upload files via GRPC
- it's using 🔥 Firebase storage for uploading files but can use any storage that extends the [uploader package](https://github.com/omaressameldin/lazy-panda-upload-service/tree/master/core/pkg/uploader)
- the `core` directory contains the required packages for uploading and any shared upload code
- the `services` directoy contains separate services for each lazy panda service that needs to upload files (currently only `user-service` need to upload)

## How to run
- make sure you have **docker version: 18.x+** installed
- make sure firebase config is present in each of the services [follow this tutorial to get firebase config](https://www.youtube.com/watch?v=9rN29jENirI)
- run `docker-compose up --build` to launch service
- the **user upload service** will be available at port `6501`

## Taking the service for a spin
**Note1:** Please, make sure that the service is running before testing any of the following snippets

**Note2:** The service is not complete so no kind of authorization is implemented!

**Note3:** To be able to test the service you need to have a file to upload!
- you can check [.core/proto-gen/proto/v1/uploador.proto](.core/proto-gen/proto/v1/uploador.proto) file for available fields

- To connect to service:
```golang
import (
	"log"
	v1 "github.com/omaressameldin/lazy-panda-upload-service/core/pkg/api/v1"
	"google.golang.org/grpc"
)
func main() {
	client, err := grpc.Dial("localhost:6501", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v", err)
	}
  c := v1.NewUploadServiceClient(client)
}
```

- To upload a file
```golang
import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	v1 "github.com/omaressameldin/lazy-panda-upload-service/core/pkg/api/v1"
	"google.golang.org/grpc"
)
func main() {
  // Open file
	file, err := os.Open("./[file_path]")
	if err != nil {
		log.Println(err)
	}

// start upload
	stream, err := c.Upload(context.Background())
	if err != nil {
		log.Println(err)
	}
	buf := make([]byte, 1024)

// upload file metadata first
	meta := v1.UploadFileRequest_Meta{
		Meta: &v1.Metadata{
			FileName: "test",
			FileType: "[file_extension]",
		},
	}
	stream.Send(&v1.UploadFileRequest{FileData: &meta})

// start writing file to buffer
	writing := true
	for writing {
		n, err := file.Read(buf)

		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}

			err = fmt.Errorf("errored while copying from file to buf %v", err)
			log.Println(err)
			return
		}
    // send buffer through stream
		stream.Send(&v1.UploadFileRequest{
			FileData: &v1.UploadFileRequest_Content{Content: buf[:n]},
		})
	}

  // close stream and log upload status
	status, err := stream.CloseAndRecv()
	log.Printf("%v, %v", status, err)
}
```


## Technologies used
- Golang
- GRPC
- firebase
- Docker
- Docker-compose
