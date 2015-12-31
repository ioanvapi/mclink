package main

import (
	"time"
)


var (
	Scheduler *scheduler
)

func init() {
	Scheduler = newScheduler()
}

type scheduler struct {
	sch *time.Timer
	// time when scheduler
	t   *time.Time
}


func newScheduler() *scheduler {
	return new(scheduler)
}

func (s *scheduler) Reset(minutes int, f func()) {
	var dm time.Duration = time.Duration(minutes) * time.Minute

	if s.sch != nil {
		s.sch.Stop()
	}

	s.sch = time.AfterFunc(dm, func() {
		defer s.null()
		f()
	})
	newTime := time.Now().Add(dm)
	s.t = &newTime
}

func (s *scheduler) When() *time.Time {
	return s.t
}

func (s *scheduler) null() {
	s.sch = nil
	s.t = nil
}