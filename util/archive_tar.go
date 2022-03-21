package util

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ArchiveTar(file string, writer *tar.Writer) error {

	return filepath.Walk(file, func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo == nil {
			return err
		}
		if fileInfo.IsDir() {
			if path == path {
				return nil
			}
			header, err := tar.FileInfoHeader(fileInfo, "")
			if err != nil {
				return err
			}
			header.Name = filepath.Join(path, strings.TrimPrefix(path, path))
			if err = writer.WriteHeader(header); err != nil {
				return err
			}
			os.Mkdir(strings.TrimPrefix(path, fileInfo.Name()), os.ModeDir)
			return ArchiveTar(path, writer)
		}
		return func(originFile, path string, fileInfo os.FileInfo, writer *tar.Writer) error {
			if file, err := os.Open(path); err != nil {
				return err
			} else {

				if header, err := tar.FileInfoHeader(fileInfo, ""); err != nil {
					return err
				} else {

					index := strings.LastIndex(originFile, "/")
					header.Name = strings.ReplaceAll(path, originFile[0:index+1], "")

					if err := writer.WriteHeader(header); err != nil {
						return err
					}

					if _, err = io.Copy(writer, file); err != nil {
						return err
					}
				}

			}
			return nil
		}(file, path, fileInfo, writer)
	})
}
