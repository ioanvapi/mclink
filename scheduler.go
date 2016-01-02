package main

import (
    "time"
    "log"
)


var (
    Scheduler *scheduler
)

func init() {
    Scheduler = newScheduler()
}

type scheduler struct {
    // timer for executing a provided function
    sch *time.Timer
    // timer for an alert message
    alert *time.Timer
    // time when scheduler will execute
    t   *time.Time
}


func newScheduler() *scheduler {
    return new(scheduler)
}


//Reset sets the internal state of scheduler.
//Sch timer will execute the schFunc function after minutes.
//Alert timer will popup an alert message one minute before Sch timer.
func (s *scheduler) Reset(minutes int, schFunc func(), alertMsg string) {
    s.setSchTimer(minutes, schFunc)

    if minutes > 0 {
        s.setAlertTimer(minutes - 1, alertMsg)
    }
}

//When says when scheduled function will be executed.
func (s *scheduler) When() *time.Time {
    return s.t
}

func (s *scheduler) setSchTimer(minutes int, schFunc func()) {
    var dm time.Duration = time.Duration(minutes) * time.Minute

    if s.sch != nil {
        s.sch.Stop()
    }

    s.sch = time.AfterFunc(dm, func() {
        defer s.null()
        schFunc()
    })

    // set the time when scheduler will execute
    newTime := time.Now().Add(dm)
    s.t = &newTime
}

func (s *scheduler) setAlertTimer(minutes int, alertMsg string) {
    var dm time.Duration = time.Duration(minutes) * time.Minute

    if s.alert != nil {
        s.alert.Stop()
    }

    s.alert = time.AfterFunc(dm, func() {
        err := popupWindow(alertMsg)
        if err != nil {
            log.Println("Error when popup alert message: ", err.Error())
        }
    })
}


func (s *scheduler) null() {
    s.sch = nil
    s.t = nil
}