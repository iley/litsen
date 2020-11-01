package bot

import "time"

type watch struct {
	period   time.Duration
	callback func()
	done     chan struct{}
}

func newWatch(period time.Duration, callback func()) *watch {
	return &watch{
		period:   period,
		callback: callback,
		done:     make(chan struct{}, 1),
	}
}

func (w *watch) start() {
	go func() {
		for {
			select {
			case <-w.done:
				return
			case <-time.After(w.period):
				w.callback()
			}
		}
	}()
}

func (w *watch) stop() {
	w.done <- struct{}{}
}
