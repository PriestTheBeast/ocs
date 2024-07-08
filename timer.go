package main

import "time"

type Timer struct {
	Name                  string
	lastStartTime         time.Time
	durationBeforePausing time.Duration
	Running               bool
}

func (t *Timer) Start() {
	if t.Running {
		return
	}
	t.lastStartTime = time.Now()
	t.Running = true
}

func (t *Timer) Pause() {
	if !t.Running {
		return
	}
	t.durationBeforePausing += time.Since(t.lastStartTime)
	t.lastStartTime = time.Time{}
	t.Running = false
}

func (t *Timer) ElapsedTime() time.Duration {
	if !t.Running {
		return t.durationBeforePausing
	}

	return t.durationBeforePausing + time.Since(t.lastStartTime)
}
