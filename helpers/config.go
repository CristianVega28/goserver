package helpers

import "github.com/CristianVega28/goserver/core/db"

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
)

var KeyCfg = ConfigKeyContext{}

func (cfg *ConfigServerApi) ReturnMetadataTable() []db.MetadataTable {
	metadata := []db.MetadataTable{}
	for key, value := range cfg.Schema {
		var metaType string

		if typeValue, ok := value.(string); ok {

			metaType, _, _ = db.PublicParserColumnsFields(typeValue)
		}

		metadata = append(metadata, db.MetadataTable{
			Type:  metaType,
			Field: key,
		})
	}

	return metadata
}
