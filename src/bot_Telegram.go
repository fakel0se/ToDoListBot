package main

//go run bot_Telegram.go GalendarBot.go
import (
	"log"
	"./GalendarBot"
	"github.com/Syfaro/telegram-bot-api"
	"fmt"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("559800435:AAE_aExKTPXbcwEto2qsHTHux_Wlh5McQic")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	// инициализируем канал, куда будут прилетать обновления от API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	// читаем обновления из канала
	updates, err := bot.GetUpdatesChan(ucfg)
	if err != nil {
		log.Panic(err)
	}
	for update := range updates {
		// Пользователь, который написал боту
		UserName := update.Message.From.UserName
		userID := update.Message.From.ID
		
		// ID чата/диалога.
		// Может быть идентификатором как чата с пользователем
		// (тогда он равен UserNameID) так и публичного чата/канала
		ChatID := update.Message.Chat.ID
		var reply string
		var Text string
		if update.Message.Command() != "" {
			reply = GalendarBot.ParseCommand(update.Message.Command(), fmt.Sprint(userID))
		} else {
			//Получем сообщение и парсим его
			Text = update.Message.Text
			reply = GalendarBot.ParseText(Text, fmt.Sprint(userID))
		}

		log.Printf("[%s] %d %s", UserName, ChatID, Text)

		// Созадаем сообщение
		msg := tgbotapi.NewMessage(ChatID, reply)
		// и отправляем его
		bot.Send(msg)
	}
}
