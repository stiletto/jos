import hashlib, bisect

class ConsistentHashRing(object):
    def __init__(self, copies=3, nodes=None, hashfun=hashlib.md5):
        self.hashfun = hashfun
        self.copies = copies
        self._nodes = {}
        self._node_keys = []

    def add_node(self, nodeid, payload=None, priority=128):
        node_key = chr(priority) + self.hashfun(nodeid).digest()[:7]
        if node_key in self._nodes:
            raise ValueError('fuck')
        self._nodes[node_key] = [nodeid, payload]
        bisect.insort(self._node_keys, node_key)

    def get_key_nodes(self, key):
        res = []
        key_hash = self.hashfun(key).digest()
        divr = 2**32 / len(self._node_keys)
        int_key = (ord(key_hash[0]) << 24) + (ord(key_hash[1]) << 16) + (ord(key_hash[2]) << 8) + ord(key_hash[3])
        #int_key = 2**32-1
        for copy in range(self.copies):
            nodenum = (int_key + (copy << 32) / self.copies ) % 2**32
            nodenum = (nodenum / divr) % len(self._node_keys)
            res.append(self._nodes[self._node_keys[nodenum]])
        return res

    def get_nodes(self):
        pass

if __name__=="__main__":
    ring = ConsistentHashRing()
    for x in range(30000):
        ring.add_node('node%d' % x, x)

    for key in range(16):
        key = chr(key+ord('A'))
        print repr(key),':', repr(ring.get_key_nodes(key))

    nodec = {}
    cnt = 3
    for a in range(2**(8*cnt)):
        if a%65536==0:
            print '%08x/%08x' %(a, 2**(8*cnt))
        key = ''
        for x in range(cnt):
            key = chr((a >> (x*8)) % 0x100 ) + key
        #print a, repr(key)
        nodes = ring.get_key_nodes(key)
        for node,payload in nodes:
            nodec[node] = nodec.get(node,0) + 1
    print 'len(nodec) ==',len(nodec)#, repr(nodec)
    minc, maxc = 2**32-1, 0
    sumc = 0
    for k,node in nodec.iteritems():
        minc = min(minc, node)
        maxc = max(maxc, node)
        sumc += node
    print 'minc ==', minc
    print 'maxc ==', maxc
    avgc = float(sumc)/len(nodec)
    print 'avgc ==', avgc
    mdev = 0
    for k,node in nodec.iteritems():
        mdev += (avgc-node)**2
    import math
    mdev = math.sqrt(mdev /len(nodec))
    print 'mdev ==', mdev
