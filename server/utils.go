package server

func checkTypesForResponse(value any) []map[string]any {
	var response = make([]map[string]any, 0)
	if arr, ok := value.([]any); ok {
		for _, v := range arr {
			if m, ok := v.(map[string]any); ok {
				response = append(response, m)
			}
		}
	}

	return response

}
