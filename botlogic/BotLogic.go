package botlogic

import (
	"kubete_torrentBot/strgred"
	"log"
)

// массив сообщений, на которые мы можем дать ответ
var OKprocessing []string = []string{"start", "status", "help", "login", "login_type", "add_me", "logout", "logot_alltrue"}

func Processing(request string, chat_id int64) string {

	status := get_status(chat_id)

	var msg string

	switch request {
	case "start":
		msg = start()
	case "help":
		msg = help()
	case "status":
		msg = "Сейчас вы " + get_status(chat_id) + " пользователь"
	case "add_me": // временно
		strgred.Redis_add(chat_id /*status*/, "Авторизованный")
		msg = "Пользователь внесён в базу. Текущий статус: " + get_status(chat_id)
	case "login":
		msg = login(status, chat_id)
	case "login_type":
		msg = login_type(status)
	case "logout":
		msg = logout(chat_id)
	case "logout_alltrue":
		msg = logout_all(chat_id)
	default:
		// мы уже отсеяли иные команды, но на всякий случай
		msg = "Нет такой команды"
	}

	return msg
}

func start() string {
	return "Hello, telega user!\n\nЕсли не знаешь с чего начать, введи /help"
}

func help() string {
	var commands_text string
	for _, command := range OKprocessing {
		commands_text += "/" + command + "\n"
	}

	commands_text = "Команды бота:\n" + commands_text + "По всем вопросам обращаться к @voraxas"

	return commands_text
}

func get_status(chat_id int64) string {
	// получаем статус из бд
	status := strgred.Redis_get(chat_id)

	if status == "nil" {
		status = "Неизвестный"
	}

	return status
}

func login(status string, chat_id int64) string {
	switch status {
	case "Неизвестный":
		// предлагаем авторизацию
		return "Вы не залогинены. \nАвторизуйтесь через Github, Яндекс ID либо через код"
	case "Анонимный":
		//! достаём токен входа из redis
		//! запрос модулю авторизации - проверка токена входа
		code := 1 // возвращаемый код (временно)
		switch code {
		case 1: // не опознанный/истёкший токен:
			strgred.Redis_delete(chat_id)
			return login("Неизвестный", chat_id)
		case 2: // в доступе отказано:
			strgred.Redis_delete(chat_id)
			return "Неудачная авторизация."
		case 3: // доступ предоставлен
			//! получаем jwt-токен доступа и токен обновления
			//! сохраняем оба токена и статус Авторизованный в базу
			strgred.Redis_delete(chat_id)
			strgred.Redis_add(chat_id, "Авторизованный")
			//! обрабатываем запрос по статусу "Авторизованный"
		default:
			log.Panic("Error: Unknown autorization code")
			return ""
		}

	case "Авторизованный":
		return "Вы уже авторизованы!"
	}
	return ""
}

func login_type(status string) string {
	switch status {
	case "Неизвестный":
		// генерируем токен входа
		// заносим в редис chat_id и токен входа со статусом
		// отправляем токен входа модулю авторизации
		// получаем ответ и отправляем
		return "Модуль авторизации не подключен.\nComing soon..."
	case "Анонимный":
		// генерируем токен входа
		// заменяем токен входа в redis
		// отправляем токен входа модулю авторизации
		// получаем ответ и отправляем
		return "Модуль авторизации не подключен.\nComing soon..."
	case "Авторизованный":
		return "Вы уже авторизованы!"
	default:
		log.Panic("Error: Unknown status")
	}
	return ""
}

func logout(chat_id int64) string {
	//! удаляем chat_id из redis
	ok := strgred.Redis_delete(chat_id)
	if !ok {
		return "Сеанс уже завершён ранее"
	}
	return "Сеанс завершён."
}

func logout_all(chat_id int64) string {
	logout(chat_id)
	//! запрос модулю авторизации /logout, отправляем токен обновления
	return "Сеанс завершён на всех устройствах."
}

func Entry() map[int64]string {
	replies := make(map[int64]string)
	chat_ids := strgred.GetSomeIDs("Анонимный")

	for _, user_id0 := range chat_ids {
		user_id := int64(user_id0)
		//! запрос модулю авторизации - проверка токена входа
		code := 1 // возвращаемый код (временно)
		switch code {
		case 1: // неопознаный токен или время действия закончилось
			strgred.Redis_delete(user_id)
		case 2: // в доступе отказано
			strgred.Redis_delete(user_id)
			replies[user_id] = "Статус входа: неудачная авторизация"
		case 3: // доступ предоставлен
			//! получаем jwt-токен доступа и токен обновления
			//! сохраняем оба токена и статус Авторизованный в базу
			strgred.Redis_delete(user_id)
			strgred.Redis_add(user_id, "Авторизованный")
			replies[user_id] = "Статус входа: успешная авторизация"
		default:
			log.Panic("Error: Unknown autorization code")
		}
	}
	return replies
}

func Alert() map[int64]string {
	notifications := make(map[int64]string)

	chat_ids := strgred.GetSomeIDs("Авторизованный")

	for _, user_id := range chat_ids {
		//! запрос главному модулю на URL /notification по токену доступа
		notic := "тестовое уведомление" // уведомление (временно)
		notifications[int64(user_id)] = notic
		//! запрос главному модулю на удаление уведомлений по JWT токену доступа
	}
	return notifications
}
