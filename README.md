jos
===

This is JarOfShit. Sometimes consistent, wanna be distributed, not yet fault-tolerant and replicated BLOB storage.

jos is used as avatars storage for [BnW](http://bnw.im).

# How to use:

    $ export GOPATH=$HOME/.go`
    $ go get github.com/stiletto/jos
    $ mkdir datadir
    $ $HOME/.go/bin/jos ./datadir ./datadir/meta 127.0.0.1:9999

# jos HTTP API

Storing object:

    $ curl -i --data-binary 'Hello, world' -H 'Content-Type: text/plain' -X PUT http://127.0.0.1:9999/put/hello.txt
    HTTP/1.1 200 OK
    Date: Sat, 09 Aug 2014 20:16:49 GMT
    Content-Length: 32
    Content-Type: text/plain; charset=utf-8
    
    bc6e6f16b8a077ef5fbc8d59d0b931b9

Retreiving object:

    $ curl -i http://127.0.0.1:9999/get/hello.txt 
    HTTP/1.1 200 OK                                       
    Accept-Ranges: bytes
    Content-Length: 12
    Content-Type: text/plain
    Etag: bc6e6f16b8a077ef5fbc8d59d0b931b9
    Last-Modified: Sat, 09 Aug 2014 20:16:49 GMT
    X-Created: Sat, 09 Aug 2014 20:16:49 GMT
    Date: Sat, 09 Aug 2014 20:17:53 GMT
    
    Hello, world

All your fancy Etag stuff is supported of course:

    $ curl -i -H 'If-None-Match: bc6e6f16b8a077ef5fbc8d59d0b931b9' http://127.0.0.1:9999/get/hello.txt 
    HTTP/1.1 304 Not Modified                                                         
    Etag: bc6e6f16b8a077ef5fbc8d59d0b931b9
    Last-Modified: Sat, 09 Aug 2014 20:16:49 GMT
    X-Created: Sat, 09 Aug 2014 20:16:49 GMT
    Date: Sat, 09 Aug 2014 20:18:57 GMT
    

Deleting object:

    $ curl -i -X DELETE http://127.0.0.1:9999/delete/hello.txt                                  
    HTTP/1.1 200 OK
    Date: Sat, 09 Aug 2014 20:21:01 GMT
    Content-Length: 7
    Content-Type: text/plain; charset=utf-8
    
    deleted
