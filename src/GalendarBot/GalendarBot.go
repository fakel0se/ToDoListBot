package GalendarBot

//1. криво добавляет имя события, пробел и время
//2. добавить событие к уже существующему json

import (
	"log"
	"os"
	"strings"
	"time"

	"encoding/json"
	"io/ioutil"

	"./GoogleWrap"
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

func WorkMake(Word []string) event {
	//TODO:Нужно сделать хоть как-то
	var Work event
	//позиция даты
	k := 0
	for i := 1; i < len(Word)-1; i++ {
		if len(strings.Split(Word[i], ".")) == 3 {
			k = i
			Work.DateTimeStart = DateParse(Word[k], Word[k+1])
			DateTimeEnd, _ := time.Parse("2006-01-02T15:04:05+08:00", Work.DateTimeStart)
			DateTimeEnd = DateTimeEnd.Add(time.Hour)
			Work.DateTimeEnd = DateTimeEnd.Format(time.RFC3339)
			Work.DateTimeLastChange = (time.Now()).Format(time.RFC3339) //Время создания/последнего изменения
		} else {
			Work.EventName += Word[i] + " "
			//log.Printf("Имя события:\n" + Work.Name + "\n")
		}
	}
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
	if len(Word) >= 5 {
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

func ParseText(Text string, UserName string) string {
	//TODO: Разделить строку на соответствующие части
	//Добавить Купить еды 15.06.12 16.45 Важно Домашние дела.
	var Work JsonStruct
	Word := strings.Split(Text, " ") //"Дело, на дату, время, с пометкой важности")
	//действие по первому ключевому слову
	Work = ReadJson(UserName)
	// log.Println(Work.Areas[])
	if Work.Areas == nil {
		ar := area{"Неразмечено", nil}
		Work.Areas = append(Work.Areas, ar)
	}
	Work.Areas[0].Events = append(Work.Areas[0].Events, WorkMake(Word))
	// log.Printf(Work.Date.Format(time.RFC3339))
	switch Word[0] {
	case "Код":
		if GoogleWrap.SaveToken(Word[1], UserName) {
			return "Авторизирован"
		}
		break
	case "Добавить":
		if WriteJson(UserName, Work) {
			// 	if CreateWork(Work, UserName) {
			return "Добавлено." //successful
			// 	} else {
			// 		break //go to the Failed
			// 	}
		} else {
			break //go to the Failed
		}
		if CreateWork(Work, UserName) {
			return "Добавлено." //successful
		} else {
			break //go to the Failed
		}

	case "Изменить":
		if ChangeWork(Work, UserName) {
			return "Изменено." //successful
		} else {
			break //go to the Failed
		}
	case "Удалить":
		if DeleteWork(Work, UserName) {
			return "Удалено." //successful
		} else {
			break //go to the Failed
		}
	case "Помоги", "Help":
		return GetHelp()
	case "Дела":
		if len(Word) > 1 {
			if Word[1] == "на" {
				return "Дела на завтра:\n" + GetWorksOnTomorrow(UserName) //функцию
			}
			//TODO: Распасрсить дату(((((
			//_, err:=time.Parse("01.01.1970",Word[1])
			//if err!=nil{
			if TryDateParse(Word[1]) {
				log.Printf("Ошибка:" + Word[1])
				Date := DateParse(Word[1], "00.00")
				return "Дела на " + Word[1] + ":\n" + GetWorksOnDate(Date, UserName) //функцию
			}
		}
		return "Дела на сегодня:\n" + GetWorks(UserName) //функцию
	default:
		return "У нас проблема("
	}
	return "Failed"
}

//Структура json-файла пользователя
type JsonStruct struct {
	UserName string   //Имя пользователя
	Settings settings //Настройки
	Areas    []area   //Сфры дейтельности

}

//Структура сфер дейтельности
type area struct {
	AreaName string
	Events   []event //Дела
}

//Структура дел	*-по Google API
type event struct {
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
	Created       string //Время создания события* RFC3339
	Updated       string //Время изменения события* RFC3339
}

//Структура настроек
type settings struct {
	CountWorkPerDay int
}

func ReadJson(UserName string) JsonStruct {
	file, _ := os.Open(UserName + ".json")
	defer file.Close()

	JsonFile, _ := ioutil.ReadAll(file)
	var JsonStructIn JsonStruct
	json.Unmarshal(JsonFile, &JsonStructIn)
	return JsonStructIn
}

func CreateUser(UserName string) bool {
	_, err := os.Create(UserName + ".json")
	if err != nil {
		return false
	}
	var Work JsonStruct
	Work.UserName = UserName
	WriteJson(UserName, Work)
	return true
}

//надо придумать название
func ddd(UserName string) {
	if CreateUser(UserName) {
		log.Println("Успешно создан файл пользователя " + UserName)
	} else {
		log.Println("Не удалось создать файл пользователя " + UserName)
	}
	// var Work JsonStruct
	// js, _ := json.Marshal(Work)
	// ioutil.WriteFile(UserName+".json", js, 0644)

}

func WriteJson(UserName string, Work JsonStruct) bool {
	js, err := json.Marshal(Work)
	if err != nil {
		log.Println("Failed Json Marshal")
		return false
	}
	ioutil.WriteFile(UserName+".json", js, 0644)
	return true
}

func ParseCommand(Command string, UserName string) string {
	switch Command {
	case "start":
		//регистрация
		ddd(UserName)
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

func Auth(UserName string) string {
	//запрос к bot_Galendar
	return "URL:edrfyujilkuydfasthd"
}

func GetWorks(UserName string) string {
	GoogleWrap.Auth(UserName)
	return GoogleWrap.ShowEvents()
	//запрос к bot_Galendar
	return "1. Написать бота"
}

func GetWorksOnDate(Dates string, UserName string) string {
	//запрос к bot_Galendar
	return "1. Написать бота на дату " //+ Dates.Format(time.RFC822)
}

func GetWorksOnTomorrow(UserName string) string {
	//запрос к bot_Galendar
	return "1. Написать бота завтра"
}

func CreateWork(Work JsonStruct, UserName string) bool {
	GoogleWrap.Auth(UserName)
	Workevent := ConverToDOtoEvent(Work.Areas[0].Events[0])
	log.Printf("\n11111111111111111111111\n" + Work.UserName)
	return GoogleWrap.AddEvent(Workevent)
	//запрос к bot_Galendar
	return true
}

func ChangeWork(Work JsonStruct, UserName string) bool {
	//запрос к bot_Galendar
	return true
}

func DeleteWork(Work JsonStruct, UserName string) bool {
	//запрос к bot_Galendar
	return true
}

func ConverToDOtoEvent(Event event) *calendar.Event {
	event := &calendar.Event{
		Summary:     Event.EventName,
		Description: "Event.Name",
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
