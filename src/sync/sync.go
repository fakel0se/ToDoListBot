package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	gb "../GalendarBot"
	GoogleWrap "../GalendarBot/GoogleWrap"
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

// // вызов каждый день в 10 утра функцию myfunc
// func main() {
// 	err := callAt(10, 0, 0, myfunc)
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 	}
// 	//для дальнейшей работы программы.
// 	time.Sleep(time.Hour * 24)
// }

func main() {
	//worker выполняется через 1сек
	//не троготь цикл, он для таймера.
	// go worker(time.NewTicker(1 * time.Hour))
	// for {
	// }
	// worker выполняется через 1 час
	// не троготь цикл, он для таймера.
	// 2 минуты
	go SyncEventsFromJson(time.NewTicker(2 * time.Minute))
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

func SyncEventsFromJson(ticker *time.Ticker) {
	for range ticker.C {

		files, _ := filepath.Glob("../*.json") //место хранения файлов с событиями пользователей
		for _, file := range files {
			changed := false
			filenames := strings.Split(file, ".json")
			log.Println("Получили файл " + filenames[0])
			Work := gb.ReadJson(filenames[0])
			// for event := 1; event < len(Word)-1; event++ {
			for area := 0; area < len(Work.Areas); area++ {
				// var area *gb.Area
				// for _, *area = range Work.Areas {
				log.Println("Получили сферу " + Work.Areas[area].AreaName)
				// log.Println("Получили сферу " + area.AreaName)
				for event := 0; event < len(Work.Areas[area].Events); event++ {
					// var event *gb.Event
					// for _, *event = range area.Events {
					log.Println("Получили событие " + Work.Areas[area].Events[event].EventName)
					// log.Println("Получили событие " + event.EventName)
					//Если не размечено
					if !Work.Areas[area].Events[event].InTheCalendar {
						// if !event.InTheCalendar {
						//добавляем в календарь
						// GoogleWrap.Auth(Work.UserName)// так как авторизация внутри
						if GoogleWrap.AddEvent(gb.ConverToDOtoEvent(Work.Areas[area].Events[event]), Work.UserName) {
							// if GoogleWrap.AddEvent(gb.ConverToDOtoEvent(*event), Work.UserName) {
							log.Println("Дело " + Work.Areas[area].Events[event].EventName + " добавлено в " + Work.Areas[area].AreaName)
							//помечаем что размечено
							Work.Areas[area].Events[event].InTheCalendar = true
							changed = true
							log.Println("Размечено = ", Work.Areas[area].Events[event].InTheCalendar)
						} else {
							log.Println("Дело " + Work.Areas[area].Events[event].EventName + " не добавлено!!!")
						}
					}
				}
			}
			changed2 := false
			Work, changed2 = SyncEventsFromGoogleCalendar(Work)
			log.Println("!!!!Имя последнего события ", Work.Areas[0].Events[len(Work.Areas[0].Events)-1].EventName)
			//Work.Areas[0].Events[0].InTheCalendar = false
			//Если были изменения, то записываем
			if changed || changed2 {
				if gb.WriteJson(filenames[0], Work) {
					log.Println(time.Now().Format(time.RFC3339) + " Дела синхронизированиы успешно в " + filenames[0] + ".json")
				}
			}

		}
	}
}

func SyncEventsFromGoogleCalendar(Work gb.JsonStruct) (gb.JsonStruct, bool) {
	changed := false
	if GoogleWrap.Auth(Work.UserName) {
		t := time.Now().Format(time.RFC3339)
		events, err := GoogleWrap.Srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalf("Не удается получить следующие десять событий пользователя. %v", err)
		}

		// buffer := ""
		if len(events.Items) > 0 {
			for _, event := range events.Items {
				matched := false
				for eve := 0; eve < len(Work.Areas[0].Events); eve++ {
					log.Println("Номер события ", eve)
					if event.Summary == Work.Areas[0].Events[eve].EventName {
						matched = true
					}
				}
				if !matched {
					var Jevent gb.Event
					Jevent.EventName = event.Summary
					Jevent.Description = event.Description
					Jevent.DateTimeStart = event.Start.DateTime
					Jevent.DateTimeEnd = event.End.DateTime
					Jevent.InTheCalendar = true
					Jevent.DateTimeLastChange = time.Now().Format(time.RFC3339)
					Jevent.ColorId = event.ColorId
					// var when string
					// //если время пустая строка событие является на весь день, так как доступна только дата
					// if event.Start.DateTime != "" {
					// 	when = event.Start.DateTime
					// } else {
					// 	when = event.Start.Date
					// }
					// fmt.Printf("%s (%s) %s\n", event.Summary, when, event.Id)
					// buffer = fmt.Sprint(buffer + event.Summary + " " + when + " " + event.Id + " " + "\n")
					// if Work.Areas == nil {
					// 	ar := Area{"Неразмечено", nil}
					// 	Work.Areas = append(Work.Areas, ar)
					// }
					Work.Areas[0].Events = append(Work.Areas[0].Events, Jevent)
					changed = true
					log.Println("Имя последнего события ", Work.Areas[0].Events[len(Work.Areas[0].Events)-1].EventName)
				}
			}
		}
	}
	return Work, changed
}
