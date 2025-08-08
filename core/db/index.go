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

	CreateTable(migration)
}

func CreateTable(mgn Migration) {
	conn := models.Connect()

	defer conn.Close()

	parserFieldsToSql(mgn.TableName, mgn.Fields)
}

/*
	key: type,size|constraint
*/

func parserFieldsToSql(tableName string, fields map[string]string) string {
	var tableSQL strings.Builder

	tableSQL.WriteString(fmt.Sprintf("\nCREATE TABLE IF NOT EXISTS %s (\n", tableName))

	for key, value := range fields {
		// arrType[0] ->  type
		// arrType[1] ->  size
		//separateAtribute[1] -> constraint
		separateAtribute := strings.Split(value, "|")
		arrType := strings.Split(separateAtribute[0], ",")
		var size string = ""
		var constraint string = ""

		if conditionSize := len(arrType) > 1; conditionSize {
			size = arrType[1]
		}

		if conditionConstraint := len(separateAtribute) > 1; conditionConstraint {
			constraint = separateAtribute[1]
		}

		format := valueInKey(key, arrType[0], size, constraint)

		tableSQL.WriteString(format)

	}
	tableSQL.WriteString(")")

	return tableSQL.String()
}

func valueInKey(key string, typeDb string, size string, constraint string) string {

	if typeDb == "primary_key" {
		return fmt.Sprintf("\t%s INTEGER PRIMARY KEY AUTOINCREMENT, \n", key)
	} else {
		if size != "" {
			return fmt.Sprintf("\t%s %s(%s) %s, \n", key, strings.ToUpper(typesByDatabase(typeDb)), size, constraint)
		} else if constraint != "" {
			return fmt.Sprintf("\t%s %s %s, \n", key, strings.ToUpper(typesByDatabase(typeDb)), constraintByDatabase(constraint))
		}
		return fmt.Sprintf("\t%s %s, \n", key, strings.ToUpper(typesByDatabase(typeDb)))

	}
}

func typesByDatabase(typedb string) string {
	switch typedb {
	case "datetime":
		return "DATETIME"
	case "url":
		return "TEXT"
	default:
		return typedb
	}
}

func constraintByDatabase(constraint string) string {
	switch constraint {
	case "not_null":
		return "NOT NULL"
	case "unique":
		return "UNIQUE"
	default:
		return ""
	}
}
