package main

import (
	"fmt"
    "path/filepath"
    "bytes"
    "os"
    "log"
)

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

//AppendFileInfo adds a file to DirInfo under its appropriate parent directory, returns the new file's depth
func (d *DirInfo)AppendFileInfo(newFile *FileInfo) int {
    //parse the filepath for its parent directory so we can check to see if we have it listed
    parentPath := filepath.Dir(newFile.Path)
    
    //if there are no files listed, add it to the top level
    if len(d.Files) == 0 {
        d.Files = append(d.Files, *newFile)
        return 0
    }
    
    //otherwise find its parent in the tree and add it underneath its parent
    for i := range d.Files {
        parentPtr, depth, err := d.Files[i].ScanTreeForPath(parentPath, 1)
        if err == nil {
            parentPtr.Children = append(parentPtr.Children, *newFile)
            return depth
        }
    }
    return 0
}

//WalkAndBuildFileInformation wraps filepath.WalkFunc in a closure
//This allows us to take advantage of WalkFunc while building up our DirInfo struct.
//It also allows us to pass in the recursive flag and take advantage of it.
func (d *DirInfo) WalkAndBuildFileInformation(recursive bool) filepath.WalkFunc {
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
		if fileInfo.IsLink {
			symlink, _ := os.Readlink(fileInfo.Path)
			fileInfo.LinksTo = symlink
		}
        
        depth := d.AppendFileInfo(&fileInfo)
        
        if !recursive && fileInfo.IsDir && depth > 0  {
            return filepath.SkipDir
        }
        
        return nil
    })
}