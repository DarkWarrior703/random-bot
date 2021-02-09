package youtube

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/DarkWarrior703/random-bot/utility"
	"github.com/bwmarrin/discordgo"
)

var running = false

// Query class
type Query struct {
	ID         string
	GuildID    string
	listofopts [][]string
}

var (
	listOfQ = []Query{}
)

var (
	queue       [][]string
	length      int = 0
	stopChannel chan bool
)

// QueryUserYtb asks user about what opts
func QueryUserYtb(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	list := strings.Split(m.Content, " ")
	if len(list) < 2 {
		return
	}
	query := ""
	if list[0] == "-play" || list[0] == "-p" {
		query = strings.Join(list[1:], "+")
	} else {
		return
	}
	if isNumeric(query) {
		return
	}
	url := "https://www.youtube.com/results?search_query=" + query + "&sp=EgIQAQ%253D%253D"
	body, err := utility.GetData(url)
	if err != nil {
		return
	}
	reg := regexp.MustCompile("(/watch\\?v=.{11})|(title\":{\"runs\":[[]{\"text\":\").+?\"")
	q := reg.FindAllString(string(body), -1)
	ques := "What video do you want?\n"
	listopt := [][]string{}
	size := 10
	if len(q) < 10 {
		size = len(q)
	}
	for i := 0; i < size; {
		title := ""
		if strings.Split(q[i], "\"")[0] == "title" {
			tmpReg := regexp.MustCompile("t\":.*")
			title = tmpReg.FindString(q[i])[3:]
		} else {
			break
		}
		ques += fmt.Sprintf("%d. ", i/2+1)
		ques += title
		ques += "\n"
		listopt = append(listopt, []string{title, q[i+1]})
		i += 2
	}
	s.ChannelMessageSend(m.ChannelID, ques)
	listOfQ = append(listOfQ, Query{m.Author.ID, m.GuildID, listopt})
}

// HandleChoice gets the choice of user
func HandleChoice(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	if len(list) < 2 {
		return
	}
	if (list[0] != "-play") && (list[0] != "-p") {
		return
	}
	choice, err := strconv.ParseFloat(strings.Join(list[1:], "+"), 64)
	if err != nil {
		return
	}
	fmt.Println("Hi")
	title := ""
	id := ""
	for i, q := range listOfQ {
		if q.ID == m.Author.ID && q.GuildID == m.GuildID {
			if int(choice) <= len(q.listofopts) {
				c := q.listofopts[int(choice)-1]
				title = c[0]
				id = c[1][9:]
				listOfQ = remove(listOfQ, i)
				break
			} else {
				continue
			}
		}
	}
	err = youtubeDownloader(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	filename := "D:\\tmp\\" + id + ".m4a"
	queue = append(queue, []string{title, filename})
	length++
	if !running {
		fmt.Println("WTF!")
		err = playSong(s, m)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func youtubeDownloader(id string) error {
	filename := "'D:/tmp/%(id)s.%(ext)s'"
	command := "youtube-dl.exe -f 'bestaudio[ext=m4a]' -o " + filename + " " + id
	cmd := exec.Command("powershell.exe", "-Command", command)
	err := cmd.Run()
	return err
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func playSong(s *discordgo.Session, m *discordgo.MessageCreate) error {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return err
	}
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		return err
	}
	ch := ""
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			ch = vs.ChannelID
		}
	}
	vc, err := s.ChannelVoiceJoin(g.ID, ch, false, true)
	if err != nil {
		return err
	}
	title := queue[0][0]
	filename := queue[0][1]
	queue = removeQuery(queue)
	stopChannel := make(chan bool)
	running = true
	s.ChannelMessageSend(m.ChannelID, "```\nNow playing "+title+"```\n")
	playAudioFile(vc, filename, stopChannel)
	vc.Speaking(false)
	running = false
	playNextSong(s, m, vc)
	return nil
}

func playNextSong(s *discordgo.Session, m *discordgo.MessageCreate, vc *discordgo.VoiceConnection) {
	fmt.Println(length)
	if length == 0 {
		vc.Disconnect()
		return
	}
	title := queue[0][0]
	filename := queue[0][1]
	queue = removeQuery(queue)
	stopChannel := make(chan bool)
	running = true
	s.ChannelMessageSend(m.ChannelID, "```\nNow playing "+title+"```\n")
	playAudioFile(vc, filename, stopChannel)
	running = false
	playNextSong(s, m, vc)
	return
}

func remove(s []Query, i int) []Query {
	if len(s) == 1 {
		return []Query{}
	} else if i == 0 {
		return s[1:]
	}
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removeQuery(s [][]string) [][]string {
	length--
	if len(s) > 1 {
		return s[1:]
	}
	return [][]string{}
}

// SkipHandler handles -skip command
func SkipHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	list := strings.Split(m.Content, " ")
	if len(list) > 2 {
		return
	}
	if list[0] == "-skip" {
		running = false
	}
}

// ClearQueueHandler handles -clearqueue
func ClearQueueHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	list := strings.Split(m.Content, " ")
	if len(list) > 2 {
		return
	}
	if list[0] == "-clearqueue" {
		queue = [][]string{}
		length = 0
	}
}
