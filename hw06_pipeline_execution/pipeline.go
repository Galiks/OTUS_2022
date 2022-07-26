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
			continue
		}
		out = stageHandler(done, s(out))
	}
	return out
}

func stageHandler(done, data In) Out {
	stageChan := make(Bi)
	go func(data In, stageResultChan Bi) {
		defer close(stageResultChan)
		for {
			select {
			case <-done:
				return
			case item, ok := <-data:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case stageResultChan <- item:
				}
			}
		}
	}(data, stageChan)
	return stageChan
}
