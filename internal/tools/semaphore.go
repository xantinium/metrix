package tools

// Semaphore структура семафора.
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore создает новый семафор.
func NewSemaphore(max int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, max),
	}
}

// Acquire увеличивает значение семафора на 1.
// Если текущее значение превышает max, исполнение горутрины
// приостановится, в ожидании вызова Release другой горутиной.
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Acquire уменьшает значение семафора на 1.
func (s *Semaphore) Release() {
	<-s.ch
}
