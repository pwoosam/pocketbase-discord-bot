package discord

import (
	"github.com/bwmarrin/discordgo"
)

func CreateBot(authToken string) error {
	discord, err := discordgo.New("Bot " + authToken)
	if err != nil {
		return err
	}

	err = discord.Open()
	if err != nil {
		return err
	}

	commands := []*discordgo.ApplicationCommand{&helloCommand, &rpsCommand}
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		helloCommand.Name:      helloHandler,
		rpsCommand.Name:        rpsHandler,
		rpsJoinCommandName:     rpsJoinHandler,
		rpsRockCommandName:     rpsChoiceHandler,
		rpsPaperCommandName:    rpsChoiceHandler,
		rpsScissorsCommandName: rpsChoiceHandler,
	}

	_, err = discord.ApplicationCommandBulkOverwrite(discord.State.Application.ID, "", commands)
	if err != nil {
		return err
	}

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		} else if i.Type == discordgo.InteractionMessageComponent {
			if h, ok := commandHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	return nil
}
