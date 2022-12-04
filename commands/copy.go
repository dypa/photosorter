package commands

import (
	"fmt"
	"path/filepath"
	"photosorter/utils"
	"regexp"
	"strings"
	"time"
)

// // TODO H-i-s_YYYY-MM-DD
// const FileFormat = "15-04-05_2006-01-02"

// // TODO YYYY/MM-YYYY/DD-MM-YYYY
// const DirFormat = "2006/01-2006/02-01-2006"

// format  YYYY-MM-DD H-i-s
const FileFormat = "2006-01-02 15-04-05"

// format YYYY/YYYY-MM-DD
const DirFormat = "2006/2006-01-02"

type CopyCommand struct {
	ArchivePath string
	SourcePath  string
}

func (currentCommand CopyCommand) getDirectoryForIteration() string {
	return currentCommand.SourcePath
}

func (currentCommand CopyCommand) CanRun(file string) bool {
	return utils.IsAllowedExtension(file)
}

func (currentCommand CopyCommand) Run(file string) error {
	isImage := utils.IsExifImage(file)

	return currentCommand.copyFile(file, isImage)
}

func (currentCommand CopyCommand) copyFile(filePath string, useExif bool) error {
	hash, _ := utils.FileContentHash(filePath)

	var err error
	var datetime time.Time

	if datetime, err = determineDatetime(filePath, useExif); err != nil {
		return nil
	}

	newPath, err := generateNewPath(currentCommand, datetime, hash, filePath)
	if err != nil {
		return err
	}

	if !utils.IsDirOrFileExists(newPath) {
		//TODO файлы .MOV имеют не верное время изменения
		err := utils.CopyFileContents(filePath, newPath)
		if err != nil {
			fmt.Println("[FAIL]", filePath)
			return err
		}

		// checksum, _ := utils.FileContentHash(newPath)
		// if hash != checksum {
		// 	fmt.Println("[FAIL]", filePath, "Expected:", hash, "Current:", checksum)
		// } else {
		fmt.Println("[COPY]", filePath, "to", newPath)
		// }
	} else {
		//TODo before skip checksum, if not - write log
		//TODO реализовать проверку на размер файла, если меньше - то можно перезаписывать файл

		size1 := utils.ReadFileSize(filePath)
		size2 := utils.ReadFileSize(newPath)
		if size1 != size2 {
			fmt.Println("[FAIL]", filePath, "size", size1, "not same", newPath, "size", size2)
		} else {
			fmt.Println("[SKIP]", filePath, "found in", newPath)
		}

	}

	return nil
}

func generateNewPath(currentCommand CopyCommand, datetime time.Time, hash string, filePath string) (string, error) {
	newPath := currentCommand.ArchivePath + "/" + datetime.Format(DirFormat)
	err := utils.CreateDirectoryIfNotExists(newPath)
	if err != nil {
		return "", err
	}

	newPath = newPath + "/" + datetime.Format(FileFormat) + "." + hash + strings.ToUpper(utils.GetExtension(filePath))

	return newPath, nil
}

func determineDatetime(filePath string, useExif bool) (time.Time, error) {
	var err error
	var datetime time.Time

	if useExif {
		datetime, err = utils.ExifParseFileDateTime(filePath)
		if err != nil {
			fmt.Println("[INFO]", filePath, err)
		}
	}

	if datetime.IsZero() {
		re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s\d{2}-\d{2}-\d{2}`)
		base := filepath.Base(filePath)
		if re.Match([]byte(base)) {
			datetime, _ = time.Parse(FileFormat, re.FindString(base))
		}
	}

	if datetime.IsZero() {
		if datetime, err = utils.ReadLastModFile(filePath); err != nil {
			return time.Time{}, err
		}

	}

	return datetime, nil
}
