package helpers

import (
	"fmt"
	"os"

	"github.com/CristianVega28/goserver/core/db"
	"github.com/CristianVega28/goserver/core/models"
	"github.com/CristianVega28/goserver/utils"
	"github.com/samber/lo"
)

type (
	ConfigServerApi struct {
		Request       []string       `json:"request"`
		MiddlewareApi MiddlewareApi  `json:"middleware"`
		Response      any            `json:"response"`
		Schema        map[string]any `json:"schema"`
		Env           map[string]any `json:"env"`
	}

	MiddlewareApi struct {
		Auth     string   `json:"auth"`
		Logging  bool     `json:"logging"`
		Security []string `json:"security"`
	}

	ConfigKeyContext struct{}

	ResponseConfig struct {
		Path          string         `json:"path"`
		Request       []string       `json:"request"`
		MiddlewareApi MiddlewareApi  `json:"middleware"`
		Schema        map[string]any `json:"schema"`
	}
	ConfigServerStatistics struct {
		TotalRequests int
		TotalTables   int
		TotalRecords  int
	}
)

var KeyCfg = ConfigKeyContext{}

func (cfgS *ConfigServerStatistics) Loader(data map[string]any) {
	if len(data) != 0 {
		cfgS.TotalRequests = len(data)
	}

	cnn := db.Connect()

	defer cnn.Close()

}

func (cfg *ConfigServerApi) ReturnMetadataTable() []db.MetadataTable {
	metadata := []db.MetadataTable{}
	for key, value := range cfg.Schema {
		var metaType string

		if key == "table_name" {
			fks := db.ForeignKeysTable(value.(string))
			for _, v := range fks {
				metadata = append(metadata, db.MetadataTable{
					Type:  "INTEGER",
					Field: v,
				})
			}
			continue
		}

		if typeValue, ok := value.(string); ok {
			if typeValue == "primary_key" {
				metaType = "INTEGER"
			} else {
				metaType, _, _ = db.PublicParserColumnsFields(typeValue)
			}
		}

		if _, ok := value.(map[string]any); ok {
			continue
		}
		metadata = append(metadata, db.MetadataTable{
			Type:  metaType,
			Field: key,
		})
	}

	return metadata
}

// Helper to know if the schema (database) exist
func (cfg *ConfigServerApi) ExistSchema() bool {
	var tableNames []string
	var TrueDatabases []bool

	if cfg.Schema != nil {
		tableNames = extractTableName(cfg.Schema)
	}

	conn := db.Connect()

	defer conn.Close()

	for _, v := range tableNames {
		existTable, _ := db.CheckAndTableInDatabase(v, conn)
		TrueDatabases = append(TrueDatabases, existTable)
	}

	return lo.EveryBy(TrueDatabases, func(item bool) bool {
		return item
	})
}

func extractTableName(mapSchema map[string]any) []string {
	var tableNameVar []string
	var tableNameInternal []string
	for index, v := range mapSchema {
		if index == "table_name" {
			tableNameVar = append(tableNameVar, v.(string))
			break
		}

		if vMap, ok := v.(map[string]any); ok {
			tableNameInternal = extractTableName(vMap)
		}
	}

	tableNameVar = append(tableNameVar, tableNameInternal...)
	return tableNameVar
}

func (cfg *ConfigServerApi) PreLoader() {

	// Set enviroment variables

	for i, v := range cfg.Env {

		err := os.Setenv(i, fmt.Sprintf("%v", v))
		if err != nil {
			utils.Log.Fatal(err.Error())
		}
	}

	for _, v := range cfg.MiddlewareApi.Security {
		if v == "rate_limit" {
			exist, _ := db.CheckAndTableInDatabase("rate_limits", db.Connect())

			if !exist {
				rate_limit := models.RateLimit{}
				rate_limit.SeederTable()
			}
		}
	}

}
