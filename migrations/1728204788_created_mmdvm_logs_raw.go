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
			"id": "0fjg5ahesyz3afi",
			"created": "2024-10-06 08:53:08.354Z",
			"updated": "2024-10-06 08:53:08.354Z",
			"name": "mmdvm_logs_raw",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "pmw4pel5",
					"name": "value",
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

		collection, err := dao.FindCollectionByNameOrId("0fjg5ahesyz3afi")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
