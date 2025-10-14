package clipboard

import (
	"context"
	"time"
)

type Watcher struct {
	lastContent string
	interval    time.Duration
}

func NewWatcher(interval time.Duration) *Watcher {
	return &Watcher{
		interval: interval,
	}
}

func (w *Watcher) Start(ctx context.Context, onChange func(string)) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			current, err := Read()
			if err != nil {
				continue
			}

			if current != w.lastContent && current != "" && len(current) > 0 {
				w.lastContent = current
				onChange(current)
			}
		}
	}
}
