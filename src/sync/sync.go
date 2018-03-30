package sync

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	gb "./GalendarBot"
	GoogleWrap "../GoogleWrap"
)

func initSync() {
	//выполняется через 1 час
	// не троготь цикл, он для таймера.
	// 2 минуты
	go SyncEventsFromJson(time.NewTicker(10 * time.Second))
	for {
	}
}

func SyncEventsFromJson(ticker *time.Ticker) {
	for range ticker.C {
		log.Println("Синхронизация сервера с Google-Календарь:")
		files, _ := filepath.Glob("./*.json") //место хранения файлов с событиями пользователей
		for _, file := range files {
			changed := false
			filenames := strings.Split(file, ".json")
			log.Println("Получили файл " + filenames[0])
			Work := gb.ReadJson(filenames[0])
			for area := 0; area < len(Work.Areas); area++ {
				log.Println("Получили сферу " + Work.Areas[area].AreaName)
				for event := 0; event < len(Work.Areas[area].Events); event++ {
					log.Println("Получили событие " + Work.Areas[area].Events[event].EventName)
					//--------На добавление
					//Если не размечено
					if !Work.Areas[area].Events[event].InTheCalendar {
						// if !event.InTheCalendar {
						//добавляем в календарь
						// GoogleWrap.Auth(Work.UserID)// так как авторизация внутри
						Id, Successed := GoogleWrap.AddEvent(gb.ConverToDOtoEvent(Work.Areas[area].Events[event]), Work.UserID)
						if Successed {
							// if GoogleWrap.AddEvent(gb.ConverToDOtoEvent(*event), Work.UserID) {
							log.Println("Дело " + Work.Areas[area].Events[event].EventName + " добавлено в " + Work.Areas[area].AreaName)
							//помечаем что размечено
							Work.Areas[area].Events[event].EventId = Id
							Work.Areas[area].Events[event].InTheCalendar = true
							//!!!!
							changed = true
							log.Println("Размечено = ", Work.Areas[area].Events[event].InTheCalendar)
							log.Println("Id события: " + Work.Areas[area].Events[event].EventId)
						} else {
							log.Println("Дело " + Work.Areas[area].Events[event].EventName + " не добавлено!!!")
						}
					} else {
						log.Println("Событие " + Work.Areas[area].Events[event].EventName + " уже есть в Google-Календаре.")
					}
					//------------
					//---------На изменение
					if Work.Areas[area].Events[event].Changed {
						if GoogleWrap.UpdateEvent(gb.ConverToDOtoEvent(Work.Areas[area].Events[event]), (Work.Areas[area].Events[event].EventId), Work.UserID) {
							changed = true
							Work.Areas[area].Events[event].Changed = false
							log.Println("Дело " + Work.Areas[area].Events[event].EventName + " изменено в " + Work.Areas[area].AreaName)
						} else {
							log.Println("Дело " + Work.Areas[area].Events[event].EventName + " НЕ изменено в " + Work.Areas[area].AreaName)
						}
					}
					//---------------
					//-----------На удаление
					//Если удалено
					if Work.Areas[area].Events[event].Deleted {
						//Если удалено успешно
						if GoogleWrap.DeleteEvent(Work.Areas[area].Events[event].EventId, Work.UserID) {
							log.Println("Дело " + Work.Areas[area].Events[event].EventName + " удалено из " + Work.Areas[area].AreaName)
							Work.Areas[area].Events = gb.Remv(Work.Areas[area].Events, event)
							changed = true
						} else {
							log.Println("Дело " + Work.Areas[area].Events[event].EventName + " НЕ удалено из " + Work.Areas[area].AreaName)
						}
					}
					//--------------------
				}
			}
			log.Println("-----Синхронизация сервера-----")
			changed2 := false
			Work, changed2 = SyncEventsFromGoogleCalendar(Work)
			//log.Println("!!!!Имя последнего события ", Work.Areas[0].Events[len(Work.Areas[0].Events)-1].EventName)
			//Work.Areas[0].Events[0].InTheCalendar = false
			//Если были изменения, то записываем
			if changed || changed2 {
				if gb.WriteJson(filenames[0], Work) {
					log.Println(time.Now().Format(time.RFC3339) + " Дела синхронизированиы успешно в " + filenames[0] + ".json")
				} else {
					log.Println(time.Now().Format(time.RFC3339) + " Ошибка синхронизации " + filenames[0] + ".json")
				}

			}

		}
	}
}

