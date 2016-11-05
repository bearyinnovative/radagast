package bearychat

import (
	"context"

	bc "github.com/bcho/bearychat.go"
)

type RTMMessage struct {
	Text        string                  `json:"text"`
	VchannelId  string                  `json:"vchannel"`
	IsMarkdown  bool                    `json:"markdown,omitempty"`
	Attachments []bc.IncomingAttachment `json"attachments,omitempty"`
}

func SendToVchannel(c context.Context, rtmClient *bc.RTMClient, message RTMMessage) error {
	_, err := rtmClient.Post("message", message, nil)
	return err
}
