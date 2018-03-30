package main

import (
	"log"
	"os"
	"strings"
	"time"

	"encoding/json"
	"io/ioutil"

	"./handler"
	"./GoogleWrap"
	//"./sync"
	calendar "google.golang.org/api/calendar/v3"
)

// type ToDo struct {
// 	Name      string
// 	Date      string
// 	Time      string
// 	Important string
// 	Area      string
// }

func TryDateParse(Date string) bool {
	// dt, err := DateParse(Date)
	// if err != nil {
	// 	return false
	// }
	// log.Printf("Дата: " + dt.Format(time.RFC3339))
	return true
}

func DateParse(Date string, Time string) string {
	Dates := strings.Split(Date, ".")
	Times := strings.Split(Time, ".")
	var Day, Month, Year, Hour, Min string
	if len(Dates) == 3 {
		Day = Dates[0]
		Month = Dates[1]
		Year = Dates[2]
	}
	if len(Times) == 2 {
		Hour = Times[0]
		Min = Times[1]
	}
	log.Println("День: " + Day)
	log.Println("Месяц: " + Month)
	log.Println("Год: " + Year)
	log.Println("Час: " + Hour)
	log.Println("Минуты: " + Min)
	//RFC3339     = "2006-01-02T15:04:05+08:00"
	times := Year + "-" + Month + "-" + Day + "T" + Hour + ":" + Min + ":00+08:00"
	log.Println("Время: " + times)
	//dt, err := time.Parse(time.RFC3339, times) //z часовой пояс
	//log.Printf("\nДата: " + dt.Format(time.RFC3339))
	return times
}
func DateTime(oldTime string, Type string) string {
	DateTimeZone := strings.Split(oldTime, "T")
	Date := strings.Split(DateTimeZone[0], "-")
	TimeZone := strings.Split(DateTimeZone[1], "+")
	Time := strings.Split(TimeZone[0], ":")
	switch Type {
	case "Year":
		return Date[0]
	case "Month":
		return Date[1]
	case "Day":
		return Date[2]
	case "Hour":
		return Time[0]
	case "Min":
		return Time[1]
	case "Zone":
		return TimeZone[1]
	case "Date":
		return Date[2] + "." + Date[1] + "." + Date[0]
	case "Time":
		return Time[0] + "." + Time[1]
	}
	return Type
}

func CreateEvent(Word []string) Event {
	//TODO:Нужно сделать хоть как-то
	var Work Event
	//позиция даты
	k := 0
	for i := 1; i < len(Word)-1; i++ {
		if len(strings.Split(Word[i], ".")) == 3 {
			k = i
			Work.DateTimeStart = DateParse(Word[k], Word[k+1])
			log.Println("Время старта события:\n" + Work.DateTimeStart)
			DateTimeEnd, _ := time.Parse(time.RFC3339, Work.DateTimeStart)
			log.Println("Время конца события:\n" + DateTimeEnd.Format(time.RFC3339))
			DateTimeEnd = DateTimeEnd.Add(time.Hour)
			log.Println("Время точного конца события:\n" + DateTimeEnd.Format(time.RFC3339))
			Work.DateTimeEnd = strings.Replace(DateTimeEnd.Format(time.RFC3339), "Z", "+", 1)
			log.Println("Время ОКОНЧАТЕЛЬНОГО конца события:\n" + Work.DateTimeEnd)
			Work.DateTimeLastChange = (time.Now()).Format(time.RFC3339) //Время создания/последнего изменения
		} else {
			Work.EventName += Word[i] + " "
			//log.Printf("Имя события:\n" + Work.Name + "\n")
		}
	}
	Work.EventName = strings.Trim(Work.EventName, " ")
	//если длина текста больше 3 слов
	if k == 0 {
		//Work.Date = ""
	}
	//если длина текста больше 4 слов
	// if len(Word) >= 4 {
	// 	Work.Time = Word[k]
	// 	k++
	// } else {
	// 	Work.Time = ""
	// }
	//если длина текста больше 5 слов

	if len(Word) >= k {

		Work.Importance = Word[k]
		k++
	} else {
		Work.Importance = ""
	}
	//если длина текста больше 6 слов
	// if len(Word) >= 6 {
	// 	for j := k; j < len(Word); j++ {
	// 		Work.Area += Word[j]
	// 	}
	// } else {
	// 	Work.Area = ""
	// }
	return Work
}

