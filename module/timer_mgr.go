package module

import (
	"Galia/log"
	"Galia/util/timer"
	"sync/atomic"
	"time"
)

type TimerMgr struct {
	dispatcher       *timer.Dispatcher
	mapActiveTimer   map[timer.ITimer]struct{}
	mapActiveIdTimer map[uint64]timer.ITimer
}

var timerSeedId uint32

func (m *TimerMgr) GenTimerId() uint64 {
	for {
		newTimerId := uint64(atomic.AddUint32(&timerSeedId, 1))
		if _, ok := m.mapActiveIdTimer[newTimerId]; ok == true {
			continue
		}
		return newTimerId
	}
}

func (m *TimerMgr) Release() {
	for pTimer := range m.mapActiveTimer {
		pTimer.Cancel()
		delete(m.mapActiveTimer, pTimer)
	}

	for id, t := range m.mapActiveIdTimer {
		t.Cancel()
		delete(m.mapActiveIdTimer, id)
	}
}
func (m *TimerMgr) OnCloseTimer(t timer.ITimer) {
	t.Cancel()
	delete(m.mapActiveIdTimer, t.GetId())
	delete(m.mapActiveTimer, t)
}

func (m *TimerMgr) OnAddTimer(t timer.ITimer) {
	if t != nil {
		if m.mapActiveTimer == nil {
			m.mapActiveTimer = map[timer.ITimer]struct{}{}
		}
		m.mapActiveTimer[t] = struct{}{}
	}
}

func (m *TimerMgr) AfterFunc(d time.Duration, cb func(*timer.Timer)) *timer.Timer {
	if m.mapActiveTimer == nil {
		m.mapActiveTimer = map[timer.ITimer]struct{}{}
	}
	return m.dispatcher.AfterFunc(d, nil, cb, m.OnCloseTimer, m.OnAddTimer)
}

func (m *TimerMgr) CronFunc(cronExpr *timer.CronExpr, cb func(*timer.Cron)) *timer.Cron {
	if m.mapActiveTimer == nil {
		m.mapActiveTimer = map[timer.ITimer]struct{}{}
	}
	return m.dispatcher.CronFunc(cronExpr, nil, cb, m.OnCloseTimer, m.OnAddTimer)
}

func (m *TimerMgr) NewTicker(d time.Duration, cb func(*timer.Ticker)) *timer.Ticker {
	if m.mapActiveTimer == nil {
		m.mapActiveTimer = map[timer.ITimer]struct{}{}
	}
	return m.dispatcher.TickerFunc(d, nil, cb, m.OnCloseTimer, m.OnAddTimer)
}
func (m *TimerMgr) SafeAfterFunc(timerId *uint64, d time.Duration, AdditionData interface{}, cb func(uint64, interface{})) {
	if m.mapActiveIdTimer == nil {
		m.mapActiveIdTimer = map[uint64]timer.ITimer{}
	}

	if *timerId != 0 {
		m.CancelTimerId(timerId)
	}

	*timerId = m.GenTimerId()
	t := m.dispatcher.AfterFunc(d, cb, nil, m.OnCloseTimer, m.OnAddTimer)
	t.AdditionData = AdditionData
	t.Id = *timerId
	m.mapActiveIdTimer[*timerId] = t
}

func (m *TimerMgr) SafeCronFunc(cronId *uint64, cronExpr *timer.CronExpr, AdditionData interface{}, cb func(uint64, interface{})) {
	if m.mapActiveIdTimer == nil {
		m.mapActiveIdTimer = map[uint64]timer.ITimer{}
	}

	*cronId = m.GenTimerId()
	c := m.dispatcher.CronFunc(cronExpr, cb, nil, m.OnCloseTimer, m.OnAddTimer)
	c.AdditionData = AdditionData
	c.Id = *cronId
	m.mapActiveIdTimer[*cronId] = c
}

func (m *TimerMgr) SafeNewTicker(tickerId *uint64, d time.Duration, AdditionData interface{}, cb func(uint64, interface{})) {
	if m.mapActiveIdTimer == nil {
		m.mapActiveIdTimer = map[uint64]timer.ITimer{}
	}

	*tickerId = m.GenTimerId()
	t := m.dispatcher.TickerFunc(d, cb, nil, m.OnCloseTimer, m.OnAddTimer)
	t.AdditionData = AdditionData
	t.Id = *tickerId
	m.mapActiveIdTimer[*tickerId] = t
}

func (m *TimerMgr) CancelTimerId(timerId *uint64) bool {
	if m.mapActiveIdTimer == nil {
		log.SError("mapActiveIdTimer is nil")
		return false
	}

	t, ok := m.mapActiveIdTimer[*timerId]
	if ok == false {
		log.SError("cannot find timer id ", timerId)
		return false
	}

	t.Cancel()
	*timerId = 0
	return true
}
