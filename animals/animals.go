package animals

import (
	"encoding/json"
	"strings"

	"github.com/DarkWarrior703/random-bot/utility"
	"github.com/bwmarrin/discordgo"
)

// DogHandler handles -cat command
func DogHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	list := strings.Split(m.Content, " ")
	if list[0] != "-dog" {
		return
	}
	body, err := utility.GetData("https://random.dog/woof.json")
	if err != nil {
		return
	}
	type tmp struct {
		URL string `json:"url"`
	}
	tmpdata := &tmp{}
	err = json.Unmarshal(body, tmpdata)
	if err != nil {
		return
	}
	s.ChannelMessageSend(m.ChannelID, tmpdata.URL)
}

// CatHandler handles -cat command
func CatHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	list := strings.Split(m.Content, " ")
	if list[0] != "-cat" {
		return
	}
	body, err := utility.GetData("https://aws.random.cat/meow")
	if err != nil {
		return
	}
	type tmp struct {
		File string `json:"file"`
	}
	tmpdata := &tmp{}
	err = json.Unmarshal(body, tmpdata)
	if err != nil {
		return
	}
	s.ChannelMessageSend(m.ChannelID, tmpdata.File)
}

func FoxHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	list := strings.Split(m.Content, " ")
	if list[0] != "-fox" {
		return
	}
	body, err := utility.GetData("https://randomfox.ca/floof/")
	if err != nil {
		return
	}
	type tmp struct {
		Image string `json:"image"`
	}
	tmpdata := &tmp{}
	err = json.Unmarshal(body, tmpdata)
	if err != nil {
		return
	}
	s.ChannelMessageSend(m.ChannelID, tmpdata.Image)
}
