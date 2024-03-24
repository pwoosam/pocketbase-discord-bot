package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "udl9dmclhim8sbr",
			"created": "2024-03-23 02:25:50.179Z",
			"updated": "2024-03-23 02:25:50.179Z",
			"name": "rps_game_interaction",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "tmh7h59v",
					"name": "game_id",
					"type": "text",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				}
			],
			"indexes": [],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("udl9dmclhim8sbr")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
