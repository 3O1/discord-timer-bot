package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	Token     = "your-token"      // replace with your bot token
	ChannelID = "your-channel-id" // replace with your channel id
	EndDate   = "2030-10-14T00:00:00Z"
)

func main() {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Open a websocket connection to Discord
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// update the channel name immediately and then every 24 hours
	updateChannelName(discord)
	ticker := time.NewTicker(24 * time.Hour)

	go func() {
		for {
			select {
			case <-ticker.C:
				updateChannelName(discord)
			}
		}
	}()

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	<-make(chan struct{})
}

func updateChannelName(s *discordgo.Session) {
	channel, err := s.Channel(ChannelID)
	if err != nil {
		fmt.Println("error fetching channel,", err)
		return
	}

	now := time.Now()
	t, err := time.Parse(time.RFC3339, EndDate)
	if err != nil {
		fmt.Println("error parsing end date,", err)
		return
	}

	daysLeft := int(t.Sub(now).Hours() / 24)
	newName := fmt.Sprintf("Launching in: %d days", daysLeft)

	_, err = s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Name: newName,
	})
	if err != nil {
		fmt.Println("error editing channel,", err)
		return
	}

	fmt.Println("Successfully updated channel name to:", newName)
}
