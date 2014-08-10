package main;

import (
    "sync/atomic"
    "encoding/asn1"
    "fmt"
    "strings"
    "time"
//    "crypto/sha256"
    leveldb "code.google.com/p/go-leveldb"
)

type MetaStorageLDB struct {
    Root string
    opts *leveldb.Options
    cache *leveldb.Cache
    db *leveldb.DB
    uro *leveldb.ReadOptions
    uwo *leveldb.WriteOptions
    itcounter uint32
}

func NewMetaLDB(Root string) (*MetaStorageLDB, error) {
    storage := &MetaStorageLDB{}
    opts := leveldb.NewOptions()
    cache := leveldb.NewLRUCache(10<<20)
    opts.SetCache(cache)
    opts.SetCreateIfMissing(true);
    db, err := leveldb.Open(Root, opts)
    if err != nil {
        cache.Close()
        opts.Close()
        return nil, err
    }
    storage.db = db
    storage.itcounter = 0
    storage.cache = cache
    storage.opts = opts
    storage.uro = leveldb.NewReadOptions()
    storage.uwo = leveldb.NewWriteOptions()

    it := db.NewIterator(storage.uro)
    defer it.Close()
    it.SeekToFirst()
    for it = it; it.Valid(); it.Next() {
        fmt.Printf("key: %s\n", it.Key())
    }
    return storage, nil
}

func mlog(t time.Time, itc uint32) string {
    return fmt.Sprintf("l|%017X|%08X", t.Unix(), itc)
}
func (storage *MetaStorageLDB) Set(b *BlobInfo) error {
    meta_key := "m|"+b.Name
    itcounter := atomic.AddUint32(&storage.itcounter, 1)
    log_key := mlog(b.ModifiedLocal, itcounter)
    asn_b, err := asn1.Marshal(*b)
    if err != nil {
        return err
    }
    wb := leveldb.NewWriteBatch()
    defer wb.Close()
    wb.Put([]byte(log_key), []byte(b.Name))
    wb.Put([]byte(meta_key), asn_b)
    return storage.db.Write(storage.uwo, wb)
}

func (storage *MetaStorageLDB) Get(name string) (*BlobInfo, error) {
    data, err := storage.db.Get(storage.uro, []byte("m|"+name))
    if data == nil || err != nil {
        return nil, err
    }
    bi := &BlobInfo{}
    _, err = asn1.Unmarshal(data, bi)
    if err != nil {
        return nil, err
    }
    return bi, nil
}

func (storage *MetaStorageLDB) Delete(b *BlobInfo) error {
    panic("WTF")
}

type LogIteratorLDB struct {
    storage *MetaStorageLDB
    snapshot *leveldb.Snapshot
    ro *leveldb.ReadOptions
    iterator *leveldb.Iterator
}

func (storage *MetaStorageLDB) GetLogIterator(since *time.Time) (LogIterator, error) {
    it := &LogIteratorLDB{}
    snapshot := storage.db.NewSnapshot()
    ro := leveldb.NewReadOptions()
    ro.SetSnapshot(snapshot)
    iterator := storage.db.NewIterator(ro)
    if since != nil {
        iterator.Seek([]byte(mlog(*since,0)))
    } else {
        iterator.Seek([]byte("l|-----------------|00000000"))
    }

    it.storage = storage
    it.iterator = iterator
    it.snapshot = snapshot
    it.ro = ro
    return it, nil
}

func (it *LogIteratorLDB) GetNext() (*BlobInfo, error) {
    if !it.iterator.Valid() {
        err := it.iterator.GetError()
        return nil, err
    }
    k := string(it.iterator.Key())
    if !strings.HasPrefix(k, "l|") {
        return nil, nil
    }
    bi, err := it.storage.Get(string(it.iterator.Value()))
    it.iterator.Next()
    return bi, err
}

func (it *LogIteratorLDB) Close() (error) {
    it.iterator.Close()
    it.ro.Close()
    it.storage.db.ReleaseSnapshot(it.snapshot)
    return nil
}
