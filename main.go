package main

import (
	"1841/pkg/semaphor"
	"fmt"
	"sync"
	"time"
)

// 5 воркеров должны отработать 10 запросов,
// но одновременно могут обрабатывать запрос только 3 воркера
const (
	maxWorkers     = 5  // Максимальное кол-во горутин-воркеров
	requestsNumber = 10 // Кол-во запросов, которое необходимо обработать воркерам
	maxResources   = 3  // Максимальное кол-во ресурсов, которое могут использовать воркеры
)

// Данные запроса
type data struct{}

// worker - горутина принимает запросы из канала и обрабатывает их.
// Семафор ограничивает кол-во одновременно обрабатываемых запросов
func worker(id int, c <-chan data, s *semaphor.Semaphore, wg *sync.WaitGroup) {
	defer wg.Done()

	makeStr := func(msg string) string {
		return fmt.Sprintf("Worker #%d -> %s", id, msg)
	}

	// Получаем запросы
	for range c {
		// Захватываем ресурс
		if err := s.Acquire(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(makeStr("захватили ресурс"))
		}

		time.Sleep(time.Second)

		// Освобождаем ресурс
		if err := s.Release(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(makeStr("освободили ресурс"))
		}
	}
}

func main() {
	// Канал, по которому будут поступать запросы к воркерам
	c := make(chan data, requestsNumber)

	// Семафор, ограничивающий кол-во одновременного использования ресурсов воркерами
	var s = semaphor.NewSemaphore(0, maxResources, time.Second*10)

	// Запускаем пул воркеров
	var wg sync.WaitGroup
	wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go worker(i, c, s, &wg)
	}

	// Заполняем канал запросами
	for i := 0; i < requestsNumber; i++ {
		c <- data{}
	}
	close(c)

	// Ожидаем завершения воркеров
	wg.Wait()
}
