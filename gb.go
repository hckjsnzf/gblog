
package main

import (
    "os"
    "bufio"
    "bytes"
    //"os/exec"
    "log"
    "net/http"
    "fmt"
    "regexp"
    //"io/ioutil"
    "flag"
    "path/filepath"
    "time"
    "strconv"
)

var httpport int
var workpath string

type rec struct {
    pname string
}
var src []rec
var articlen int

type bcont struct {
    name string
    title string
    context bytes.Buffer
}

var slblog []bcont



var pos int

var bmap map[string] bytes.Buffer


func httprmote(w http.ResponseWriter, r *http.Request) {
    path := r.URL.String()
    fmt.Println("from client: ", "map len", len(bmap))
    fmt.Println("--",path,"--")
    b, ok := bmap[path]
    if ok {
        b.WriteTo(w)
    } else {
        fmt.Fprintf(w, "zz NOT FOUND")
    }
}




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

    src = make([]rec, articlen)
    slblog = make([]bcont, articlen)

    pos = 0
    filepath.Walk(workpath, setname)

    fmt.Println("range src")
    i := 0
    for _, v := range src {
        fmt.Println(v.pname)
        catfile(workpath, v.pname, &(slblog[i]))
        i++
    }

    bmap = make(map[string] bytes.Buffer, articlen+1)
    for _, v := range slblog {
        bmap[v.name] =  v.context
        fmt.Println(v.name)
        fmt.Println("-------")
        http.HandleFunc(v.name, httprmote)
    }
    fmt.Println("bmap len", len(bmap))

    port := strconv.Itoa(httpport)
    port = ":"+port
    err := http.ListenAndServe(port, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }

    return
}


func catfile(path string, name string, thisb *bcont) {
    filename := path+name
    fmt.Println(filename+"'s context:")

    thisb.name = "/"+name


    fd, err := os.Open(filename)
    if(err != nil) {
        fmt.Println(err)
        os.Exit(1)
    }

    defer fd.Close()

    tohtml_init()

    f := bufio.NewReader(fd)
    scanner := bufio.NewScanner(f)
    i := 0
    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(line)
        if i==0 {
            thisb.title = line
            i = 1
            tohtml_line("@title "+line)
        }

        _,err := tohtml_line(line)
        _,err := thisb.context.WriteString(line)
        if err != nil {
            fmt.Println(err)
            return
        }

    }

    tohtml_end(thisb)

    if err:= scanner.Err(); err != nil {
        fmt.Println(err)
    }

}

type htmlt struct {
    status int
    context bytes.Buffer


func tohtml_init() {

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
            fmt.Println("pos now is", pos)
            src[pos].pname = fname
            pos++
            fmt.Println("ADD ONE")
        }
    }
    return nil
}




