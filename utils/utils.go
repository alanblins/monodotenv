/*
Copyright © 2024 Alan Lins <alanblins@gmail.com>
*/
package utils

import (
	"errors"
	"log"
	"os"
)

func IsFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetValue(value string, key string, source string, userfile map[string]string, extendValue string) (string, error) {
	if value == "" && extendValue == "" {
		return "", errors.New("not found value for the key: " + key)
	}

	if value == "" && extendValue != "" {
		value = extendValue
	}

	if source == "value" {
		return value, nil
	}
	if source == "user" {
		if userfile == nil {
			return "", errors.New("no .monodotenv.user.yaml file found")
		}
		val, ok := userfile[value]
		if !ok {
			return "", errors.New(value + "not found in .monodotenv.user.yaml file.")
		}
		return val, nil
	}
	return "", errors.New("some error")
}

func CheckConfigFile() {
	_, err := os.Stat("monodotenv.yaml")
	if os.IsNotExist(err) {
		log.Fatalln("monodotenv.yaml file not found")
	}
}
