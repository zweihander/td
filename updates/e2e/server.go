package e2e

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/gotd/td/tg"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type Server struct {
	common   map[int]tg.MessageClass
	secret   map[int]tg.EncryptedMessageClass
	channels map[int]map[int]tg.MessageClass

	// Current client remote state.
	pts        int
	qts        int
	date       int
	seq        int
	channelPts map[int]int

	log *zap.Logger
	mux sync.Mutex
}

func NewServer(
	common []tg.MessageClass,
	secret []tg.EncryptedMessageClass,
	channels map[int][]tg.MessageClass,
	log *zap.Logger,
) *Server {
	s := &Server{
		common:     map[int]tg.MessageClass{},
		secret:     map[int]tg.EncryptedMessageClass{},
		channels:   map[int]map[int]tg.MessageClass{},
		channelPts: map[int]int{},
		pts:        1,
		qts:        1,
		log:        log,
	}

	for i, msg := range common {
		s.common[i] = msg
	}
	for i, msg := range secret {
		s.secret[i] = msg
	}
	for channelID, msgs := range channels {
		events := make(map[int]tg.MessageClass)
		for i, msg := range msgs {
			events[i] = msg
		}

		s.channels[channelID] = events
		s.channelPts[channelID] = 1
	}

	return s
}

func (s *Server) UpdatesGetState(ctx context.Context) (*tg.UpdatesState, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	return &tg.UpdatesState{
		Pts:  s.pts,
		Qts:  s.qts,
		Date: s.date,
		Seq:  s.seq,
	}, nil
}

func (s *Server) UpdatesGetDifference(ctx context.Context, request *tg.UpdatesGetDifferenceRequest) (tg.UpdatesDifferenceClass, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if request.Pts > s.pts {
		return nil, xerrors.Errorf("pts is bigger than remote pts")
	}

	if request.Qts > s.qts {
		return nil, xerrors.Errorf("qts is bigger than remote qts")
	}

	var (
		ptsUpdates []tg.MessageClass
		qtsUpdates []tg.EncryptedMessageClass
		users      []tg.UserClass
		chats      []tg.ChatClass
	)

	for i := request.Pts + 1; i <= s.pts; i++ {
		msg, ok := s.common[i-1]
		if !ok {
			break
		}
		ptsUpdates = append(ptsUpdates, msg)
	}

	for i := request.Qts + 1; i <= s.qts; i++ {
		msg, ok := s.secret[i-1]
		if !ok {
			break
		}
		qtsUpdates = append(qtsUpdates, msg)
	}

	if len(ptsUpdates) == 0 && len(qtsUpdates) == 0 {
		return &tg.UpdatesDifferenceEmpty{
			Date: s.date,
			Seq:  s.seq,
		}, nil
	}

	return &tg.UpdatesDifference{
		NewMessages:          ptsUpdates,
		NewEncryptedMessages: qtsUpdates,
		Users:                users,
		Chats:                chats,
		State: tg.UpdatesState{
			Pts:  s.pts,
			Qts:  s.qts,
			Date: s.date,
			Seq:  s.seq,
		},
	}, nil
}

func (s *Server) UpdatesGetChannelDifference(ctx context.Context, request *tg.UpdatesGetChannelDifferenceRequest) (tg.UpdatesChannelDifferenceClass, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	channel, ok := request.Channel.(*tg.InputChannel)
	if !ok {
		return nil, xerrors.Errorf("unexpected channel input: %T", request.Channel)
	}

	pts, ok := s.channelPts[channel.ChannelID]
	if !ok {
		return nil, xerrors.Errorf("unknown channelID")
	}

	if request.Pts > pts {
		return nil, xerrors.Errorf("pts is bigger than remote pts")
	}

	var (
		newMessages  []tg.MessageClass
		otherUpdates []tg.UpdateClass
		users        []tg.UserClass
		chats        []tg.ChatClass
	)

	for i := request.Pts + 1; i <= pts; i++ {
		msg, ok := s.channels[channel.ChannelID][i-1]
		if !ok {
			break
		}

		newMessages = append(newMessages, msg)
	}

	if len(newMessages) == 0 {
		e := &tg.UpdatesChannelDifferenceEmpty{
			Pts: pts,
		}
		e.SetFinal(true)
		return e, nil
	}

	r := &tg.UpdatesChannelDifference{
		Pts:          pts,
		NewMessages:  newMessages,
		OtherUpdates: otherUpdates,
		Users:        users,
		Chats:        chats,
	}
	r.SetFinal(true)
	return r, nil
}

func (s *Server) Run(ctx context.Context, send func(tg.UpdatesClass) error) error {
	rand.Seed(time.Now().Unix())

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			update, ok := s.next()
			if !ok {
				s.log.Info("Send nops")
				return send(s.nops())
			}

			if rand.Intn(2) == 1 {
				s.log.Info("Update loss")
				continue
			}

			s.log.Info("Send updates", zap.Any("update", update))
			if err := send(update); err != nil {
				return err
			}

			time.Sleep(time.Millisecond * 20)
		}
	}
}

func (s *Server) next() (tg.UpdatesClass, bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	u := new(tg.Updates)
	if msg, ok := s.common[s.pts-1]; ok {
		u.Updates = append(u.Updates, &tg.UpdateNewMessage{
			Message:  msg,
			Pts:      s.pts,
			PtsCount: 1,
		})
		s.pts++
	}

	if msg, ok := s.secret[s.qts-1]; ok {
		u.Updates = append(u.Updates, &tg.UpdateNewEncryptedMessage{
			Message: msg,
			Qts:     s.qts,
		})
		s.qts++
	}

	for channelID, events := range s.channels {
		pts := s.channelPts[channelID]
		if msg, ok := events[pts-1]; ok {
			u.Updates = append(u.Updates, &tg.UpdateNewChannelMessage{
				Message:  msg,
				Pts:      pts,
				PtsCount: 1,
			})
			s.channelPts[channelID] = pts + 1
		}
	}

	if len(u.Updates) == 0 {
		return nil, false
	}

	return u, true
}

func (s *Server) nops() tg.UpdatesClass {
	s.mux.Lock()
	defer s.mux.Unlock()

	u := &tg.Updates{
		Updates: []tg.UpdateClass{
			&tg.UpdateWebPage{
				Pts:      s.pts,
				PtsCount: 1,
			},
			&tg.UpdateBotStopped{
				Qts: s.qts,
			},
		},
	}

	for channelID, pts := range s.channelPts {
		u.Updates = append(u.Updates, &tg.UpdateChannelWebPage{
			Pts:       pts,
			PtsCount:  1,
			ChannelID: channelID,
		})
	}

	return u
}
