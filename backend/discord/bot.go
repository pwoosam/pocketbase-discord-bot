package discord

import (
	"github.com/bwmarrin/discordgo"
)

func CreateBot(authToken string) error {
	discord, err := discordgo.New("Bot " + authToken)
	if err != nil {
		return err
	}

	if err = discord.Open(); err != nil {
		return err
	}

	commands := []*discordgo.ApplicationCommand{&rpsCommand}
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		rpsCommand.Name:        rpsHandler,
		rpsJoinCommandName:     rpsJoinHandler,
		rpsRockCommandName:     rpsChoiceHandler,
		rpsPaperCommandName:    rpsChoiceHandler,
		rpsScissorsCommandName: rpsChoiceHandler,
	}

	if _, err = discord.ApplicationCommandBulkOverwrite(discord.State.Application.ID, "", commands); err != nil {
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
