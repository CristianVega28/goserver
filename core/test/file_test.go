package utils_test

import (
	"testing"
)

func TestExtractDataByExtension(t *testing.T) {
	// available extension
	// .json, .yaml, .yml, xml

	// var TargetPath string = filepath.Join("api", "api.json")
	// var result map[string]any = map[string]any{
	// 	"user": []map[string]any{{
	// 		"id":    "1",
	// 		"name":  "John Do",
	// 		"email": "admin@gmail.com",
	// 	}},
	// }

	// // file := File{}

	// response, err := file.ExtractData(TargetPath)
	// jsonDataResult, _ := json.Marshal(result)
	// jsonDataResponse, _ := json.Marshal(response)

	// // fmt.Println(result)
	// // fmt.Println(response)

	// if err != nil {
	// 	t.Errorf(err.Error())

	// }

	// if string(jsonDataResult) != string(jsonDataResponse) {
	// 	t.Errorf("The data is not the same")
	// }

}
