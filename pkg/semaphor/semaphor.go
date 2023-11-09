package semaphor

import (
	"fmt"
	"log"
	"time"
)

var errAcquire = fmt.Errorf("Не удалось захватить семафор")

// Semaphore — структура двоичного семафора
type Semaphore struct {
	// Семафор — абстрактный тип данных,
	// в нашем случае в основе его лежит канал
	sem chan struct{}
	// Время ожидания основных операций с семафором, чтобы не
	// блокировать
	// операции с ним навечно (необязательное требование, зависит от
	// нужд программы)
	timeout time.Duration
}

// Acquire — метод захвата семафора
func (s *Semaphore) Acquire() error {
	select {
	case s.sem <- struct{}{}:
		return nil
	case <-time.After(s.timeout):
		return errAcquire
	}
}

// Release — метод освобождения семафора
func (s *Semaphore) Release() error {
	select {
	case <-s.sem:
		return nil
	case <-time.After(s.timeout):
		return errAcquire
	}
}

// NewSemaphore — функция создания семафора
func NewSemaphore(initialCount, maxCount int, timeout time.Duration) *Semaphore {
	// Проверяем initialCount и maxCount на граничные значения
	switch {
	case initialCount < 0:
		log.Fatalln("initialCount не может быть отрицательным")
	case maxCount < 1:
		log.Fatalln("maxCount должен быть положительным")
	case maxCount < initialCount:
		log.Fatalln("maxCount не может быть меньше initialCount")
	}

	s := &Semaphore{
		sem:     make(chan struct{}, maxCount),
		timeout: timeout,
	}

	// Занимаем начальное кол-во ресурсов
	for i := 0; i < initialCount; i++ {
		s.sem <- struct{}{}
	}

	return s
}
