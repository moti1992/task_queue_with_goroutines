# task_queue_with_goroutines

1. Push the tasks to channel
2. Tasks will be picked by goroutine, process it and write to channel for other goroutine to check
3. Concurrently other goroutine will pick the processed data from channel and check if its completed or not
