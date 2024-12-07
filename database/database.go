package database

import (
	"Goddess/datatype/hash"
	"Goddess/datatype/list"
)

const (
	DataDictSize = 1 << 16
)

type Database struct {
	index int
	data  *hash.Dict
}

func initDatabase() *Database {
	return nil
}

func (db *Database) Close()                                    {}
func (db *Database) getOrInitList(key string) (list list.List) {}
