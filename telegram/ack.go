package telegram

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/gotd/td/internal/mt"
)

func (c *Client) ackLoop(ctx context.Context) {
	c.wg.Add(1)
	defer c.wg.Done()

	log := c.log.Named("ack")

	var buf []int64
	send := func() {
		defer func() { buf = buf[:0] }()

		if err := c.writeServiceMessage(ctx, &mt.MsgsAck{MsgIds: buf}); err != nil {
			c.log.Error("Failed to ACK", zap.Error(err))
			return
		}

		log.Debug("ACK", zap.Int64s("message_ids", buf))
	}

	ticker := time.NewTicker(c.ackInterval) // TODO: remove side-effect
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if len(buf) > 0 {
				send()
			}
		case msgID := <-c.ackSendChan:
			buf = append(buf, msgID)
			if len(buf) >= c.ackBatchSize {
				send()
				ticker.Reset(c.ackInterval)
			}
		}
	}
}

// waitACK waits to receive ACK from the server until the context is canceled.
// If ACK was received - returns nil.
// If the context was canceled before the ACK was received, it returns a context error.
func (c *Client) waitACK(ctx context.Context, msgID int64) error {
	got := make(chan struct{})

	c.ackMux.Lock()
	c.ack[msgID] = func() { got <- struct{}{} }
	c.ackMux.Unlock()

	defer func() {
		c.ackMux.Lock()
		delete(c.ack, msgID)
		c.ackMux.Unlock()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-got:
		return nil
	}
}
