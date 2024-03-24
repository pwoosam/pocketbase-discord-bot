package gamedata

import (
	"myapp/service"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

func CreateCollections() error {
	return service.App.Dao().RunInTransaction(func(dao *daos.Dao) error {
		if dao.IsCollectionNameUnique(rpsCollection.Name) {
			rpsCollection.MarkAsNotNew()
		}
		if err := dao.SaveCollection(rpsCollection); err != nil {
			return err
		}

		if dao.IsCollectionNameUnique(rpsInteractionCollection.Name) {
			rpsInteractionCollection.MarkAsNotNew()
		}
		if err := dao.SaveCollection(rpsInteractionCollection); err != nil {
			return err
		}

		return nil
	})
}

var rpsCollection = &models.Collection{
	Type: models.CollectionTypeBase,
	Name: "rps_game",
	Schema: schema.NewSchema(
		&schema.SchemaField{
			Name:     "status",
			Type:     schema.FieldTypeNumber,
			Required: true,
		},
		&schema.SchemaField{
			Name:     "player1_id",
			Type:     schema.FieldTypeText,
			Required: true,
		},
		&schema.SchemaField{
			Name:     "player2_id",
			Type:     schema.FieldTypeText,
			Required: false,
		},
		&schema.SchemaField{
			Name:     "player1_choice",
			Type:     schema.FieldTypeNumber,
			Required: false,
		},
		&schema.SchemaField{
			Name:     "player2_choice",
			Type:     schema.FieldTypeNumber,
			Required: false,
		},
		&schema.SchemaField{
			Name:     "player_id_winner",
			Type:     schema.FieldTypeText,
			Required: false,
		},
		&schema.SchemaField{
			Name:     "player_id_loser",
			Type:     schema.FieldTypeText,
			Required: false,
		},
	),
}

var rpsInteractionCollection = &models.Collection{
	Type: models.CollectionTypeBase,
	Name: "rps_game_interaction",
	Schema: schema.NewSchema(
		&schema.SchemaField{
			Name:     "game_id",
			Type:     schema.FieldTypeText,
			Required: true,
		},
	),
}
