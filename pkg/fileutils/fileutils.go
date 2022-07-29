package fileutils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codeclysm/extract/v3"
)

func CopyFile(src, dst string) error {
	// open src
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// create dst
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// copy content
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// get src permissions
	fileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	// give same permissions
	if err := os.Chmod(dst, fileStat.Mode()); err != nil {
		return err
	}

	return out.Close()
}

func DoesFileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func GetAbsolutePath(relativePath string) string {

	fullPath, err := filepath.Abs(relativePath)
	if err != nil {
		panic(err)
	}
	return fullPath
}

func GetPathToCurrentBinary() (string, error) {
	currentFilePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	// resolve any symlinks
	resolvedFilePath, err := filepath.EvalSymlinks(currentFilePath)
	if err != nil {
		return "", err
	}

	return resolvedFilePath, nil
}

func HasWritePermissionToFile(filePath string) (bool, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0744)
	if err != nil {
		if errors.Is(err, fs.ErrPermission) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	return true, nil
}

func ExtractTarGzFile(sourceFile, target string) error {
	ctx := context.Background()
	data, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(data)
	err = extract.Gz(ctx, buffer, target, func(s string) string { return s })
	if err != nil {
		return err
	}

	return nil
}

func SafeMoveFile(source, target string, showLogs bool) (err error) {
	// Resolve symlinks and use actual path for critical operation
	if source, err = filepath.EvalSymlinks(source); err != nil {
		return err
	}
	if target, err = filepath.EvalSymlinks(target); err != nil {
		return err
	}

	fileExists, err := DoesFileExists(target)
	if err != nil {
		return err
	}

	// if file exists, make a backup and remove
	// file will always exist in case of updates, check to support generic usage
	if fileExists {
		if showLogs {
			fmt.Printf("> Creating backup of existing file (%s)\n", source)
		}
		backupTargetFile := filepath.Dir(source) + filepath.Base(source) + "-backup"
		if err = CopyFile(target, backupTargetFile); err != nil {
			return err
		}

		if err = os.Remove(target); err != nil {
			return err
		}

		defer func() {
			// remove backup: when no panic
			if panicErr := recover(); panicErr != nil {
				err = fmt.Errorf("failed to move file, failed to restore backup: %v", panicErr)
			} else {
				if showLogs {
					fmt.Println("> Removing backup file")
				}
				os.Remove(backupTargetFile)
			}
		}()
	}

	err = os.Rename(source, target)
	if err != nil {
		// if file existed initially, place it back
		if fileExists {
			if showLogs {
				fmt.Println("> Failed to move updated file, restoring from backup")
			}

			backupTargetFile := filepath.Dir(source) + filepath.Base(source) + "-backup"
			if backupErr := os.Rename(backupTargetFile, target); backupErr != nil {
				fmt.Println()
				fmt.Printf("Unable to restore original file \nKindly move %s to %s \n", backupTargetFile, target)
				fmt.Println()
				// panic to skip deletion of
				panic(backupTargetFile)
			}
		}
		return err
	}
	if showLogs {
		fmt.Println("> Move successful")
	}

	return nil
}