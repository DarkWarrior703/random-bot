package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DarkWarrior703/anime-bot/anime"
	"github.com/DarkWarrior703/anime-bot/jokes"
	"github.com/DarkWarrior703/anime-bot/youtube"
	"github.com/bwmarrin/discordgo"
)

const token = ""

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session", err)
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("Error getting sound", err)
		os.Exit(1)
	}
	dg.AddHandler(youtube.QueryUserYtb)
	dg.AddHandler(youtube.HandleChoice)
	dg.AddHandler(youtube.SkipHandler)
	dg.AddHandler(anime.AnimeHandler)
	dg.AddHandler(anime.MangaHandler)
	dg.AddHandler(anime.ImageHandler)
	dg.AddHandler(jokes.KanyeHandler)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
