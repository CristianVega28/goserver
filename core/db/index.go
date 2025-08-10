package db

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"github.com/CristianVega28/goserver/core/models"
	"github.com/CristianVega28/goserver/utils"
	"github.com/samber/lo"
)

type (
	Migration struct {
		TableName string
		Fields    map[string]string // field and type
	}
	ValuesKey struct {
		TableName       string // name of the table
		Key             string // name field in table
		TypeDb          string // field's type in table
		Size            string // field's size in table
		Constraint      string // field's constraint in table
		ExistsTable     bool   // boolean it verfied the existence of the table
		IndexCurrent    int    // current index in the loop
		LenMissingArray int    // length of the array missing columns
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
			continue
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

	sqlString := parserFieldsToSql(mgn.TableName, mgn.Fields, conn)
	_, err := conn.Exec(sqlString)

	if err != nil {
		log.Fatal("Error creating table: " + err.Error())
	}
}

/*
	key: type,size|constraint
*/

func parserFieldsToSql(tableName string, fields map[string]string, conn *sql.DB) string {
	var tableSQL strings.Builder

	existsTable, columns := models.CheckAndTableInDatabase(tableName, conn)

	if !existsTable {
		tableSQL.WriteString(fmt.Sprintf("\nCREATE TABLE IF NOT EXISTS %s (\n", tableName))
		var index int = 0

		for key, value := range fields {
			// arrType[0] ->  type
			// arrType[1] ->  size
			//separateAtribute[1] -> constraint
			typeV, size, constraint := parserColumnsFields(value)

			var cfgValuesKey ValuesKey = ValuesKey{
				Key:         key,
				ExistsTable: existsTable,
				TypeDb:      typeV,
				Size:        size,
				Constraint:  constraint,
			}

			format := valueInKey(cfgValuesKey)

			if index == (len(fields) - 1) {
				format = strings.Replace(format, ", \n", " \n", 1)
			}

			tableSQL.WriteString(format)

			index++

		}
		tableSQL.WriteString(");")

	} else {
		var keys []string
		for k := range fields {
			keys = append(keys, k)
		}

		missingColumns := lo.Filter(keys, func(col string, _ int) bool {
			return !slices.Contains(columns, col)
		})

		for _, col := range missingColumns {
			value := fields[col]
			typeV, size, constraint := parserColumnsFields(value)
			var cfgValuesKey ValuesKey = ValuesKey{
				TableName:   tableName,
				Key:         col,
				ExistsTable: existsTable,
				TypeDb:      typeV,
				Size:        size,
				Constraint:  constraint,
			}

			format := valueInKey(cfgValuesKey)

			tableSQL.WriteString(format)

		}
	}

	return tableSQL.String()
}

func valueInKey(cfg ValuesKey) string {

	var toSql strings.Builder
	if cfg.TypeDb == "primary_key" {
		return fmt.Sprintf("\t%s INTEGER PRIMARY KEY AUTOINCREMENT, \n", cfg.Key)
	}

	if !cfg.ExistsTable {
		toSql.WriteString(fmt.Sprintf("\t%s %s", cfg.Key, strings.ToUpper(typesByDatabase(cfg.TypeDb))))

		if cfg.Size != "" {
			toSql.WriteString(fmt.Sprintf("(%s)", cfg.Size))
		}

		if cfg.Constraint != "" {
			toSql.WriteString(fmt.Sprintf(" %s", constraintByDatabase(cfg.Constraint)))
		}

		toSql.WriteString(", \n")

	} else {

		toSql.WriteString(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", cfg.TableName, cfg.Key, strings.ToUpper(typesByDatabase(cfg.TypeDb))))
		if cfg.Size != "" {
			toSql.WriteString(fmt.Sprintf("(%s)", cfg.Size))
		}

		if cfg.Constraint != "" {
			toSql.WriteString(fmt.Sprintf(" %s", constraintByDatabase(cfg.Constraint)))
		}

		toSql.WriteString(";")

	}

	return toSql.String()

}

/*
	return (
		type string,
		size string,
		constraint string
	)
*/

func parserColumnsFields(value string) (string, string, string) {
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

	return arrType[0], size, constraint

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
