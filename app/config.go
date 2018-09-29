package app

import (
	"bufio"
	"os"
	"reflect"
	"strings"
)

// Config структура представляет config.json файл.
type Config struct {
	AppVersion string

	ServerAddress string

	DbAddress string
	DbName    string
}

// LoadConfig загружает json файл и возвращает Config структуру.
func LoadConfig(path string) *Config {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	var config Config

	// Reflect on config.
	typReflect := reflect.TypeOf(&config).Elem()
	valReflect := reflect.ValueOf(&config).Elem()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineArr := strings.Split(scanner.Text(), "=")
		if len(lineArr) == 2 {
			key := strings.TrimSpace(lineArr[0])
			val := strings.TrimSpace(lineArr[1])
			key = strings.Title(strings.Replace(key, ".", " ", -1))
			key = strings.Replace(key, " ", "", -1)
			_, ok := typReflect.FieldByName(key)
			if ok {
				valReflect.FieldByName(key).SetString(val)
			}
		}
	}

	return &config
}
