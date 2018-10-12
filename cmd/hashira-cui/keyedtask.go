package main

import "github.com/pankona/hashira/service"

// KeyedTask is service.Task with Key() function.
type KeyedTask service.Task

// Key returns task's ID.
// This function is for satisfying orderedmap.Keyer interface.
func (kt *KeyedTask) Key() string {
	return kt.Id
}
