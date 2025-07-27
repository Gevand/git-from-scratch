package commands

import (
	"fmt"
	"geo-git/lib"

	repostatus "geo-git/lib/repository"
	"geo-git/lib/utils"
	"slices"
)

var ColorMap = map[string]string{
	"green": "\033[32m",
	"red":   "\033[31m",
	"reset": "\033[0m",
}

var ShortStatusMap = map[repostatus.Status]string{
	repostatus.Deleted:  "D",
	repostatus.Modified: "M",
	repostatus.Added:    "A",
}
var LongStatusMap = map[repostatus.Status]string{
	repostatus.Deleted:  "deleted",
	repostatus.Modified: "modified",
	repostatus.Added:    "new file",
}

func RunStatus(repo *lib.Respository, cmd *Command) error {

	statusTracking := repostatus.NewStatusTracking()
	err := statusTracking.GenerateStatus(repo)
	if err != nil {
		return err
	}

	if utils.Contains(cmd.Args, "--porcelain") {
		printResultsPorcelain(statusTracking)
	} else {
		printResultsLongs(statusTracking)
	}
	return nil
}

func printResultsLongs(statusTracking *repostatus.RepositoryStatusTracking) {
	printChanges("Changes to be committed", statusTracking.IndexChanges, "green")
	printChanges("Changes not staged for commit", statusTracking.WorkSpaceChanges, "red")
	printUntrackedChanges("Untracked files", statusTracking.Untracked, "red")
	printCommitStatus(statusTracking)
}

func printCommitStatus(statusTracking *repostatus.RepositoryStatusTracking) {
	if len(statusTracking.IndexChanges) > 0 {
		return
	}

	if len(statusTracking.WorkSpaceChanges) > 0 {
		fmt.Println("no changes added to commit")
	} else if len(statusTracking.Untracked) > 0 {
		fmt.Println("nothing added to commit but untracked files present")
	} else {
		fmt.Println("nothing to commit, working tree clean")
	}
}

func printChanges(message string, changes map[string]repostatus.Status, color string) {
	if len(changes) == 0 {
		return
	}
	fmt.Println(message)
	fmt.Println("")
	color_code, ok := ColorMap[color]
	reset := ColorMap["reset"]
	for path, status := range changes {
		if ok {
			fmt.Printf("\t\t%s%s %s\n%s", color_code, LongStatusMap[status], path, reset)
		} else {
			fmt.Printf("%s %s\n", LongStatusMap[status], path)
		}
	}
	fmt.Println("")
}
func printUntrackedChanges(message string, changes []string, color string) {
	if len(changes) == 0 {
		return
	}
	fmt.Println(message)
	fmt.Println("")
	color_code, ok := ColorMap[color]
	reset := ColorMap["reset"]
	for _, path := range changes {
		if ok {
			fmt.Printf("\t\t%s%s\n%s", color_code, path, reset)
		} else {
			fmt.Printf("\t\t%s\n", path)
		}
	}
	fmt.Println("")
}
func printResultsPorcelain(statusTracking *repostatus.RepositoryStatusTracking) {
	for _, path := range statusTracking.Changed {
		status := statusForPath(path, statusTracking)
		fmt.Printf("%s %s\r\n", status, path)
	}
	for _, file := range slices.Compact(statusTracking.Untracked) {
		fmt.Printf("?? %s\r\n", file)
	}
}

func statusForPath(path string, statusTracking *repostatus.RepositoryStatusTracking) string {
	var left string
	var right string
	left_status, ok := statusTracking.IndexChanges[path]
	if !ok {
		left = ""
	} else {
		left = ShortStatusMap[left_status]
	}

	right_status, ok := statusTracking.WorkSpaceChanges[path]
	if !ok {
		right = ""
	} else {
		right = ShortStatusMap[right_status]
	}

	return left + right
}
