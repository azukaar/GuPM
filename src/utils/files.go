package utils 

import (
	"os"
	"strings"
	// "fmt"
	"path/filepath"
	"github.com/otiai10/copy"
)

var EmptyFileStructure = FileStructure{}

type FileStructure struct {
	Children map[string]FileStructure
	Name string
	Content []byte
	Filetype int
}

func Dir(path string) (matches []string, err error) {
	return filepath.Glob(path)
}

func RemoveFiles(files []string) error {
	for _, file := range files {
		return os.RemoveAll(file)
	}
	return nil
}

func CopyFiles(files []string, destination string) error {
	for _, file := range files {
		return copy.Copy(file, destination)
	}
	return nil
}

func (g *FileStructure) getOrCreate(path string, options FileStructure) FileStructure {
	var folders = strings.Split(path, "/")
	var folder = folders[0]
	var child, _ = g.Children[folder]
	
	if(child.Name == "") {
		if(len(folders) > 1) {
			g.Children[folder] = FileStructure{
				Children: make(map[string]FileStructure),
				Name: folder,
				Filetype: 0,
			}
		} else {			
			g.Children[folder] = FileStructure{
				Children: make(map[string]FileStructure),
				Name: folder,
				Filetype: options.Filetype,
				Content: options.Content,
			}
		}
		child, _ = g.Children[folder]
	}

	if(len(folders) > 1) {
		next := folders[1:len(folders)]
		return child.getOrCreate(strings.Join(next[:], "/"), options)
	} else {
		return child
	}
}

func (g *FileStructure) SaveSelfAt(path string) error {
	if(g.Filetype == 0) {
		newPath := Path(path+"/"+g.Name)
		os.MkdirAll(newPath, os.ModePerm);
		for _, child := range g.Children {
			child.SaveSelfAt(newPath)
		} 
	} else {
		filePath := path
		if(g.Name != "") {
			filePath = Path(filePath + "/" +g.Name)
		}
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.FileMode(0777))

		if err != nil {
			f.Close()
			return  err
		}
		
		if _, err := f.Write(g.Content); err != nil {
			f.Close()
			return  err
		}
		
		f.Close()
	}
	return nil
}

func (g *FileStructure) SaveAt(path string) error {
	if(g.Filetype == 0) {
		os.MkdirAll(Path(path), os.ModePerm);
		for _, child := range g.Children {
			child.SaveSelfAt(Path(path))
		} 
	}
	if(g.Filetype == 1) {
		g.SaveSelfAt(Path(path))
	}
	return nil
}
