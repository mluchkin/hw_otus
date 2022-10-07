package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	env := make(Environment, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		env[filename], err = getFileEnvValue(filepath.Join(dir, filename))
		if err != nil {
			return nil, err
		}
	}

	return env, nil
}

func getFileEnvValue(filename string) (EnvValue, error) {
	file, err := os.Open(filename)
	if err != nil {
		return EnvValue{}, err
	}
	defer file.Close()

	scanner := bufio.NewReader(file)
	lineContent, err := scanner.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return EnvValue{}, err
	}

	lineContent = strings.TrimRight(lineContent, "\t\n ")
	lineContent = strings.ReplaceAll(lineContent, "\x00", "\n")

	return EnvValue{lineContent, len(lineContent) == 0}, nil
}
