package main

import (
	"log"
	"strings"
	"time"
)

type ToDo struct {
	Name      string
	Date      string
	Time      string
	Important string
	Area      string
}

func TryDateParse(Date string) bool {
	_, err := DateParse(Date)
	if err != nil {
		return true
	}
	return false
}

func DateParse(Date string) (time.Time, error) {
	Dates := strings.Split(Date, ".")
	var Day, Month, Year string
	if len(Dates) == 3 {
		Day = Dates[0]
		Month = Dates[1]
		Year = Dates[2]
	}
	//RFC3339     = "2006-01-02T15:04:05Z07:00"
	return time.Parse(time.RFC3339, Year+"-"+Month+"-"+Day+"T00:00:00Z08:00") //z часовой пояс
}

func ToDoMake(Word []string) ToDo {
	//TODO:Нужно сделать хоть как-то
	var Work ToDo
	//позиция даты
	k := 0
	for i := 1; i < len(Word); i++ {
		if TryDateParse(Word[i]) {
			k = i
			Work.Date = Word[k]
		} else {
			Work.Name += Word[i]
		}
	}
	//если длина текста больше 3 слов
	if k == 0 {
		Work.Date = ""
	}
	//если длина текста больше 4 слов
	if len(Word) >= 4 {
		Work.Time = Word[k]
		k++
	} else {
		Work.Time = ""
	}
	//если длина текста больше 5 слов
	if len(Word) >= 5 {
		Work.Important = Word[k]
		k++
	} else {
		Work.Important = ""
	}
	//если длина текста больше 6 слов
	if len(Word) >= 6 {
		for j := k; j < len(Word); j++ {
			Work.Area += Word[j]
		}
	} else {
		Work.Area = ""
	}
	return Work
}

func ParseText(Text string, UserName string) string {
	//TODO: Разделить строку на соответствующие части
	//Добавить Купить еды 15.06.12 16.45 Важно Домашние дела.

	Word := strings.Split(Text, " ") //"Дело, на дату, время, с пометкой важности")
	//действие по первому ключевому слову
	Work := ToDoMake(Word)
	switch Word[0] {
	case "Добавить":
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
				Date, _ := DateParse(Word[1])
				return "Дела на " + Word[1] + ":\n" + GetWorksOnDate(Date, UserName) //функцию
			}
		}
		return "Дела на сегодня:\n" + GetWorks(UserName) //функцию
	default:
		return "У нас проблема("
	}
	return "Failed"
}

func ParseCommand(Command string, UserName string) string {
	switch Command {
	case "start":
		//регистрация
		return "Привет. Что бы начать работать тебе нужно перейти по ссылке " + Auth(UserName) + " и авторизоваться в Google-Календаре." //Заглушка
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
	//запрос к bot_Galendar
	return "1. Написать бота"
}

func GetWorksOnDate(Dates time.Time, UserName string) string {
	//запрос к bot_Galendar
	return "1. Написать бота на дату " //+ Dates.Format(time.RFC822)
}

func GetWorksOnTomorrow(UserName string) string {
	//запрос к bot_Galendar
	return "1. Написать бота завтра"
}

func CreateWork(Work ToDo, UserName string) bool {
	//запрос к bot_Galendar
	return true
}

func ChangeWork(Work ToDo, UserName string) bool {
	//запрос к bot_Galendar
	return true
}

func DeleteWork(Work ToDo, UserName string) bool {
	//запрос к bot_Galendar
	return true
}
