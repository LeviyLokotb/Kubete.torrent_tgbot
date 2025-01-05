package main

import (
	"kubete_torrentBot/botlogic"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var msgs []tgbotapi.MessageConfig

var mutex sync.Mutex

func main() {
	// запускаем таймер в отдельной горутине
	go Timer()

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
		// отсылаем сообщения по таймеру
		// блокируем чтобы msgs не изменился в процессе

		mutex.Lock()
		if len(msgs) > 0 {
			for n, msg := range msgs {
				bot.Send(msg)
				msgs = append(msgs[:n], msgs[n+1:]...)
			}
		}
		mutex.Unlock()

		// обрабатываем остальные сообщения
		if update.Message != nil {

			// вывод логов
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// основная часть
			var msg tgbotapi.MessageConfig

			chat_id := update.Message.Chat.ID

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
			//msg.ParseMode = "markdown"
			bot.Send(msg)
			//отправляем ответ

		}
	}
}

func Timer() {
	for {
		log.Println("◘ Timer run ◘")
		newtimer := time.NewTimer(10 * time.Second)

		<-newtimer.C
		log.Println("• Timer stop •")

		// выполняем действия по таймеру
		// проверка входа

		replies := botlogic.Entry()
		for id1, reply := range replies {
			msg := tgbotapi.NewMessage(id1, reply)
			msgs = append(msgs, msg)
		}

		//проверка уведомлений
		notifications := botlogic.Alert()
		for id2, alert := range notifications {
			msg := tgbotapi.NewMessage(id2, alert)
			msgs = append(msgs, msg)
		}
		log.Println("Timer messages: ", msgs)
	}
}
