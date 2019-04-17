package utils 

import (
	"os"
	"strings"
)

type FileStructure struct {
	Children map[string]FileStructure
	// parent *FileStructure
	Name string
	Content string
	Filetype int
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
		newPath := path+"/"+g.Name
		os.MkdirAll(newPath, os.ModePerm);
		for _, child := range g.Children {
			child.SaveSelfAt(newPath)
		} 
	} else {
		filePath := path+"/"+g.Name
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.FileMode(0777))

		if err != nil {
			return  err
		}
		
		if _, err := f.WriteString(g.Content); err != nil {
			return  err
		}
		
		f.Close()
	}
	return nil
}

func (g *FileStructure) SaveAt(path string) error {
	if(g.Filetype == 0) {
		for _, child := range g.Children {
			child.SaveSelfAt(path)
		} 
	}
	return nil
}