package main

import "github.com/serialx/hashring"

type hashKey string

func (hk hashKey) Less(other hashring.HashKey) bool {
	return hk < other.(hashKey)
}
