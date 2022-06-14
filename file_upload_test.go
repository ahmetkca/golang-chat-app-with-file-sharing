package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"gocloud.dev/blob/fileblob"
)

func setupTestBucket(filesDir string) error {
	u, err := url.Parse(fmt.Sprintf("file:///%s", filesDir))
	if err != nil {
		return err
	}

	opts := fileblob.Options{
		URLSigner: fileblob.NewURLSignerHMAC(
			u,
			[]byte("super secret password nobody knows"),
		),
		CreateDir: true,
	}
	filesBucket, err = fileblob.OpenBucket(filesDir, &opts)
	if err != nil {
		return err
	}
	return nil
}

func TestFileUploadTest(t *testing.T) {
	tmpDir, err := os.MkdirTemp("./", "*")
	if err != nil {
		t.Fatalf("error could not create a temporary directory")
	}
	setupTestBucket(tmpDir)
	ctx := context.Background()
	mux := http.NewServeMux()
	setupHandlers(mux, ctx)       // setup the endpoints for multiplexer
	ts := httptest.NewServer(mux) // create a new test server with the multiplexer
	defer ts.Close()

	wsUrl := "ws" + strings.TrimPrefix(ts.URL, "http") + "/chat?username=testuser" // websocket endpoint
	dialer := websocket.DefaultDialer                                              // create a test client for websocket connection
	ws_conn, _, err := dialer.DialContext(
		ctx, wsUrl, nil,
	)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	for i := 1; i < 2; i++ {
		_, p, err := ws_conn.ReadMessage()
		if err != nil {
			t.Fatalf("%v\n", err)
		}
		expectedFirstMessage := "Welcome to support. My name is Ahmet. How can I help you today?"
		t.Logf("\nExpected Message: \n%s", expectedFirstMessage)
		if len(string(p)) != len(expectedFirstMessage) {
			t.Fatalf("Expected first mesage from websocket: %v. Got: %v\n", expectedFirstMessage, string(p))
		}
		t.Logf("\nFirst Message received: \n%s", string(p))
		for x := 0; x < len(expectedFirstMessage); x++ {
			// t.Logf("%c", string(p)[x])
			if string(p)[x] != expectedFirstMessage[x] {
				t.Fatalf("Expected first mesage from websocket: %v. Got: %v\n", expectedFirstMessage, string(p))
			}
		}
	}

	for i := 1; i < 2; i++ {
		// osFile, err := os.OpenFile("test.pdf", os.O_RDONLY, os.ModeTemporary)
		// if err != nil {
		// 	t.Fatalf("error could not open the temporary pdf test file")
		// }

		var dta []byte
		// osFile.Read(dta)
		dta, err = os.ReadFile("./test.pdf")
		if err != nil {
			t.Fatalf("error could not read test pdf file's contents")
		}

		ws_conn.WriteMessage(websocket.BinaryMessage, dta)

		_, p, err := ws_conn.ReadMessage()

		if err != nil {
			t.Fatalf("error could not read message after uploading a valid file")
		}

		if !strings.Contains(string(p), "file:///") {
			t.Fatalf("Expected a signed url which should contains \"file:///\" but Got: %v", string(p))
		}
	}
}

func TestInvalidFileUploadTest(t *testing.T) {
	tmpDir, err := os.MkdirTemp("./", "*")
	if err != nil {
		t.Fatalf("error could not create a temporary directory")
	}
	setupTestBucket(tmpDir)
	ctx := context.Background()
	mux := http.NewServeMux()
	setupHandlers(mux, ctx)       // setup the endpoints for multiplexer
	ts := httptest.NewServer(mux) // create a new test server with the multiplexer
	defer ts.Close()

	wsUrl := "ws" + strings.TrimPrefix(ts.URL, "http") + "/chat?username=testuser" // websocket endpoint
	dialer := websocket.DefaultDialer                                              // create a test client for websocket connection
	ws_conn, _, err := dialer.DialContext(
		ctx, wsUrl, nil,
	)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	for i := 1; i < 2; i++ {
		_, p, err := ws_conn.ReadMessage()
		if err != nil {
			t.Fatalf("%v\n", err)
		}
		expectedFirstMessage := "Welcome to support. My name is Ahmet. How can I help you today?"
		t.Logf("\nExpected Message: \n%s", expectedFirstMessage)
		if len(string(p)) != len(expectedFirstMessage) {
			t.Fatalf("Expected first mesage from websocket: %v. Got: %v\n", expectedFirstMessage, string(p))
		}
		t.Logf("\nFirst Message received: \n%s", string(p))
		for x := 0; x < len(expectedFirstMessage); x++ {
			// t.Logf("%c", string(p)[x])
			if string(p)[x] != expectedFirstMessage[x] {
				t.Fatalf("Expected first mesage from websocket: %v. Got: %v\n", expectedFirstMessage, string(p))
			}
		}
	}

	for i := 1; i < 2; i++ {
		osFile, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			t.Fatalf("error could not open the temporary plain text test file")
		}
		nBytes, err := osFile.WriteString("invalid test data")
		if err != nil || nBytes == 0 {
			t.Fatalf("%v", err)
		}

		var dta []byte

		// s := "invalid test file content"
		// dta = []byte(s)
		osFile.Read(dta)
		// dta, err = os.ReadFile("./test.pdf")
		if err != nil {
			t.Fatalf("error could not read test pdf file's contents")
		}

		ws_conn.WriteMessage(websocket.BinaryMessage, dta)

		_, p, err := ws_conn.ReadMessage()

		if err != nil {
			t.Fatalf("error could not read message after uploading a valid file")
		}
		t.Logf("\nMessage received back from websocket-server: \n%s", string(p))
		if !strings.Contains(string(p), "Invalid data type received only (pdf, png, or jpeg) allowed") {
			t.Fatalf("Expected a error message saying it is invalid file type but Got: %v", string(p))
		}

	}
}
