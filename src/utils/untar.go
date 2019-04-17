package utils

import "archive/tar"
import "compress/gzip"

import (
	"io"
	"strings"
	"bytes"
)

func Untar(r string) (FileStructure, error) {
	gzr, err := gzip.NewReader(strings.NewReader(r))
	root := FileStructure{
		Children: make(map[string]FileStructure),
		Name : "/",
		Filetype: 0,
	}

	if err != nil {
		return FileStructure{}, err
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
			s := buf.String() 
			_ = s

			root.getOrCreate(header.Name, FileStructure{
				Filetype: 1,
				Content: s,
			})
		}	
	}

	return root, nil
}