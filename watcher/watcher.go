package watcher

type IWatcher interface {
	Stop() error
	SetReloadData(f func(file_identifier string, buf []byte) error) error
	LoadFile(file_identifier string) error
	LoadFiles(func(file_identifier string) bool) error
}
