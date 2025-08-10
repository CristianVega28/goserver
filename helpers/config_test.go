package helpers

import (
	"encoding/json"
	"testing"

	"github.com/CristianVega28/goserver/utils"
)

var looger = utils.Logger{}
var log = looger.Create()

func TestConfigStruct(t *testing.T) {
	var config ConfigServerApi

	var jsonMock string = `{
    
        "request": ["POST", "GET"],
        "middleware": {
            "auth": "bearer",
            "logging": true,
            "db": true,
            "security": ["csrf", "xss"]
        },
        "response": [
            {
                "id": "100001_500001",
                "created_time": "2025-07-26T18:03:21+0000",
                "message": "¡Hola a todos! Este es mi primer post 😄",
                "from": {
                  "id": "100001",
                  "name": "Lucía Fernández"
                },
                "permalink_url": "https://facebook.com/100001/posts/500001"
            },
            {
                "id": "100002_500002",
                "created_time": "2025-07-26T18:05:10+0000",
                "message": "¡Feliz viernes! ¿Qué planes tienen para el fin de semana? 🎉",
                "from": {
                  "id": "100002",
                  "name": "Carlos Pérez"
                },
                "permalink_url": "https://facebook.com/100002/posts/500002"
            },
            {
                "id": "100003_500003",
                "created_time": "2025-07-26T18:07:45+0000",
                "message": "¡Saludos desde la playa! 🏖️",
                "from": {
                  "id": "100003",
                  "name": "Ana Torres"
                },
                "permalink_url": "https://facebook.com/100003/posts/500003"
            }
        ],
        "schema": {
            "table_name": "posts",
            "id": "primary_key",
            "created_time": "datetime",
            "message": "text",
            "from": {
                "table_name": "users",
                "id": "primary_key",
                "email": "varchar,255|unique",
                "name": "varchar,255"
            },
            "permalink_url": "url|not_null"
        }
        
    
}`

	err := json.Unmarshal([]byte(jsonMock), &config)
	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	log.Structs("Config Struct", config)

}
