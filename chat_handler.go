package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var allowedContentTypes = []string{"application/pdf", "image/png", "image/jpeg"}

func isValidData(data []byte) (string, bool) {
	contentType := http.DetectContentType(data)
	log.Printf("Detected content type: %s", contentType)

	for _, act := range allowedContentTypes {
		if contentType == act {
			return contentType, true
		}
	}
	return "", false
}

type Username string

func loginAndUpgradeToWebSocket(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")
	log.Printf("username = %v\n", username)
	if len(username) == 0 || username == "" {
		log.Println("Error: username is not supplied")
		http.Error(w, "Error: username is not supplied", http.StatusBadRequest)
		// err = ws_conn.Close()
		// if err != nil {
		// 	log.Println(err)
		// }
		http.Error(w, "Unable to verify login details", http.StatusInternalServerError)
		return
	}

	ws_conn, err := upgrader.Upgrade(w, r, nil) // Try to establish a WebSocket Connection with the client Web Browser
	if err != nil {                             // non-nil err no WebSocket connection can be establish
		http.Error(w, "Error: WebSocket connection cannot be established", http.StatusInternalServerError)
		return
	}

	sessionId := base64.StdEncoding.EncodeToString([]byte(username))
	chats := supportChats.Read(sessionId)
	if len(chats) == 0 {
		log.Printf("No chat history for username=%s found", username)
	} else {
		if err := ws_conn.WriteMessage(1, []byte(strings.Join(chats, "<p>"))); err != nil {
			log.Println("Error sending chat history")
		}
	}

	const welcomeMessage = "Welcome to support. My name is Ahmet. How can I help you today?"
	if err := ws_conn.WriteMessage( // WebSocket connection established send back message
		1, []byte(welcomeMessage),
	); err != nil {
		defer ws_conn.Close() // if there is erro while sending message back close the WebSocket connection
		http.Error(w, "Error: there was an error while send data back", http.StatusInternalServerError)
		return
	}

	for {

		messageType, p, err := ws_conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage:
			log.Printf("Message received: %s\n", string(p))

			supportChats.Write(sessionId, string(p))

			if err != nil {
				ws_conn.Close()
				// http.Error(w, "Error: There was an error while reading messages", http.StatusInternalServerError)
				return
			}
			capitalized := strings.ToUpper(string(p))
			if err := ws_conn.WriteMessage(messageType, []byte(capitalized)); err != nil {
				ws_conn.Close()
				// http.Error(w, "Error: cannot write message back to the client", http.StatusInternalServerError)
				return
			}
		case websocket.BinaryMessage:
			//TODO: check the type of the file content (PDF, PNG, or, JPEG etc.)
			contentType, isValid := isValidData(p)
			if isValid && len(contentType) > 0 {
				//TODO: upload the file to the files-bucket of MinIO
				uploadedFileUrl, err := uploadFile(ctx, sessionId, username, p, contentType)
				if err != nil {
					log.Println(err)
					ws_conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				} else {
					ws_conn.WriteMessage(websocket.TextMessage, []byte(uploadedFileUrl))
				}
			} else {
				log.Printf("invalid data type submitted only (pdf, png, or jpeg) allowed")
				response := fmt.Sprintf("Invalid data type received only (pdf, png, or jpeg) allowed")
				ws_conn.WriteMessage(websocket.TextMessage, []byte(response))
			}
		default:
			response := fmt.Sprintf("Invalid data received: %d", messageType)
			ws_conn.WriteMessage(websocket.TextMessage, []byte(response))
		}
	}

	err = ws_conn.Close()
	if err != nil {
		log.Println(err)
	}
}
