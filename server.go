package main

import (
	"context"
	"embed"
	"errors"
	"log"
	"mime"
	"net/http"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"
)

var filesBucket *blob.Bucket

const (
	bucketName = "files-bucket"
	location   = "local"
)

// Initialize the blob Need to create bucket before running
func initBlob(ctx context.Context) error {
	var err error
	filesBucket, err = blob.OpenBucket(
		ctx, "s3://"+bucketName+"?"+
			"endpoint=127.0.0.1:9000&"+
			"region="+location+"&"+
			"disableSSL=true&"+
			"s3ForcePathStyle=true",
	)
	return err
}

// Uploads data to the MinIO server specifically to the bucket 'files-bucket'
// returns:
//		-	(string) 	presigned URL to the file in the bucket
//		- 	(error)		err, if any.
func uploadFile(ctx context.Context, sessionId string, username string, data []byte, contentType string) (string, error) {
	// check the extension of supplied data
	extensions, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return "", errors.New("error obtaining file extensions for content type")
	}

	objectKey := sessionId + extensions[0]
	log.Printf("Uploading %s", objectKey)
	w, err := filesBucket.NewWriter(ctx, objectKey, nil)
	if err != nil {
		return "", err
	}
	nBytes, err := w.Write(data)
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", nil
	}
	log.Printf("Uploaded %s - %d bytes", objectKey, nBytes)
	return filesBucket.SignedURL(ctx, objectKey, nil)
}

//go:embed static
var staticFiles embed.FS

func setupHandlers(mux *http.ServeMux, ctx context.Context) *http.ServeMux {
	if mux == nil {
		mux = http.NewServeMux() // create a brand new multiplexer
	}
	var staticFS = http.FS(staticFiles) // creates a http.FileSystem out of embeded files of directory
	fs := http.FileServer(staticFS)     // creates a file server with the given FileSystem
	mux.Handle("/static/", fs)          // register FileServer at the /static/ url path

	mux.HandleFunc("/", index) // index HandlerFunc will handle root url path

	ctxHandler := &ContextAdapter{
		ctx:     ctx,
		handler: middleware(ContextHandlerFunc(loginAndUpgradeToWebSocket)),
	}

	mux.Handle("/chat", ctxHandler)
	return mux
}

func main() {
	ctx := context.Background()
	var err error

	// endpoint := "localhost:9000"
	// accessKeyId := "olan"
	// secretAccessKey := "ovWVGMHV5TE3bW3x"
	// useSSL := false
	// initGlobalMinioClient(endpoint, accessKeyId, secretAccessKey, useSSL)
	err = initBlob(ctx)
	if err != nil {
		log.Fatalf("error cannot initialize the MinIO server")
	}

	mux := http.NewServeMux() // create a brand new multiplexer

	setupHandlers(mux, ctx)

	server := http.Server{ // Server struct with the port 8080 and Handler as created mutiplexer
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	err = server.ListenAndServe() // Starts the server and listen on given port (8080 by default)
	if err != nil {
		server.ErrorLog.Panic("Error: cannot start the server on port 8080")
	}
}
