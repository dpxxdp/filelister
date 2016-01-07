package main

import (
	"fmt"
	"time"
    "path/filepath"
    "errors"
    "bytes"
	"strings"
    "os"
    "log"
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

//ScanTreeForPath scans the FileInfo and its children for a pointer to the given path, if it doesn't exist it returns an error
//useful for finding a "parent" directory to attach a child under
func (f *FileInfo) ScanTreeForPath(path string) (*FileInfo, error) {
    
    //if this file is the one we're looking for, return it
    if(path == f.Path) {
        return f, nil
    }
    
    //if not, scan the children recursively
    for i := range f.Children {
        if(path == f.Children[i].Path) {
            return &f.Children[i], nil
        }
        
        ptr,_ := f.Children[i].ScanTreeForPath(path)
        if(ptr != nil) {
            return ptr, nil
        }
    }
    
    //if we make it all the way through, the path we are looking for is not in this tree
    return nil, errors.New("these are not the droids you are looking for")
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

//DirInfo is a slice of FileInfo plus metadata TODO
type DirInfo struct {
    Files   []FileInfo
}

func (d *DirInfo) String() string {
    var buffer bytes.Buffer

    for i := range d.Files {
        buffer.WriteString(d.Files[i].BuildStringGivenDepth(0))
    }

    fmt.Println(buffer.String())
    return buffer.String()
}

//AppendFileInfo adds a file to DirInfo under its appropriate parent directory
func (d *DirInfo)AppendFileInfo(newFile *FileInfo) {
    
    //parse the filepath for its parent directory so we can check to see if we have it listed
    parentPath := filepath.Dir(newFile.Path)
    
    //if there are no files listed, add it to the top level
    if len(d.Files) == 0 {
        d.Files = append(d.Files, *newFile)
    }
    
    //otherwise find its parent in the tree and add it underneath its parent
    for i := range d.Files {
        parentPtr, err := d.Files[i].ScanTreeForPath(parentPath)
        if err == nil {
            parentPtr.Children = append(parentPtr.Children, *newFile)
        }
    }
}

//WalkAndBuildFileInformation wraps filepath.WalkFunc in a closure
//This allows us to take advantage of WalkFunc while building up our DirInfo struct.
//It also allows us to pass in the recursive flag and take advantage of it.
func (dirPtr *DirInfo) WalkAndBuildFileInformation(recursive bool) filepath.WalkFunc {
	return filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		fileInfo := FileInfo {
			Path: 			path,
			ModTime: 		info.ModTime(),
			IsLink: 		info.Mode() & os.ModeSymlink == os.ModeSymlink,
			IsDir: 			info.IsDir(),
			Size:			info.Size(),
			Name:			info.Name(),
		}
		if(fileInfo.IsLink) {
			symlink, _ := os.Readlink(fileInfo.Path)
			fileInfo.LinksTo = symlink
		}
        
        dirPtr.AppendFileInfo(&fileInfo)
        
        if(!recursive) {
            return filepath.SkipDir
        }
        
        return nil
    })
}