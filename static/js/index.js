// borrowed from https://github.com/mdn/samples-server/tree/master/s/websocket-chat
// and changed to only include relevant functionality
"use strict";
var connection = null;
var clientID = 0;

function connect() {
  var serverUrl;
  var scheme = "ws";

  var username = document.getElementById("username").value;
  serverUrl = scheme + "://" + document.location.hostname + ":8080" + "/chat" + "?username=" + username;

  connection = new WebSocket(serverUrl);
  console.log("***CREATED WEBSOCKET");

  connection.onopen = function(evt) {    
    document.getElementById("text").disabled = false;
    document.getElementById("send").disabled = false;
  };
  
  connection.onmessage = function(evt) {
    var f = document.getElementById("chatbox").contentDocument;    
    var msg = evt.data;
    console.log("Message received: ");
    console.dir(msg);    
    f.write('<p>'+ msg);
  };  
}

function send() {  
  var msg = document.getElementById("text").value;
  connection.send(msg);
  document.getElementById("text").value = "";
}

function sendFile() {
  var file = document.getElementById('filename').files[0];    
    var reader = new FileReader();
    var rawData = new Blob();
    reader.loadend = function() {
    }
    reader.onload = function(e) {
        rawData = e.target.result;
        connection.binaryData = "blob";
        connection.send(rawData);
    }
    reader.readAsArrayBuffer(file);
}

function handleKey(evt) {
  if (evt.keyCode === 13 || evt.keyCode === 14) {
    if (!document.getElementById("send").disabled) {
      send();
    }
  }
}