package main

import "github.com/pankona/hashira/service"

type keyedTask service.Task

func (kt *keyedTask) Key() string {
	return kt.Id
}
