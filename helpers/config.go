package helpers

import (
	"github.com/CristianVega28/goserver/core/db"
	"github.com/samber/lo"
)

type (
	ConfigServerApi struct {
		Request       []string       `json:"request"`
		MiddlewareApi MiddlewareApi  `json:"middleware"`
		Response      any            `json:"response"`
		Schema        map[string]any `json:"schema"`
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
)

var KeyCfg = ConfigKeyContext{}

func (cfg *ConfigServerApi) ReturnMetadataTable() []db.MetadataTable {
	metadata := []db.MetadataTable{}
	for key, value := range cfg.Schema {
		var metaType string

		if key == "table_name" {
			continue
		}

		if typeValue, ok := value.(string); ok {
			if typeValue == "primary_key" {
				metaType = "INTEGER"
			} else {
				metaType, _, _ = db.PublicParserColumnsFields(typeValue)
			}
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
