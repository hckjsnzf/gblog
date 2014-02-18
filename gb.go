
package main

import (
    //"os"
    "log"
    "fmt"
    //"io/ioutil"
    "flag"
)

var httpport int
var workpath string

func main () {
    flag.StringVar(&workpath, "workp", "./", "default work path is current dir")
    flag.IntVar(&httpport, "port", 6969, "default http port is 6969")
    flag.Parse()

    log.Println("after Parse")
    fmt.Println("port is ", httpport)
    fmt.Println("work path is ", workpath)

    return
}

