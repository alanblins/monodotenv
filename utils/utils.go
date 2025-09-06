/*
Copyright © 2024 Alan Lins <alanblins@gmail.com>
*/
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/alanblins/monodotenv/models"
	"gopkg.in/yaml.v3"
)

func IsFileExist(pathval string, myOs models.OsI) (bool, error) {
	_, error := myOs.Stat(pathval)
	if error == nil {
		return true, nil
	}
	return false, error

}

func ValidateEnvironmanetVariable(environmentVariable models.EnvironmentVariable, environment string, extendEnvironment string, userFile map[string]string) (bool, error) {

	return true, nil
}

func GetValue(environmentVariable models.EnvironmentVariable, extendEnvironment string, environment string, userfile map[string]string, secretsFile models.SecretsYaml) (string, error) {
	environmentFinal := environment
	value, environmentExist := environmentVariable.Environments[environment]
	if !environmentExist {
		if extendEnvironment != "" {
			value, environmentExist = environmentVariable.Environments[extendEnvironment]
			if !environmentExist {
				return "", errors.New("not found value for the key: " + environmentVariable.Key + " and environment: " + environment)
			}
			environmentFinal = extendEnvironment
		} else {
			return "", errors.New("not found value for the key: " + environmentVariable.Key + " and environment: " + environment)
		}
	}

	if environmentVariable.Source == "" {
		return value, nil
	}

	if environmentVariable.Source == "value" {
		return value, nil
	}

	if environmentVariable.Source == "user" {
		if userfile == nil {
			return "", errors.New("no .monodotenv.user.yaml file found")
		}
		val, ok := userfile[value]
		if !ok {
			return "", errors.New(value + "not found in .monodotenv.user.yaml file.")
		}
		return val, nil
	}

	if environmentVariable.Source == "aes-gcm" {
		if secretsFile.Secrets == nil {
			return "", errors.New("no .monodotenv.secrets.yaml file found")
		}
		keyHex := secretsFile.Secrets[environmentFinal][environmentVariable.Key][0]
		nounceHex := secretsFile.Secrets[environmentFinal][environmentVariable.Key][1]
		value = GCMDecrypter(keyHex, value, nounceHex)
		if value == "" {
			return "", errors.New(value + "not found in .monodotenv.user.yaml file.")
		}
		return value, nil
	}
	return "", errors.New("source unknown")
}

func ReadYaml[T models.ConfigYaml | map[string]string | models.SecretsYaml](filepath string, yamlDecoded *T) error {
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &yamlDecoded)
	if err != nil {
		return err
	}
	return nil
}

func WriteContent(environmentVariable models.EnvironmentVariable, configYaml *models.ConfigYaml, environment string, outputEnvMap map[string]string, pathval string, userFileYaml map[string]string, secretsFileYaml models.SecretsYaml) {
	extendEnvironment := configYaml.Extends[environment]
	value, errorReadValue := GetValue(environmentVariable, extendEnvironment, environment, userFileYaml, secretsFileYaml)
	if errorReadValue != nil {
		log.Fatalln(errorReadValue)
	}
	content := outputEnvMap[pathval] + environmentVariable.Key + "=" + value + "\n"
	outputEnvMap[pathval] = content
}

func WriteContentDocLine(contents [][]string, environmentVariable models.EnvironmentVariable, configYaml *models.ConfigYaml, environments []string, pathval string, userFileYaml map[string]string, secretsFileYaml models.SecretsYaml) [][]string {
	content := []string{}
	content = append(content, environmentVariable.Key)
	content = append(content, environmentVariable.Name)
	content = append(content, environmentVariable.Description)
	content = append(content, pathval)
	for _, environment := range environments {
		value, errorReadValue := GetValue(environmentVariable, environment, environment, userFileYaml, secretsFileYaml)
		if errorReadValue != nil {
			log.Fatalln(errorReadValue)
		}
		content = append(content, value)
	}
	contents = append(contents, content)
	return contents
}

func WriteFile(pathval string, outputBytes []byte, suffix string) error {
	envFile := ".env"
	if suffix != "" {
		envFile += "." + suffix
	}
	finalPath := path.Join(pathval, envFile)
	return os.WriteFile(finalPath, outputBytes, 0644)
}

func DryWriteFile(pathval string, outputBytes []byte, suffix string) {
	envFile := ".env"
	if suffix != "" {
		envFile += "." + suffix
	}
	finalPath := path.Join(pathval, envFile)
	var myString = string(outputBytes[:])
	fmt.Println(finalPath)
	fmt.Println(myString)
}

func GCMEncrypter(keyString string, textString string, nonceHex string) (string, string) {
	// AES-128 or AES-256.
	key := []byte(keyString)
	plaintext := []byte(textString)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, 12)
	if nonceHex == "" {
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err.Error())
		}
		nonceHex = hex.EncodeToString(nonce)
	} else {
		nonce, err = hex.DecodeString(nonceHex)
		if err != nil {
			panic(err.Error())
		}
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	cipherHex := hex.EncodeToString(ciphertext)
	return cipherHex, nonceHex
}

func GCMDecrypter(keyString string, ciphertextHex string, nonceHex string) string {
	// AES-128 or AES-256.
	key := []byte(keyString)
	ciphertext, _ := hex.DecodeString(ciphertextHex)

	nonce, _ := hex.DecodeString(nonceHex)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return string(plaintext)
}
