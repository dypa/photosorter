package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func IsDirOrFileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func DirectoryIterator(dir string) ([]string, error) {
	fileList := make([]string, 0)
	err := filepath.Walk(dir, func(filePath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		fileList = append(fileList, filePath)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func CreateDirectoryIfNotExists(path string) error {
	return os.MkdirAll(path, 0777)
}

func CopyFileContents(src string, dst string) (err error) {
	fdSrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fdSrc.Close()

	dstTmp := dst + ".temp"

	fdDstTemp, err := os.Create(dstTmp)
	if err != nil {
		return err
	}
	defer fdDstTemp.Close()

	if _, err = io.Copy(fdDstTemp, fdSrc); err != nil {
		return err
	}

	// https://www.joeshaw.org/dont-defer-close-on-writable-files/
	fdSrc.Close()
	fdDstTemp.Close()

	stat, err := os.Stat(src)
	if err != nil {
		return err
	}

	//sync if > 10Mb
	if stat.Size() > 10000000 {
		fdDstTemp.Sync()
	}

	if os.Chtimes(dstTmp, stat.ModTime(), stat.ModTime()); err != nil {
		return err
	}
	if os.Chmod(dstTmp, stat.Mode()); err != nil {
		return err
	}

	//sometimes rename return error, but file renamed
	if err = os.Rename(dstTmp, dst); err != nil {
		fmt.Println("[DEBUG] rename error", err.Error())
		if !IsDirOrFileExists(dst) {
			return err
		}
	}

	return nil
}

func ReadLastModFile(filepath string) (time.Time, error) {
	file, err := os.Stat(filepath)

	if err != nil {
		return time.Time{}, err
	}

	return file.ModTime(), nil
}

func ReadFileSize(filepath string) (size int64) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return -1
	}

	return stat.Size()
}
