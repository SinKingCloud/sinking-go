package job

import (
	"sync"
)

type Task struct {
	Thread   int
	Producer func(channel chan interface{})
	Consumer func(param interface{})
}

// SetThread 设置协程数量
func (task *Task) SetThread(num int) *Task {
	task.Thread = num
	return task
}

// SetProducer 设置生产者
func (task *Task) SetProducer(fun func(fun chan interface{})) *Task {
	task.Producer = fun
	return task
}

// SetConsumer 设置消费者
func (task *Task) SetConsumer(fun func(fun interface{})) *Task {
	task.Consumer = fun
	return task
}

func (task *Task) Run() {
	wg := &sync.WaitGroup{}
	tasks := make(chan interface{})
	for i := 0; i < task.Thread; i++ {
		wg.Add(1)
		go func(group *sync.WaitGroup, tasks chan interface{}) {
			for task2 := range tasks {
				if task.Consumer != nil && task2 != nil {
					task.Consumer(task2)
				}
			}
			group.Done()
		}(wg, tasks)
	}
	if task.Producer != nil {
		task.Producer(tasks)
	}
	tasks <- nil
	wg.Wait()
}
