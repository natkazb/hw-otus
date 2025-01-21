package hw06pipelineexecution

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

		// Запускаем горутину для обработки текущего этапа
		go func(in In, stage Stage, out Bi) { // число горутин == числу stages
			defer close(out)

			for v := range in {
				select {
					case <-done: return
					default: 
						// Создаем временный канал для передачи значения в stage
						tmpCh := make(Bi, 1)
						tmpCh <- v
						close(tmpCh)

						// Получаем результат этапа и отправляем его дальше
						result := stage(tmpCh)
						// Пытаемся отправить результат
						out <- <-result
				}
			}
		}(current, stage, intermediate)

		current = intermediate
	}

	return current
}

func ExecutePipelineChat(in In, done In, stages ...Stage) Out {
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

		// Запускаем горутину для обработки текущего этапа
		go func(in In, stage Stage, out Bi) { // число горутин == числу stages
			defer close(out)

			for v := range in {
				// Проверяем сигнал done
				/* select {
				case <-done:
					return
				default:
				} */

				// Создаем временный канал для передачи значения в stage
				tmpCh := make(chan interface{}, 1)
				tmpCh <- v
				close(tmpCh)

				// Получаем результат этапа и отправляем его дальше
				result := stage(tmpCh)

				// Пытаемся отправить результат
				select {
				case out <- <-result:
				case <-done:
					return
				}
			}
		}(current, stage, intermediate)

		current = intermediate
	}

	return current
}
