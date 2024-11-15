package main

import (
	commands "dendenmushi/commands/info"
	"dendenmushi/handlers"
	"log"
	"os"
	"os/signal"
	"flag"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	SyncCommands = flag.Bool("sync", false, "Whether to sync slash commands or not")
)

func main() {
	flag.Parse() 
	
	if !*SyncCommands {
		log.Printf("Sync is disabled")
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file with error: %v", err)
	}

	s, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	if err != nil {
		log.Fatalf("Failed to create discord client with error error: %v", err)
	}

	s.AddHandler(func (s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v", s.State.User.Username)
	})

	if err := s.Open(); err != nil {
		log.Fatalf("Failed to open session with error: %v", err)
	}

	cmdHandler := handlers.NewCmdHandler(s)
	cmdHandler.Add([]handlers.SlashCmd{&commands.PingCmd{}}, *SyncCommands)

	s.AddHandler(func (s *discordgo.Session, i *discordgo.InteractionCreate) {
		cmdHandler.Handle(i.Interaction)
	})

	defer s.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop 
}