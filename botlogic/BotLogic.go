package botlogic

import (
	"context"
	"kubete_torrentBot/strgred"
	"log"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// массив сообщений, на которые мы можем дать ответ
var OKprocessing []string = []string{"start", "status", "help", "login", "login_type", "add_me", "logout", "logot_alltrue"}

func Processing(request, chat_id string) string {

	status := get_status(chat_id)

	var msg string

	switch request {
	case "start":
		msg = start()
	case "help":
		msg = help()
	case "status":
		msg = "Сейчас вы " + status + " пользователь"
	case "add_me": // временно
		add_to_redis(chat_id, status)
		msg = "Пользователь внесён в базу. Текущий статус: " + get_status(chat_id)
	case "login":
		msg = login(status)
	case "login_type":
		msg = login_type(status)
	case "logout":
		msg = logout()
	case "logout_alltrue"
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
	commands_text := ""
	for _, command := range OKprocessing {
		commands_text += "/" + command + "\n\n"
	}

	commands_text += "По всем вопросам обращаться к @voraxas"

	return "Команды бота:\n" + commands_text
}

func get_status(chat_id string) string {
	// подключаемся к redis
	cfg := strgred.Config{
		Addr:        "localhost:6380",
		Password:    "ylp3QnB(VR0v>oL<Y3heVgsdE)+O+RZ",
		User:        "leosah",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}

	db, err := strgred.NewClient(context.Background(), cfg)
	if err != nil {
		log.Panic("db creating fail: ", err)
	}

	// получаем статус из бд
	status, err := db.Get(context.Background(), chat_id).Result()

	if err == redis.Nil {
		// пользователя нет в базе
		status = "Неизвестный"
	}

	return status

}

func login(status string) string {
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
			//! удаляем из redis chat_id
			return login("Неизвестный")
		case 2: // в доступе отказано:
			//! удаляем из redis chat_id
			return "Неудачная авторизация."
		case 3: // доступ предоставлен
			// получаем jwt-токен доступа и токен обновления
			// сохраняем оба токена и статус Авторизованный в базу
			// обрабатываем запрос по статусу "Авторизованный"
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
		// генерируем токенн входа
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

func logout() string {
	//! удаляем chat_id из redis
	return "Сеанс завершён."
}

func logout_all() string{
	logout()
	//! запрос модулю авторизации /logout, отправляем токен обновления
	return "Сеанс завершён на всех устройствах."
}

func add_to_redis(key, value string) {
	// подключаемся к redis
	cfg := strgred.Config{
		Addr:        "localhost:6380",
		Password:    "ylp3QnB(VR0v>oL<Y3heVgsdE)+O+RZ",
		User:        "leosah",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}

	db, err := strgred.NewClient(context.Background(), cfg)
	if err != nil {
		log.Println("db creating fail: ", err)
	}

	db.Set(context.Background(), key, value, 0)
}
