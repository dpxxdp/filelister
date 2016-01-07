package main

import (
	"fmt"
	"time"
    "path/filepath"
    "errors"
    "bytes"
	"strings"
    "os"
)

//FileInfo represents information about a single file or directory including a slice of its subdirectories
type FileInfo struct {
	Path           string      `json:"-" yaml:"-"`
    Name           string      `json:"name" yaml:"name"`
	ModTime        time.Time   `json:"modified" yaml:"modified"`
	IsLink         bool        `json:"isLink" yaml:"isLink"`
	IsDir          bool        `json:"isDir" yaml:"isDir"`
	LinksTo        string      `json:"linksTo" yaml:"linksTo"`
	Size           int64       `json:"size" yaml:"size"`
	Children       []FileInfo  `json:"children" yaml:"children"`
}

//ScanTreeForPath scans the FileInfo and its children for a pointer to the given path and its depth in the tree
//if it doesn't exist it returns an error
//useful for finding a "parent" directory to attach a child under
func (f *FileInfo) ScanTreeForPath(path string, depth int) (*FileInfo, int, error) {
    //if this file is the one we're looking for, return it
    if path == f.Path {
        return f, depth, nil
    }
    
    //if not, scan the children recursively
    for i := range f.Children {
        if(path == f.Children[i].Path) {
            depth++
            return &f.Children[i],depth, nil
        }
        
        ptr,depth,_ := f.Children[i].ScanTreeForPath(path, depth+1)
        if(ptr != nil) {
            return ptr,depth, nil
        }
    }
    
    //if we make it all the way through, the path we are looking for is not in this tree
    return nil, 0, errors.New("these are not the droids you are looking for")
}

//BuildStringGivenDepth prints the filepath with the given number of tabs in front to simulate directory structure
func (f *FileInfo) BuildStringGivenDepth(depth int) string {
    var buffer bytes.Buffer
    buffer.WriteString(strings.Repeat("\t",depth))
    
    if f.IsLink {
    	buffer.WriteString(fmt.Sprintf("%v* (%v)", filepath.Base(f.Path), f.LinksTo))
    } else {
    	buffer.WriteString(fmt.Sprintf("%v", filepath.Base(f.Path)))
    }
    if f.IsDir {
    	buffer.WriteRune(os.PathSeparator)
    }
    buffer.WriteString("\n")

    for i := range f.Children {
        buffer.WriteString(f.Children[i].BuildStringGivenDepth(depth+1))
    }

    return buffer.String()
}