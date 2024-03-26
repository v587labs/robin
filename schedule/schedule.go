package schedule

import (
	"codnect.io/chrono"
	"context"
	"github.com/spf13/cobra"
	"github.com/v587labs/robin/rlog"
	"sync"
	"time"
)

type Task func(ctx context.Context)

type scheduler struct {
	id            string
	name          string
	task          Task
	scheduler     chrono.TaskScheduler
	schedulerTask chrono.ScheduledTask
	scheduleFn    func() (chrono.ScheduledTask, error)

	scheduled bool
	lck       sync.Mutex
}

func (s *scheduler) run(ctx context.Context) {
	ctx, logger := rlog.With(ctx, "task", s.name, "id", "")
	logger.Debug("start")
	s.task(ctx)
	logger.Debug("complete")
}

func (s *scheduler) Schedule(options ...chrono.Option) *scheduler {
	s.lck.Lock()
	defer s.lck.Unlock()
	if s.scheduled {
		return s
	}
	s.scheduleFn = func() (chrono.ScheduledTask, error) {
		return s.scheduler.Schedule(s.run, options...)
	}
	return s
}
func (s *scheduler) ScheduleWithCron(expression string, options ...chrono.Option) *scheduler {
	s.lck.Lock()
	defer s.lck.Unlock()
	if s.scheduled {
		return s
	}
	s.scheduleFn = func() (chrono.ScheduledTask, error) {
		return s.scheduler.ScheduleWithCron(s.run, expression, options...)
	}
	return s
}
func (s *scheduler) ScheduleWithFixedDelay(delay time.Duration, options ...chrono.Option) *scheduler {
	s.lck.Lock()
	defer s.lck.Unlock()
	if s.scheduled {
		return s
	}
	s.scheduleFn = func() (chrono.ScheduledTask, error) {
		return s.scheduler.ScheduleWithFixedDelay(s.run, delay, options...)
	}
	return s
}
func (s *scheduler) ScheduleAtFixedRate(period time.Duration, options ...chrono.Option) *scheduler {
	s.lck.Lock()
	defer s.lck.Unlock()
	if s.scheduled {
		return s
	}
	s.scheduleFn = func() (chrono.ScheduledTask, error) {
		return s.scheduler.ScheduleAtFixedRate(s.run, period, options...)
	}
	return s
}
func (s *scheduler) IsShutdown() bool {
	return s.scheduler.IsShutdown()
}
func (s *scheduler) Shutdown() chan bool {
	return s.scheduler.Shutdown()
}
func (s *scheduler) startSchedule() error {
	s.lck.Lock()
	defer s.lck.Unlock()
	if s.scheduled || s.scheduleFn == nil {
		return nil
	}
	s.scheduled = true
	t, err := s.scheduleFn()
	s.schedulerTask = t
	return err
}

func (s *scheduler) runECmd(cmd *cobra.Command, args []string) error {
	s.run(CommandContext(cmd))
	return nil
}
