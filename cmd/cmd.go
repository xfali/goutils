// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package cmd

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "os"
    "os/exec"
)

func ExecCommand(command string, args ...string) error {
    cmd := exec.Command(command, args...)
    if cmd == nil {
        return errors.New("exec command nil")
    }

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    cmd.Start()

    err := cmd.Wait()
    if err != nil {
        fmt.Errorf("exec error %s %s\n", err.Error(), command)
    }

    return nil
}

func printPipe(pipe io.Reader, out io.Writer) {
    if pipe == nil {
        return
    }

    buf := bufio.NewReader(pipe)
    for {
        line, err := buf.ReadString('\n')
        fmt.Fprintln(out, line)
        if err != nil {
            if err == io.EOF {
                break
            }
            break
        }
    }
}
