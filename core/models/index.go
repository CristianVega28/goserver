package models

import (
	"database/sql"
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
		conn      *sql.DB
		TableName string
		Fields    []db.MetadataTable
	}

	ModelsI[T any] interface {
		Select(id string) T
		Insert(m any) error
		Init() Models[T]
		SetMetadataTable(fields []db.MetadataTable)
		SetTableName(name string)
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
func (model *Models[T]) Insert(m any) error {
	var rawSql string
	if mapsInsert, ok := m.([]map[string]any); ok {
		rawSql = db.InsertIntoTableRawSql(model.TableName, mapsInsert, model.Fields)
	}

	_, err := model.conn.Exec(rawSql)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil

}

func (model *Models[T]) Select(id string) T {
	var a T

	return a
}

func (model *Models[T]) SetMetadataTable(fields []db.MetadataTable) {
	model.Fields = fields
}
func (model *Models[T]) SetTableName(name string) {
	model.TableName = name
}
