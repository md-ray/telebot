package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/magiconair/properties"
)

func main() {
	// Check command-line param first
	if len(os.Args) < 2 {
		fmt.Println("command-line argument for properties file not provided")
		os.Exit(1)
	}

	// Check properties file first
	cfgFile := os.Args[1]
	// init from a file
	p := properties.MustLoadFile(cfgFile, properties.UTF8)
	bottoken := p.MustGet("bot.token")
	bot, err := tgbotapi.NewBotAPI(bottoken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Transmission param init
	tranmissionPair := p.MustGet("transmission.user") + ":" + p.MustGet("transmission.pass")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		words := strings.Fields(update.Message.Text)
		keyword := words[0]

		if strings.EqualFold(keyword, "torrent-add") {
			if len(words) < 2 {
				keyword = "Command not complete. Missing parameter"
			} else {

				param := "--auth " + tranmissionPair + " -a \"" + words[1] + "\""
				log.Printf("param : " + param)
				cmd := exec.Command("/usr/bin/transmission-remote", "--auth", tranmissionPair, "-a", words[1])
				//cmd := exec.Command("dir")
				out, err := cmd.CombinedOutput()
				if err != nil {
					log.Printf("Command finished without error")
					keyword = "Successfully executed. " + string(out)
				} else {
					log.Printf("Command finished with error: %v", err)
					keyword = "Error executed. " + string(out)
				}
			}
		} else if strings.EqualFold(keyword, "torrent-list") {
			param := "--auth transmission:transmission -a \"" + words[1] + "\""
			log.Printf("param : " + param)
			cmd := exec.Command("/usr/bin/transmission-remote", "--auth", tranmissionPair, "-l")
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("Command finished without error")
				keyword = "Successfully executed. " + string(out)
			} else {
				log.Printf("Command finished with error: %v", err)
				keyword = "Error executed. " + string(out)
			}
		} else {
			keyword = "Command Unknown"
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyword)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

	}
}
