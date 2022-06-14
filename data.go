package main

import "sync"

// server-wide global in-memory store
var supportChats = chatHistory{
	mutex: sync.Mutex{},
	data:  make(map[string][]string),
}

// Contains all of the users chat history data
type chatHistory struct {
	mutex sync.Mutex
	data  map[string][]string
}

func (ch *chatHistory) Read(sessionId string) []string {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	return ch.data[sessionId]
}

func (ch *chatHistory) Write(sessionId string, chatData string) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	ch.data[sessionId] = append(ch.data[sessionId], chatData)
}
