package main

import (
	"encoding/json"
	"os"
)

func ReadJson() map[string]Port {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	jsonData := make(map[string]Port)

	decodedFile := json.NewDecoder(file)
	err = decodedFile.Decode(&jsonData)
	if err != nil {
		panic(err)
	}

	return jsonData
}
