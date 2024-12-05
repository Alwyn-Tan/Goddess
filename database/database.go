package database

import (
	"Goddess/datastruct/dict"
	"Goddess/datastruct/list"
)

const (
	DataDictSize = 1 << 16
)

type Database struct {
	index int
	data  *dict.ConcurrentDict
}

func initDatabase() *Database {
	return nil
}

func (db *Database) Close()                                    {}
func (db *Database) getOrInitList(key string) (list list.List) {}
