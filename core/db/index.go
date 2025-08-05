package db

import (
	"fmt"
	"strings"

	"github.com/CristianVega28/goserver/core/models"
	"github.com/CristianVega28/goserver/utils"
)

type (
	Migration struct {
		TableName string
		Fields    map[string]string // field and type
	}
)

var logs = utils.Logger{}
var log = logs.Create()

func MigrateSchema(schema map[string]any) {

	var migration Migration
	fields := make(map[string]string)
	for key, value := range schema {
		log.Msg("Migrating schema for: " + key)

		if key == "table_name" {
			if tablename, ok := value.(string); ok {
				migration.TableName = tablename
			}
		}

		if nested, ok := value.(map[string]any); ok {
			MigrateSchema(nested)
		}

		if typeValue, ok := value.(string); ok {
			fields[key] = typeValue
		}
	}

	migration.Fields = fields

	// log.Structs("Migration schema", migration)

}

func CreateTable(mgn Migration) {
	conn := models.Connect()

	defer conn.Close()

}

/*
	key: type,size|constraint
*/

func parserFields(fields map[string]string) map[string]string {
	fieldsnew := make(map[string]string)

	for key, value := range fields {
		separateAtribute := strings.Split(value, "|")
		arrType := strings.Split(separateAtribute[0], ",")
		formtType := fmt.Sprintf("%s (%s) %s", arrType[0], arrType[1], separateAtribute[1])
		fieldsnew[key] = formtType

	}

	return fieldsnew
}
