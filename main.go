package main

import (
	"context"
	"kubete_torrentBot/strgred"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	redis "github.com/redis/go-redis"
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

	// подключение к redis (пока не готово)
	cfg := strgred.Config{
		Addr:        "localhost:6379",
		Password:    "test1234",
		User:        "testuser",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}

	db, err := strgred.NewClient(context.Background(), cfg)

	// обработка сообщений
	for update := range updates {
		if update.Message != nil {

			// вывод логов
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// основная часть
			var msg, chat_id tgbotapi.MessageConfig

			chat_id = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			var status string

			// запрос к redis
			///////////////////
			_, err := db.Get(context.Background(), chat_id).Result()
			if err == redis.Nil {
				// сценарий неизвестного пользователя
				status = "Неизвестный"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваш статус: "+status+".\nЗарегистрируйтесь:")
				bot.Send(msg)
				// запрос на авторизацию
				/*...*/
			} else if err != nil {
				// ошибки
				log.Printf("failed to get value, error: %v\n", err)
			} else {
				// сценарий анонимного пользователя
				status = "Анонимный"
				// требуем вход
				/*...*/
			}
			///////////////////

			// комманды

			commands_text := "/start\n- начало работы\n" + "/help\n- информация\n" + "/status\n- статус пользователя"

			switch update.Message.Command() {
			case "start":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, telega user! I useless now...")
			case "help":
				help_text := "Вот что я могу:\n" + commands_text + "\nэто не так уж и плохо, как по мне"
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, help_text)
			case "status":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ваш статус: "+status+".")
			default:
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Нет такой команды")
			}

			//отправляем ответ
			bot.Send(msg)
		}
	}
}
