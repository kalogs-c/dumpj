package filemanager

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func UnzipFile(path string, dest string) error {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		fpath := filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, file.Mode())
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), file.Mode()); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			if _, err = io.Copy(outFile, fileReader); err != nil {
				return err
			}
		}
	}

	return nil
}
