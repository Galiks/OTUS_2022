package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, s := range stages {
		if s == nil {
			// continue - see pipeline_test.go:94
			break
		}
		stageChan := make(Bi)
		go func(out In, stageResultChan Bi) {
			defer close(stageResultChan)
			for {
				select {
				case <-done:
					return
				case item, ok := <-out:
					if !ok {
						return
					}
					stageResultChan <- item
				}
			}
		}(out, stageChan)
		out = s(stageChan)
	}
	return out
}
