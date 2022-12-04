package commands

import (
	"errors"
	"fmt"
	"path/filepath"
	"photosorter/utils"
	"strings"
)

type CheckCommand struct {
	ArchivePath string
}

func (currentCommand CheckCommand) getDirectoryForIteration() string {
	return currentCommand.ArchivePath
}

func (currentCommand CheckCommand) CanRun(file string) bool {
	return utils.IsAllowedExtension(file)
}

func (currentCommand CheckCommand) Run(file string) error {
	return currentCommand.checkFile(file)
}

func (currentCommand CheckCommand) checkFile(filePath string) error {
	array := strings.Split(filepath.Base(filePath), ".")
	arrayLen := len(array)
	slice := array[arrayLen-2 : arrayLen-1]

	if len(slice) == 0 {
		return errors.New("hash in file name not found")
	}

	hashFromFileName := slice[0]

	hash, _ := utils.FileContentHash(filePath)

	if hashFromFileName != hash {
		fmt.Println("[FAIL]", filePath, "Expected:", hashFromFileName, "Current:", hash)
	} else {
		fmt.Println("[OK]", filePath)
	}

	return nil
}
