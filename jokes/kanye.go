package jokes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// KanyeHandler handler -kanye command
func KanyeHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	if m.Content[0] != '-' {
		return
	}
	list := strings.Split(m.Content, " ")
	if list[0] != "-kanye" {
		return
	}
	type tmp struct {
		Quote string `json:"quote"`
	}
	url := "https://api.kanye.rest"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	tmpdata := &tmp{}
	err = json.Unmarshal(body, tmpdata)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Kanye said: "+tmpdata.Quote)
}
