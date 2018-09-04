package main

import (
	"fmt"
	nats "github.com/nats-io/go-nats"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"gopkg.in/mailgun/mailgun-go.v1"
	"log"
	"strconv"
	"strings"
	"time"
)

type config struct {
	NatUser       string `env:"NAT_USER"`
	NatPass       string `env:"NAT_PASS"`
	AppName       string `env:"APP_NAME" envDefault:"Notification monitor - Golang"`
	NatHost       string `env:"NAT_HOST"`
	MailgunDomain string `env:"MAILGUN_DOMAIN"`
	MailgunKey    string `env:"MAILGUN_PRIVATE_KEY"`
	MailgunPublic string `env:"MAILGUN_PUBLIC_KEY"`
	EmailFrom     string `env:"EMAIL_FROM" envDefault:"postmaster@ci.connectapp.biz"`
	EmailTo       string `env:"EMAIL_TO" envDefault:"tech@connectapp.biz"`
}

//Email send , used only when load is high.
func sendMessage(mg mailgun.Mailgun, sender, subject, body, recipient string) {
	message := mg.NewMessage(sender, subject, body, recipient)
	resp, id, err := mg.Send(message)

	if err != nil {
		log.Fatal(err)
	}

	//Logs for email queue
	log.Printf("ID: %s Resp: %s\n", id, resp)

}
func usage() {
	log.Fatalf("Environment varibles for NAT, MAILGUN and EMAIL not set \n NAT_USER\n NAT_PASS\n APP_NAME\n NAT_HOST\n MAINGUN_DOMANI\n MAILGUN_PRIVATE_KEY\n MAILGUN_PUBLIC_KEY\n EMAIL_FROM\n EMAIL_TO\n")
}

func main() {

	//Load env variables
	err := godotenv.Load()

	//deault config
	cfg := config{}
	err = env.Parse(&cfg)
	if err != nil || len(cfg.NatUser) == 0 || len(cfg.NatPass) == 0 || len(cfg.NatHost) == 0 || len(cfg.MailgunDomain) == 0 || len(cfg.MailgunKey) == 0 || len(cfg.MailgunDomain) == 0 {
		usage()
		panic("Cannot pass environment variables")
	}
	// Connect to a server
	options := nats.Options{
		User:           cfg.NatUser,
		Password:       cfg.NatPass,
		Name:           cfg.AppName,
		Url:            cfg.NatHost,
		AllowReconnect: true,
		MaxReconnect:   10,
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
	}

	mg := mailgun.NewMailgun(cfg.MailgunDomain, cfg.MailgunKey, cfg.MailgunPublic)
	conn, err := options.Connect()
	if err != nil {
		panic(err)
	} else {
		conn.Subscribe("stats.loadaverage", func(m *nats.Msg) {
			//fmt.Printf("%v+\n", time.Now())
			msg := strings.Split(string(m.Data), " ")
			host := msg[0]
			p, _ := strconv.ParseFloat(msg[1], 64)
			cpu, _ := strconv.ParseFloat(msg[2], 64)
			load := p / cpu //strconv.ParseFloat(msg[1], 64) / strconv.ParseFloat(msg[2], 64)

			if load > 0.60 {

				subject := fmt.Sprintf(" Warning! Server is Overloaded: %s ", host)
				body := fmt.Sprintf("Load average : %f", load)

				sendMessage(mg, "postmaster@ci.connectapp.biz", subject, body, "hitesh.joshi@ziploan.in")

				//fmt.Printf("%v Received a message :%s %f\n ", time.Now(), host, load)
			}
		})
	}

	defer conn.Close()

	go forever()
	select {} // block forever
}

func forever() {
	for {
		time.Sleep(time.Second)
	}
}
