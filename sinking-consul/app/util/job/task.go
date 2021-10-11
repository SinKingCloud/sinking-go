package job

import (
	"sync"
)

type Task struct {
	Thread   int
	Producer func(channel chan string)
	Consumer func(param string)
}

// SetThread 设置协程数量
func (task *Task) SetThread(num int) *Task {
	task.Thread = num
	return task
}

// SetProducer 设置生产者
func (task *Task) SetProducer(fun func(fun chan string)) *Task {
	task.Producer = fun
	return task
}

// SetConsumer 设置消费者
func (task *Task) SetConsumer(fun func(fun string)) *Task {
	task.Consumer = fun
	return task
}

func (task *Task) Run() {
	wg := &sync.WaitGroup{}
	tasks := make(chan string)
	for i := 0; i < task.Thread; i++ {
		wg.Add(1)
		go func(group *sync.WaitGroup, tasks chan string) {
			for task2 := range tasks {
				if task.Consumer != nil && task2 != "" {
					task.Consumer(task2)
				}
			}
			group.Done()
		}(wg, tasks)
	}
	if task.Producer != nil {
		task.Producer(tasks)
	}
	tasks <- ""
	wg.Wait()
}
