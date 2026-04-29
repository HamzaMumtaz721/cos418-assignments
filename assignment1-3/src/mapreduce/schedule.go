package mapreduce

import (
	"sync"
	)

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int
	switch phase {
		case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
		case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	var wg sync.WaitGroup

	for i := 0; i < ntasks; i++ {
			wg.Add(1)
			taskNum := i
			go func() {
				defer wg.Done()
				for {
					worker := <-mr.registerChannel
					args := DoTaskArgs{
						JobName:       mr.jobName,
						Phase:         phase,
						TaskNumber:    taskNum,
						NumOtherPhase: nios,
					}
					if phase == mapPhase {
							args.File = mr.files[taskNum]
						}
					ok := call(worker, "Worker.DoTask", args, new(struct{}))
					if ok {
						go func() {
								mr.registerChannel <- worker
							}()
						return
						}
					}
				}()
		}

	wg.Wait()
	debug("Schedule: %v phase done\n", phase)
	}
