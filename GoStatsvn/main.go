package main

import(
    "fmt"
	"flag"
	"log"
	"os"
	"util"
	"strconv"
)

var svnXmlFile *string = flag.String("f", "", "svn log with xml format")
var svnDir *string = flag.String("d", "", "code working directory")

type AuthorStat struct {
	AppendLines int
	RemoveLines int
}

func main() {
	flag.Parse()

	//判断有没有指定svnlog.xml文件
	if *svnXmlFile == "" {
		log.Fatal("-f cannot be empty, -f svnlog.xml")
		return
	}

	//判断有没有指定svnlog.xml文件
	if *svnDir == "" {
		log.Fatal("-d cannot be empty, -d svnWorkDir")
		return
	}

	//判断文件是否存在
	if _,err := os.Stat(*svnXmlFile); os.IsNotExist(err) {
		log.Fatalf("svn log file '%s' not exists.", *svnXmlFile)
	}

	//获取svn root目录
	svnRoot, err := util.GetSvnRoot(*svnDir);

	svnXmlLogs, err := util.ParaseSvnXmlLog(*svnXmlFile)
//	fmt.Printf("%v", svnXmlLogs)
	util.CheckErr(err)

	AuthorStats := make(map[string]AuthorStat)

	for _, svnXmlLog := range svnXmlLogs.Logentry {
		newRev, _ := strconv.Atoi(svnXmlLog.Revision)
		fmt.Printf("svn diff on r%d ,\n", newRev)
		for _, path := range svnXmlLog.Paths {
			if path.Action == "M" && path.Kind == "file" {
				stdout, err := util.CallSvnDiff(newRev-1, newRev, svnRoot+path.Path)
				if err == nil {
					//fmt.Println("stdout ",stdout)
				} else {
					fmt.Println("err ", err.Error())
				}
				appendLines, removeLines, err := util.GetLineDiff(stdout)
				fmt.Printf("\t%s on r%d +%d -%d,\n", path.Path, newRev, appendLines, removeLines)
				if err == nil {
					Author, ok := AuthorStats[svnXmlLog.Author]
					if ok {
						Author.AppendLines += appendLines
						Author.RemoveLines += removeLines
						AuthorStats[svnXmlLog.Author] = Author
					} else {
						Author.AppendLines = appendLines
						Author.RemoveLines = removeLines
						AuthorStats[svnXmlLog.Author] = Author
					}
					//fmt.Println(appendLines, removeLines, AuthorStats)
				}
			}
		}
	}
	//输出结果
	ConsoleOutPutTable(AuthorStats)

}

//console输出结果
func ConsoleOutPutTable(AuthorStats map[string]AuthorStat) {/*{{{*/
	fmt.Printf(" ==user== \t==lines==\n")
	for author, val := range AuthorStats {
		fmt.Printf("%10s\t%d\n", author, val.AppendLines+val.RemoveLines)
	}
}/*}}}*/
