package main

import (
	"kubete_torrentBot/botlogic"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// создаём бота по зарегестрированному API
	bot, err := tgbotapi.NewBotAPI("7798824633:AAFd4IkF6Rfvs-Fh0h2rSytWCdKLSQJTwaM")
	if err != nil {
		log.Panic(err)
	}

	// режим отладки для логов
	bot.Debug = true

	// канал обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// список новых сообщений
	updates := bot.GetUpdatesChan(u)

	// обработка сообщений
	for update := range updates {
		if update.Message != nil {

			// вывод логов
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// основная часть
			var msg tgbotapi.MessageConfig

			chat_id := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text).Text

			// коммады
			notok := true
			for _, command := range botlogic.OKprocessing {
				if command == update.Message.Command() {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, botlogic.Processing(command, chat_id))
					notok = false
				}
			}
			if notok {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Нет такой комманды")
			}
			// режим markdown для форматирования текста
			msg.ParseMode = "markdown"
			bot.Send(msg)
			//отправляем ответ

		}
	}
}