func ChangeEvent(Events []Event, Word []string) []Event {
	ChangedName := Word[1]
	//изменение имени события
	NewName := Word[2]
	var k int
	for event := 0; event < len(Events); event++ {

		log.Println("Получили событие " + Events[event].EventName)

		//Если нашли событие
		if Events[event].EventName == ChangedName {
			Events[event].EventName = NewName
			k = event
		}
	}
	//парсер на дату время
	for wrd := 2; wrd < len(Word); wrd++ {
		// var oldDate string
		// var newTime []string
		//Если Дата
		if len(strings.Split(Word[wrd], ".")) == 3 {
			// oldHour := DateTime(Events[k].DateTimeStart, "Hour")
			// oldMin := DateTime(Events[k].DateTimeStart, "Min")
			// log.Println("!!Час: " + oldHour)
			// log.Println("!!Минуты: " + oldMin)
			Events[k].DateTimeStart = DateParse(Word[wrd], DateTime(Events[k].DateTimeStart, "Time"))
			DateTimeEnd, _ := time.Parse("2006-01-02T15:04:05+08:00", Events[k].DateTimeStart)
			DateTimeEnd = DateTimeEnd.Add(time.Hour)
			Events[k].DateTimeEnd = DateTimeEnd.Format(time.RFC3339)
			Events[k].DateTimeLastChange = (time.Now()).Format(time.RFC3339) //Время создания/последнего изменения
		}
		//Если время
		if len(strings.Split(Word[wrd], ".")) == 2 {
			oldDate := DateTime(Events[k].DateTimeStart, "Date")
			log.Println("!Дата: " + oldDate)
			// newTime = strings.Split(Word[wrd], ".")
			Events[k].DateTimeStart = DateParse(oldDate, Word[wrd])
			DateTimeEnd, _ := time.Parse(time.RFC3339, Events[k].DateTimeStart)
			DateTimeEnd = DateTimeEnd.Add(time.Hour)
			log.Println("Время ОКОНЧАТЕЛЬНОГО конца события:\n" + DateTimeEnd.Format(time.RFC3339))
			Events[k].DateTimeEnd = strings.Replace(DateTimeEnd.Format(time.RFC3339), "Z", "+", 1)
			Events[k].DateTimeLastChange = (time.Now()).Format(time.RFC3339) //Время создания/последнего изменения
		}

		if strings.ToLower(Word[wrd]) != Events[k].Importance {
			Events[k].Importance = strings.ToLower(Word[wrd])
		}
	}

	Events[k].Changed = true
	return Events
}

func DeleteEvent(Events []Event, Word []string) []Event {
	DeletedName := Word[1]
	var MatchedEventId []int
	for event := 0; event < len(Events); event++ {

		log.Println("Получили событие " + Events[event].EventName)

		//Если нашли событие
		if Events[event].EventName == DeletedName {
			MatchedEventId = append(MatchedEventId, event)
			log.Println("Нашли событие " + Events[event].EventName)
		}
	}
	if len(MatchedEventId) > 1 {
		//Предлагаем на удаление
	} else {
		//удаляем 1 событие

		//Удаляем через синхронизацию
		Events[MatchedEventId[0]].Deleted = true

		//Удаление здесь
		//Events = Remv(Events, MatchedEventId[0])
	}
	return Events
}

//удаление элемента из массива без сортировки
//то есть просто последний перемещаем на место удаляемого и возвращаем на 1 короче
func Remv(events []Event, id int) []Event {
	log.Println("Удаляемое " + events[id].EventName)
	log.Println("Вставляемое " + events[len(events)-1].EventName)
	events[id] = events[len(events)-1]
	log.Println("Получили " + events[id].EventName)
	return events[:len(events)-1]
}

