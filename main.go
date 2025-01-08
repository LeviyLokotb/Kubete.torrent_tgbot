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
	for {
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
				status := botlogic.Get_status(chat_id)
				command := update.Message.Command()
				args := update.Message.CommandArguments()

				// не даём выполнять команды неищвестным и анонимным
				if status != "Авторизованный" {
					if command != "start" && command != "help" && command != "status" && command != "login" {
						command = "login"
						msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, "Сначала пройдите авторизацию")
						bot.Send(msg1)
					}
				}

				// коммады
				switch command {
				case "start":
					start_text :=
						`Это телеграм-бот для проведения массовых опросов и тестирований.
					Узнайте о коммандах введя /help.
					Не забудьте пройти авторизацию: /login`
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, start_text)
				case "help":
					head :=
						`Вот что умеет этот бот:
					переключайтесь между страницами (1-7):
					/help <номер страницы>` + "\n"

					help1 :=
						`1. Команды бота и авторизация
					/start - начало работы
					/help - возможности бота
					/status - статус пользователя (Неизвестный/Анонимный/Авторизованный)
					/login - вход
					/login type - регистрация
					/logout - завершение сеанса
					/logout all=true - на всех устройствах`

					help2 :=
						`2. Команды, связанные с пользователями 
					(<id> - id пользователя. Не вводите если операция направлена на вас)
					/userList - список пользователей
					/name <id> - посмотреть ФИО пользователя
					/nameChange <id> - изменить ФИО пользователя
					/userData <id> - информация о пользователе (курсы, оценки, тесты)
					/role <id> - посмотреть роль пользователя (студент/преподаватель/админ)
					/roleChange <id> <role> - изменить роль пользователя
					/blockInfo <id> - заблокирован ли пользователь
					/block <id> - заблокировать пользователя
					/unblock <id> - разблокировать пользователя`

					help3 :=
						`3. Комманды, связанные с дисциплинами
					(<id> - id дисциплины)
					/courseList - список дисциплин
					/course <id> - информация о дисциплине (Название, Описание, ID преподавателя)
					/courseChangeName <id> - изменить название дисциплины
					/courseChangeInfo <id> - изменить описание дисциплины
					/courseTestList <id> - список тестов дисциплины
					/testActive <id> - активен ли тест
					/testActivate <id> - активировать тест
					/testDeactivate <id> - деактивировать тест
					/testAdd <id> - добавить пустой тест
					/testDelete <id> - удалить тест
					/courseStudentList <id> - список студентов дисциплины
					/courseStudentAdd <id дисциплины> <id пользователя> - записать студента на дисциплину
					/courseStudentDelete <id дисциплины> <id пользователя> - отчислить студента с дисциплины
					/courseAdd - создать дисциплину
					/courseDelete - удалить дисциплину`

					help4 :=
						`4. Команды, связанные с вопросами
					(<id> - id вопроса)
					/questList - список всех ваших вопросов
					/questInfo <id> - информация о вопросе (Название, Текст, id, Ответ)
					/questUpdate <id> - изменить вопрос
					/questCreate - создать вопрос
					/questDelete <id> - удалить тест`

					help_text := head
					switch args {
					case "2":
						help_text += help2
					case "3":
						help_text += help3
					case "4":
						help_text += help4
					default:
						help_text += help1
					}

					msg = tgbotapi.NewMessage(update.Message.Chat.ID, help_text)
				case "status":
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Сейчас вы "+status+" пользователь")
				case "login":
					switch args {
					case "type":
						switch status {
						case "Авторизованный":
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Вы уже авторизованы!")
						default:
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, botlogic.Login_type(chat_id))
						}
					default:
						switch status {
						case "Неизвестный":
							keyboard := tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonURL("🐱Github", "https://github.com"),
								),
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonURL("📕Яндекс ID", "https://yandex.ru"),
								),
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData("🧾Код", "button_code"),
								),
							)
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Авторизуйтесь через:")
							msg.ReplyMarkup = keyboard
						case "Авторизованный":
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Вы уже авторизованы!")
						case "Анонимный":
							msg = tgbotapi.NewMessage(update.Message.Chat.ID, botlogic.Login(chat_id))
						}
					}
				case "logout":
					if args == "all=true" {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, botlogic.Logout_all(chat_id))
					} else {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, botlogic.Logout(chat_id))
					}
				default:
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Нет такой команды")
				}

				// режим markdown для форматирования текста
				//msg.ParseMode = "markdown"
				bot.Send(msg)
				//отправляем ответ

			} else if update.CallbackQuery != nil {
				// обработка нажатия на кнопку
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
				if _, err := bot.Request(callback); err != nil {
					log.Panic("Ошибка при обработке callback:", err)
				}

				// Определяем, какая кнопка была нажата
				switch update.CallbackQuery.Data {
				case "button_code":
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
						botlogic.Login_type(update.CallbackQuery.Message.Chat.ID))
					bot.Send(msg)
				}
			}
		}
	}
}

func Timer() {
	for {
		log.Println("◘ Timer run ◘")
		newtimer := time.NewTimer(60 * time.Second)

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
