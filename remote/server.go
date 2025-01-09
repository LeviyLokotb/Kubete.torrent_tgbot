package remote

import (
	"bufio"
	"log"
	"net"
)

func SendMain(request string) string {
	// открываем соединение
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Panic("Error in remote: ", err)
	}
	defer conn.Close()
	// отправляем данные
	if _, err = conn.Write([]byte(request)); err != nil {
		log.Panic(err)
	}
	connReader := bufio.NewReader(conn)
	message, _ := connReader.ReadString('\n')
	return message
}

func SendAu(request string) string {
	// открываем соединение
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Panic("Error in remote: ", err)
	}
	defer conn.Close()
	// отправляем данные
	if _, err = conn.Write([]byte(request)); err != nil {
		log.Panic(err)
	}
	connReader := bufio.NewReader(conn)
	message, _ := connReader.ReadString('\n')
	return message
}
