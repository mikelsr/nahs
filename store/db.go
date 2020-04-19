package store

import (
	"bytes"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/mikelsr/bspl"
	"github.com/syndtr/goleveldb/leveldb"
)

// DB is a wrapper for leveldb.DB
type DB struct {
	path string
	db   *leveldb.DB
	snap *leveldb.Snapshot
}

// NewDB is the default constructor for DB
func NewDB(path string) (*DB, error) {
	d := DB{path: path}
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	d.db = db
	d.snap, err = db.GetSnapshot()
	return &d, err
}

// Put a Value in the DB
func (d *DB) Put(k Key, v Value) {
	d.db.Put(k.dump(), v.dump(), nil)
}

// Get a Value from the DB Snapshot
func (d *DB) Get(k peer.ID) Value {
	v, err := d.snap.Get(Key(k).dump(), nil)
	if err != nil {
		return nil
	}
	return loadValue(v)
}

// Has checks if a key exists in the DB Snapshot
func (d *DB) Has(k peer.ID) bool {
	has, err := d.snap.Has(Key(k).dump(), nil)
	return err == nil && has
}

// Update the DB snapshot
func (d *DB) Update() error {
	d.snap.Release()
	var err error
	d.snap, err = d.db.GetSnapshot()
	return err
}

// Open the DB file
func (d *DB) Open() {
	var err error
	d.db, err = leveldb.OpenFile(d.path, nil)
	if err != nil {
		panic(err)
	}
}

// Close the DB file
func (d *DB) Close() {
	d.snap.Release()
	d.db.Close()
}

// Key is a peer.ID mapped to a list of protocols
type Key peer.ID

// dump a Key to bytes
func (k Key) dump() []byte {
	return []byte(k)
}

// Value of the Peer info database
type Value []bspl.Protocol

// dump Value to bytes
func (v Value) dump() []byte {
	buff := bytes.NewBuffer(nil)
	n := len(v)
	for i, p := range v {
		b := []byte(p.String())
		_, err := buff.Write(b)
		if err != nil {
			panic(err)
		}
		if i < n-1 {
			buff.Write(separator)
		}
	}
	return buff.Bytes()
}

// loadValue from bytes
func loadValue(b []byte) Value {
	bProtos := bytes.Split(b, separator)
	v := make(Value, len(bProtos))
	for i, bp := range bProtos {
		reader := bytes.NewReader(bp)
		protocol, err := bspl.Parse(reader)
		if err != nil {
			panic(err)
		}
		v[i] = protocol
	}
	return v
}
