package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDirNoR(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.Contains(name, "=") {
			continue
		}

		filePath := filepath.Join(dir, name)
		fmt.Println(filePath)
		value, err := readEnvFile(filePath)
		if err != nil {
			return nil, err
		}

		env[name] = value
	}

	return env, nil
}

func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, errWalk error) error {
			if errWalk != nil {
				return errWalk
			}
			if info.Mode().IsDir() {
				return nil
			}
			name := info.Name()
			if strings.Contains(name, "=") {
				return nil
			}
			value, err := readEnvFile(path)
			if err != nil {
				return err
			}
			env[name] = value
			return nil
		})
	return env, err
}

func readEnvFile(path string) (EnvValue, error) {
	info, err := os.Stat(path)
	if err != nil {
		return EnvValue{}, err
	}

	// Если файл пустой, помечаем переменную на удаление
	if info.Size() == 0 {
		return EnvValue{NeedRemove: true}, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return EnvValue{}, err
	}
	defer file.Close()

	// Читаем только первую строку
	content, err := io.ReadAll(file)
	if err != nil {
		return EnvValue{}, err
	}

	// Находим первую строку (до \n или до конца файла)
	value := string(bytes.Split(content, []byte("\n"))[0])

	// Заменяем нулевые байты на \n
	value = strings.ReplaceAll(value, string(byte(0)), "\n")

	// Убираем пробелы и табуляции в конце
	value = strings.TrimRight(value, " \t")

	return EnvValue{
		Value:      value,
		NeedRemove: false,
	}, nil
}
