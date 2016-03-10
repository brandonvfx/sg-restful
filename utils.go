package main

import "encoding/json"

func StructToString(s interface{}) string {
	j, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(j)
}
