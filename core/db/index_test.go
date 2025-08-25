package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserFieldsForCreatingTable(t *testing.T) {

	t.Run("Create fields for table users and posts", func(t *testing.T) {
		conn := Connect()
		defer conn.Close()
		var migrationSchema Migration = Migration{
			TableName: "users",
			Fields: map[string]string{
				"id":            "primary_key",
				"created_at":    "datetime",
				"message":       "text",
				"permalink_url": "url|not_null",
			},
		}

		var fields map[string]string = map[string]string{
			"id":            "	id INTEGER PRIMARY KEY AUTOINCREMENT, \n",
			"created_at":    "	created_at DATETIME, \n",
			"message":       "	message TEXT, \n",
			"permalink_url": "	permalink_url TEXT NOT NULL, \n",
		}

		checkFields(migrationSchema.TableName, fields, migrationSchema.Fields, t, false)
		var migrationSchema2 Migration = Migration{
			TableName: "posts",
			Fields: map[string]string{
				"id":      "primary_key",
				"title":   "text",
				"author":  "varchar,125",
				"content": "text|not_null",
			},
		}

		var fieldsWithSizeAndConstrained map[string]string = map[string]string{
			"id":      "	id INTEGER PRIMARY KEY AUTOINCREMENT, \n",
			"title":   "	title TEXT, \n",
			"author":  "	author VARCHAR(125), \n",
			"content": "	content TEXT NOT NULL, \n",
		}

		checkFields(migrationSchema2.TableName, fieldsWithSizeAndConstrained, migrationSchema2.Fields, t, false)

	})

	t.Run("Alter fields for existing table", func(t *testing.T) {

		var migrationSchema Migration = Migration{
			TableName: "users",
			Fields: map[string]string{
				"email": "varchar,255",
				"name":  "varchar,255",
			},
		}

		var fieldsAlter map[string]string = map[string]string{
			"email": "ALTER TABLE users ADD COLUMN email VARCHAR(255);",
			"name":  "ALTER TABLE users ADD COLUMN name VARCHAR(255);",
		}

		checkFields(migrationSchema.TableName, fieldsAlter, migrationSchema.Fields, t, true)
	})
}

func TestVerifyCreateTableWithValueKeyFunction(t *testing.T) {

	t.Run("Create table users with valueKey function", func(t *testing.T) {
		conn := Connect()
		defer conn.Close()
		var migrationSchema Migration = Migration{
			TableName: "users",
			Fields: map[string]string{
				"id":            "primary_key",
				"created_at":    "datetime",
				"message":       "text",
				"permalink_url": "url|not_null",
			},
		}

		sqlString := parserFieldsToSql(migrationSchema.TableName, migrationSchema.Fields, conn)

		_, err := conn.Exec(sqlString)

		if err != nil {
			t.Errorf("Error creating table: %v", err)
		}

		existTable, _ := CheckAndTableInDatabase(migrationSchema.TableName, conn)

		assert.Equal(t, existTable, true, "Rows affected should be 0 for table creation")

	})

	// Above create the table users, now we test when it'll put new columns in existing table
	t.Run("Create columns in existing table", func(t *testing.T) {
		conn := Connect()
		defer conn.Close()
		var migrationSchema Migration = Migration{
			TableName: "users",
			Fields: map[string]string{
				"id":            "primary_key",
				"created_at":    "datetime",
				"message":       "text",
				"permalink_url": "url|not_null",
				"email":         "varchar,80",
				"name":          "varchar,255",
			},
		}

		sqlString := parserFieldsToSql(migrationSchema.TableName, migrationSchema.Fields, conn)
		t.Log("SQL String", sqlString)

		_, err := conn.Exec(sqlString)

		if err != nil {
			t.Errorf("Error creating new columns: %v", err)
		}

	})

}

func TestRawSqlForInsertIntoTable(t *testing.T) {

	t.Run("Insert data when it's one record", func(t *testing.T) {
		data := []map[string]any{{
			"id":            1,
			"created_at":    "2023-10-01 12:00:00",
			"message":       "Hello World",
			"permalink_url": "http://example.com",
			"is_active":     true,
		}}

		metadata := []MetadataTable{
			{
				Type:  "INTEGER",
				Field: "id",
			},
			{
				Type:  "DATETIME",
				Field: "created_at",
			},
			{
				Type:  "TEXT",
				Field: "message",
			},
			{
				Type:  "TEXT",
				Field: "permalink_url",
			},
			{
				Type:  "BOOLEAN",
				Field: "is_active",
			},
		}
		raw := InsertIntoTableRawSql("users", data, metadata)
		assert.Equal(t, raw, "INSERT INTO users (id ,created_at ,message ,permalink_url ,is_active ) VALUES (1 ,'2023-10-01 12:00:00' ,'Hello World' ,'http://example.com' ,true );", "Raw SQL should match expected format")

	})

	t.Run("Insert data when it's multiple records", func(t *testing.T) {
		data := []map[string]any{
			{
				"id":            1,
				"created_at":    "2023-10-01 12:00:00",
				"message":       "Hello World",
				"permalink_url": "http://example.com",
			},
			{
				"id":            2,
				"created_at":    "2023-10-02 13:00:00",
				"message":       "Second Message",
				"permalink_url": "http://example.org",
			},
		}

		metadata := []MetadataTable{
			{
				Type:  "INTEGER",
				Field: "id",
			},
			{
				Type:  "DATETIME",
				Field: "created_at",
			},
			{
				Type:  "TEXT",
				Field: "message",
			},
			{
				Type:  "TEXT",
				Field: "permalink_url",
			},
		}
		raw := InsertIntoTableRawSql("users", data, metadata)
		assert.Equal(t, raw, "INSERT INTO users (id ,created_at ,message ,permalink_url ) VALUES (1 ,'2023-10-01 12:00:00' ,'Hello World' ,'http://example.com' ),(2 ,'2023-10-02 13:00:00' ,'Second Message' ,'http://example.org' );", "Raw SQL should match expected format")
	})
}

func checkFields(tablename string, expected map[string]string, fields map[string]string, t *testing.T, exists bool) {
	expectBoolFormat := make([]bool, len(fields))

	var index int = 0
	for k, v := range fields {
		typeV, size, constraint := parserColumnsFields(v)

		var cfg ValuesKey = ValuesKey{
			TableName:   tablename,
			Key:         k,
			TypeDb:      typeV,
			Size:        size,
			Constraint:  constraint,
			ExistsTable: exists,
		}
		format := valueInKey(cfg)

		if assert.Equal(t, expected[k], format, "Fields SQL") {
			expectBoolFormat[index] = true
		} else {
			expectBoolFormat[index] = false
		}

		index++
	}

	var paseedFormatSql int = 0

	for _, v := range expectBoolFormat {
		if v {
			paseedFormatSql++
		}
	}

}
