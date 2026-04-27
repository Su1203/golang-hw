package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(wrapWithDone(in, done))
	}
	return wrapWithDone(in, done)
}

func wrapWithDone(in In, done In) Out {
	if done == nil {
		return in
	}
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}
