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
			"id": "2pzccfvlzhajatw",
			"created": "2024-10-06 08:53:08.353Z",
			"updated": "2024-10-06 08:53:08.353Z",
			"name": "mmdvm_logs_fm",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "wpu5pxjo",
					"name": "state",
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

		collection, err := dao.FindCollectionByNameOrId("2pzccfvlzhajatw")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
