package main

import (
	"bufio"
	"bytes"
	"os"
	//"os/exec"
	"fmt"
	"log"
	"net/http"
	"regexp"
	//"io/ioutil"
	"flag"
	"path/filepath"
	"strconv"
	"time"
)

var httpport int
var workpath string

type rec struct {
	pname string
}

var src []rec
var articlen int

type bcont struct {
	name    string
	title   string
	context string
}

var slblog []bcont

var pos int

var bmap map[string]string

func httprmote(w http.ResponseWriter, r *http.Request) {
	path := r.URL.String()
	fmt.Println("from client: ", "map len", len(bmap))
	fmt.Println("--", path, "--")
	b, ok := bmap[path]
	fmt.Println(b)
	if ok {
		fmt.Fprintf(w, b)
	} else {
		fmt.Fprintf(w, "zz NOT FOUND")
	}
}

func main() {

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
	time.Sleep(500 * time.Millisecond)

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
		if i == 1 {
			fmt.Println("index 0 again\n" + slblog[0].context)
		}
		fmt.Println("in range\n" + slblog[i].context)
		i++
	}

	bmap = make(map[string]string, articlen+2)
	/* index and css */
	for _, v := range slblog {
		bmap[v.name] = v.context
		fmt.Println(v.name)
		fmt.Println("-------")
		fmt.Println(v.context)
		http.HandleFunc(v.name, httprmote)
	}

	setcss()

	port := strconv.Itoa(httpport)
	port = ":" + port
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return
}

func setcss() {
	bmap["/style.css"] = "body {\nbackground-color: #BBBBB; color: #552800;\n" +
        "padding-top: 10px; padding-left: 20px; padding-right: 20px;}\n" +
		"a.link {color:#0000A0;}\n" +
		"a.visited {color: #A000A0;}\n" +
		"a.active {color: #00A000;}\n" +
		"img.thumb {height: 100px; width: 100px;}\n" +
		"img.big {height: 900px; width: 900px;}\n" +
		"img.xref {vertical-align: middle; padding-left:6px; padding-right:10px;}\n"

	http.HandleFunc("/style.css", httprmote)
}

func catfile(path string, name string, thisb *bcont) {
	filename := path + name
	fmt.Println(filename + "'s context:")

	thisb.name = "/" + name

	fd, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer fd.Close()

	f := bufio.NewReader(fd)
	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if i == 0 {
			thisb.title = line
			i = 1
			tohtml_init(line)
		} else {
			tohtml_line(line)
		}
		/*_,err := thisb.context.WriteString(line)
		  if err != nil {
		      fmt.Println(err)
		      return
		  }*/

	}

	thisb.context = tohtml_end()
	fmt.Println("thisb.context")
	fmt.Println(thisb.context)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

type htmlt struct {
	status   int
	inaction int
	title    string
	context  bytes.Buffer
}

var thtml htmlt

func tohtml_init(title string) {
	thtml.status = 1
	thtml.inaction = 0
	thtml.context.Truncate(0)
	//thtml.context = bytes.NewBuffer(nil)

	thtml.context.WriteString("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\"" +
		"\"http://www.w3.org/TR/html4/loose.dtd\">\n")
	thtml.context.WriteString("<html>\n<head>\n" +
		"<link rel=\"stylesheet\" type=\"text/css\" href=\"style.css\" />\n" +
		"<title>" + title + "</title></head>\n<body>\n")
	thtml.context.WriteString("<h1 align=center>" + title + "</h1>")
}

func tohtml_line(line string) {
	thtml.context.WriteString("<p>" + line + "</p>" + "</hr>\n")
}

func tohtml_end() string {
	thtml.context.WriteString("</body></html>")

	return string(thtml.context.Bytes())
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
		}
	}
	return nil
}
