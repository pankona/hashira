package main

import "github.com/pankona/hashira/service"

type KeyedTask service.Task

func (kt *KeyedTask) Key() string {
	return kt.Id
}
