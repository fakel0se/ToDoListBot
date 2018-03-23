package main

import (
	"fmt"
	"log"
	"time"
)

// Вызов переданной функции раз в сутки в указанное время.
func callAt(hour, min, sec int, f func()) error {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return err
	}

	// Вычисляем время первого запуска.
	now := time.Now().Local()
	firstCallTime := time.Date(
		now.Year(), now.Month(), now.Day(), hour, min, sec, 0, loc)
	if firstCallTime.Before(now) {
		// Если получилось время раньше текущего, прибавляем сутки.
		firstCallTime = firstCallTime.Add(time.Hour * 24)
	}

	// Вычисляем временной промежуток до запуска.
	duration := firstCallTime.Sub(time.Now().Local())

	go func() {
		time.Sleep(duration)
		for {
			f()
			// Следующий запуск через сутки.
			time.Sleep(time.Hour * 24)
		}
	}()

	return nil
}

//функция вызываемая каждый день.
func myfunc(eventTime string) {
	nowTime := time.Now() //.Format("2006.01.02T15:04:05+08:00")
	eventTime = "2018-03-24T09:00:00+08:00"

	eventT, err := time.Parse(time.RFC3339, eventTime)
	if err != nil {
		log.Fatalf("%v", err)
	}

	//оповещение будет если до события 2 часа и менее
	if int(nowTime.Sub(eventT)/time.Hour) <= 24 {
		fmt.Println("скоро событие")
	}
	//fmt.Println(int(nowTime.Sub(eventT) / time.Hour))
}

// вызов каждый день в 10 утра функцию myfunc
func main() {
	err := callAt(10, 0, 0, myfunc)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	//для дальнейшей работы программы.
	time.Sleep(time.Hour * 24)
}

/*
func main() {
	//worker выполняется через 1сек
	//не троготь цикл, он для таймера.
	go worker(time.NewTicker(1 * time.Second))
	for {
	}
}

func worker(ticker *time.Ticker) {

	for range ticker.C {
		nowTime := time.Now() //.Format("2006.01.02T15:04:05+08:00")
		eventTime := "2018-03-24T09:00:00+08:00"

		eventT, err := time.Parse(time.RFC3339, eventTime)
		if err != nil {
			log.Fatalf("%v", err)
		}

		//оповещение будет если до события 2 часа и менее
		if int(nowTime.Sub(eventT)/time.Hour) < 2 {
			fmt.Println("скоро событие")
		}
		//fmt.Println(int(nowTime.Sub(eventT) / time.Hour))
	}
}
*/
