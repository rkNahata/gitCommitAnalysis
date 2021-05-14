package main

import (
	"flag"
	"fmt"
)

func main() {
	var folder, email string
	flag.StringVar(&folder, "add", "", "add a new folder to scan for git repos")
	flag.StringVar(&email, "email", "your@email.com", "enter ur git account email")
	flag.Parse()

	if folder != "" {
		scanner(folder)
		return
	}
	statistics(email)

}

func scanner(folder string) {
	fmt.Println("Found Folders")

	repositories := RecursiveScanFolders(folder)
	filePath := getDotFile()
	addNewSliceElementsFile(filePath, repositories)
	fmt.Println("\n Successfully added files")
}

func statistics(email string) {
	fmt.Println(email)
	commits := processRepositories(email)
	printCommitStats(commits)
}
