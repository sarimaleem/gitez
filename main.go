package main

import (
    "fmt"
    "io"
    "os"
    "os/exec"
    "strings"
    "time"
)

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

func main() {
    filtered := withFilter("fzf -m", func(in io.WriteCloser) {
        fmt.Fprintln(in, "switch branch");
        fmt.Fprintln(in, "commit code");
        fmt.Fprintln(in, "commit everything");
        time.Sleep(5* time.Millisecond)
        // for i := 0; i < 1000; i++ {
        //     fmt.Fprintln(in, i)
        //     time.Sleep(5 * time.Millisecond)
        // }
    })
    test := strings.Join(filtered, " ");
    fmt.Println(test);
    // fmt.Println(filtered + "what's up")
}
