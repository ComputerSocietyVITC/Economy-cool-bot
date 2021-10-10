package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	err   error
)

func init() {
	Token = os.Getenv("BOT_TOKEN")
}

func main() {
	hello()
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error in creating session : ", err)
		return
	}

	session.AddHandler(messageCreate)
	session.AddHandler(messageReg)

	session.Identify.Intents = discordgo.IntentsGuildMessages

	err = session.Open()
	if err != nil {
		fmt.Println("Error in opening connection : ", err)
		return
	}

	fmt.Println("Bot in now running. CTRL-C to exit.")
	sc := make(chan os.Signal, 100)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	session.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	// author := strings.FieldsFunc(m.Author.String(), func(r rune) bool {
	// 	if r == '#' {
	// 		return true
	// 	}
	// 	return false
	// })

	if m.Content == "hi!" {

		_, err = s.ChannelMessageSend(m.ChannelID, "hello boyo")
		if err != nil {
			fmt.Println(err)
		}
	}

	if m.Content[0] == '!' {
		if m.Content == "!diss me" {
			message := diss("<@!" + m.Author.ID + ">")
			_, err = s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}

		} else if strings.Contains(m.Content[:5], "!diss") {
			author := strings.Split(m.Content, " ")[1]
			fmt.Println(author)
			message := diss(author)
			_, err = s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func diss(author string) string {
	message := author + " is a piece of sh*t"
	return message
}
