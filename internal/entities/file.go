package entities

import "bytes"

type FileInfo struct {
	Name        string
	ContentType string
	Size        int64
}

type File struct {
	Buffer bytes.Buffer
	FileInfo
}
