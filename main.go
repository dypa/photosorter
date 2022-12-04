package main

import (
	"fmt"
	"os"
	"photosorter/commands"
	"photosorter/utils"
	"time"
)

// TODO add verbose
func main() {
	benchmark := time.Now()

	if len(os.Args) < 2 {
		println("Enter command")
		printUsage()
	}

	switch os.Args[1] {
	case "check":
		if len(os.Args) < 3 {
			printUsage()
		}

		archivePath := os.Args[2]

		if !utils.IsDirOrFileExists(archivePath) {
			println("Archive path:", archivePath, "not exists")
			printUsage()
		}

		commands.RunCommand(commands.CheckCommand{
			ArchivePath: archivePath,
		})

	case "copy":
		if len(os.Args) < 4 {
			printUsage()
		}

		archivePath := os.Args[3]
		sourcePath := os.Args[2]

		if archivePath == sourcePath {
			println("Archive path must be different with sourcePath")
			printUsage()
		}

		if !utils.IsDirOrFileExists(sourcePath) {
			println("Source path:", sourcePath, " not exists")
			printUsage()
		}

		if !utils.IsDirOrFileExists(archivePath) {
			println("Archive path:", archivePath, "not exists")
			printUsage()
		}

		commands.RunCommand(commands.CopyCommand{
			SourcePath:  sourcePath,
			ArchivePath: archivePath,
		})

	default:
		printUsage()
	}

	fmt.Println("[INFO] Benchmark:", time.Since(benchmark))
}

func printUsage() {
	println("Usage:")
	println("./photosort copy source-path archive-path")
	println("./photosort check archive-path")
	os.Exit(1)
}
