#!/usr/bin/env python2
import os,requests,datetime,json
import dateutil.parser
from wsgiref.handlers import format_date_time
from time import mktime
from os import path

def parsedate(d):
    #return datetime.datetime.strptime(d, "%Y-%m-%dT%H:%M:%S.%f%zO")
    return dateutil.parser.parse(d)
if __name__=="__main__":
    import sys
    datadir = sys.argv[1]
    baseurl = sys.argv[2]
    for fname in os.listdir(datadir):
        if not fname.endswith('.meta'):
            continue
        dfname = fname.split('.')[0]
        meta = json.load(file(path.join(datadir,fname)))
        print meta['Name']
        if meta['Deleted']:
            print meta['Name'],'is deleted, skipping'
            continue
        data = file(path.join(datadir,dfname))
        mod = parsedate(meta['Modified'])
        r = requests.put(baseurl+'put/'+meta['Name'],
            data=data, headers={
                'Content-Type': meta['MimeType'],
                'X-Last-Modified': format_date_time(mktime(mod.timetuple()))
            })
        print r, r.text
