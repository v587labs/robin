package schedule

import (
	"codnect.io/chrono"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/v587labs/robin/rcobra"
	"github.com/v587labs/robin/rlog"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
)

var schedulers = map[string]*scheduler{}
var schedulerLck = &sync.Mutex{}
var startOnce = &sync.Once{}

func Add(name string, task Task) *scheduler {
	_name := strings.TrimSpace(strings.ToLower(name))
	if len(_name) < 1 {
		panic(errors.New("task name error"))
	}
	schedulerLck.Lock()
	defer schedulerLck.Unlock()

	if _, ok := schedulers[_name]; ok {
		panic(errors.New("task " + name + " exists "))
	}

	s := &scheduler{
		id:        _name,
		name:      name,
		task:      task,
		scheduler: chrono.NewDefaultTaskScheduler(),
		lck:       sync.Mutex{},
	}
	schedulers[_name] = s
	return s
}

func Schedule(ctx context.Context) {
	startOnce.Do(func() {
		go _startAll(ctx)
	})

}
func Shutdown(ctx context.Context) error {
	schedulerLck.Lock()
	defer schedulerLck.Unlock()
	wg := errgroup.Group{}

	for _, s := range schedulers {
		(func(s *scheduler) {
			wg.Go(func() error {
				<-s.Shutdown()
				return nil
			})
		})(s)
	}
	return wg.Wait()

}

type commandContextKey struct {
}

func CommandContext(cmd *cobra.Command) context.Context {
	return context.WithValue(cmd.Context(), commandContextKey{}, cmd)
}

func _startAll(ctx context.Context) {
	for {
		(func() {
			schedulerLck.Lock()
			defer schedulerLck.Unlock()

			for _, s := range schedulers {
				err := s.startSchedule()
				if err != nil {
					panic(err)
				}
			}
		})()
	}
}

func Commands() []*cobra.Command {

	cmds := []*cobra.Command{}
	schedulerLck.Lock()
	defer schedulerLck.Unlock()

	for _, s := range schedulers {
		cmd := &cobra.Command{
			Use:   s.id,
			Short: fmt.Sprintf("run task :%s ", s.name),
			RunE:  s.runECmd,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				rlog.L(cmd.Context()).Info("task")
			},
		}
		cmds = append(cmds, cmd)
	}
	return cmds
}

func AddExecCommand(parent *cobra.Command) {
	sub := &cobra.Command{
		Use:   "exec",
		Short: "execute schedule task",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			rlog.L(cmd.Context()).Info("exec")
			return nil
		},
	}
	rcobra.AddCommands(sub, Commands()...)
	rcobra.AddCommands(parent, sub)
}
