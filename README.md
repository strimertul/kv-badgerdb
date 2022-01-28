# BadgerDB driver for Kilovolt

Simple BadgerDB driver for Kilovolt.

## Usage

Usage is literally a function call to wrap an existing BadgerDB instance in a Kilovolt driver interface and then passing it over.

```go
package example

import (
	"github.com/dgraph-io/badger/v3"
	kv "github.com/strimertul/kilovolt/v7"
	badger_driver "github.com/strimertul/kv-badgerdb"
)

func main() {
	// Initialize your database 
	options := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(options)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create driver using database instance
	driver := badger_driver.NewBadgerBackend(db)

	// Pass it to Kilovolt
	hub, err := kv.NewHub(driver, kv.HubOptions{}, nil)
	if err != nil {
		panic(err)
	}
	go hub.Run()

	// etc.
}
```