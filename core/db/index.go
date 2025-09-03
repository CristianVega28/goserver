package db

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"github.com/CristianVega28/goserver/utils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
)

type (
	//struct in charge to create/modify a table
	Migration struct {
		TableName string
		Fields    map[string]string // field and type
	}
	MetadataTable struct {
		Type  string
		Field string
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

func ExecSqlTable(mgn Migration) {
	conn := Connect()

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

	existsTable, columns := CheckAndTableInDatabase(tableName, conn)

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

func InsertIntoTableRawSql(tableName string, data []map[string]any, metadataTable []MetadataTable) string {

	var insertSql strings.Builder
	insertSql.WriteString(fmt.Sprintf("INSERT INTO %s (", tableName))

	for index, value := range metadataTable {

		insertSql.WriteString(fmt.Sprintf("%s ", value.Field))
		if len(metadataTable)-1 != index {
			insertSql.WriteString(",")
		}
	}

	insertSql.WriteString(") VALUES ")

	count := reviewLengthValues(tableName)

	if count == 0 {
		return insertSqlFunc(&insertSql, data, metadataTable)
	} else {
		dataBk := data[count:]
		return insertSqlFunc(&insertSql, dataBk, metadataTable)
	}
}

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", "file:database.db?cache=shared&mode=rwc")

	if err != nil {
		log.Fatal("Failed to connect to the database: " + err.Error())
	}

	errPing := db.Ping()
	if errPing != nil {
		log.Fatal("Failed to ping the database: " + errPing.Error())
	} else {
		log.Msg("Connected to the database (SQLite)")
	}

	return db
}

/*
return (

	existTable bool,
	columns []string

)
*/
func CheckAndTableInDatabase(name string, conn *sql.DB) (bool, []string) {

	var existTable bool = false

	rows, err := conn.Query(fmt.Sprintf("PRAGMA table_info(%s)", name))

	if err != nil {
		log.Fatal(err.Error())
	}

	defer rows.Close()

	existTable = rows.Next()
	if !existTable {
		return existTable, nil
	}

	rowsCol, err := conn.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 0", name))
	if err != nil {
		log.Fatal(err.Error())
	}

	defer rowsCol.Close()

	cols, err := rowsCol.Columns()
	if err != nil {
		log.Fatal(err.Error())
	}

	return existTable, cols
}
