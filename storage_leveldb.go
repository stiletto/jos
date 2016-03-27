package main;

import (
    "sync/atomic"
    "encoding/asn1"
    "fmt"
    "strings"
    "time"
//    "crypto/sha256"
    //leveldb "code.google.com/p/go-leveldb"
    leveldb "github.com/syndtr/goleveldb/leveldb"
    leveldb_cache "github.com/syndtr/goleveldb/leveldb/cache"
    leveldb_opt "github.com/syndtr/goleveldb/leveldb/opt"
    leveldb_iterator "github.com/syndtr/goleveldb/leveldb/iterator"
)

type MetaStorageLDB struct {
    Root string
    opts *leveldb_opt.Options
    cache *leveldb_cache.Cache
    db *leveldb.DB
    uro *leveldb_opt.ReadOptions
    uwo *leveldb_opt.WriteOptions
    itcounter uint32
}

func NewMetaLDB(Root string) (*MetaStorageLDB, error) {
    storage := &MetaStorageLDB{}
    opts := &leveldb_opt.Options{}
    //cache := leveldb_cache.NewLRU(10<<20)
    //opts.SetCache(cache)
    //opts.SetCreateIfMissing(true);
    db, err := leveldb.OpenFile(Root, opts)
    if err != nil {
        //cache.Close()
        //opts.Close()
        return nil, err
    }
    storage.db = db
    storage.itcounter = 0
    //storage.cache = cache
    storage.opts = opts
    storage.uro = &leveldb_opt.ReadOptions{}
    storage.uwo = &leveldb_opt.WriteOptions{}

    it := db.NewIterator(nil, storage.uro)
    defer it.Release()
    //it.SeekToFirst()
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
    wb := new(leveldb.Batch)
    wb.Put([]byte(log_key), []byte(b.Name))
    wb.Put([]byte(meta_key), asn_b)
    return storage.db.Write(wb, storage.uwo)
}

func (storage *MetaStorageLDB) Get(name string) (*BlobInfo, error) {
    data, err := storage.db.Get([]byte("m|"+name), storage.uro)
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

type MetaIteratorLDB struct {
    storage *MetaStorageLDB
    snapshot *leveldb.Snapshot
    ro *leveldb_opt.ReadOptions
    iterator leveldb_iterator.Iterator
    list bool
    seekto string
}

func (storage *MetaStorageLDB) GetLogIterator(since *time.Time) (LogIterator, error) {
    var seekto string
    if since != nil {
        seekto = mlog(*since,0)
    } else {
        seekto = "l|-----------------|00000000"
    }
    it := NewMetaIteratorLDB(storage, seekto)
    it.list = false
    return it, nil
}

func NewMetaIteratorLDB(storage *MetaStorageLDB, seekto string) (*MetaIteratorLDB) {
    it := &MetaIteratorLDB{}
    snapshot, err := storage.db.GetSnapshot()
    if err != nil {
	panic(err.Error())
    }	
    ro := &leveldb_opt.ReadOptions{}
    //ro.SetSnapshot(snapshot)
    iterator := snapshot.NewIterator(nil, ro)
    iterator.Seek([]byte(seekto))
    it.seekto = seekto

    it.storage = storage
    it.iterator = iterator
    it.snapshot = snapshot
    it.ro = ro
    it.list = false
    return it
}

func (it *MetaIteratorLDB) GetNext() (bi *BlobInfo, err error) {
    if !it.iterator.Valid() {
        err := it.iterator.Error()
        return nil, err
    }
    k := string(it.iterator.Key())
    if it.list {
        if !strings.HasPrefix(k, it.seekto) {
            return nil, nil
        }
        bi = &BlobInfo{}
        _, err = asn1.Unmarshal(it.iterator.Value(), bi)
    } else {
        if !strings.HasPrefix(k, "l|") {
            return nil, nil
        }
        bi, err = it.storage.Get(string(it.iterator.Value()))
    }
    it.iterator.Next()
    return bi, err
}

func (it *MetaIteratorLDB) Close() (error) {
    it.iterator.Release()
    it.snapshot.Release()
    return nil
}

func (storage *MetaStorageLDB) GetListIterator(prefix string) (LogIterator, error) {
    it := NewMetaIteratorLDB(storage, "m|"+prefix)
    it.list = true
    return it, nil

}
