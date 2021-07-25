package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
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

func sendMail(p *Person) {
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
		fmt.Printf("发送失败：\n")
		printPerson(p)
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")

}

func initServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	http.HandleFunc("/api/send_mail", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		query := url.Query()

		name := query["name"][0]
		phone := query["phone"][0]
		region := query["region"][0]
		address := query["address"][0]

		p := Person{name, phone, region, address}
		fmt.Printf("加入队列：\n")
		printPerson(&p)
		writeFile(&p)
		sendMail(&p)
	})
	fmt.Printf("localhost:9999 监听中。。。")
	http.ListenAndServe("localhost:9999", nil)

}

func printPerson(p *Person) {
	fmt.Printf(("姓名：%v\t"), p.name)
	fmt.Printf("电话：%v\t", p.phone)
	fmt.Printf("地区：%v\t", p.region)
	fmt.Printf("地址：%v\n", p.address)

}
func writeFile(p *Person) {

	data := time.Now().Format("01-02-2006 15:04:05 Mon") + "\t" + p.name + "\t" + p.phone + "\t" + p.region + "\t" + p.address

	file, err := os.OpenFile("backup.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	datawriter.WriteString(data + "\n")

	datawriter.Flush()
	file.Close()
}

func main() {
	initServer()
}
