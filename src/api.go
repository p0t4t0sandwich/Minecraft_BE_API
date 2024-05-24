package server

import (
	"encoding/json"
	"net/http"

	"github.com/p0t4t0sandwich/minecraft-be-websocket-api/src/protocol"
	"github.com/p0t4t0sandwich/minecraft-be-websocket-api/src/protocol/commands"
)

// CMDHandler - Handle the CMD route
func CMDHandler(wss *WebSocketServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		// Get command from JSON body
		var body struct {
			Command string `json:"cmd"`
		}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		msg := commands.NewCommandPacket(body.Command)
		data, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = wss.Send(id, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		wss.AddCommand(msg.Header.RequestId, commands.FromString(body.Command))
	}
}

// EventSubscribeHandler - Handle the EventSubscribe route
func EventSubscribeHandler(wss *WebSocketServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		eventName := r.PathValue("name")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}
		if eventName == "" {
			http.Error(w, "Event Name is requried", http.StatusBadRequest)
			return
		}

		msg := protocol.NewEventSubPacket(eventName, protocol.SubscribeType)
		data, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = wss.Send(id, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// EventUnsubscribeHander - Handle the Unsubscribe route
func EventUnsubscribeHander(wss *WebSocketServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		eventName := r.PathValue("name")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}
		if eventName == "" {
			http.Error(w, "Event Name is requried", http.StatusBadRequest)
			return
		}

		msg := protocol.NewEventSubPacket(eventName, protocol.UnsubscribeType)
		data, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = wss.Send(id, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}