package commands

import (
	"github.com/google/uuid"
	"github.com/p0t4t0sandwich/minecraft-be-websocket-api/src/protocol"
)

// CommandTextResponse - The body of a command text response message
type CommandTextResponse string

// Command Origin Types
type OriginType string

const (
	OriginTypePlayer OriginType = "player"
)

// CommandOrigin - The origin of a command message
type CommandOrigin struct {
	Type OriginType `json:"type"`
}

// CommandBody - The body of a command message
type CommandBody struct {
	Version     int           `json:"version"`
	Origin      CommandOrigin `json:"origin"`
	CommandLine string        `json:"commandLine"`
}

// NewCommandPacket - Create a new command packet
func NewCommandPacket(command string) *protocol.Packet {
	return &protocol.Packet{
		Header: protocol.Header{
			RequestId:      uuid.New(),
			MessagePurpose: protocol.CommandRequestType,
			MessageType:    protocol.CommandRequestType,
			Version:        1,
		},
		Body: CommandBody{
			Version:     1,
			Origin:      CommandOrigin{Type: OriginTypePlayer},
			CommandLine: command,
		},
	}
}

// CommandResponseBody - The body of a command response message
type CommandResponseBody struct {
	Message       string `json:"message,omitempty"`
	StatusCode    int    `json:"statusCode"`
	StatusMessage string `json:"statusMessage,omitempty"`
}

// CommandResponse - The body of a command response message
type CommandResponse struct {
	*protocol.Packet
	Body CommandResponseBody `json:"body"`
}

// NewCommandResponse - Create a new command response packet
func NewCommandResponse(packet *protocol.Packet) *CommandResponse {
	bodyMap := packet.Body.(map[string]interface{})
	if _, ok := bodyMap["message"]; !ok {
		bodyMap["message"] = ""
	}
	if _, ok := bodyMap["statusMessage"]; !ok {
		bodyMap["statusMessage"] = ""
	}

	return &CommandResponse{
		Packet: packet,
		Body: CommandResponseBody{
			Message:       bodyMap["message"].(string),
			StatusCode:    bodyMap["statusCode"].(int),
			StatusMessage: bodyMap["statusMessage"].(string),
		},
	}
}