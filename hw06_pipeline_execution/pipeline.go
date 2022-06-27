package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	go func() {
		for _, s := range stages {
			select {
			case <-done:
				return
			default:
				s := s
				in = s(in)
			}
		}
		go func() {
			defer close(out)
			for {
				select {
				case <-done:
					return
				case out <- <-in:
					continue
				default:
					return
				}
			}
		}()
	}()

	return out
}
