package main

import (
	"fmt"
	"time"
)

type FileInformation struct {
	Path string
	ModTime time.Time
	IsLink bool
	IsDir bool
	LinksTo string
	Size int64
	Name string
	Children []*FileInformation
}

func (f *FileInformation) String() string {
	return fmt.Sprintf("%v\n\t%v\n\t%v\n\t%v\n\t%v\n\t%v\n\t%v\n\n", f.Path, f.ModTime, f.IsLink, f.IsDir, f.LinksTo, f.Size, f.Name)
}