
package main

import (
    "os"
    "bufio"
    //"os/exec"
    "log"
    "fmt"
    "regexp"
    //"io/ioutil"
    "flag"
    "path/filepath"
    "time"
)

var httpport int
var workpath string

type rec struct {
    pname string
}

var src []rec
var articlen int

var pos int

func main () {

    //if os.Getppid() != 1 {
    //    //filePath,_ := filepath.Abs(os.Args[0])
    //    cmd := exec.Command(os.Args[0], os.Args[1:]...)
    //    cmd.Stdin = os.Stdin
    //    cmd.Stdout = os.Stdout
    //    cmd.Stderr = os.Stderr
    //    log.Println("before Start")
    //    err := cmd.Start()
    //    if err != nil {
    //        log.Println("Start err ", err)
    //        os.Exit(0)
    //    }
    //    log.Println("after Start")
    //    time.Sleep(10*time.Second)
    //    return
    //}

    flag.StringVar(&workpath, "workp", "./", "default work path is current dir")
    flag.IntVar(&httpport, "port", 6969, "default http port is 6969")
    flag.Parse()

    log.Println("after Parse")
    fmt.Println("port is ", httpport)
    fmt.Println("work path is ", workpath, "\n")
    time.Sleep(500*time.Millisecond)

    articlen = 0
    filepath.Walk(workpath, checkname)

    fmt.Println("src's len is", len(src), "before make")
    src = make([]rec, articlen+1)
    pos = 0
    filepath.Walk(workpath, setname)

    fmt.Println("src's cnt", pos, len(src))
    fmt.Println("range src")
    for _, v := range src {
        fmt.Println(v.pname)
        catfile(workpath, v.pname)
    }

    return
}

func catfile(path string, name string) {
    filename := path+name
    fmt.Println(filename+"'s contexts:")

    fd, err := os.Open(filename)
    if(err != nil) {
        fmt.Println(err)
        os.Exit(1)
    }

    defer fd.Close()

    f := bufio.NewReader(fd)
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }

    if err:= scanner.Err(); err != nil {
        fmt.Println(err)
    }

}

func checkname(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    fname := info.Name()
    if info.IsDir() {
        fmt.Println("Just a Dir", fname)
        if path != workpath {
            return filepath.SkipDir
        }
    } else {
        matched, err := regexp.MatchString("r20[0-9]{12}z", fname)
        if err != nil {
            fmt.Println("Match err", err)
            return err
        }
        if matched == true {
            fmt.Println("YES, ", fname)
            articlen++
        } else {
            fmt.Println("NO , ", fname)
        }
    }
    return nil
}

func setname(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    fname := info.Name()
    if info.IsDir() {
        if path != workpath {
            return filepath.SkipDir
        }
    } else {
        matched, err := regexp.MatchString("r20[0-9]{12}z", fname)
        if err != nil {
            fmt.Println("Match err", err)
            return err
        }
        if matched == true {
            src[pos].pname = fname
            pos++
            fmt.Println("ADD ONE")
        }
    }
    return nil
}




