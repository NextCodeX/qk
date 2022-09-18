package core

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (this *InternalFunctionSet) Zip(path, dst string) {
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	if !dirExists(path) {
		runtimeExcption(path, " is not found!")
		return
	}
	dstZipFileName := getDstZipFileName(path, dst)
	dstZipFileDir := filepath.Dir(dstZipFileName)

	mkdirIfNotExists(dstZipFileDir)

	fmt.Println("creating zip archive:", dstZipFileName)
	archive, err := os.Create(dstZipFileName)
	if err != nil {
		panic(err)
	}
	defer ioclose(archive)

	zipWriter := zip.NewWriter(archive)

	relativeIndex := len(path)
	if isDir(path) {
		var srcFiles []*FileInfo
		doScanForInfo(path, &srcFiles)
		for _, srcFile := range srcFiles {
			if srcFile.IsDir() {
				continue
			}

			fileAbsPath := srcFile.Path()
			relativePath := fileAbsPath[relativeIndex:]
			writeToZip(zipWriter, fileAbsPath, relativePath)
		}
	} else {
		relativePath := path[relativeIndex:]
		writeToZip(zipWriter, path, relativePath)
	}

	ioclose(zipWriter)
	fmt.Printf("create zip archive: %v successfully!\n", dstZipFileName)
}

func writeToZip(zipWriter *zip.Writer, fileAbsPath string, relativePath string) {
	fobj, err := os.Open(fileAbsPath)
	if err != nil {
		panic(err)
	}
	defer ioclose(fobj)

	zw, err := zipWriter.Create(relativePath)
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(zw, fobj); err != nil {
		panic(err)
	}
}

func ioclose(archive io.Closer) {
	err := archive.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func getDstZipFileName(path string, dst string) string {
	dst = strings.TrimSpace(dst)
	if dst != "" {
		if filepath.IsAbs(dst) {
			if isDir(dst) {
				return filepath.Join(dst, filepath.Base(path)+".zip")
			} else {
				return dst
			}
		}
		return filepath.Join(filepath.Dir(path), dst)
	}
	if isDir(path) {
		return filepath.Join(filepath.Dir(path), filepath.Base(path)+".zip")
	}
	return path + ".zip"
}

func (this *InternalFunctionSet) Unzip(path, dst string) {
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	if !dirExists(path) {
		runtimeExcption(path, " is not found!")
		return
	}
	unzipDir := getUnzipDir(path, dst)

	mkdirIfNotExists(unzipDir)

	archive, err := zip.OpenReader(path)
	if err != nil {
		panic(err)
	}
	defer ioclose(archive)

	for _, f := range archive.File {
		filePath := filepath.Join(unzipDir, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(unzipDir)+string(os.PathSeparator)) {
			fmt.Println("invalid file path:", filePath)
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("mkdir:", filePath)
			mkdirIfNotExists(filePath)
			continue
		}

		mkdirIfNotExists(filepath.Dir(filePath))

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		ioclose(dstFile)
		ioclose(fileInArchive)
	}
}

func getUnzipDir(path string, dst string) string {
	dst = strings.TrimSpace(dst)
	if dst != "" {
		if filepath.IsAbs(dst) {
			return dst
		}
		return filepath.Join(filepath.Dir(path), dst)
	}
	if strings.HasSuffix(path, ".zip") {
		return path[:len(path)-4]
	} else {
		return path + "_unzip"
	}
}
