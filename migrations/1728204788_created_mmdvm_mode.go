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
			"id": "t91kzxxouy6l3nk",
			"created": "2024-10-06 08:53:08.354Z",
			"updated": "2024-10-06 08:53:08.354Z",
			"name": "mmdvm_mode",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "bshtto24",
					"name": "mode",
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
			"listRule": "",
			"viewRule": "",
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

		collection, err := dao.FindCollectionByNameOrId("t91kzxxouy6l3nk")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
