package main

import (
	"fmt"
	"github.com/aaxiba/grocksdb"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/kv/get", kvGet)
	r.GET("/kv/set", kvSet)
	r.Run(":9999") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func kvGet(c *gin.Context) {
	k := c.Query("k")

	cache := grocksdb.NewLRUCache(3 << 10)
	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(cache)

	opts := grocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)

	db, err := grocksdb.OpenDb(opts, "/tmp/rocksdb-test")
	defer db.Close()
	fmt.Println(db, err)
	value, err := db.Get(grocksdb.NewDefaultReadOptions(), []byte(k))
	defer value.Free()
	fmt.Printf("value.Data: %+v, err: %v\n", string(value.Data()), err)
	res := make(map[string]interface{}, 0)
	res["data"] = string(value.Data())
	c.JSON(0, res)
}

func kvSet(c *gin.Context) {
	k := c.Query("k")
	v := c.Query("v")

	cache := grocksdb.NewLRUCache(3 << 10)
	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(cache)

	opts := grocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)

	db, err := grocksdb.OpenDb(opts, "/tmp/rocksdb-test")
	wo := grocksdb.NewDefaultWriteOptions()
	err = db.Put(wo, []byte(k), []byte(v))
	defer db.Close()
	fmt.Println(db, err)
	res := make(map[string]interface{}, 0)
	res["err"] = err
	c.JSON(0, res)
}