func ParseText(Text string, UserID string) string {
	//TODO: Разделить строку на соответствующие части
	//Добавить Купить еды 15.06.12 16.45 Важно Домашние дела.
	var Work JsonStruct
	Word := strings.Split(Text, " ") //"Дело, на дату, время, с пометкой важности")
	//действие по первому ключевому слову
	Work = ReadJson(UserID)
	// log.Println(Work.Areas[])

	//Work.Areas[0].Events = append(Work.Areas[0].Events, CreateEvent(Word))
	// log.Printf(Work.Date.Format(time.RFC3339))
	switch Word[0] {
	case "Код", "код":
		if GoogleWrap.SaveToken(Word[1], UserID) {
			return "Авторизирован"
		}
		break
	case "Добавить":
		Work.Areas[0].Events = append(Work.Areas[0].Events, CreateEvent(Word))
		if WriteJson(UserID, Work) {
			// 	if CreateWork(Work, UserID) {
			return "Добавлено." //successful
			// 	} else {
			// 		break //go to the Failed
			// 	}
		} else {
			break //go to the Failed
		}
		if CreateWork(Work, UserID) {
			return "Добавлено." //successful
		} else {
			break //go to the Failed
		}

	case "Изменить":
		// if ChangeWork(Work, UserID) {
		Work.Areas[0].Events = ChangeEvent(Work.Areas[0].Events, Word)
		if WriteJson(UserID, Work) {
			return "Изменено." //successful
		} else {
			break //go to the Failed
		}
	case "Удалить":
		Work.Areas[0].Events = DeleteEvent(Work.Areas[0].Events, Word)
		if WriteJson(UserID, Work) {
			return "Удалено." //successful
		} else {
			break //go to the Failed
		}
	case "Помоги", "Help":
		return GetHelp()
	case "Дела":
		if len(Word) > 1 {
			if Word[1] == "на" {
				return "Дела на завтра:\n" + GetWorksOnTomorrow(UserID) //функцию
			}
			//TODO: Распасрсить дату(((((
			//_, err:=time.Parse("01.01.1970",Word[1])
			//if err!=nil{
			if TryDateParse(Word[1]) {
				log.Printf("Ошибка:" + Word[1])
				Date := DateParse(Word[1], "00.00")
				return "Дела на " + Word[1] + ":\n" + GetWorksOnDate(Date, UserID) //функцию
			}
		}
		return "Дела на сегодня:\n" + GetWorks(UserID) //функцию
	default:
		return "У нас проблема("
	}
	return "Failed"
}

//Структура json-файла пользователя
type JsonStruct struct {
	UserID   string   //Имя пользователя
	Settings settings //Настройки
	Areas    []Area   //Сфры дейтельности

}

//Структура сфер дейтельности
type Area struct {
	AreaName string
	Events   []Event //Дела
}

//Структура дел	*-по Google API
type Event struct {
	EventId            string                 //ИД события*
	EventName          string                 //Название дела
	DateTimeLastChange string                 //Дата и время последнего изменения
	Description        string                 //Описания*
	DateTimeStart      string                 //Время начала дела
	Start              calendar.EventDateTime //Время начала дела*
	DateTimeEnd        string                 //+1час к старту  //Время завершения дела   (:=TimeStart.Add(time.Hour))
	End                calendar.EventDateTime //Время завершения дела*
	Importance         string                 //Важность
	ColorId            string                 //Важность, но в цвете*
	//Area               string                 //Сфера дейтельности
	InTheCalendar bool   //добавлено ли в G-календарь?
	Deleted       bool   //Удалён?
	Changed       bool   //Изменён
	Created       string //Время создания события* RFC3339
	Updated       string //Время изменения события* RFC3339
}

//Структура настроек
type settings struct {
	CountWorkPerDay int
}

func ReadJson(UserID string) JsonStruct {
	file, _ := os.Open(UserID + ".json")
	defer file.Close()

	JsonFile, _ := ioutil.ReadAll(file)
	var JsonStructIn JsonStruct
	json.Unmarshal(JsonFile, &JsonStructIn)
	return JsonStructIn
}

