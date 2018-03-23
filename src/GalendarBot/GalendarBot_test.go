package GalendarBot

import (
	"log"
	"strings"
	"testing"
)

func WorkMake_test(t *testing.T) {
	var v Sevent
	v = WorkMake(strings.Split("Добавить сломать бота 23.05.2018 22.16", " "))
	log.Println(v.EventName)
}
