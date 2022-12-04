package commands

import (
	"photosorter/utils"
	"sync"
)

type CommandInterface interface {
	getDirectoryForIteration() string
	CanRun(file string) bool
	Run(file string) error
}

//TODO make counter for total and done gorutines
func RunCommand(currentCommand CommandInterface) error {
	var wg sync.WaitGroup

	//TODO check 32, 64 etc
	maxNbConcurrentGoroutines := 16
	concurrentGoroutines := make(chan struct{}, maxNbConcurrentGoroutines)

	fileList, err := utils.DirectoryIterator(currentCommand.getDirectoryForIteration())

	if err != nil {
		return err
	}

	for _, file := range fileList {
		if !currentCommand.CanRun(file) {
			continue
		}

		wg.Add(1)
		go func(file string, currentCommand CommandInterface) {
			defer wg.Done()
			concurrentGoroutines <- struct{}{}

			err := currentCommand.Run(file)

			if err != nil {
				panic(err)
			}

			<-concurrentGoroutines
		}(file, currentCommand)

	}
	wg.Wait()

	return nil
}
