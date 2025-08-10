package models

import (
	"database/sql"
	"fmt"
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
		conn      *sql.DB
		TableName string
		Fields    []string
	}

	ModelsI[T any] interface {
		Select(id string) T
		Create(model T) T
		Init() Models
	}
	DB struct {
		Conn *sql.DB
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
