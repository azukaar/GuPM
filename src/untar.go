package main

import "archive/tar"
import "compress/gzip"

import (
	"os"
	"io"
	"strings"
	"path/filepath"
	// "fmt"
)

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(dst string, r string) error {
	gzr, err := gzip.NewReader(strings.NewReader(r))
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created

		// fmt.Println("--",  header.Name, header.Typeflag, tar.TypeDir)

		if(header.Name != "package") {	
			var target = filepath.Join(dst, header.Name[8:])

			// the following switch could also be done using fi.Mode(), not sure if there
			// a benefit of using one vs. the other.
			// fi := header.FileInfo()
			
			
			// check the file type
			switch header.Typeflag {

			// if its a dir and it doesn't exist create it
			case tar.TypeDir: {
				os.MkdirAll(target, os.ModePerm);
			}

			// if it's a file create it
			case tar.TypeReg:
				var folders = strings.Split(target, "/")
				var folder = folders[:len(folders)-1]
				
				os.MkdirAll(strings.Join(folder[:], "/"), os.ModePerm);

				f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
				if err != nil {
					return err
				}

				// copy over contents
				if _, err := io.Copy(f, tr); err != nil {
					return err
				}
				
				// manually close here after each file operation; defering would cause each file close
				// to wait until all operations have completed.
				f.Close()
			}
		}		
	}
}