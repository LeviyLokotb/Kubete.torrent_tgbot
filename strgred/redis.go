package strgred

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// структура для подключения к redis
type Config struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	User        string        `yaml:"user"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

// Функция подключения к redis
func NewClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	// проверка соединения (пингуем по пустому контексту)
	if err := db.Ping(ctx).Err(); err != nil {
		log.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	return db, nil
}

// подключение к редису
func Connection() *redis.Client {
	cfg := Config{
		Addr:        "localhost:6380",
		Password:    "ylp3QnB(VR0v>oL<Y3heVgsdE)+O+RZ",
		User:        "leosah",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}

	db, err := NewClient(context.Background(), cfg)
	if err != nil {
		log.Panic("db creating fail: ", err)
	}

	return db
}

func Redis_add(key, value any) {
	// подключаемся к redis
	cfg := Config{
		Addr:        "localhost:6380",
		Password:    "ylp3QnB(VR0v>oL<Y3heVgsdE)+O+RZ",
		User:        "leosah",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}

	db, err := NewClient(context.Background(), cfg)
	if err != nil {
		log.Panic("db creating fail: ", err)
	}
	// вносим значение
	err2 := db.Set(context.Background(), fmt.Sprint(key), value, 0).Err()
	if err2 != nil {
		log.Panic("ERROR in redis_add: ", err)
	}
}

func Redis_get(key any) string {
	// подключаемся к redis
	cfg := Config{
		Addr:        "localhost:6380",
		Password:    "ylp3QnB(VR0v>oL<Y3heVgsdE)+O+RZ",
		User:        "leosah",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}

	db, err := NewClient(context.Background(), cfg)
	if err != nil {
		log.Panic("db creating fail: ", err)
	}
	// получаем значение из бд
	value, err := db.Get(context.Background(), fmt.Sprint(key)).Result()
	if err == redis.Nil {
		return "nil"
	}
	return value
}

func Redis_delete(key any) bool {
	cfg := Config{
		Addr:        "localhost:6380",
		Password:    "ylp3QnB(VR0v>oL<Y3heVgsdE)+O+RZ",
		User:        "leosah",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}

	db, err := NewClient(context.Background(), cfg)
	if err != nil {
		log.Panic("db creating fail: ", err)
	}

	ok, err := db.Del(context.Background(), fmt.Sprint(key)).Result()
	if err != nil {
		log.Panic("ERROR in redis_delete: ", err)
	}

	return ok > 0
}

// поиск ключей по значению find
func GetSomeIDs(find string) []int64 {
	cfg := Config{
		Addr:        "localhost:6380",
		Password:    "ylp3QnB(VR0v>oL<Y3heVgsdE)+O+RZ",
		User:        "leosah",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}

	db, err := NewClient(context.Background(), cfg)
	if err != nil {
		log.Panic("db creating fail: ", err)
	}

	var results []int64

	keys, err2 := db.Keys(context.Background(), "*").Result()
	if err2 != nil {
		log.Panic("ERROR in GetSomeIDs: ", err2)
	}
	//log.Println(keys)

	for _, key := range keys {
		status := Redis_get(key)
		intkey, errstr := strconv.Atoi(key)
		if errstr != nil {
			log.Panic("ERROR in GetSomeIDs: strconv error: ", errstr)
		}
		if status == find {
			results = append(results, int64(intkey))
		}
	}
	return results
}
