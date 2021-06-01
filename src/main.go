package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
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
    "view commits",
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

func fuzzyFind(list []string, options ...string) string  {
    command := "fzf"
    for _, val := range(options) {
        command += " " + val
    }
	filtered := withFilter(command, func(in io.WriteCloser) {
		for _, val := range list {
			fmt.Fprintln(in, val)
		}
		time.Sleep(5 * time.Millisecond)
	})
	output := strings.TrimSpace(strings.Join(filtered, " "))
    return output
}

func getCommandOutput(command string) (string, string){
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		shell = "sh"
	}
	cmd := exec.Command(shell, "-c", command)
    var outb, errb bytes.Buffer
    cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
    return outb.String(), errb.String()
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
	case "view commits":
		viewCommits()
	}
    
}

func switchBranch() {
    r, _ := regexp.Compile("[\\w\\d].+")
    res, _ := getCommandOutput("git branch")
    branches := r.FindAllString(res, -1);
    targetBranch := fuzzyFind(branches)
    getCommandOutput("git checkout " + targetBranch)
}

func createBranch() {
	fmt.Println("TODO: do the create branch feature")
    fmt.Println("Name of the branch: ")
    var branch string
    fmt.Scanln(&branch)
    // figure out what happens when an error is logged
    getCommandOutput("git checkout -b " + branch)
}

func stageChange() {
	fmt.Println("TODO: do the stage changes feature")
}

func commitChanges() {
	fmt.Println("TODO: do the commit changes feature")
}

func viewCommits()  {
    res, _ := getCommandOutput("git log --format=\"%h author: %an (%ar)\"")
    fzfInput := strings.Split(res, "\n")
    fuzzyFind(fzfInput, "--preview=\"git show {1}\"")
}
