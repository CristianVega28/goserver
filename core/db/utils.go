package db

import (
	"fmt"
	"strings"
)

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
func PublicParserColumnsFields(value string) (string, string, string) {
	return parserColumnsFields(value)
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

func castValueByType(value any, typeDb string) any {
	switch strings.ToUpper(typeDb) {
	case "INTEGER", "REAL", "BOOLEAN":
		return value
	default:
		return fmt.Sprintf("'%v'", value)
	}
}

func reviewLengthValues(tableName string) int {
	db := Connect()

	defer db.Close()

	existable, _ := CheckAndTableInDatabase(tableName, db)

	if !existable {
		return 0
	}

	var count int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)

	if err != nil {
		logs.Fatal(err.Error())
	}

	return count
}

func insertSqlFunc(stringBk *strings.Builder, data []map[string]any, metadataTable []MetadataTable) string {
	for i := 0; i < len(data); i++ {
		info := data[i]
		// var columnRaqSql strings.Builder
		stringBk.WriteString("(")

		for index, value := range metadataTable {
			if _, ok := info[value.Field].(map[string]any); !ok {
				cast := castValueByType(info[value.Field], value.Type)
				stringBk.WriteString(fmt.Sprintf("%v ", cast))
			}

			if index == len(metadataTable)-1 {
				stringBk.WriteString(")")
			} else {
				stringBk.WriteString(",")
			}
		}

		if i != len(data)-1 {
			stringBk.WriteString(",")
		} else {
			stringBk.WriteString(";")
		}

	}

	return stringBk.String()
}
func dropTable(tableName string) error {
	db := Connect()
	defer db.Close()
	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName))
	return err
}
