package e2e

import (
	"fmt"
	"sync"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/updates"
	"go.uber.org/zap"
)

type State struct {
	common     []tg.MessageClass
	secret     []tg.EncryptedMessageClass
	channels   map[int][]tg.MessageClass
	pts        int
	qts        int
	seq        int
	channelPts map[int]int

	stateless []tg.UpdateClass
	mux       sync.Mutex
	log       *zap.Logger
}

func NewState() *State {
	return &State{
		channels:   map[int][]tg.MessageClass{},
		channelPts: map[int]int{},
	}
}

func (s *State) HandleDiff(newState tg.UpdatesState, diff updates.DiffUpdate) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.pts = newState.Pts
	s.qts = newState.Qts
	s.seq = newState.Seq

	s.common = append(s.common, diff.NewMessages...)
	s.secret = append(s.secret, diff.NewEncryptedMessages...)
	return nil
}

func (s *State) HandleChannelDiff(channelID int, pts int, diff updates.ChannelDiffUpdate) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.channelPts[channelID] = pts
	s.channels[channelID] = append(s.channels[channelID], diff.NewMessages...)
	return nil
}

func (s *State) HandleStatelessUpdates(updates []tg.UpdateClass) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.stateless = append(s.stateless, updates...)
	return nil
}

func (s *State) HandlePtsUpdate(pts int, update updates.PtsUpdate) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, u := range update.Updates {
		switch u := u.(type) {
		case *tg.UpdateNewMessage:
			s.common = append(s.common, u.Message)
		default:
			s.log.Debug(fmt.Sprintf("unexpected update type: %T", u))
		}
	}

	s.pts = pts
	return nil
}

func (s *State) HandleQtsUpdate(qts int, update updates.QtsUpdate) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, u := range update.Updates {
		switch u := u.(type) {
		case *tg.UpdateNewEncryptedMessage:
			s.secret = append(s.secret, u.Message)
		default:
			s.log.Debug(fmt.Sprintf("unexpected update type: %T", u))
		}
	}

	s.qts = qts
	return nil
}

func (s *State) HandleSeqUpdate(seq int, update updates.SeqUpdate) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.stateless = append(s.stateless, update.Updates...)
	s.seq = seq
	return nil
}

func (s *State) HandleChannelUpdate(channelID int, pts int, update updates.PtsUpdate) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, u := range update.Updates {
		switch u := u.(type) {
		case *tg.UpdateNewChannelMessage:
			s.channels[channelID] = append(s.channels[channelID], u.Message)
		default:
			s.log.Debug(fmt.Sprintf("unexpected update type: %T", u))
		}
	}

	return nil
}

func (s *State) ChannelTooLong(channelID int) {
	panic("not implemented")
}
