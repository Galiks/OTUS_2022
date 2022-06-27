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
		defer close(out)
		var i int = 0
		for {
			select {
			case <-done:
				return
			default:
				in = stages[i](in)
				i++
				if len(stages) == i {
					for i := range in {
						out <- i
					}
					return
				}
			}
		}
	}()
	return out
}

// go func() {
// 	for _, s := range stages {
// 		in = s(in)
// 	}
// 	go func() {
// 		defer close(out)
// 		for {
// 			select {
// 			case <-done:
// 				return
// 			case out <- <-in:
// 				continue
// 			default:
// 				return
// 			}
// 		}
// 	}()
// }()
