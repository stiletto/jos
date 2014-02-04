package main;

import (
    "encoding/hex"
    "encoding/json"
//    "crypto/sha256"
    "crypto/sha1"
    "path"
    "io"
    "io/ioutil"
    "os"
)

type DataStorageFS struct {
    Root string
}

func getPath(root string, b string) string {
    h := sha1.New()
    hname := hex.EncodeToString(h.Sum([]byte(b)))
    return path.Join(root,hname);
}

func (storage *DataStorageFS) Save(b *BlobInfo, src io.Reader) (int64, error) {
    fname := getPath(storage.Root, b.Name)
    f, err := os.OpenFile(fname+".tmp", os.O_CREATE | os.O_EXCL | os.O_WRONLY, 0666)
    if err != nil {
        return 0, err
    }

    num, err := io.Copy(f, src)
    if err != nil {
        return 0, err
    }
    f.Close()

    err = os.Rename(fname+".tmp", fname)
    return num, err
}

func (storage *DataStorageFS) Open(b *BlobInfo) (io.ReadSeeker, error) {
    fname := getPath(storage.Root, b.Name)
    f, err := os.OpenFile(fname, os.O_RDONLY, 0666)
    return f, err
}

func (storage *DataStorageFS) Delete(b *BlobInfo) error {
    fname := getPath(storage.Root, b.Name)
    os.Remove(fname+".tmp")
    return os.Remove(fname)
}

type MetaStorageFS struct {
    Root string
}

func (storage *MetaStorageFS) Set(b *BlobInfo) error {
    fname := getPath(storage.Root, b.Name)
    f, err := os.OpenFile(fname+".meta.tmp", os.O_CREATE | os.O_EXCL | os.O_WRONLY, 0666)
    if err != nil {
        return err
    }

    bm, err := json.Marshal(b)
    if err != nil {
        return err
    }
    _, err = f.Write(bm)
    if err != nil {
        return err
    }
    f.Close()

    err = os.Rename(fname+".meta.tmp", fname+".meta")
    return err
}

func (storage *MetaStorageFS) Get(name string) (*BlobInfo, error) {
    fname := getPath(storage.Root, name)
    f, err := os.OpenFile(fname+".meta", os.O_RDONLY, 0666)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    bm, err := ioutil.ReadAll(f)
    if err != nil {
        return nil, err
    }

    b:= new(BlobInfo)
    err = json.Unmarshal(bm, &b)

    return b, err
}

func (storage *MetaStorageFS) Delete(b *BlobInfo) error {
    fname := getPath(storage.Root, b.Name)
    os.Remove(fname+".meta.tmp")
    return os.Remove(fname+".meta")
}
