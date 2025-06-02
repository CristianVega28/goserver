package models

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/CristianVega28/goserver/utils"
)

type (
	BaseModel struct {
		Id         uint
		Created_at time.Time
		Updated_at time.Time
	}

	Models struct {
		conn *sql.DB
	}

	ModelsI[T any] interface {
		Select(id string) T
		Create(model T) T
		Init() Models
	}
)

var looger = utils.Logger{}
var log = looger.Create()

func (base *Models) Init() Models {
	return Models{
		conn: Connect(),
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

func Migration() {

}
