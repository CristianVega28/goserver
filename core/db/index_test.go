package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserFieldsForCreatingTable(t *testing.T) {
	var migrationSchema Migration = Migration{
		TableName: "users",
		Fields: map[string]string{
			"id":            "primary_key",
			"created_at":    "datetime",
			"message":       "text",
			"permalink_url": "url|not_null",
		},
	}
	parsedFields := parserFieldsToSql(migrationSchema.TableName, migrationSchema.Fields)

	var expectedFields string = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	created_at DATETIME, 
	message TEXT, 
	permalink_url TEXT NOT NULL, 
)`

	assert.Equal(t, expectedFields, parsedFields, "The parsed fields should match the expected SQL table creation string")
}
