package rb

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func BackupLog(fileNum int, sourcePath, targetPath string) {
	fmt.Printf("file #%d (%s -> %s)\n", fileNum, sourcePath, targetPath)
}

func Backup(sourceFilePath, logPath, sourcePathRoot, targetPathRoot string, i int, startTime time.Time) (string, string, time.Time, error) {
	relativePath, err := filepath.Rel(sourcePathRoot, sourceFilePath)
	if err != nil {
		return "", "", time.Unix(0, 0), err
	}
	targetFilePath, err := filepath.Abs(fmt.Sprintf("%s/%s", targetPathRoot, relativePath))
	if err != nil {
		return "", "", time.Unix(0, 0), err
	}
	BackupLog(i, sourceFilePath, targetFilePath)
	modTime, err := Copy(sourceFilePath, targetFilePath, targetPathRoot)
	if err != nil {
		return "", "", time.Unix(0, 0), err
	}
	return sourceFilePath, targetFilePath, modTime, nil
}

func Copy(sourcePath, targetPath string, targetPathRoot string) (time.Time, error) {
	fileStatSource, err := os.Stat(sourcePath)
	if err != nil {
		WaitForDirectory(targetPathRoot)
		return Copy(sourcePath, targetPath, targetPathRoot)
	}
	if !fileStatSource.Mode().IsRegular() {
		return time.Unix(0, 0), fmt.Errorf("%s is not a regular file", sourcePath)
	}
	_, err = os.Stat(targetPath)
	if err != nil {
		err := os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			WaitForDirectory(targetPathRoot)
			return Copy(sourcePath, targetPath, targetPathRoot)
		}
	}
	src, err := os.Open(sourcePath)
	if err != nil {
		return time.Unix(0, 0), err
	}
	defer src.Close()

	dest, err := os.Create(targetPath)

	if err != nil {
		WaitForDirectory(targetPathRoot)
		return Copy(sourcePath, targetPath, targetPathRoot)
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	return fileStatSource.ModTime(), err
}

func WaitForDirectory(path string) {
	fmt.Printf("Waiting for directory %s to be available...\n", path)
	var searching = true
	for searching {
		_, err := os.Stat(path)
		if err != nil {
			time.Sleep(2 * time.Second)
		} else {
			searching = false
		}
	}
}
