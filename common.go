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
    ModifiedLocal time.Time
    MimeType string
    Tag string
    Deleted bool
    Nodes []string
    WriteLocked bool
    DataStorage string
    LogRecord string
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
    GetLogIterator(since *time.Time) (LogIterator, error)
    GetListIterator(prefix string) (LogIterator, error)
}

type LogIterator interface {
    GetNext() (*BlobInfo, error)
    Close() error
}
