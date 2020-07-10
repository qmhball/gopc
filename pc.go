package pc

import (
	"sync"
)

type Task interface {
}

type AbstructPC struct {
	ConsumerNum int
	ChanLen     int
	Tasks       chan Task
	Consumer    func(Task)
	Producer    func(chan Task)
	wg          sync.WaitGroup
}

func (a *AbstructPC) Init(cl int, cn int, p func(chan Task), c func(Task)) {
	a.ChanLen = cl
	a.ConsumerNum = cn
	a.Tasks = make(chan Task, a.ConsumerNum)
	a.Producer = p
	a.Consumer = c
}

func (a *AbstructPC) producerDispatch() {
	defer close(a.Tasks)

	a.Producer(a.Tasks)
}

func (a *AbstructPC) consumerDispatch() {
	for i := 0; i < a.ConsumerNum; i++ {
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			for task := range a.Tasks {
				a.Consumer(task)
			}
		}()
	}

	a.wg.Wait()
}

func (a *AbstructPC) Run() {
	go a.producerDispatch()
	a.consumerDispatch()
}
