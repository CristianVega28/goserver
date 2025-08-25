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
