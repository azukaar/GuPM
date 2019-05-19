package utils

import "archive/tar"
import "archive/zip"
import "compress/gzip"

import (
	"io"
	"io/ioutil"
	"strings"
	"bytes"
)

func Tar(files []string) (FileStructure, error) {
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)
	finalList := make([]string, 0)

	for _, file := range files {
		if(IsDirectory(file)) {
			finalList = append(finalList, RecursiveFileWalkDir(file)...)
		} else {
			finalList = append(finalList, file)
		}
	}
	
	for _, file := range finalList {
		content, err := ioutil.ReadFile(file)
		if(err != nil) {
			return EmptyFileStructure, err
		}
	
		hdr := &tar.Header{
			Name: file,
			Mode: 0740,
			Size: int64(len(content)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return EmptyFileStructure, err
		}

		if _, err := tw.Write([]byte(content)); err != nil {
			return EmptyFileStructure, err
		}
	}

	if err := tw.Close(); err != nil {
		return EmptyFileStructure, err
	}

	if err := gzw.Close(); err != nil {
		return EmptyFileStructure, err
	}

	root := FileStructure{
		Content: buf.Bytes(),
		Filetype: 1,
	}

	return root, nil
}

func Untar(r string) (FileStructure, error) {
	gzr, err := gzip.NewReader(strings.NewReader(r))
	root := FileStructure{
		Children: make(map[string]FileStructure),
		Name : "/",
		Filetype: 0,
	}

	if err != nil {
		return EmptyFileStructure, err
	}

	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		case err == io.EOF:
			return root, nil

		case err != nil:
			return root, err

		case header == nil:
			continue
		}

		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir: {
			root.getOrCreate(header.Name, FileStructure{
				Filetype: 0,
			})
		}
		
		// if it's a file create it
		case tar.TypeReg:
			buf := new(bytes.Buffer)
			buf.ReadFrom(tr)

			root.getOrCreate(header.Name, FileStructure{
				Filetype: 1,
				Content: buf.Bytes(),
			})
		}	
	}

	return root, nil
}

func Unzip(r string) (FileStructure, error) {
	root := FileStructure{
		Children: make(map[string]FileStructure),
		Name : "/",
		Filetype: 0,
	}

	tr, _ := zip.NewReader(strings.NewReader(r), int64(len(r)))

	for _, f := range tr.File {
		// if its a dir and it doesn't exist create it
		if(f.FileInfo().IsDir()) {
			root.getOrCreate(f.Name, FileStructure{
				Filetype: 0,
			})
		} else {
			rc, _ := f.Open()
			buf := new(bytes.Buffer)
			buf.ReadFrom(rc)

			root.getOrCreate(f.Name, FileStructure{
				Filetype: 1,
				Content: buf.Bytes(),
			})
		}	
	}
	return root, nil
}