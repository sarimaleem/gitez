package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var startingOptions = [...]string{
	"create branch",
	"switch branch",
	"stage changes",
	"commit changes",
}

func withFilter(command string, input func(in io.WriteCloser)) []string {
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		shell = "sh"
	}
	cmd := exec.Command(shell, "-c", command)
	cmd.Stderr = os.Stderr
	in, _ := cmd.StdinPipe()
	go func() {
		input(in)
		in.Close()
	}()
	result, _ := cmd.Output()
	return strings.Split(string(result), "\n")
}

func fuzzyFind(list []string) string  {
	filtered := withFilter("fzf -m", func(in io.WriteCloser) {
		for _, val := range list {
			fmt.Fprintln(in, val)
		}
		time.Sleep(5 * time.Millisecond)
	})
	output := strings.TrimSpace(strings.Join(filtered, " "))
    return output
}

func getCommandOutput(command string) string{
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		shell = "sh"
	}
	cmd := exec.Command(shell, "-c", command)
	cmd.Stderr = os.Stderr
	result, _ := cmd.Output()
    return string(result)
}

func main() {
	start()
}

func start() {
    output := fuzzyFind(startingOptions[:])

	switch output {
	case "switch branch":
		switchBranch()
	case "create branch":
		createBranch()
	case "stage changes":
		stageChange()
	case "commit changes":
		commitChanges()
	}
}

func switchBranch() {
    r, _ := regexp.Compile("[\\w\\d].+")
    cmdOutput := getCommandOutput("git branch")
    branches := r.FindAllString(cmdOutput, -1);
    targetBranch := fuzzyFind(branches)
    getCommandOutput("git checkout " + targetBranch)
}

func createBranch() {
	fmt.Println("TODO: do the create branch feature")
}

func stageChange() {
	fmt.Println("TODO: do the stage changes feature")
}

func commitChanges() {
	fmt.Println("TODO: do the commit changes feature")
}
