package utility

import (
	"os"

	"muzzammil.xyz/jsonc"
)

func ReadJsonFile(path string) (r map[string]interface{}) {
	// read json file
	d, err := os.ReadFile(path)
	if err != nil {
		panic("Unable to read json file")
	}
	err = jsonc.Unmarshal(d, &r)
	if err != nil {
		panic("Unable to parse json file")
	}
	return r
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
