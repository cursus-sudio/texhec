package internal

import (
	"engine/modules/batcher"
	"engine/services/clock"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"
	"time"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type Service struct {
	EventsBuilder events.Builder `inject:"1"`
	Clock         clock.Clock    `inject:"1"`
	Logger        logger.Logger  `inject:"1"`

	workers    int
	timeBudget time.Duration
	tasks      []batcher.Task
}

func NewService(
	c ioc.Dic,
	workers int,
	frameLoadingBudget time.Duration,
) *Service {
	s := ioc.GetServices[*Service](c)
	s.workers = workers
	s.timeBudget = frameLoadingBudget
	return s
}

func (s *Service) NewTask() batcher.TaskFactory { return NewTaskFactory(s.workers) }
func (s *Service) Queue(task batcher.Task)      { s.tasks = append(s.tasks, task) }
func (s *Service) Progress() float32 {
	if len(s.tasks) != 0 {
		return s.tasks[0].Progress()
	}
	return -1
}

func (s *Service) System() batcher.System {
	return ecs.NewSystemRegister(func() error {
		events.Listen(s.EventsBuilder, s.Listen)
		return nil
	})
}

func (s *Service) Listen(frames.FrameEvent) {
	if len(s.tasks) == 0 {
		return
	}
	task := s.tasks[0]

	start := s.Clock.Now()
	for s.Clock.Now().Sub(start) < s.timeBudget {
		task.Step()
		if task.Progress() != 1 {
			continue
		}
		s.tasks = s.tasks[1:]
		if len(s.tasks) == 0 {
			break
		}
		task = s.tasks[0]
	}
}
