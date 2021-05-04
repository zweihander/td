package e2e

import (
	"strings"

	"github.com/gotd/td/tg"
)

type SuiteBuilder struct {
	common   []tg.MessageClass
	secret   []tg.EncryptedMessageClass
	channels map[int][]tg.MessageClass
	pcounter int
	peerIDs  map[string]int

	initialChannelPts map[int]int
}

func NewSuite() *SuiteBuilder {
	return &SuiteBuilder{
		channels:          map[int][]tg.MessageClass{},
		peerIDs:           map[string]int{},
		initialChannelPts: map[int]int{},
	}
}

func (e *SuiteBuilder) getPeerID(name string) int {
	name = strings.TrimSpace(name)
	peerID, ok := e.peerIDs[name]
	if !ok {
		e.peerIDs[name] = e.pcounter
		e.pcounter++
	}

	return peerID
}

func (e *SuiteBuilder) createDialog(dialog ...string) []tg.MessageClass {
	if len(dialog)%2 != 0 {
		panic("bad dialog")
	}

	var messages []tg.MessageClass
	for i := 0; i < len(dialog); {
		from, text := dialog[i], dialog[i+1]
		messages = append(messages, &tg.Message{
			FromID: &tg.PeerUser{
				UserID: e.getPeerID(from),
			},
			Message: text,
		})
		i += 2
	}

	return messages
}

func (e *SuiteBuilder) createChannelDialog(channelID int, dialog ...string) []tg.MessageClass {
	if len(dialog)%2 != 0 {
		panic("bad dialog")
	}

	var messages []tg.MessageClass
	for i := 0; i < len(dialog); {
		from, text := dialog[i], dialog[i+1]
		messages = append(messages, &tg.Message{
			FromID: &tg.PeerUser{
				UserID: e.getPeerID(from),
			},
			PeerID: &tg.PeerChannel{
				ChannelID: channelID,
			},
			Message: text,
		})
		i += 2
	}

	return messages
}

func (e *SuiteBuilder) EncryptedMessage(from, text string) {
	e.secret = append(e.secret, &tg.EncryptedMessage{
		ChatID: e.getPeerID(from),
		Bytes:  []byte(text),
	})
}

func (e *SuiteBuilder) CommonDialog(dialog ...string) {
	e.common = append(e.common, e.createDialog(dialog...)...)
}

func (e *SuiteBuilder) ChannelDialog(channel string, dialog ...string) {
	channelID := e.getPeerID(channel)
	e.channels[channelID] = append(e.channels[channelID], e.createChannelDialog(channelID, dialog...)...)
	e.initialChannelPts[channelID] = 0
}
