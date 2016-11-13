package mapreduce

import "fmt"

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//

	barrier := make(chan struct{})
	for taskId := 0; taskId < ntasks; taskId++ {
		// Get a worker
		go func(taskId int) {
			worker := <-mr.registerChannel
			if call(worker, "Worker.DoTask", &DoTaskArgs{
				mr.jobName,
				mr.files[taskId],
				phase,
				taskId,
				nios,
			}, new(struct{})) {
				// Run task successfully
				barrier <- struct{}{}
				// Return the worker!!!
				mr.registerChannel <- worker
				fmt.Printf("Finish task %d.\n", taskId)
			}
		}(taskId)
	}

	for taskId := 0; taskId < ntasks; taskId++ {
		<- barrier
	}

	fmt.Printf("Schedule: %v phase done\n", phase)
}
