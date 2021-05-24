package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

// Variables used for command line parameters
var (
	Token     string
	ChannelID string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token | DISCORD_TOKEN")
	flag.StringVar(&ChannelID, "c", "", "Channel ID | DISCORD_CHANNEL_ID")
	flag.Parse()

	if Token == "" {
		Token = os.Getenv("DISCORD_TOKEN")
	}
	if ChannelID == "" {
		ChannelID = os.Getenv("DISCORD_CHANNEL_ID")
	}

	if Token == "" || ChannelID == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {

	// Create a new Discord session using the provided login information.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	go schedule(dg)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func randomSentence(category string) string {
	red := []string{
		"Coucou c'est l'equipe ROUGE, est ce qu'on peut me sortir ?",
	}
	green := []string{
		"Tout habillé de VERT, avec les VERRES, je fais mes prieres",
	}
	out := []string{
		"TocToc on peut rentrer ?",
	}
	rand.Seed(time.Now().Unix())
	switch category {
	case "red":
		fmt.Println("RED OUT")
		return red[rand.Intn(len(red))]
	case "green":
		fmt.Println("GREEN OUT")
		return green[rand.Intn(len(green))]
	default:
		fmt.Println("IN")
		return out[rand.Intn(len(out))]
	}
}

func schedule(s *discordgo.Session) {
	c := cron.New()
	c.AddFunc("CRON_TZ=Europe/Paris 0 20 * * SUN,THU", func() {
		s.ChannelMessageSend(ChannelID, randomSentence("red"))
	})
	c.AddFunc("CRON_TZ=Europe/Paris 0 20 * * TUE", func() {
		s.ChannelMessageSend(ChannelID, randomSentence("green"))
	})
	c.AddFunc("CRON_TZ=Europe/Paris 0 12 * * MON,WED,FRI", func() {
		s.ChannelMessageSend(ChannelID, randomSentence("out"))
	})
	fmt.Println("Schedule are set")
	c.Start()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "!bot getChannelID":
		s.ChannelMessageSend(m.ChannelID, m.ChannelID)
	default:
		return
	}
}
