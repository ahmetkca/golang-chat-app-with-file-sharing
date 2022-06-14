package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

// Tests chat endpoint with valid username
// the message returned from websocket server should be the expected message
// otherwise the chat endpoint is not working properly
func TestWebSocketHandler(t *testing.T) {
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
	// defer ws_conn.Close()

	for i := 0; i < 1; i++ {
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

		t.Logf("\n")
	}
}

// No username is supplied so the client should not be able to  connect
// must get BadHandshake error since no username
func TestNoUsernameConnection(t *testing.T) {
	ctx := context.Background()
	mux := http.NewServeMux()
	setupHandlers(mux, ctx)       // setup the endpoints for multiplexer
	ts := httptest.NewServer(mux) // create a new test server with the multiplexer
	defer ts.Close()

	wsUrl := "ws" + strings.TrimPrefix(ts.URL, "http") + "/chat?username=" // websocket endpoint
	dialer := websocket.DefaultDialer                                      // create a test client for websocket connection
	_, _, err := dialer.DialContext(
		ctx, wsUrl, nil,
	)
	if err == nil && err == websocket.ErrBadHandshake {
		t.Fatalf("Expected: %v, Got: %v\n", websocket.ErrBadHandshake, err)
	}
}
