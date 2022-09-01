package main

import (
	"fmt"
	"github.com/margleb/booking/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"io/ioutil"
	"log"
	"strings"
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
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email send")
	}

}
