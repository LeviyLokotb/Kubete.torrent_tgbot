package main

import (
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

	updates := bot.GetUpdatesChan(u)

	// обработка сообщений
	for update := range updates {
		if update.Message != nil {

			// вывод логов
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// основная часть
			var msg /*, chat_id*/ tgbotapi.MessageConfig

			//chat_id = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			// комманды

				// текст команды /help
				commands_text := "/start\n- начало работы\n" + "/help\n- информация\n" + "/status\n- статус пользователя"

			switch update.Message.Command() {
			case "start":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, telega user! I useless now...")
			case "help":
				help_text := "Вот что я могу:\n" + commands_text + "\nэто не так уж и плохо, как по мне"
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, help_text)
			default:
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Нет такой команды")
			}

			//отправляем ответ
			bot.Send(msg)
		}
	}
}
