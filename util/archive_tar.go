/*
 * Copyright 2025 The ChaosBlade Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
