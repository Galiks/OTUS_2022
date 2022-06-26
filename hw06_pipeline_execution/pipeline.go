package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	temp := in
	go func() {
		defer close(out)
		for _, s := range stages {
			temp = s(temp)
		}
		out <- <-temp
	}()

	return out
}
