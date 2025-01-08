package botlogic

import (
	"kubete_torrentBot/strgred"
	"log"
)

func Get_status(chat_id int64) string {

	// получаем статус из бд
	status := strgred.Redis_get(chat_id)

	if status == "nil" {
		status = "Неизвестный"
	}

	return status
}

func Login(chat_id int64) string {
	//! достаём токен входа из redis
	//! запрос модулю авторизации - проверка токена входа
	code := 1 // возвращаемый код (временно)
	switch code {
	case 1: // не опознанный/истёкший токен:
		strgred.Redis_delete(chat_id)
		return "Вы не вошли либо время входа истекло."
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
	return ""
}

func Login_type(chat_id int64) string {
	status := Get_status(chat_id)
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
	default:
		log.Panic("Error: Unknown status")
	}
	return "0"
}

func Logout(chat_id int64) string {
	//! удаляем chat_id из redis

	if !strgred.Redis_delete(chat_id) {
		return "Сеанс уже завершён ранее"
	}
	return "Сеанс завершён."
}

func Logout_all(chat_id int64) string {
	Logout(chat_id)
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
