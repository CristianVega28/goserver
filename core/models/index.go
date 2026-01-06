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
		conn        *sql.DB
		TableName   string
		Fields      []db.MetadataTable
		PrimaryKey  string
		OtherModels []Models[T]
		Response    []map[string]any
	}

	ModelsI[T any] interface {
		Select(id string, columns []string) ([]map[string]any, error)
		Insert(m []map[string]any, meta []db.MetadataTable) error
		Update(m map[string]any, primaryKey string) error
		InsertMigration(isInsert bool) error
		SelectAll() []T
		Init() Models[T]
		SetMetadataTable(fields []db.MetadataTable)
		SetTableName(name string)
		GetTableName() string
		SetPrimaryKey(key string)
		GetPrimaryKey() string
		ValidateFields(bodyR any) map[string]any
		AddModels(m Models[T])
		GetResponse() []map[string]any
		SetResponse(res any)
		ParserColumn(arr db.Migration) []string
		GenerateMetadata(model any) []db.MetadataTable
	}
	DB struct {
		Conn *sql.DB
	}
)

func (base *Models[T]) Init() Models[T] {
	return Models[T]{
		conn:        db.Connect(),
		OtherModels: []Models[T]{},
	}
}
func (model *Models[T]) InsertMigration(isInsert bool) error {
	var rawSql string
	response := model.GetResponse()
	rawSql = db.InsertIntoTableRawSql(model.TableName, response, model.Fields, isInsert)

	if rawSql != "" {
		_, err := model.conn.Exec(rawSql)
		if err != nil {
			utils.Log.Fatal(err.Error())
			return err
		}

	}

	return nil
}
func (model *Models[T]) Insert(data []map[string]any, meta []db.MetadataTable) error {
	var rawSql string
	conn := db.Connect()
	verifed, _ := db.CheckAndTableInDatabase(model.TableName, conn)

	if !verifed {
		var errorT = fmt.Errorf("the table %s does not exist in the database", model.TableName)
		return errorT
	}

	rawSql = db.InsertIntoTableRawSql(model.TableName, data, meta, true)

	if rawSql != "" {
		_, err := conn.Exec(rawSql)
		if err != nil {
			utils.Log.Fatal(err.Error())
			return err
		}

	}

	return nil
}

func (model *Models[T]) Update(data map[string]any, primaryKey string) error {
	var rawSql string
	conn := db.Connect()
	verifed, _ := db.CheckAndTableInDatabase(model.TableName, conn)

	if !verifed {
		var errorT = fmt.Errorf("the table %s does not exist in the database", model.TableName)
		return errorT
	}
	rawSql = db.UpdateIntoTableRawSql(model.TableName, data, primaryKey)

	if rawSql != "" {
		_, err := conn.Exec(rawSql)
		if err != nil {
			utils.Log.Fatal(err.Error())
			return err
		}

	}

	return nil
}

func (model *Models[T]) Select(id string, columns []string) ([]map[string]any, error) {

	conn := db.Connect()
	verifed, _ := db.CheckAndTableInDatabase(model.TableName, conn)

	if !verifed {
		var errorT = fmt.Errorf("the table %s does not exist in the database", model.TableName)
		return nil, errorT
	}

	sqlRow := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?;", "*", model.TableName, model.PrimaryKey)
	rows, errQuery := conn.Query(sqlRow, id)
	if errQuery != nil {
		return nil, errQuery
	}

	if errQuery == sql.ErrNoRows {
		return nil, nil
	}

	defer rows.Close()

	results := []map[string]any{}
	_column, errColumn := rows.Columns()
	if errColumn != nil {
		return nil, errColumn
	}

	for rows.Next() {
		values := make([]interface{}, len(_column))
		valuePtrs := make([]interface{}, len(_column))
		for i := range _column {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		rowMap := make(map[string]any)

		for i, col := range _column {
			v := values[i]
			// MySQL devuelve []byte para strings
			if b, ok := v.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = v
			}
		}

		results = append(results, rowMap)

	}

	return results, nil
}

func (model *Models[T]) SelectAll() []map[string]any {

	var valuesArr []map[string]any
	rows, err := model.conn.Query("SELECT * FROM " + model.TableName)

	if err != nil {
		utils.Log.Fatal(err.Error())
	}

	defer rows.Close()

	cols, err := rows.Columns()

	if err != nil {
		utils.Log.Fatal(err.Error())
	}

	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range cols {
			valuePtrs[i] = &values[i]
		}

		// 2. Escanear la fila
		if err := rows.Scan(valuePtrs...); err != nil {
			utils.Log.Fatal(err.Error())
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

func (model *Models[T]) AddModels(m Models[T]) {
	model.OtherModels = append(model.OtherModels, m)
}

func (model *Models[T]) GetResponse() []map[string]any {
	return model.Response
}

func (model *Models[T]) SetResponse(res any) {
	response, err := utils.CheckTypesForResponse(res)
	if err != nil {
		utils.Log.Fatal(err.Error())
	}
	model.Response = response
}
func (model *Models[T]) GenerateMetadata(modelAny any) []db.MetadataTable {
	var metadata []db.MetadataTable
	mapMeta := utils.ReturnMetadataTable(modelAny, "db")

	for _, v := range mapMeta {
		meta := db.MetadataTable{
			Field: v["Field"],
			Type:  v["Type"],
		}
		metadata = append(metadata, meta)

	}
	return metadata
}

func (model *Models[T]) ParserColumn(arr db.Migration) []string {
	var columns []string
	for field, _ := range arr.Fields {
		columns = append(columns, field)
	}
	return columns
}
