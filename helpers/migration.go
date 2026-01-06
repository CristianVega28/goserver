package helpers

import (
	"github.com/CristianVega28/goserver/core/db"
	"github.com/CristianVega28/goserver/core/models"
)

type ()

// MigrateSchema creates the tables in database based on the schema provided
// It returns a ModelsI interface which can be used to interact with the database
func MigrateSchema(schema map[string]any) models.ModelsI[map[string]any] {

	var migration db.Migration
	var base models.ModelsI[map[string]any] = &models.Models[map[string]any]{}
	model := base.Init()
	metadata := []db.MetadataTable{}
	var otherForeign []db.ForeignKey
	fields := make(map[string]string)
	for key, value := range schema {
		var metaType string

		if key == "table_name" {
			if tablename, ok := value.(string); ok {
				migration.TableName = tablename
				model.SetTableName(tablename)
			}
			continue
		}
		if typeValue, ok := value.(string); ok {
			fields[key] = typeValue
			metaType, _, _ = db.PublicParserColumnsFields(typeValue)
			if metaType == "primary_key" {
				model.SetPrimaryKey(key)
			}
		}

		if nested, ok := value.(map[string]any); ok {
			om := MigrateSchema(nested)
			foreginId := om.GetTableName() + "_id"

			otherForeign = append(otherForeign, db.ForeignKey{
				Field:          foreginId,
				ReferenceTable: om.GetTableName(),
				ReferenceField: om.GetPrimaryKey(),
			})

			metadata = append(metadata, db.MetadataTable{
				Type:  "integer",
				Field: foreginId,
			})
			fields[foreginId] = "integer"
			model.AddModels(om.Init())
			continue
		}

		metadata = append(metadata, db.MetadataTable{
			Type:  metaType,
			Field: key,
		})
	}

	migration.Fields = fields
	model.SetMetadataTable(metadata)

	migration.Foreigns = otherForeign

	db.ExecSqlTable(migration)

	return &model
}
