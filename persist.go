package dag

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"log"
)

const (
	KeyFormatterVertex   = "vertex:%d"
	KeyFormatterParents  = "parents:%s:%d"
	KeyFormatterChildren = "children:%s:%d"
)

func (dag *DAG) Store(path string) {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	txn := db.NewTransaction(true)

	vertexIndex := 0
	dag.Vertexes.Range(func(key, value interface{}) bool {
		v, err := json.Marshal(value)
		if err != nil {
			log.Println(err.Error())
			return false
		}
		log.Println(string(v))
		err = txn.Set([]byte(fmt.Sprintf(KeyFormatterVertex, vertexIndex)), v)

		switch err {
		case badger.ErrTxnTooBig:
			_ = txn.Commit()
			txn = db.NewTransaction(true)
			_ = txn.Set([]byte(fmt.Sprintf(KeyFormatterVertex, vertexIndex)), v)
		case nil:
			break
		default:
			return false
		}

		// 存储parents
		vert := value.(*Vertex)
		flag := 0
		if vert.Parents.Len() > 0 {
			for e := vert.Parents.Front(); e != nil; e = e.Next() {
				err := txn.Set(
					[]byte(fmt.Sprintf(KeyFormatterParents, vert.Hash, flag)),
					[]byte(e.Value.(*Vertex).Hash),
				)
				switch err {
				case badger.ErrTxnTooBig:
					_ = txn.Commit()
					txn = db.NewTransaction(true)
					_ = txn.Set([]byte(fmt.Sprintf(KeyFormatterVertex, vertexIndex)), v)
				case nil:

				default:
					return false
				}

				flag++
			}
		}

		// 存储children
		flag = 0
		if vert.Children.Len() > 0 {
			for e := vert.Children.Front(); e != nil; e = e.Next() {
				err := txn.Set(
					[]byte(fmt.Sprintf(KeyFormatterChildren, vert.Hash, flag)),
					[]byte(e.Value.(*Vertex).Hash),
				)

				switch err {
				case badger.ErrTxnTooBig:
					_ = txn.Commit()
					txn = db.NewTransaction(true)
					_ = txn.Set([]byte(fmt.Sprintf(KeyFormatterVertex, vertexIndex)), v)
				case nil:
					break
				default:
					return false
				}
			}

			flag++
		}

		vertexIndex++
		return true
	})

	_ = txn.Commit()

}

func (dag *DAG) Load(path string) {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dag.badgerDB = db

	for index := 0; ; index++ {
		if !dag.loadVertex(fmt.Sprintf(KeyFormatterVertex, index)) {
			break
		}
	}

	dag.Vertexes.Range(func(key, value interface{}) bool {

		dag.loadChildren(key.(string))

		return true
	})

}

func (dag *DAG) loadVertex(key string) bool {
	log.Println(key)
	db := dag.badgerDB
	err := db.View(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			log.Println(string(val))
			v := &Vertex{}
			err := json.Unmarshal(val, v)
			if err != nil {
				return err
			}
			dag.Vertexes.Store(v.Hash, v)
			dag.Length++

			return nil
		})
	})

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (dag *DAG) loadChildren(hash string) bool {

	db := dag.badgerDB
	err := db.View(func(txn *badger.Txn) error {
		for flag := 0; ; flag++ {
			item, err := txn.Get([]byte(fmt.Sprintf(KeyFormatterChildren, hash, flag)))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return nil
				}

				return err
			}

			toHash := ""
			err = item.Value(func(val []byte) error {
				toHash = string(val)
				return nil
			})
			if err != nil {
				return err
			}

			dag.AddEdge(
				NewVertex("TX", hash, "transaction"), // {Hash: fmt.Sprintf("%d", from), Type: "TX", Value: "transaction"},
				NewVertex("TX", toHash, "transaction"),
			)

		}

	})

	if err != nil {
		log.Println(err)
		return false
	}

	return true

}
