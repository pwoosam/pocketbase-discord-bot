package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var helloCommand = discordgo.ApplicationCommand{
	Name:        "hello",
	Description: "Say hello",
}

func helloHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hello, World!",
		},
	})

	if err != nil {
		log.Println("failed to send hello response:", err)
	}
}
