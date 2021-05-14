package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"sort"
	"time"
)

const (
	daysInLastSixMonths  = 183
	outOfRange           = 9999
	weeksInLastSixMonths = 26
)

type column []int

func processRepositories(email string) map[int]int {
	filePath := getDotFile()
	repos := parseFileLinesToSlice(filePath)
	daysInMap := daysInLastSixMonths
	commits := make(map[int]int, daysInMap)

	for i := daysInMap; i > 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = fillCommits(email, path, commits)
	}
	return commits
}

func fillCommits(email, path string, commits map[int]int) map[int]int {
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}
	ref, err := repo.Head()
	if err != nil {
		//fmt.Println(err)
		return commits

	}

	iterator, err := repo.Log(&git.LogOptions{
		From: ref.Hash(),
	})
	if err != nil {
		panic(err)
	}

	offset := calcOffset()

	err = iterator.ForEach(func(commit *object.Commit) error {
		daysAgo := countDaysSinceDate(commit.Author.When) + offset

		if commit.Author.Email != email {
			return nil
		}
		if daysAgo != outOfRange {
			commits[daysAgo]++
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return commits
}

func countDaysSinceDate(date time.Time) int {
	days := 0
	now := getBeginningOfTheDay(time.Now())
	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++
		if days > outOfRange {
			return outOfRange
		}
	}
	return days
}

func getBeginningOfTheDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

func calcOffset() int {
	var offset int
	weekDay := time.Now().Weekday()

	switch weekDay {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1

	}
	return offset
}

func printCommitStats(commits map[int]int) {
	keys := sortMapIntoSlice(commits)
	cols := buildColumns(keys, commits)
	printCells(cols)
}

func sortMapIntoSlice(commits map[int]int) []int {
	var keys []int
	for k, _ := range commits {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func buildColumns(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)
	col := column{}

	for _, k := range keys {
		week := k / 7
		dayOfWeek := k % 7
		if dayOfWeek == 0 {
			col = column{}
		}
		col = append(col, commits[k])
		if dayOfWeek == 6 {
			cols[week] = col
		}
	}
	return cols
}

func printMonths() {
	week := getBeginningOfTheDay(time.Now().Add(-(daysInLastSixMonths * time.Hour * 24)))
	month := week.Month()
	fmt.Printf("         ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}
		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

func printDayCol(day int) {
	out := "     "
	switch day {
	case 0:
		out = "	Sun "
	case 1:
		out = "	Mon "
	case 2:
		out = "	Tue "
	case 3:
		out = "	Wed "
	case 4:
		out = "	Thu "
	case 5:
		out = "	Fri "
	case 6:
		out = "	Sat "
	}
	fmt.Printf(out)
}

func printCell(val int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 3:
		escape = "\033[1;30;43m"
	case val >= 3 && val < 6:
		escape = "\033[1;30;46m"
	case val >= 6:
		escape = "\033[1;30;42m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

func printCells(cols map[int]column) {
	printMonths()
	for j := 6; j >= 0; j-- {
		for i := weeksInLastSixMonths + 1; i >= 0; i-- {
			if i == weeksInLastSixMonths+1 {
				printDayCol(j)
			}
			if col, ok := cols[i]; ok {
				if i == 0 && j == calcOffset()-1 {
					printCell(col[j], true)
					continue
				} else {
					if len(col) > j {
						printCell(col[j], false)
						continue
					}
				}
			}
			printCell(0, false)
		}
		fmt.Printf("\n")
	}

}