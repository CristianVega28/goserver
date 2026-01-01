package models

import (
	"strconv"
	"time"

	"github.com/CristianVega28/goserver/core/db"
	"github.com/CristianVega28/goserver/utils"
)

var env utils.Env = utils.Env{}

/*
Algorithm: Sliding Window Counter
*/
type RateLimit struct {
	CurrentCount   int       `db:"current_count"`
	LastCount      int       `db:"last_count"`
	TimestampStart time.Time `db:"timestamp_start"`
	Ip             string    `db:"ip"`
	Models[map[string]any]
}

func (r *RateLimit) SeederTable() {

	mgn := r.GetMigration()
	db.ExecSqlTable(mgn)

}

func (r *RateLimit) InsertData() error {
	// colums := r.Models.ParserColumn(r.GetMigration())
	parserMaps := utils.StructToMap(r, "db")

	meta := r.GenerateMetadata(r)

	log.Structs("RateLimit InsertData Meta", meta)
	r.SetTableName("rate_limits")
	r.Insert([]map[string]any{parserMaps}, meta)

	return nil
}

func (r *RateLimit) UpdateData(ip string) error {
	// colums := r.Models.ParserColumn(r.GetMigration())
	return nil
}

func (r *RateLimit) GetEnvTime() int {
	v, _v := env.GetEnv("rate_limit_time")

	if v {
		time, _ := strconv.Atoi(_v)
		return time
	} else {
		return 60 // default seconds
	}

}

func (r *RateLimit) GetEnvLimit() int {
	v, _v := env.GetEnv("rate_limit_requests")

	if v {
		limit, _ := strconv.Atoi(_v)
		return limit
	} else {
		return 60 // default limit
	}

}

func (r *RateLimit) GetMigration() db.Migration {
	var mgn db.Migration = db.Migration{
		TableName: "rate_limits",
		Fields: map[string]string{
			"current_count":   "integer",
			"last_count":      "integer",
			"timestamp_start": "text",
			"ip":              "string",
		},
		Foreigns: []db.ForeignKey{},
	}

	return mgn
}
