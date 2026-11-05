package repos

import (
	"context"
	"errors"
	"time"
)

// ErrRepoBusy is retryable preparation backpressure, not an analysis failure.
var ErrRepoBusy = errors.New("repository prepare lock busy")

type repoLockWaitKey struct{}

// WithRepoLockWait bounds only time waiting for the writer lock. It does not
// shorten an acquired holder's fetch deadline or release a live writer's lock.
// Zero means try once. A negative duration preserves the caller's deadline.
func WithRepoLockWait(ctx context.Context, wait time.Duration) context.Context {
	return context.WithValue(ctx, repoLockWaitKey{}, wait)
}

func lockWaitContext(ctx context.Context) (context.Context, context.CancelFunc, bool) {
	wait, ok := ctx.Value(repoLockWaitKey{}).(time.Duration)
	if !ok || wait < 0 {
		return ctx, func() {}, false
	}
	child, cancel := context.WithTimeout(ctx, wait)
	return child, cancel, true
}
