package main;

import (
//    "github.com/hoisie/web"
    "io"
    "errors"
    "hash"
    "encoding/hex"
    "crypto/md5"
    "net/http"
    "fmt"
//    "strconv"
    "os"
    "time"
    "sync"
)

type BlobJar struct {
    Data DataStorage
    Meta MetaStorage
    MutexMap
}

func hello(val string) string {
    return "hello " + val
}

var bj *BlobJar

//func get(ctx *web.Context, val string) {
func get(w http.ResponseWriter, r *http.Request) {
    val := r.URL.Path[len("/get/"):]
    if r.Method != "GET" {
        w.WriteHeader(405)
        return
    }

    meta, err := bj.Meta.Get(val)
    if (err != nil) || meta.Deleted {
        fallback := r.URL.Query().Get("fallback")
        if fallback != "" {
            meta, err = bj.Meta.Get(fallback)
        }
    }
    if (err != nil) || meta.Deleted {
        header := w.Header()
        if err==nil {
            header.Set("Last-Modified", meta.Modified.UTC().Format(http.TimeFormat))
            header.Set("ETag", "deleted")
        } else {
            header.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
            header.Set("ETag", "notfound")
        }
        w.WriteHeader(404)
        io.WriteString(w,"Blob not found\n")
        return
    }

    data, err := bj.Data.Open(meta)
    if err != nil {
        w.WriteHeader(500)
        io.WriteString(w, fmt.Sprintf("Has meta but no data. %#v\n", err))
        return
    }

    header := w.Header()
    header.Set("X-Created", meta.Created.UTC().Format(http.TimeFormat))
    header.Set("Content-Type", meta.MimeType)
    header.Set("ETag", meta.Tag)
    http.ServeContent(w, r, "", meta.Modified, data)
}

func log(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        w.WriteHeader(405)
        return
    }

    sstart := time.Now()
    iterator, err := bj.Meta.GetLogIterator(nil)
    if err != nil {
        w.WriteHeader(500)
        io.WriteString(w,"Failed to get log iterator\n")
        return
    }
    defer iterator.Close()
    fmt.Fprintf(w, "%s\n", sstart.UTC())
    for {
        bi, err := iterator.GetNext()
        if err != nil {
            fmt.Printf("Error while iterating: %#v\n", err)
            break
        }
        if bi == nil {
            break
        }
        fmt.Fprintf(w, "%s|%s\n", bi.Modified.UTC(), bi.Name)
    }
}

type HashingReader struct {
    R io.Reader
    H hash.Hash
}
func (hr *HashingReader) Read(p []byte) (int, error) {
    n, errR := hr.R.Read(p)
    _, errH := hr.H.Write(p[:n])
    //fmt.Printf("Read: %s. %#v %#v",p[:n], errR, errH)
    if errR != nil {
        return n, errR
    }
    return n, errH
}

var olderThanStored = errors.New("New blob is older than stored one")

func (bj *BlobJar) Update(newMeta *BlobInfo, r io.Reader) error {
    bj.Lock(newMeta.Name)
    defer bj.Unlock(newMeta.Name)

    newMeta.ModifiedLocal = time.Now()

    oldMeta, err := bj.Meta.Get(newMeta.Name)
    if err == nil {
        newMeta.Created = oldMeta.Created
    } else {
        newMeta.Created = time.Now()
        oldMeta = nil
    }

    if oldMeta != nil && !oldMeta.Modified.Before(newMeta.Modified) {
        return olderThanStored
    }

    if newMeta.Deleted {
        newMeta.Size = 0
        newMeta.Tag = "deleted"
        bj.Data.Delete(newMeta)
    } else {
        hr := &HashingReader{r, md5.New()}
        nums, err := bj.Data.Save(newMeta, hr)
        if err != nil {
            return err
        }
        //if newMeta.Size == -1 {
        newMeta.Size = nums
        //}
        newMeta.Tag = hex.EncodeToString(hr.H.Sum(nil))
    }
    err = bj.Meta.Set(newMeta)
    if err != nil {
        return err
    }
    return nil
}

func modify(w http.ResponseWriter, r *http.Request, val string, delete bool) {
    //fmt.Printf("FUUUUUUUUUUUUUUUUCK: %#v\n", r)
    if ((r.Method != "PUT") && !delete)||((r.Method != "DELETE") && delete) {
        w.WriteHeader(405)
        return
    }

    var newMeta BlobInfo

    newMeta.Name = val
    newMeta.MimeType = r.Header.Get("Content-Type")
    if newMeta.MimeType=="" {
        newMeta.MimeType = "application/octet-stream"
    }
    var err error
    if newMeta.Modified, err = time.Parse(http.TimeFormat, r.Header.Get("X-Last-Modified")); err != nil {
        newMeta.Modified = time.Now()
    }

    newMeta.Deleted = delete
    err = bj.Update(&newMeta, r.Body)
    fmt.Printf("Update %s -> (%s, %s, %d) %#v\n", val, newMeta.Modified, newMeta.Tag, newMeta.Size, err)
    if err != nil {
        if err == olderThanStored {
            w.WriteHeader(412)
            io.WriteString(w,"Your blob is older than mine\n")
            return
        }
        w.WriteHeader(500)
        fmt.Printf("Error while saving: %#v\n", err)
        io.WriteString(w, "Unable to save data.\n")
        return
    }

    w.WriteHeader(200)
    io.WriteString(w,newMeta.Tag)
}

//func put(ctx *web.Context, val string) {
func put(w http.ResponseWriter, r *http.Request) {
    val := r.URL.Path[len("/put/"):]
    modify(w,r,val,false)
}

//func put(ctx *web.Context, val string) {
func del(w http.ResponseWriter, r *http.Request) {
    val := r.URL.Path[len("/delete/"):]
    modify(w,r,val,true);
}

func main() {
    if len(os.Args) < 3 {
        fmt.Printf("Usage: %s <dirname> <host:port>\n", os.Args[0])
        return
    }
    bj = &BlobJar{}
    bj.Data = &DataStorageFS{os.Args[1]}
    var err error
    bj.Meta, err = NewMetaLDB(os.Args[2])
    if err != nil {
        fmt.Printf("Error: %#v\n");
        return
    }
    bj.mutexes = map[string]*sync.Mutex{}
    //bj.lock = sync.RWMutex{}
    //lol.Lock()
    //lol.Unlock()
    http.HandleFunc("/get/", get)
    http.HandleFunc("/put/", put)
    http.HandleFunc("/delete/", del)
    http.HandleFunc("/log/", log)
    fmt.Printf("%#v\n", http.ListenAndServe(os.Args[3], nil))
    /* web.Get("/get/(.*)", get) //"127.0.0.1:9999"
    web.Put("/put/(.*)", put)
    web.Run("0.0.0.0:9999") */
}
