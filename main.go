package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
)

type Person struct {
	name    string
	phone   string
	region  string
	address string
}

type Mail struct {
	From     string   `json:"from"`
	To       []string `json:"to"`
	Password string   `json:"password"`
}

func send_mail(p *Person) {

	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	mail := Mail{}
	errDecode := decoder.Decode(&mail)
	if errDecode != nil {
		fmt.Println("error:", errDecode)
	}

	// Sender data.
	from := mail.From
	password := mail.Password

	// Receiver email address.
	to := mail.To

	// smtp server configuration.
	smtpHost := "smtp.qq.com"
	smtpPort := "587"

	// Message.

	message := []byte("To: " + to[0] + "\r\n" +
		"Subject: 来自" + p.name + "的新订单！\r\n" +
		"\r\n" +
		"姓名：" + p.name + "\r\n" +
		"电话：" + p.phone + "\r\n" +
		"地区：" + p.region + "\r\n" +
		"地址：" + p.address + "\r\n")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")

}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	http.HandleFunc("/send_mail", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		query := url.Query()

		name := query["name"][0]
		phone := query["phone"][0]
		region := query["region"][0]
		address := query["address"][0]

		fmt.Printf(("姓名：%v\t"), name)
		fmt.Printf("电话：%v\t", phone)
		fmt.Printf("地区：%v\t", region)
		fmt.Printf("地址：%v\n", address)
		send_mail(&Person{name, phone, region, address})
		w.Write([]byte("ok"))
	})

	http.ListenAndServe("localhost:9999", nil)
}
