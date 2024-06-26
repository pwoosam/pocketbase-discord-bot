package discord

import (
	"fmt"
	"log"
	"myapp/gamedata"
	"myapp/service"

	"github.com/bwmarrin/discordgo"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

var rpsCommand = discordgo.ApplicationCommand{
	Name:        "rps",
	Description: "Rock, Paper, Scissors",
}

func rpsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userID string
	if i.Interaction.User != nil {
		userID = i.Interaction.User.ID
	} else {
		userID = i.Interaction.Member.User.ID
	}

	var (
		newGame *models.Record
		err     error
	)
	if err = service.App.Dao().RunInTransaction(func(dao *daos.Dao) error {
		var err error
		newGame, err = gamedata.RPSCreateGame(dao, userID)
		if err != nil {
			return err
		}
		if err = gamedata.RPSUpsertGameInteraction(dao, i.Interaction.ID, newGame.Id); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("failed to upsert rps interaction:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprint("Failed to create interaction: ", err),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("<@!%s> created game. Waiting for players to join.", userID),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Join",
							CustomID: rpsJoinCommandName,
							Style:    discordgo.PrimaryButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "👍",
							},
						},
					},
				},
			}},
	}); err != nil {
		log.Println("failed to send rps response:", err)
	}
}

const (
	rpsJoinCommandName     = "rps_join"
	rpsRockCommandName     = "rps_rock"
	rpsPaperCommandName    = "rps_paper"
	rpsScissorsCommandName = "rps_scissors"
)

func rpsCommandNameToChoice(commandName string) gamedata.RPSChoice {
	switch commandName {
	case rpsRockCommandName:
		return gamedata.Rock
	case rpsPaperCommandName:
		return gamedata.Paper
	case rpsScissorsCommandName:
		return gamedata.Scissors
	default:
		return gamedata.Undecided
	}
}

func rpsChoiceToString(choice gamedata.RPSChoice) string {
	switch choice {
	case gamedata.Rock:
		return "rock"
	case gamedata.Paper:
		return "paper"
	case gamedata.Scissors:
		return "scissors"
	default:
		return "undecided"
	}
}

func rpsChoiceToEmoji(choice gamedata.RPSChoice) string {
	switch choice {
	case gamedata.Rock:
		return "🪨"
	case gamedata.Paper:
		return "📄"
	case gamedata.Scissors:
		return "✂️"
	default:
		return ""
	}
}

func rpsJoinHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var originalUserID string
	if i.Message.Interaction.User != nil {
		originalUserID = i.Message.Interaction.User.ID
	} else {
		originalUserID = i.Message.Member.User.ID
	}

	var (
		game *models.Record
		err  error
	)
	if game, err = gamedata.RPSGetGameByInteractionId(service.App.Dao(), i.Message.Interaction.ID); err != nil {
		log.Println("failed to get rps game from db: ", err)
		return
	} else {
		log.Println(game)
	}

	var userID string
	if i.Interaction.User != nil {
		userID = i.Interaction.User.ID
	} else {
		userID = i.Interaction.Member.User.ID
	}

	if err = gamedata.RPSJoinGame(service.App.Dao(), game.Id, userID); err != nil {
		log.Println("failed to join rps game:", err)
		if err = InteractionRespondNewMessageEphemeral(s, i, fmt.Sprint("Failed to join game: ", err), []discordgo.MessageComponent{}); err != nil {
			log.Println("failed to send rps join response:", err)
		}
		return
	}

	if err = InteractionRespondUpdateMessage(s, i, fmt.Sprintf("<@!%s> created game. <@!%s> joined.", originalUserID, userID), []discordgo.MessageComponent{}); err != nil {
		log.Println("failed to send rps join response:", err)
		return
	}

	var newMsg *discordgo.Message
	if newMsg, err = InteractionFollowupMessage(s, i, fmt.Sprintf("<@!%s> vs. <@!%s>", originalUserID, userID),
		[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Rock",
						CustomID: rpsRockCommandName,
						Style:    discordgo.PrimaryButton,
						Emoji: discordgo.ComponentEmoji{
							Name: "🪨",
						},
					},
					discordgo.Button{
						Label:    "Paper",
						CustomID: rpsPaperCommandName,
						Style:    discordgo.PrimaryButton,
						Emoji: discordgo.ComponentEmoji{
							Name: "📄",
						},
					},
					discordgo.Button{
						Label:    "Scissors",
						CustomID: rpsScissorsCommandName,
						Style:    discordgo.PrimaryButton,
						Emoji: discordgo.ComponentEmoji{
							Name: "✂️",
						},
					},
				},
			},
		}); err != nil {
		log.Println("failed to send rps followup message:", err)
		return
	}

	if err = gamedata.RPSUpsertGameInteraction(service.App.Dao(), newMsg.ID, game.Id); err != nil {
		log.Println("failed to update rps game:", err)
		if err = InteractionRespondUpdateMessageEphemeral(s, i, fmt.Sprint("Failed to update game: ", err), []discordgo.MessageComponent{}); err != nil {
			log.Println("failed to send rps join response:", err)
		}
	}
}

func rpsChoiceHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var (
		choice gamedata.RPSChoice
		game   *models.Record
		err    error
	)

	if choice = rpsCommandNameToChoice(i.MessageComponentData().CustomID); choice == gamedata.Undecided {
		log.Println("invalid rps choice:", choice)
		if err = InteractionRespondNewMessageEphemeral(s, i, fmt.Sprintf("Invalid choice %s!", rpsChoiceToString(choice)), []discordgo.MessageComponent{}); err != nil {
			log.Println("failed to send rps invalid choice response:", err)
		}
		return
	}

	if game, err = gamedata.RPSGetGameByInteractionId(service.App.Dao(), i.Message.ID); err != nil {
		if err = InteractionRespondNewMessageEphemeral(s, i, fmt.Sprint("Failed to get game: ", err), []discordgo.MessageComponent{}); err != nil {
			log.Println("failed to send rps choice response:", err)
		}
		return
	}
	if err = gamedata.RPSMakeChoice(service.App.Dao(), game.Id, i.Interaction.Member.User.ID, choice); err != nil {
		if err = InteractionRespondNewMessageEphemeral(s, i, fmt.Sprint("Failed to make choice: ", err), []discordgo.MessageComponent{}); err != nil {
			log.Println("failed to send rps choice response:", err)
		}
		return
	}

	if game, err = gamedata.RPSGetGameByInteractionId(service.App.Dao(), i.Message.ID); err != nil {
		if err = InteractionRespondNewMessageEphemeral(s, i, fmt.Sprint("Failed to get updated game: ", err), []discordgo.MessageComponent{}); err != nil {
			log.Println("failed to send rps choice response:", err)
		}
		return
	}
	if err = gamedata.RPSUpsertGameInteraction(service.App.Dao(), i.Message.ID, game.Id); err != nil {
		if err = InteractionRespondNewMessageEphemeral(s, i, fmt.Sprint("Failed to update game interaction: ", err), []discordgo.MessageComponent{}); err != nil {
			log.Println("failed to send rps choice response:", err)
		}
		return
	}

	if gamedata.RPSGameStatus(game.GetInt("status")) == gamedata.RPSGameStatusFinished {
		winnerId := game.GetString("player_id_winner")
		player1Id := game.GetString("player1_id")
		player2Id := game.GetString("player2_id")
		player1Emoji := rpsChoiceToEmoji(gamedata.RPSChoice(game.GetInt("player1_choice")))
		player2Emoji := rpsChoiceToEmoji(gamedata.RPSChoice(game.GetInt("player2_choice")))
		tied := winnerId == ""
		message := fmt.Sprintf("<@!%s> vs. <@!%s> - tied!", player1Id, player2Id)
		if !tied {
			message = fmt.Sprintf("<@!%s>%s vs. <@!%s>%s - <@!%s> won!", player1Id, player1Emoji, player2Id, player2Emoji, winnerId)
		}

		if err = InteractionRespondUpdateMessage(s, i, message, []discordgo.MessageComponent{}); err != nil {
			log.Println("failed to update rps scissors response:", err)
		}

		message2 := "Tied!"
		if !tied {
			if winnerId == player1Id {
				message2 = fmt.Sprintf("%s beats %s, <@!%s> won!", player1Emoji, player2Emoji, winnerId)
			} else {
				message2 = fmt.Sprintf("%s beats %s, <@!%s> won!", player2Emoji, player1Emoji, winnerId)
			}
		}
		if _, err = InteractionFollowupMessage(s, i, message2, []discordgo.MessageComponent{}); err != nil {
			log.Println("failed to send rps scissors result:", err)
		}
		return
	}

	if err = InteractionRespondNewMessageEphemeral(s, i, fmt.Sprintf("<@!%s> chose %s!", i.Interaction.Member.User.ID, rpsChoiceToString(choice)), []discordgo.MessageComponent{}); err != nil {
		log.Println("failed to send rps scissors response:", err)
	}
}
