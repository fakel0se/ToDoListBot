package main

import (
	"log"

	"github.com/Syfaro/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("559800435:AAE_aExKTPXbcwEto2qsHTHux_Wlh5McQic")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	//TODO: Ожидание входящего сообщения в цикле. А так же сделать получение именно последнего сообщения
	// инициализируем канал, куда будут прилетать обновления от API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	// updch, err := bot.GetUpdatesChan(ucfg) //updates channel
	// читаем обновления из канала
	// for {
	var updates []tgbotapi.Update
	updates, err = bot.GetUpdates(ucfg)
	if err != nil {
		log.Panic(err)
	}
	update := updates[1] //номер сообщение
	// select {
	// case <-updates:
	// Пользователь, который написал боту
	UserName := update.Message.From.UserName

	// ID чата/диалога.
	// Может быть идентификатором как чата с пользователем
	// (тогда он равен UserID) так и публичного чата/канала
	ChatID := update.Message.Chat.ID

	// Текст сообщения
	Text := update.Message.Text

	log.Printf("[%s] %d %s", UserName, ChatID, Text)

	// Ответим пользователю его же сообщением
	reply := Text
	// Созадаем сообщение
	msg := tgbotapi.NewMessage(ChatID, reply)
	// и отправляем его
	bot.Send(msg)
	// }
	// }
}
