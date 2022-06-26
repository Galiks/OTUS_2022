package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	wg := sync.WaitGroup{}
	go func() {
		temp := in
		defer close(out)
		for _, s := range stages {
			wg.Add(1)
			defer wg.Done()
			temp = s(temp)
		}
		go func() {
			for t := range temp {
				if t != nil {
					out <- t
				}
			}
		}()
		wg.Wait()
	}()

	return out
}
