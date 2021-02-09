package misc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//ShortenUrlsHandler handles -shorten and returns a shortened URL
func ShortenUrlsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	if len(m.Content) == 0 {
		return
	}
	list := strings.Split(m.Content, " ")
	if list[0] != "-shorten" {
		return
	}
	if len(list) > 2 {
		s.ChannelMessageSend(m.ChannelID, "-shorten must have two arguments")
		return
	}
	query := list[1]
	temp := strings.NewReader("url=" + query)
	req, err := http.NewRequest("POST", "https://cleanuri.com/api/v1/shorten", temp)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	type tmp struct {
		ResultURL string `json:"result_url"`
	}
	tmpdata := &tmp{}
	err = json.Unmarshal(body, tmpdata)
	s.ChannelMessageSend(m.ChannelID, "```\nYour shortened URL is "+tmpdata.ResultURL+".\n```")
}