func SyncEventsFromGoogleCalendar(Work gb.JsonStruct) (gb.JsonStruct, bool) {
	log.Println("Синхронизация Google-Календаря с сервером:")
	changed := false
	if GoogleWrap.Auth(Work.UserID) {
		t := time.Now().Format(time.RFC3339)
		events, err := GoogleWrap.Srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalf("Не удается получить следующие десять событий пользователя. %v", err)
		}
		matched := false
		if len(events.Items) > 0 {
			//получаем события из Google календарь
			for _, event := range events.Items {
				matched = false
				//Ищем совпадение в json по Id
				for are := 0; are < len(Work.Areas); are++ {
					for eve := 0; eve < len(Work.Areas[are].Events); eve++ {
						//log.Println("Номер события ", eve)
						if event.Id == Work.Areas[are].Events[eve].EventId {
							matched = true
							//Если событие было изменено в гугле,
							//то есть гугл вермя после времени в json
							GalendarTime, _ := time.Parse(time.RFC3339, event.Updated)
							JsonTime, _ := time.Parse(time.RFC3339, Work.Areas[are].Events[eve].DateTimeLastChange)
							if GalendarTime.After(JsonTime) {
								//Изменяем событие
								Work.Areas[are].Events[eve].EventName = event.Summary
								Work.Areas[are].Events[eve].Description = event.Description
								Work.Areas[are].Events[eve].DateTimeStart = event.Start.DateTime
								Work.Areas[are].Events[eve].DateTimeEnd = event.End.DateTime
								Work.Areas[are].Events[eve].InTheCalendar = true
								Work.Areas[are].Events[eve].DateTimeLastChange = time.Now().Format(time.RFC3339)
								Work.Areas[are].Events[eve].ColorId = event.ColorId
								//отметка
								changed = true
							}
						}
					}
				}
				//если совпадения нет
				if !matched {

					var Jevent gb.Event
					Jevent.EventId = event.Id
					Jevent.EventName = event.Summary
					Jevent.Description = event.Description
					Jevent.DateTimeStart = event.Start.DateTime
					Jevent.DateTimeEnd = event.End.DateTime
					Jevent.InTheCalendar = true
					Jevent.DateTimeLastChange = time.Now().Format(time.RFC3339)
					Jevent.ColorId = event.ColorId
					Work.Areas[0].Events = append(Work.Areas[0].Events, Jevent)
					changed = true
					log.Println("Добавлено событие. Имя последнего события ", Work.Areas[0].Events[len(Work.Areas[0].Events)-1].EventName)
				}
			}
			//Удаляем события в json
			//Если событие было удалено, то оно было добавлено в гугл и их удаляем
			//Ищем совпадение в json

			for are := 0; are < len(Work.Areas); are++ {
				for eve := 0; eve < len(Work.Areas[are].Events); eve++ {
					matched = false
					for _, event := range events.Items {
						if event.Id == Work.Areas[are].Events[eve].EventId {
							matched = true
							log.Println("Совпадение с  " + event.Summary)
							break
						}
					}

					//Если совпадения нет, то удаляем
					if !matched {
						log.Println("Совпадений нет с  " + Work.Areas[are].Events[eve].EventName)
						if Work.Areas[are].Events[eve].InTheCalendar {
							log.Println("Дело " + Work.Areas[are].Events[eve].EventName + " " + Work.Areas[are].Events[eve].DateTimeStart + " удалено из " + Work.Areas[are].AreaName)
							Work.Areas[are].Events = gb.Remv(Work.Areas[are].Events, eve)
							//log.Println("Проверка Дело " + Work.Areas[are].Events[eve].EventName + " " + Work.Areas[are].Events[eve].DateTimeStart + " удалено из " + Work.Areas[are].AreaName)
							matched = false
							changed = true
						}
					}

				}
			}

		}

	}
	log.Println("-----Синхронизация Google-Календаря-----")
	return Work, changed
}
