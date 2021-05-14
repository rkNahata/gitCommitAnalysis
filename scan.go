package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

func getDotFile() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	dotFile := usr.HomeDir + "/gitlocalstats"
	return dotFile
}

func RecursiveScanFolders(folder string) []string {
	fmt.Println(folder)
	return scanGitFolders(make([]string, 0), folder)
}

func scanGitFolders(folders []string, folder string) []string {
	folder = strings.TrimSuffix(folder, "/")
	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatalln(err)
	}

	var path string

	for _, file := range files {
		if file.IsDir() {
			path = folder + fmt.Sprintf("/%s", file.Name())
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" {
				continue
			}
			folders = scanGitFolders(folders, path)
		}

	}
	return folders
}

func addNewSliceElementsFile(filePath string, newRepos []string) {
	existingRepos := parseFileLinesToSlice(filePath)
	repos := joinSlices(existingRepos, newRepos)
	dumpStringSliceToFile(filePath, repos)
}

func parseFileLinesToSlice(filePath string) []string {
	f := openFile(filePath)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}
	return lines

}

func openFile(filePath string) *os.File {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	return f
}

func joinSlices(existing, new []string) []string {
	for _, i := range new {
		if !sliceContains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

func sliceContains(existing []string, new string) bool {
	for _, val := range existing {
		if val == new {
			return true
		}
	}
	return false
}

func dumpStringSliceToFile(filePath string, repos []string) {
	text := strings.Join(repos, "\n")
	ioutil.WriteFile(filePath, []byte(text), 0755)
}
