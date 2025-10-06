package models

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/CristianVega28/goserver/core/db"
	"github.com/CristianVega28/goserver/utils"
)

type (
	BaseModel struct {
		Id         uint
		Created_at time.Time
		Updated_at time.Time
	}

	Models[T any] struct {
		conn       *sql.DB
		TableName  string
		Fields     []db.MetadataTable
		PrimaryKey string
	}

	ModelsI[T any] interface {
		Select(id string) T
		Insert(m any) error
		InsertMigration(m any, isInsert bool) error
		SelectAll() []T
		Init() Models[T]
		SetMetadataTable(fields []db.MetadataTable)
		SetTableName(name string)
		GetTableName() string
		SetPrimaryKey(key string)
		GetPrimaryKey() string
		ValidateFields(bodyR any) map[string]any
	}
	DB struct {
		Conn *sql.DB
	}
)

var looger = utils.Logger{}
var log = looger.Create()

func (base *Models[T]) Init() Models[T] {
	return Models[T]{
		conn: db.Connect(),
	}
}
func (model *Models[T]) InsertMigration(m any, isInsert bool) error {
	var rawSql string
	if mapsInsert, ok := m.([]map[string]any); ok {
		rawSql = db.InsertIntoTableRawSql(model.TableName, mapsInsert, model.Fields, isInsert)
	}

	if rawSql != "" {
		_, err := model.conn.Exec(rawSql)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}

	}

	return nil
}
func (model *Models[T]) Insert(m any) error {
	// if mapsInsert, ok := m.([]map[string]any); ok {
	// 	rawSql = db.InsertIntoTableRawSql(model.TableName, mapsInsert, model.Fields)
	// }
	return nil
}

func (model *Models[T]) Select(id string) T {
	var a T
	return a
}

func (model *Models[T]) SelectAll() []map[string]any {

	var valuesArr []map[string]any
	rows, err := model.conn.Query("SELECT * FROM " + model.TableName)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer rows.Close()

	cols, err := rows.Columns()

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range cols {
			valuePtrs[i] = &values[i]
		}

		// 2. Escanear la fila
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Fatal(err.Error())
		}

		// 3. Crear map columna → valor
		rowMap := make(map[string]any)
		for i, col := range cols {
			rowMap[col] = values[i]
		}

		// 4. Mostrar resultado
		valuesArr = append(valuesArr, rowMap)
	}

	return valuesArr
}

func (model *Models[T]) ValidateFields(bodyR any) map[string]any {
	var existFields []string
	if body, ok := bodyR.(map[string]any); ok {
		for _, mt := range model.Fields {
			if _, exists := body[mt.Field]; !exists {
				existFields = append(existFields, "Missing field "+mt.Field)
			}
		}
	} else if bodyArr, ok := bodyR.([]map[string]any); ok {
		for i, body := range bodyArr {
			for _, mt := range model.Fields {
				if _, exists := body[mt.Field]; !exists {
					existFields = append(existFields, "Index #"+fmt.Sprintf("%d", i)+" Missing field "+mt.Field)
				}
			}
		}
	}

	if len(existFields) > 0 {
		return map[string]any{
			"errors": existFields,
		}
	}

	return nil
}

func (model *Models[T]) SetMetadataTable(fields []db.MetadataTable) {
	model.Fields = fields
}
func (model *Models[T]) SetTableName(name string) {
	model.TableName = name
}
func (model *Models[T]) GetTableName() string {
	return model.TableName
}

func (model *Models[T]) SetPrimaryKey(pk string) {
	model.PrimaryKey = pk
}

func (model *Models[T]) GetPrimaryKey() string {
	return model.PrimaryKey
}
