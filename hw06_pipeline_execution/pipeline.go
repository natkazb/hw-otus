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
		go func(in In, stage Stage, out Bi) {
			defer close(out)

			for v := range in {
				// fmt.Printf("stage %d in: %v \n", i, v)
				// Проверяем сигнал done
				select {
				case <-done:
					// fmt.Printf("stage %d done1 \n", i)
					return
				default:
				}

				// Создаем временный канал для передачи значения в stage
				tmpCh := make(Bi, 1)
				tmpCh <- v
				close(tmpCh)
				// Получаем результат этапа
				result := stage(tmpCh)

				// Пытаемся отправить результат
				select {
				case <-done:
					// fmt.Printf("stage %d done2 \n", i)
					return
				case out <- <-result:
					// fmt.Printf("stage %d out: (%v) \n", i, v)
				}
			}
		}(current, stage, intermediate)

		current = intermediate
	}

	return current
}
