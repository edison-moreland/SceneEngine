package main

import (
	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	watcher *fsnotify.Watcher
	changed bool
}

func WatchFile(file string) (*FileWatcher, error) {
	fileWatcher := FileWatcher{
		changed: false,
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	fileWatcher.watcher = watcher

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					fileWatcher.changed = true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				panic(err)
			}
		}
	}()

	err = watcher.Add(file)
	if err != nil {
		return nil, err
	}

	return &fileWatcher, nil
}

func (f *FileWatcher) HasChanged() bool {
	return f.changed
}

func (f *FileWatcher) ClearChange() {
	f.changed = false
}

func (f *FileWatcher) Close() error {
	return f.watcher.Close()
}
