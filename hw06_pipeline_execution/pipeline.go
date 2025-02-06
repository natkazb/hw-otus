package hw06pipelineexecution

// import "fmt"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	current := in
	for _, stage := range stages { // цикл по всем stages
		if stage == nil {
			continue
		}

		// Создаем промежуточный канал для этапа
		intermediate := make(Bi)

		// Запускаем горутину для обработки текущего stage
		go func(in In, out Bi) {
			defer close(out)

			for v := range in {
				select {
				case <-done:
					for i := range in {
						_ = i
					}
					return
				case out <- v:
				}
			}
		}(current, intermediate)

		current = stage(intermediate)
	}

	return current
}
