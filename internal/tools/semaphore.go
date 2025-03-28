package tools

// Semaphore структура семафора.
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore создает новый семафор.
// Если max = 0, то блокировок не будет.
func NewSemaphore(max int) *Semaphore {
	if max < 0 {
		panic("semaphore value cannot be negative")
	}

	if max == 0 {
		return &Semaphore{}
	}

	return &Semaphore{
		ch: make(chan struct{}, max),
	}
}

// Acquire увеличивает значение семафора на 1.
// Если текущее значение превышает max, исполнение горутрины
// приостановится, в ожидании вызова Release другой горутиной.
func (s *Semaphore) Acquire() {
	if s.ch != nil {
		s.ch <- struct{}{}
	}
}

// Acquire уменьшает значение семафора на 1.
func (s *Semaphore) Release() {
	if s.ch != nil {
		<-s.ch
	}
}
