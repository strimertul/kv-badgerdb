package badger_driver

import (
	"github.com/dgraph-io/badger/v3"
	kv "github.com/strimertul/kilovolt/v8"
)

type Driver struct {
	db *badger.DB
}

func NewBadgerBackend(db *badger.DB) Driver {
	return Driver{db}
}

func (b Driver) Get(key string) (string, error) {
	var out string
	err := b.db.View(func(tx *badger.Txn) error {
		val, err := tx.Get([]byte(key))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return kv.ErrorKeyNotFound
			}
			return err
		}
		byt, err := val.ValueCopy(nil)
		if err != nil {
			return err
		}
		out = string(byt)
		return nil
	})
	return out, err
}

func (b Driver) GetBulk(keys []string) (map[string]string, error) {
	out := make(map[string]string)
	err := b.db.View(func(tx *badger.Txn) error {
		for _, key := range keys {
			val, err := tx.Get([]byte(key))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					out[key] = ""
					continue
				}
				return err
			}
			byt, err := val.ValueCopy(nil)
			if err != nil {
				return err
			}
			out[key] = string(byt)
		}
		return nil
	})
	return out, err
}

func (b Driver) GetPrefix(prefix string) (map[string]string, error) {
	out := make(map[string]string)
	err := b.db.View(func(tx *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = []byte(prefix)
		it := tx.NewIterator(opt)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			byt, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			out[string(item.Key())] = string(byt)
		}
		return nil
	})
	return out, err
}

func (b Driver) Set(key, value string) error {
	return b.db.Update(func(tx *badger.Txn) error {
		return tx.Set([]byte(key), []byte(value))
	})
}

func (b Driver) SetBulk(kv map[string]string) error {
	return b.db.Update(func(tx *badger.Txn) error {
		for k, v := range kv {
			err := tx.Set([]byte(k), []byte(v))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (b Driver) Delete(key string) error {
	return b.db.Update(func(tx *badger.Txn) error {
		return tx.Delete([]byte(key))
	})
}

func (b Driver) List(prefix string) ([]string, error) {
	out := []string{}
	err := b.db.View(func(tx *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = []byte(prefix)
		it := tx.NewIterator(opt)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			out = append(out, string(item.Key()))
		}
		return nil
	})
	return out, err
}
