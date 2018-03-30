package handler

//go run bot_Telegram.go GalendarBot.go
import (
	"fmt"
	"log"

	"github.com/Syfaro/telegram-bot-api"
)

func Handle(msgCh chan string) {
	bot, err := tgbotapi.NewBotAPI("559800435:AAE_aExKTPXbcwEto2qsHTHux_Wlh5McQic") 
	//bot, err := tgbotapi.NewBotAPI("550460139:AAEG56gf2hI2NyjmpbeAkpxQR7hMlNdNhyU")
	if err != nil {
		log.Println(err)
	}
	
	// bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	// инициализируем канал, куда будут прилетать обновления от API
	
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	// читаем обновления из канала
	
	updates, err := bot.GetUpdatesChan(ucfg)
	
	if err != nil {
		// log.Panic(err)
		log.Println(err)
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
			msgCh <- fmt.Sprint(update.Message.Command(), ":", userID, ":")
			//reply = GalendarBot.ParseCommand(update.Message.Command(), fmt.Sprint(userID))
		} else {
			//Получем сообщение и парсим его
			Text = update.Message.Text
			msgCh <- fmt.Sprint(Text, ":", userID)
			//reply = GalendarBot.ParseText(Text, fmt.Sprint(userID))
		}
						
		log.Printf("[%s] %d %s", UserName, ChatID, Text)
		
		reply = <-msgCh
		// Созадаем сообщение
		msg := tgbotapi.NewMessage(ChatID, reply)
		// и отправляем его
		bot.Send(msg)
	}
}
