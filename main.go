package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var discord_token string

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	discord_token = os.Getenv("DISCORD_TOKEN")

	discord, err := discordgo.New("Bot " + discord_token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	discord.AddHandler(discordEventHandler)

	err = discord.Open()
	defer discord.Close()
	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return
	}

	// Keep until terminated by Ctrl+C
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Shutting down...")
}

func discordEventHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}

	discord.ChannelMessageSend(message.ChannelID, "I am not set up yet...")
}
