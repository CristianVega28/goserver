package utils

func ParserTypesByDatabases(typeStruct string) string {

	switch typeStruct {
	case "int":
		return "INTEGER"
	case "string":
		return "TEXT"
	case "float64":
		return "REAL"
	case "bool":
		return "BOOLEAN"
	case "time.Time":
		return "TEXT"
	default:
		return "TEXT"
	}

}
