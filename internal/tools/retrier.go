package tools

import "time"

var DefaulRetrier = NewRetrier(time.Second, time.Second*3, time.Second*5)

// NewRetrier создаёт новый ретраер.
func NewRetrier(pattern ...time.Duration) *Retrier {
	return &Retrier{pattern: pattern}
}

// execFuncT функция, которую необходимо ретраить.
// Если функция вернула true, необходимо продолжать ретраи.
type execFuncT = func() bool

// Retrier структура, описывающая ретраер.
// Вызывает повторный вызов переданной функции
// через временные промежутки, заданные в pattern.
type Retrier struct {
	pattern []time.Duration
}

// Exec запускает вызов функции execFunc с последующими ретраями.
func (r *Retrier) Exec(execFunc execFuncT) {
	if execFunc() {
		for _, delay := range r.pattern {
			time.Sleep(delay)

			if execFunc() {
				return
			}
		}
	}
}
