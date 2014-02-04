package main;

import (
    "io"
    "time"
)

type BlobInfo struct {
    Name string
    Size int64
    Created time.Time
    Modified time.Time
    MimeType string
    Tag string
    Deleted bool
    Nodes []string
    WriteLocked bool
}

type DataStorage interface {
    Save(b *BlobInfo, src io.Reader) (int64, error)
    Open(b *BlobInfo) (io.ReadSeeker, error)
    Delete(b *BlobInfo) error
}

type MetaStorage interface {
    Get(name string) (*BlobInfo, error)
    Set(b *BlobInfo) error
    Delete(b *BlobInfo) error
}
