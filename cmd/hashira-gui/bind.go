// -build linux,amd64

package main

// Asset is just a mock function
func Asset(name string) ([]byte, error) {
	return nil, nil
}

// RestoreAssets is just mock function
func RestoreAssets(dir, name string) error {
	return nil
}
