package main

import (
	"github.com/margleb/booking/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"time"
)

// listenForMail - функция, ожидающая получения письма из канала
func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {

	// настройки сервера
	server := mail.NewSMTPClient()
	server.Host = "localHost"
	server.Port = 1025
	server.KeepAlive = false                 // не хотим чтобы соединение было постоянным
	server.ConnectTimeout = 10 * time.Second // максимальное время соединения
	server.SendTimeout = 10 * time.Second    // максимальное время отправления

	// создаем клиента
	client, err := server.Connect()
	if err != nil {
		ErrorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, "Hello <strong>world</strong>!")

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email send")
	}

}
