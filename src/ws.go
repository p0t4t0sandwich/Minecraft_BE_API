package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/p0t4t0sandwich/minecraft-be-websocket-api/src/protocol"
	"github.com/p0t4t0sandwich/minecraft-be-websocket-api/src/protocol/commands"
	"github.com/p0t4t0sandwich/minecraft-be-websocket-api/src/protocol/events"
)

type PacketCallback func(id string, msg []byte, packetJSON map[string]interface{}, packet *protocol.Packet)
type CommandCallback func(id string, msg []byte, packetJSON map[string]interface{}, command *commands.CommandResponse)
type EventCallback func(id string, msg []byte, packetJSON map[string]interface{}, event *events.EventPacket)

var upgrader = websocket.Upgrader{}

// WebSocketServer WebSocket relay
type WebSocketServer struct {
	conns            map[string]*websocket.Conn
	commands         map[uuid.UUID]commands.CommandName
	commandListeners map[commands.CommandName][]CommandCallback
	eventListeners   map[events.EventName][]EventCallback
}

// NewWebSocketServer Create a new WebSocket relay
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		conns:            make(map[string]*websocket.Conn),
		commands:         make(map[uuid.UUID]commands.CommandName),
		commandListeners: make(map[commands.CommandName][]CommandCallback),
		eventListeners:   make(map[events.EventName][]EventCallback),
	}
}

// Add a connection to the relay
func (wss *WebSocketServer) Add(id string, conn *websocket.Conn) {
	wss.conns[id] = conn
	wss.HandleEvent(id, []byte{}, map[string]interface{}{}, &events.EventPacket{
		Header: &events.EventHeader{
			Header: protocol.Header{
				MessagePurpose: protocol.EventType,
				MessageType:    protocol.CommandResponseType,
				Version:        1,
			},
			EventName: events.WebSocketConnect,
		},
	})
}

// Remove a connection from the relay
func (wss *WebSocketServer) Remove(id string) {
	delete(wss.conns, id)
}

// Send a message to a connection
func (wss *WebSocketServer) Send(id string, msg []byte) error {
	conn, ok := wss.conns[id]
	if !ok {
		return nil
	}
	return conn.WriteMessage(websocket.TextMessage, msg)
}

// SendPacket Send a packet to a connection
func (wss *WebSocketServer) SendPacket(id string, packet *protocol.Packet) error {
	msg, err := json.Marshal(packet)
	if err != nil {
		return err
	}
	return wss.Send(id, msg)
}

// AddCommand Add a command to the relay
func (wss *WebSocketServer) AddCommand(id uuid.UUID, command commands.CommandName) {
	wss.commands[id] = command
}

// PopCommand Pop a command from the relay
func (wss *WebSocketServer) PopCommand(id uuid.UUID) (commands.CommandName, bool) {
	command, ok := wss.commands[id]
	if ok {
		delete(wss.commands, id)
	}
	return command, ok
}

// AddCommandListener Add a command listener
func (wss *WebSocketServer) AddCommandListener(command commands.CommandName, callback CommandCallback) {
	wss.commandListeners[command] = append(wss.commandListeners[command], callback)
}

// AddEventListener Add an event listener
func (wss *WebSocketServer) AddEventListener(event events.EventName, callback EventCallback) {
	wss.eventListeners[event] = append(wss.eventListeners[event], callback)
}

// HandleCommand Handle a command packet
func (wss *WebSocketServer) HandleCommand(id string, msg []byte, packetJSON map[string]interface{}, command *commands.CommandResponse, commandName commands.CommandName) {
	for _, callback := range wss.commandListeners[commandName] {
		callback(id, msg, packetJSON, command)
	}
	if wss.commandListeners[commandName] == nil {
		for _, callback := range wss.commandListeners[commands.Unknown] {
			callback(id, msg, packetJSON, command)
		}
	}
}

// HandleEvent Handle an event packet
func (wss *WebSocketServer) HandleEvent(id string, msg []byte, packetJSON map[string]interface{}, event *events.EventPacket) {
	for _, callback := range wss.eventListeners[event.Header.EventName] {
		callback(id, msg, packetJSON, event)
	}
	if wss.eventListeners[event.Header.EventName] == nil {
		for _, callback := range wss.eventListeners[events.Unknown] {
			callback(id, msg, packetJSON, event)
		}
	}
}

// WSHandler Handle a WebSocket connection
func (wss *WebSocketServer) WSHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(ws)
	defer wss.Remove(id)
	defer log.Printf("[%s] < Disconnected", id)

	wss.Add(id, ws)
	log.Printf("[%s] > Connected", id)

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		} else if msgType != websocket.TextMessage {
			log.Println("Message type is not text")
		}
		wss.HandlePacket(id, msg)
	}
}

// HandlePacket Handle a packet
func (wss *WebSocketServer) HandlePacket(id string, msg []byte) {
	packetJSON := make(map[string]interface{})
	err := json.Unmarshal(msg, &packetJSON)
	if err != nil {
		log.Println(err.Error())
		return
	}

	header, ok := packetJSON["header"].(map[string]interface{})
	if !ok {
		log.Println("Header is not an object")
		return
	}
	messagePurpose, ok := header["messagePurpose"].(string)
	if !ok {
		log.Println("MessagePurpose is not a string")
		return
	}
	switch protocol.MessageType(messagePurpose) {
	case protocol.CommandResponseType:
		command := &commands.CommandResponse{}
		err = json.Unmarshal(msg, command)
		if err != nil {
			log.Println(err.Error())
			return
		}
		commandName, ok := wss.PopCommand(command.Header.RequestId)
		if !ok {
			log.Printf("[%s] Command response: %s", id, command.Body.StatusMessage)
			return
		}
		wss.HandleCommand(id, msg, packetJSON, command, commandName)
	case protocol.EventType:
		event := &events.EventPacket{}
		err = json.Unmarshal(msg, event)
		if err != nil {
			log.Println(err.Error())
			return
		}
		wss.HandleEvent(id, msg, packetJSON, event)
	default:
		log.Printf("[%s] %s", id, messagePurpose)
		log.Println(string(msg))
	}
}
