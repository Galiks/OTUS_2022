package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, s := range stages {
		stageChan := make(Bi)
		go func(out In, stageResultChan Bi) {
			defer close(stageResultChan)
			for {
				select {
				case item, ok := <-out:
					if !ok {
						return
					}
					stageResultChan <- item
				case <-done:
					return
				default:
					continue
				}
			}
		}(in, stageChan)
		in = s(stageChan)
	}
	return in
}
