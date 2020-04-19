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

// Put a dbvalue in the DB
func (d *DB) Put(k peer.ID, v []bspl.Protocol) {
	d.db.Put(dbkey(k).dump(), dbvalue(v).dump(), nil)
}

// Get a dbvalue from the DB Snapshot
func (d *DB) Get(k peer.ID) []bspl.Protocol {
	v, err := d.snap.Get(dbkey(k).dump(), nil)
	if err != nil {
		return nil
	}
	return loadDBvalue(v)
}

// Has checks if a key exists in the DB Snapshot
func (d *DB) Has(k peer.ID) bool {
	has, err := d.snap.Has(dbkey(k).dump(), nil)
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

// dbkey is a peer.ID mapped to a list of protocols
type dbkey peer.ID

// dump a dbkey to bytes
func (k dbkey) dump() []byte {
	return []byte(k)
}

// dbvalue of the Peer info database
type dbvalue []bspl.Protocol

// dump dbvalue to bytes
func (v dbvalue) dump() []byte {
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

// loaddbvalue from bytes
func loadDBvalue(b []byte) dbvalue {
	bProtos := bytes.Split(b, separator)
	v := make(dbvalue, len(bProtos))
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
