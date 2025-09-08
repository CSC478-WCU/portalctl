package portal

import (
	"context"
	"errors"
	"time"
)

var ErrWaitTimeout = errors.New("wait: timeout")

// WaitForStatus polls experimentStatus -j until done() returns true or timeout.
func (c *Client) WaitForStatus(ctx context.Context, exp string, interval, timeout time.Duration, done func(*StatusPayload) bool) (*StatusPayload, error) {
	deadline := time.Now().Add(timeout)
	for {
		if timeout > 0 && time.Now().After(deadline) {
			return nil, ErrWaitTimeout
		}
		resp, err := c.ExperimentStatus(exp, true, false, false)
		if err == nil {
			if p, err2 := ParseStatusJSON(resp.Output); err2 == nil {
				if done(p) {
					return p, nil
				}
			}
		}
		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
