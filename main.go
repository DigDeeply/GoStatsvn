package main

import(
    "fmt"
	"flag"
	"log"
	"os"
)

var svnXmlFile *string = flag.String("f", "", "svn log with xml format")

func main() {
	flag.Parse()

	//判断有没有指定svnlog.xml文件
	if *svnXmlFile == "" {
		log.Fatal("-f cannot be empty, -f svnlog.xml")
		return
	}

	//判断文件是否存在
	if _,err := os.Stat(*svnXmlFile); os.IsNotExist(err) {
		log.Fatalf("svn log file '%s' not exists.", *svnXmlFile)
	}

    fmt.Println("hello world", *svnXmlFile);
}
