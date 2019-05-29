package app

import (
	"bufio"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Config represents a config struct to the app.
type Config struct {
	DBURI      string
	DBName     string
	DBUser     string
	DBPassword string

	MDRoot string
}

// NewConfig returns a config struct and fills it with default values if path is not provide.
func NewConfig(path string) (*Config, error) {
	if path == "" {
		return &Config{
			DBURI:      "mongodb://localhost:27017",
			DBName:     "website",
			DBUser:     "root",
			DBPassword: "",

			MDRoot: "data",
		}, nil
	}

	return LoadConfig(path)
}

// LoadConfig reads path of .conf file and loads it into config strcut.
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	var config Config

	// Reflect on config.
	valReflect := reflect.ValueOf(&config).Elem()

	// Scan every string.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty strings and comments.
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split line into key and value.
		slice := strings.Split(line, "=")
		if len(slice) != 2 {
			continue
		}

		// Process key to check if config struct contains it.
		key := strings.TrimSpace(slice[0])
		key = strings.ReplaceAll(key, ".", "")
		key = strings.ReplaceAll(key, "_", "")
		fieldValue := valReflect.FieldByNameFunc(func(field string) bool {
			return strings.ToLower(field) == strings.ToLower(key)
		})
		if fieldValue == (reflect.Value{}) {
			continue
		}

		val := slice[1]
		// Trim inline comments.
		if strings.ContainsRune(val, '#') {
			val = strings.TrimRightFunc(val, func(r rune) bool {
				return r != '#'
			})
			val = strings.TrimSuffix(val, "#")
		}
		val = strings.TrimSpace(val)

		// Check if value is string.
		if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
			val = val[1 : len(val)-1]
			fieldValue.SetString(val)
		} else {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			fieldValue.SetInt(int64(intVal))
		}
	}

	return &config, nil
}
