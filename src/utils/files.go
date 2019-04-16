package utils 

import (
	"os"
	"strings"
)

type FileStructure struct {
	children map[string]FileStructure
	// parent *FileStructure
	name string
	content string
	filetype int
}

func (g *FileStructure) getOrCreate(path string, options FileStructure) FileStructure {
	var folders = strings.Split(path, "/")
	var folder = folders[0]
	var child, _ = g.children[folder]
	
	if(child.name == "") {
		if(len(folders) > 1) {
			g.children[folder] = FileStructure{
				children: make(map[string]FileStructure),
				name: folder,
				filetype: 0,
			}
		} else {			
			g.children[folder] = FileStructure{
				children: make(map[string]FileStructure),
				name: folder,
				filetype: options.filetype,
				content: options.content,
			}
		}
		child, _ = g.children[folder]
	}

	if(len(folders) > 1) {
		next := folders[1:len(folders)]
		return child.getOrCreate(strings.Join(next[:], "/"), options)
	} else {
		return child
	}
}

func (g *FileStructure) SaveSelfAt(path string) error {
	if(g.filetype == 0) {
		newPath := path+"/"+g.name
		os.MkdirAll(newPath, os.ModePerm);
		for _, child := range g.children {
			child.SaveSelfAt(newPath)
		} 
	} else {
		filePath := path+"/"+g.name
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.FileMode(0777))

		if err != nil {
			return  err
		}
		
		if _, err := f.WriteString(g.content); err != nil {
			return  err
		}
		
		f.Close()
	}
	return nil
}

func (g *FileStructure) SaveAt(path string) error {
	if(g.filetype == 0) {
		for _, child := range g.children {
			child.SaveSelfAt(path)
		} 
	}
	return nil
}