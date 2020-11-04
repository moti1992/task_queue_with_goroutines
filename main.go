package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cast"
)

type Task struct {
	Id          string
	IsCompleted bool   // have a random function to mark the IsCompleted after a random period
	Status      string //completed, failed, timeout
	Time        string // when was the task created
	TaskData    int    // random string indicating task
	RetryCount  int
}

func dequeue(queue []Task) []Task {
	return queue[1:]
}

func createAndGetTasks() []Task {
	count := 100
	tasks := []Task{}
	for i := 0; i < count; i++ {
		newTask := Task{}
		newTask.Id = "task-" + cast.ToString(i)
		newTask.Time = time.Now().Format("2006-01-02")
		newTask.TaskData = i + 1
		newTask.RetryCount = 5
		tasks = append(tasks, newTask)
	}
	fmt.Println(tasks)
	return tasks
}

func runTasksFromQ(tasksQueue <-chan Task, out chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasksQueue {
		if task.TaskData%2 == 0 {
			fmt.Println("Even", task.TaskData)
			task.IsCompleted = true
			task.Status = "completed"
		} else {
			fmt.Println("Odd", task.TaskData)
			task.IsCompleted = true
			task.Status = "failed"
		}
		out <- task
		// time.Sleep(100 * time.Millisecond)
	}
	close(out)
}

func checkAndUpdateTasksFromQ(tasksQueue <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasksQueue {
		if task.Status == "failed" {
			if task.RetryCount == 0 {
				fmt.Println("cleaned as retried exhausted", task.TaskData)
				// dequeue(tasksQueue)
			}
		} else if task.Status == "completed" {
			fmt.Println("cleaned as completed", task.TaskData)
			// dequeue(tasksQueue)
		}
	}
	// time.Sleep(100 * time.Millisecond)
}

func main() {
	tasksQueue := createAndGetTasks()
	tasksQueueJob := make(chan Task, len(tasksQueue))
	var wg sync.WaitGroup
	wg.Add(1)
	tasksQueueJobCheck := make(chan Task, len(tasksQueue))
	for i := range tasksQueue {
		tasksQueueJob <- tasksQueue[i]
	}

	go runTasksFromQ(tasksQueueJob, tasksQueueJobCheck, &wg)
	close(tasksQueueJob)

	wg.Add(1)
	go checkAndUpdateTasksFromQ(tasksQueueJobCheck, &wg)
	wg.Wait()
}
