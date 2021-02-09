package anime

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/DarkWarrior703/anime-cli/anime"
	"github.com/DarkWarrior703/anime-cli/manga"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

// AnimeHandler defines -anime command
func AnimeHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	list := strings.Split(m.Content, " ")
	if list[0] != "-anime" {
		return
	}
	if len(list) <= 1 {
		s.ChannelMessageSend(m.ChannelID, "```\n-anime takes an argument.\n```")
		return
	}
	query := strings.Join(list[1:], "+")
	anime.SetLimit(5)
	animeList, err := anime.RetrieveAnimeData(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Title", "Episodes", "Type", "Status", "Rating"})
	for _, anime := range animeList {
		table.Append([]string{anime.Title, fmt.Sprint(anime.Episodes), anime.Type, anime.Status, anime.Rated})
	}
	table.Render()
	s.ChannelMessageSend(m.ChannelID, "```\n"+tableString.String()+"\n```")
}

// MangaHandler defines -manga command
func MangaHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	list := strings.Split(m.Content, " ")
	if list[0] != "-manga" {
		return
	}
	if len(list) <= 1 {
		s.ChannelMessageSend(m.ChannelID, "```\n-manga takes an argument.\n```")
	}
	query := strings.Join(list[1:], "+")
	manga.SetLimit(5)
	mangaList, err := manga.RetrieveMangaData(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Title", "Synopsis", "Type", "Chapters", "Volumes"})
	for _, manga := range mangaList {
		if len(manga.Synopsis) > 60 {
			table.Append([]string{manga.Title, manga.Synopsis[:59] + "...", manga.Type, fmt.Sprint(manga.Chapters), fmt.Sprint(manga.Volumes)})
		} else {
			table.Append([]string{manga.Title, manga.Synopsis, manga.Type, fmt.Sprint(manga.Chapters), fmt.Sprint(manga.Volumes)})
		}
	}
	table.Render()
	s.ChannelMessageSend(m.ChannelID, "```\n"+tableString.String()+"\n```")
}

// ImageHandler defines -image command
func ImageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	list := strings.Split(m.Content, " ")
	if list[0] != "-image" {
		return
	}
	if len(list) <= 1 {
		s.ChannelMessageSend(m.ChannelID, "```\n-image takes an argument.\n```")
	}
	type tmp struct {
		Results []struct {
			URL      string `json:"url"`
			ImageURL string `json:"image_url"`
		} `json:"results"`
	}
	query := strings.Join(list[1:], "+")
	url := "https://api.jikan.moe/v3/search/anime?q=" + query
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	tmpdata := &tmp{}
	err = json.Unmarshal(body, tmpdata)
	if err != nil {
		return
	}
	s.ChannelMessageSend(m.ChannelID, tmpdata.Results[0].ImageURL)
}
