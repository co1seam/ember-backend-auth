package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"log"
	"os"
	"reflect"
	"strings"
)

func New(path *string) (*Config, error) {
	if path != nil {
		log.Printf("Attempting to load .env file from: %s", path)
		err := godotenv.Load(*path)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("Notice: .env file not found at %s", path)
			} else {
				return nil, fmt.Errorf("error loading .env file: %v", err)
			}
		} else {
			log.Println("Successfully loaded .env file")
		}
	}

	envVars := make(map[string]interface{})
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key := pair[0]
		value := pair[1]
		envVars[key] = value
		log.Printf("Loaded env var: %s = %s", key, value)
	}

	var cfg Config
	decoderConfig := &mapstructure.DecoderConfig{
		Result:           &cfg,
		WeaklyTypedInput: true,
		ErrorUnused:      false,
		TagName:          "mapstructure",
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
				if f.Kind() != reflect.String || t.Kind() != reflect.Bool {
					return data, nil
				}
				switch strings.ToLower(data.(string)) {
				case "true", "1", "yes":
					return true, nil
				case "false", "0", "no":
					return false, nil
				default:
					return nil, fmt.Errorf("invalid boolean value: %s", data)
				}
			},
		),
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, fmt.Errorf("decoder creation failed: %v", err)
	}

	if err := decoder.Decode(envVars); err != nil {
		return nil, fmt.Errorf("decoding failed: %v", err)
	}

	log.Println("Configuration loaded successfully:")
	log.Printf("%+v\n", cfg)

	return &cfg, nil
}