func CreateUser(UserID string) bool {
	_, err := os.Create(UserID + ".json")
	if err != nil {
		return false
	}
	var Work JsonStruct
	Work.UserID = UserID
	Work.Areas = append(Work.Areas, Area{"Неразмечено", nil})
	WriteJson(UserID, Work)
	return true
}

//надо придумать название
func ddd(UserID string) {
	if CreateUser(UserID) {
		log.Println("Успешно создан файл пользователя " + UserID)
	} else {
		log.Println("Не удалось создать файл пользователя " + UserID)
	}
	// var Work JsonStruct
	// js, _ := json.Marshal(Work)
	// ioutil.WriteFile(UserID+".json", js, 0644)

}

func WriteJson(UserID string, Work JsonStruct) bool {
	js, err := json.Marshal(Work)
	if err != nil {
		log.Println("Failed Json Marshal")
		return false
	}
	ioutil.WriteFile(UserID+".json", js, 0644)
	return true
}

func ParseCommand(Command string, UserID string) string {
	switch Command {
	case "start":
		//регистрация
		ddd(UserID)
		return "Привет. Что бы начать работать тебе нужно перейти по ссылке \n" + GoogleWrap.GetTokenURL() + " \nи авторизоваться в Google-Календаре." //Заглушка

	}
	return "Упс."
}

func GetHelp() string {
	return "Существуютт следущие способы работы: \n" +
		"1) \"Добавить ИМЯ ДАТА ВРЕМЯ ВАЖНОСТЬ СФЕРА\"\n" +
		"2) \"Изменить ИМЯ ДАТА ВРЕМЯ ВАЖНОСТЬ СФЕРА\"\n" +
		"3) \"Удалить ИМЯ ДАТА\""
}

func Auth(UserID string) string {
	//запрос к bot_Galendar
	return "URL:edrfyujilkuydfasthd"
}

func GetWorks(UserID string) string {
	GoogleWrap.Auth(UserID)
	return GoogleWrap.ShowEvents()
	//запрос к bot_Galendar
	return "1. Написать бота"
}

func GetWorksOnDate(Dates string, UserID string) string {
	//запрос к bot_Galendar
	return "1. Написать бота на дату " //+ Dates.Format(time.RFC822)
}

func GetWorksOnTomorrow(UserID string) string {
	//запрос к bot_Galendar
	return "1. Написать бота завтра"
}

func CreateWork(Work JsonStruct, UserID string) bool {
	// GoogleWrap.Auth(UserID)
	//Workevent := ConverToDOtoEvent(Work.Areas[0].Events[0])
	//log.Printf("\n11111111111111111111111\n" + Work.UserID)
	//return GoogleWrap.AddEvent(Workevent, Work.UserID)
	//запрос к bot_Galendar
	return true
}

func ChangeWork(Work JsonStruct, UserID string) bool {
	//запрос к bot_Galendar
	return true
}

func DeleteWork(Work JsonStruct, UserID string) bool {
	//запрос к bot_Galendar
	return true
}

func ConverToDOtoEvent(Event Event) *calendar.Event {
	event := &calendar.Event{
		Summary:     Event.EventName,
		Description: Event.EventName,
		Start: &calendar.EventDateTime{
			DateTime: Event.DateTimeStart,
			TimeZone: "Asia/Irkutsk",
		},
		End: &calendar.EventDateTime{
			DateTime: Event.DateTimeEnd,
			TimeZone: "Asia/Irkutsk",
		},
	}

	return event
}

func prepareForParse(msg string) []string {
	return strings.Split(msg, ":")
}

func main() {

	log.Println("Starting service Calendar")
	
	/*
		initialize handler
		we sent string channel to Handler.
		handler load to this channel message from telegram's users
	*/
	var msgCh chan string	
	msgCh = make(chan string)
	go handler.Handle(msgCh)
	
	for {
		text := <-msgCh
		divided := prepareForParse(text)
		msg, UserID := divided[0], divided[1]
		
		var reply string
		if len(divided) == 3 {
			reply = ParseCommand(msg, UserID)
			log.Println(reply)
		} else {
			reply = ParseText(msg, UserID)
			log.Println(reply)
		}
		
		msgCh <- reply		
			
	}
	
}
