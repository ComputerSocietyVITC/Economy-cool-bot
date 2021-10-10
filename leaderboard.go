package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/sirupsen/logrus"
)

type users struct {
	Users []user `json:"users"`
}

type user struct {
	Username    string `json:"username"`
	Points      int    `json:"points"`
	Level       int    `json:"level"`
	ToNextLevel int    `json:"tonext"`
}

func messageReg(s *discordgo.Session, m *discordgo.MessageCreate) {
	var u users
	users := readJSON()
	con := true

	for _, i := range m.Member.Roles {
		if i == "896603878227841025" {
			con = false
			break
		}
	}
	if con {
		for i := 0; i < len(users.Users); i++ {
			if users.Users[i].Username == m.Author.String() {
				users.Users[i].Points += 10
				if users.Users[i].Points > users.Users[i].ToNextLevel {
					users.Users[i].Points = users.Users[i].Points % users.Users[i].ToNextLevel
					users.Users[i].Level++
					users.Users[i].ToNextLevel = users.Users[i].ToNextLevel * abs(users.Users[i].Level-10)
				}
				u = writeJSON(users, user{})

				goto LEAD
			}
		}
		u = writeJSON(users, user{
			Username:    m.Author.String(),
			Points:      10,
			Level:       0,
			ToNextLevel: 100,
		})

	LEAD:

		if strings.Contains(m.Content, "!leaderboard") {
			user_list := leaderBoard(u.Users)
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("LeaderBoard", strings.Join(user_list[:], "\n")))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func leaderBoard(user []user) []string {
	sort.Slice(user, func(i, j int) bool {
		return user[i].Points*(user[i].Level+1) > user[j].Points*(user[j].Level+1)
	})
	var user_list []string
	for i := 0; i < len(user); i++ {
		user_list = append(user_list, user[i].Username)
	}
	return user_list
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func readJSON() users {
	jsonFile, err := os.Open("db/economy.json")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened economy.json")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var u users

	json.Unmarshal(byteValue, &u)
	return u
}

func writeJSON(data users, newUser user) users {

	// Preparing the data to be marshalled and written.

	var u users
	t := user{}
	if newUser != t {
		u = users{Users: append(data.Users, newUser)}
	} else {
		u = data
	}
	dataBytes, err := json.Marshal(u)
	if err != nil {
		logrus.Error(err)
	}

	err = ioutil.WriteFile("db/economy.json", dataBytes, 0644)
	if err != nil {
		logrus.Error(err)
	}
	return u
}
