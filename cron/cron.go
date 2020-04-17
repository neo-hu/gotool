package cron

import (
	"sort"
	"sync"
	"time"
)

type JobID int
type Job interface {
	Run()
}

type FuncJob func()

func (f FuncJob) Run() { f() }

type Schedule interface {
	Next(time.Time) time.Time
}

type Cron struct {
	sync.RWMutex
	jobs    []*job
	nextID JobID
	running bool
	add     chan *job
	remove  chan JobID
	stop    chan struct{}
	event   chan interface{}
}
type job struct {
	id       JobID
	run      Job
	schedule Schedule
	next     time.Time
}

func (j *job) scheduleNext(now time.Time) {
	j.next = j.schedule.Next(now)
}

func New() *Cron {
	return &Cron{
		add:    make(chan *job),
		stop:   make(chan struct{}),
		remove: make(chan JobID),
		event:  make(chan interface{}, 10),
	}
}

func (c *Cron) Remove(id JobID)  {
	c.RLock()
	if c.running {
		c.RUnlock()
		c.remove <- id
	} else {
		c.RUnlock()
		c.removeJob(id)
	}
}
func (c *Cron) Add(duration time.Duration, cmd Job) JobID {
	c.Lock()
	j := &job{
		id:       c.nextID,
		schedule: NewSchedule(duration),
		run:      cmd,
	}
	c.nextID += 1
	if c.running {
		c.Unlock()
		c.add <- j
	} else {
		c.jobs = append(c.jobs, j)
		c.Unlock()
	}
	return j.id
}

func (c *Cron) Start() {
	c.Lock()
	defer c.Unlock()
	if c.running {
		return
	}
	c.event = make(chan interface{}, 10)
	c.stop = make(chan struct{})
	c.running = true
	go c.run()
}

func (c *Cron) run() {
	now := time.Now().Local()
	c.RLock()
	for _, j := range c.jobs {
		// todo 计算下一次执行的时间
		j.scheduleNext(now)
	}
	c.RUnlock()
	for {
		c.RLock()
		// todo 按下一次执行的时间运行
		sort.Sort(byTime(c.jobs))
		c.RUnlock()
		var effective time.Time
		if len(c.jobs) == 0 || c.jobs[0].next.IsZero() {
			effective = now.AddDate(10, 0, 0)
		} else {
			effective = c.jobs[0].next
		}
		select {
		case now = <-time.After(effective.Sub(now)):
			for _, j := range c.jobs {
				if j.next != effective {
					break
				}
				c.onEvent(&JobStartEvent{JobEvent{j.id}})
				go j.run.Run()
				j.scheduleNext(now)
			}
			continue
		case j := <-c.add:
			c.Lock()
			c.jobs = append(c.jobs, j)
			c.Unlock()
			j.scheduleNext(now)
			c.onEvent(&JobAddEvent{JobEvent{j.id}})
		case id := <-c.remove:
			c.removeJob(id)
		case <-c.stop:
			return
		}
		now = time.Now().Local()
	}
}

func (c *Cron) Event() <-chan interface{} {
	return c.event
}
func (c *Cron) onEvent(event interface{}) {
	select {
	case c.event <- event:
	default:
		// todo 确保不会因为event堵塞导致任务不能定期执行
	}
}

func (c *Cron) removeJob(id JobID) {
	c.Lock()
	defer c.Unlock()
	var jobs []*job
	for _, job := range c.jobs {
		if job.id != id {
			jobs = append(jobs, job)
		} else {
			c.onEvent(&JobRemoveEvent{JobEvent{job.id}})
		}
	}
	c.jobs = jobs
}

func (c *Cron) Stop() {
	c.Lock()
	defer c.Unlock()
	if !c.running {
		return
	}
	close(c.stop)
	close(c.event)
	c.running = false
}

type DelaySchedule struct {
	Delay time.Duration
}

func NewSchedule(duration time.Duration) DelaySchedule {
	if duration < time.Second {
		duration = time.Second
	}
	return DelaySchedule{
		Delay: duration - time.Duration(duration.Nanoseconds())%time.Second,
	}
}

func (schedule DelaySchedule) Next(t time.Time) time.Time {
	return t.Add(schedule.Delay - time.Duration(t.Nanosecond())*time.Nanosecond)
}

type byTime []*job

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	if s[i].next.IsZero() {
		return false
	}
	if s[j].next.IsZero() {
		return true
	}
	return s[i].next.Before(s[j].next)
}
