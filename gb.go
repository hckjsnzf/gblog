
package main

import (
    "os"
    "os/exec"
    "log"
    "fmt"
    //"io/ioutil"
    "flag"
    //"path/filepath"
    "time"
)

var httpport int
var workpath string

func main () {

    if os.Getppid() != 1 {
        //filePath,_ := filepath.Abs(os.Args[0])
        cmd := exec.Command(os.Args[0], os.Args[1:]...)
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        log.Println("before Start")
        err := cmd.Start()
        if err != nil {
            log.Println("Start err ", err)
            os.Exit(0)
        }
        log.Println("after Start")
        time.Sleep(10*time.Second)
        return
    }

    flag.StringVar(&workpath, "workp", "./", "default work path is current dir")
    flag.IntVar(&httpport, "port", 6969, "default http port is 6969")
    flag.Parse()

    log.Println("after Parse")
    fmt.Println("port is ", httpport)
    fmt.Println("work path is ", workpath)
    time.Sleep(10*time.Second)

    return
}

