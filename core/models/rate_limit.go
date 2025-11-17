package models

import "github.com/CristianVega28/goserver/core/db"

type RateLimit struct {
	Requests  int
	Time      string // in seconds (#s, $m, #h)
	Ip        string
	Scope     string
	CreatedAt string
	UpdateAt  string
	Models[map[string]any]
}

func (r *RateLimit) SeederTable() {

	var mgn db.Migration = db.Migration{
		TableName: "rate_limits",
		Fields: map[string]string{
			"requests":   "integer",
			"time":       "text",
			"ip":         "string",
			"scope":      "string",
			"created_at": "datetime",
			"updated_at": "datetime",
		},
		Foreigns: []db.ForeignKey{},
	}

	db.ExecSqlTable(mgn)

}

func (r *RateLimit) refreshData(ip string) error {

	return nil
}
