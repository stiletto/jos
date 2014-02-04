package main
import "sync"

//import "fmt"

type MutexMap struct {
    lock sync.RWMutex
    lock2 sync.Mutex
    mutexes map[string]*sync.Mutex
}

func (mm *MutexMap) Lock(name string) {
    //fmt.Printf("Locking %s\n",name)
    mm.lock.Lock()
        lock, ok := mm.mutexes[name]
        if !ok {
            lock = new(sync.Mutex)
            mm.mutexes[name] = lock
        }
        mm.lock2.Lock()
    mm.lock.Unlock()
    mm.lock.RLock()
        mm.lock2.Unlock()
        lock.Lock()
    mm.lock.RUnlock()
}

func (mm *MutexMap) Unlock(name string) {
    //fmt.Printf("Unlocking %s\no",name)
    mm.lock.Lock()
        mm.lock2.Lock()
            lock, ok := mm.mutexes[name]
            //fmt.Printf("lock: %#v, ok: %v\n", lock, ok)
            if ok {
                lock.Unlock()
            }
            //fmt.Printf("Unlocked %s\n", name)
            delete(mm.mutexes, name)
            //fmt.Printf("Deleted %s\n", name)
        mm.lock2.Unlock()
    mm.lock.Unlock()
    //fmt.Printf("Unlocked %s\n", name)
}
