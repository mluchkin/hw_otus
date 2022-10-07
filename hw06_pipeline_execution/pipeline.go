package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		stageData := StageData(done, in)
		in = stage(stageData)
	}
	return in
}

func StageData(done In, in In) Bi {
	StageData := make(Bi)
	go func() {
		defer close(StageData)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				StageData <- v
			}
		}
	}()

	return StageData
}
