package main

import(
    "fmt"
	"flag"
	"log"
	"os"
	"util"
	"encoding/xml"
	"io/ioutil"
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


	type Path struct {
		Action string `xml:"action,attr"`
		Kind string `xml:"kind,attr"`
		Path string `xml:",chardata"`
	}

	type Logentry struct {
		Revision string `xml:"revision,attr"`
		Author string `xml:"author"`
		Date string	`xml:"date"`
		Paths []Path `xml:"paths>path"`
		Msg string	`xml:"msg"`
	}

	type SvnXmlLogs struct {
		Logentry []Logentry `xml:"logentry"`
	}

	content, err := ioutil.ReadFile(*svnXmlFile)
	if err != nil {
		log.Fatal(err)
	}
	var svnXmlLogs SvnXmlLogs
	err = xml.Unmarshal(content, &svnXmlLogs)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(content))
	fmt.Printf("%v", svnXmlLogs)


	if err == nil {
		stdout, err := util.CallSvnDiff(43876, 43877, "/home/s/www/fukun/svn/BigDataPlatform/trunk/application/models/NewsTop.php")
		if err == nil {
			fmt.Println("stdout ",stdout)
		} else {
			fmt.Println("err ", err.Error())
		}
		appendLines, removeLines, err := util.GetLineDiff(stdout)
		fmt.Println(appendLines, removeLines, err)
	}

}
