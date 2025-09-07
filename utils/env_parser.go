package utils

import (
	"bufio"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	DB_PATH             string
	SESSIONS_DB_NAME    string
	TOKEN_EXPIRATION    string
	SESSION_COOKIE_NAME string
	ITEMS_PER_PAGE      string
}

var AppConfig *Config

// simple env loader
// ignores line comment but not inline
func LoadConfig(filepath string) error {
	var file *os.File
	var err error
	var scanner *bufio.Scanner
	var line string
	var partials []string
	var key string
	var value string

	file, err = os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	// load into ENV
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		partials = strings.SplitN(line, "=", 2)
		key = strings.TrimSpace(partials[0])
		value = strings.TrimSpace(partials[1])
		// Remove surrounding quotes if present
		if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
			(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
			value = value[1 : len(value)-1]
		}
		os.Setenv(key, value)
	}

	if err = scanner.Err(); err != nil {
		return err
	}
	// read from ENV into Config
	AppConfig = &Config{}
	field_values := reflect.ValueOf(AppConfig).Elem()
	field_structs := field_values.Type()

	for i := 0; i < field_values.NumField(); i++ {
		field := field_values.Field(i)
		field_name := field_structs.Field(i).Name

		if field.CanSet() && field.Kind() == reflect.String {
			val := os.Getenv(field_name)
			field.SetString(val)
		}
	}

	return nil
}
