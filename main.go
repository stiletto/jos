package main;

import (
//    "github.com/hoisie/web"
    "io"
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

var bj BlobJar

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
    header.Set("Content-Type", meta.MimeType)
    header.Set("ETag", meta.Tag)
    http.ServeContent(w, r, "", meta.Modified, data)
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

func modify(w http.ResponseWriter, r *http.Request, val string, delete bool) {
    fmt.Printf("FUUUUUUUUUUUUUUUUCK: %#v\n", r)
    if ((r.Method != "PUT") && !delete)||((r.Method != "DELETE") && delete) {
        w.WriteHeader(405)
        return
    }

    bj.Lock(val)
    defer bj.Unlock(val)

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

    oldMeta, err := bj.Meta.Get(val)
    if err == nil {
        newMeta.Created = oldMeta.Created
    } else {
        newMeta.Created = time.Now()
    }

    if oldMeta != nil && !oldMeta.Modified.Before(newMeta.Modified) {
        w.WriteHeader(412)
        io.WriteString(w,"Your blob is older than mine\n")
        return
    }

    if delete {
        newMeta.Size = 0
        newMeta.Deleted = true
        newMeta.Tag = "deleted"
        bj.Data.Delete(&newMeta)
    } else {
        newMeta.Deleted = false
        newMeta.Size, err = r.ContentLength, nil //strconv.ParseInt(ctx.Request.Header.Get("Content-Length"), 10, 64)
        if err != nil {
            newMeta.Size = -1
        }

        hr := &HashingReader{r.Body, md5.New()}
        nums, err := bj.Data.Save(&newMeta, hr)
        if err != nil {
            w.WriteHeader(500)
            io.WriteString(w, fmt.Sprintf("Unable to save data. %#v\n",err))
            return
        }
        if newMeta.Size == -1 {
            newMeta.Size = nums
        }
        newMeta.Tag = hex.EncodeToString(hr.H.Sum(nil))
    }
    err = bj.Meta.Set(&newMeta)
    if err != nil {
        w.WriteHeader(500)
        io.WriteString(w, fmt.Sprintf("Unable to save meta. %#v\n",err))
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
    bj.Data = &DataStorageFS{os.Args[1]}
    bj.Meta = &MetaStorageFS{os.Args[1]}
    bj.mutexes = map[string]*sync.Mutex{}
    //bj.lock = sync.RWMutex{}
    //lol.Lock()
    //lol.Unlock()
    http.HandleFunc("/get/", get)
    http.HandleFunc("/put/", put)
    http.HandleFunc("/delete/", del)
    http.ListenAndServe(os.Args[2], nil)
    /* web.Get("/get/(.*)", get) //"127.0.0.1:9999"
    web.Put("/put/(.*)", put)
    web.Run("0.0.0.0:9999") */
}
