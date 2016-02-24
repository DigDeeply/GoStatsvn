package main

import (
	"GoStatsvn/statStruct"
	"GoStatsvn/util"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	DEFAULT_SMALLEST_TIME_STRING = "1000-03-20T08:38:17.428370Z"
	DATE_DAY                     = "2006-01-02"
	DATE_HOUR                    = "2006-01-02 15"
	DATE_SECOND                  = "2006-01-02T15:04:05Z"
)

var svnXmlFile *string = flag.String("f", "", "svn log with xml format")
var svnDir *string = flag.String("d", "", "code working directory")
var chartTemplate *string = flag.String("t", "", "hightcharts Template file")
var chartData statStruct.ChartData

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

	//判断有没有指定画图的模版文件
	if *chartTemplate == "" {
		log.Fatal("-t cannot be empty, -t hightchartsTemplate file path")
		return
	}

	//判断文件是否存在
	if _, err := os.Stat(*svnXmlFile); os.IsNotExist(err) {
		log.Fatalf("svn log file '%s' not exists.", *svnXmlFile)
	}

	//获取svn root目录
	svnRoot, err := util.GetSvnRoot(*svnDir)

	svnXmlLogs, err := util.ParaseSvnXmlLog(*svnXmlFile)
	//	fmt.Printf("%v", svnXmlLogs)
	util.CheckErr(err)

	authorTimeStats := make(statStruct.AuthorTimeStats)

	AuthorStats := make(map[string]statStruct.AuthorStat)

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
					//综合统计
					Author, ok := AuthorStats[svnXmlLog.Author]
					if ok {
						Author.AppendLines += appendLines
						Author.RemoveLines += removeLines
					} else {
						Author.AppendLines = appendLines
						Author.RemoveLines = removeLines
					}
					AuthorStats[svnXmlLog.Author] = Author
					//todo 记录人和日期的详细log，用于细分统计
					authorTimeStat, ok := authorTimeStats[svnXmlLog.Author]
					saveTime, err := time.Parse("2006-01-02T15:04:05Z", svnXmlLog.Date)
					util.CheckErr(err)
					saveTimeStr := saveTime.Format(DATE_SECOND)
					if ok {
						Author, ok := authorTimeStat[saveTimeStr]
						if ok {
							Author.AppendLines += appendLines
							Author.RemoveLines += removeLines
						} else {
							Author.AppendLines = appendLines
							Author.RemoveLines = removeLines
						}
						authorTimeStat[saveTimeStr] = Author
					} else {
						Author.AppendLines = appendLines
						Author.RemoveLines = removeLines
						authorTimeStat = make(statStruct.AuthorTimeStat)
						authorTimeStat[saveTimeStr] = Author
					}
					authorTimeStats[svnXmlLog.Author] = authorTimeStat
					//fmt.Println(appendLines, removeLines, AuthorStats)
				}
			}
		}
	}
	//输出结果
	ConsoleOutPutTable(AuthorStats)
	//fmt.Printf("%v\n", authorTimeStats)
	minTimestamp, maxTimestamp := getMinMaxTimestamp(authorTimeStats)
	fmt.Printf("%d\t%d\n", minTimestamp, maxTimestamp)
	dayAuthorStats := StatLogByDay(authorTimeStats)
	fmt.Printf("%v\n", dayAuthorStats)
	dayAuthorStatsOutput := StatLogByFullDay(dayAuthorStats, minTimestamp, maxTimestamp)
	xaxis := util.GetXAxis(minTimestamp, maxTimestamp)
	series := util.GetSeries(dayAuthorStatsOutput)
	chartData.XAxis = xaxis
	chartData.Series = series
	fmt.Printf("%s\n%s\n", xaxis, series)
	DrawCharts()
	//输出按小时统计结果
	//ConsoleOutPutHourTable(authorTimeStats)
	//输出按周统计结果
	//ConsoleOutPutWeekTable(authorTimeStats)

}

//console输出结果
func ConsoleOutPutTable(AuthorStats map[string]statStruct.AuthorStat) { /*{{{*/
	fmt.Printf(" ==user== \t==lines==\n")
	for author, val := range AuthorStats {
		fmt.Printf("%10s\t%5d\n", author, val.AppendLines+val.RemoveLines)
	}
} /*}}}*/

//返回按天格式化好的数据
func StatLogByDay(authorTimeStats statStruct.AuthorTimeStats) (dayAuthorStats statStruct.AuthorTimeStats) { /*{{{*/
	dayAuthorStats = make(map[string]statStruct.AuthorTimeStat)
	for author, detail := range authorTimeStats {
		dayAuthorStat := make(map[string]statStruct.AuthorStat)
		_, ok := dayAuthorStats[author]
		if ok {
		} else {
			dayAuthorStats[author] = dayAuthorStat
		}
		for timeString, stats := range detail {
			//todo 找到正常转化时间的方法
			timeTime, err := time.Parse(time.RFC3339, timeString)
			util.CheckErr(err)
			timeFormat := timeTime.Format(DATE_DAY)
			//fmt.Printf("%v\t%v\n", timeString, timeTime)
			if err == nil {
				oldDayAuthorStat, ok := dayAuthorStat[timeFormat]
				var authorStat statStruct.AuthorStat
				if ok {
					authorStat.AppendLines = oldDayAuthorStat.AppendLines + stats.AppendLines
					authorStat.RemoveLines = oldDayAuthorStat.RemoveLines + stats.RemoveLines
				} else {
					authorStat.AppendLines = stats.AppendLines
					authorStat.RemoveLines = stats.RemoveLines
				}
				dayAuthorStat[timeFormat] = authorStat
			}
		}
		dayAuthorStats[author] = dayAuthorStat
	}
	return
} /*}}}*/

func StatLogByFullDay(dayAuthorStats statStruct.AuthorTimeStats, minTimestamp int64, maxTimestamp int64) (dayAuthorStatsOutput statStruct.AuthorTimeStats) { /*{{{*/
	//得到时间的开始和结束日期
	minTime := time.Unix(minTimestamp, 0)
	minDay := minTime.Format(DATE_DAY)
	minTime, _ = time.Parse(DATE_DAY, minDay)
	minDayTimestamp := minTime.Unix()
	maxTime := time.Unix(maxTimestamp, 0)
	maxDay := maxTime.Format(DATE_DAY)
	maxTime, _ = time.Parse(DATE_DAY, maxDay)
	maxDayTimestamp := maxTime.Unix()
	dayAuthorStatsOutput = make(statStruct.AuthorTimeStats)
	//遍历所有author
	for author, dayAuthorStat := range dayAuthorStats {
		fmt.Printf("====user: %s=====\n", author)
		minDayAuthor := minDay
		minTimeAuthor := minTime
		minDayTimestampAuthor := minDayTimestamp
		dayAuthorStatOutput := make(statStruct.AuthorTimeStat)
		//输出每个用户每天的信息
		for {
			authorStat, ok := dayAuthorStat[minDayAuthor]
			if ok {
				fmt.Printf("%s\t%d\n", minDayAuthor, authorStat.AppendLines+authorStat.RemoveLines)
				dayAuthorStatOutput[minDayAuthor] = authorStat
			} else {
				fmt.Printf("%s\t%d\n", minDayAuthor, 0)
				authorStat.AppendLines = 0
				authorStat.RemoveLines = 0
				dayAuthorStatOutput[minDayAuthor] = authorStat
			}
			minDayTimestampAuthor += 86400
			minTimeAuthor = time.Unix(minDayTimestampAuthor, 0)
			minDayAuthor = minTimeAuthor.Format(DATE_DAY)
			if minDayTimestampAuthor > maxDayTimestamp {
				break
			}
		}
		dayAuthorStatsOutput[author] = dayAuthorStatOutput
	}
	fmt.Printf("%v\n", dayAuthorStatsOutput)
	return
} /*}}}*/

//console 按天输出结果，空余的天按0补齐
//获取时间的最大值和最小值
func getMinMaxTimestamp(authorTimeStats statStruct.AuthorTimeStats) (minTimestamp int64, maxTimestamp int64) { /*{{{*/
	minTimestamp = 0
	maxTimestamp = 0
	//先取得时间的最大值和最小值
	for _, detail := range authorTimeStats {
		//fmt.Printf("%s\t%v\n", author, detail)
		for timeString, _ := range detail {
			timeTime, err := time.Parse(DATE_SECOND, timeString)
			if err == nil {
				if minTimestamp == 0 || minTimestamp > timeTime.Unix() {
					minTimestamp = timeTime.Unix()
				}
				if maxTimestamp < timeTime.Unix() {
					maxTimestamp = timeTime.Unix()
				}
			}
		}
		//fmt.Printf("%d\t%d\n", minTimestamp, maxTimestamp)
	}
	return
} /*}}}*/

//console按小时输出结果
//todo 此处有bug,1.没有全部按小时归并，还是按每天每小时归并的。2.显示的小时不是按24小时制
func ConsoleOutPutHourTable(authorTimeStats statStruct.AuthorTimeStats) { /*{{{*/
	defaultSmallestTime, _ := time.Parse("2006-01-02T15:04:05Z", DEFAULT_SMALLEST_TIME_STRING)
	fmt.Printf(" ==user== \t==hour==\t==lines==\n")
	//先取到时间的区间值
	for authorName, Author := range authorTimeStats {
		var minTime time.Time
		var maxTime time.Time
		for sTime, _ := range Author {
			fmtTime, err := time.Parse(DATE_HOUR, sTime)
			util.CheckErr(err)
			if minTime.Before(defaultSmallestTime) || minTime.After(fmtTime) {
				minTime = fmtTime
			}
			if maxTime.Before(defaultSmallestTime) || maxTime.Before(fmtTime) {
				maxTime = fmtTime
			}
		}
		//Todo 用户按时合并,去重
		//输出单个用户的数据
		for sTime, Sval := range Author {
			fmtTime, err := time.Parse(DATE_HOUR, sTime)
			util.CheckErr(err)
			fmt.Printf("%10s\t%5d\t%12d\n", authorName, fmtTime.Hour(), Sval.AppendLines+Sval.RemoveLines)
		}
	}
} /*}}}*/

//console按周输出结果
func ConsoleOutPutWeekTable(authorTimeStats statStruct.AuthorTimeStats) { /*{{{*/
	weekAuthorStats := make(map[string]map[string]statStruct.AuthorStat)
	for authorName, Author := range authorTimeStats {
		weekAuthorStat := make(map[string]statStruct.AuthorStat)
		_, ok := weekAuthorStats[authorName]
		if ok {
		} else {
			weekAuthorStats[authorName] = weekAuthorStat
		}
		for sTime, sAuthor := range Author {
			fmtTime, err := time.Parse(DATE_HOUR, sTime)
			util.CheckErr(err)
			week := fmtTime.Weekday().String()
			oldAuthorStat, ok := weekAuthorStat[week]
			var authorStat statStruct.AuthorStat
			if ok {
				authorStat.AppendLines = oldAuthorStat.AppendLines + sAuthor.AppendLines
				authorStat.RemoveLines = oldAuthorStat.RemoveLines + sAuthor.RemoveLines
			} else {
				authorStat.AppendLines = sAuthor.AppendLines
				authorStat.RemoveLines = sAuthor.RemoveLines
			}
			weekAuthorStat[week] = authorStat
		}
		weekAuthorStats[authorName] = weekAuthorStat
	}
	fmt.Printf(" ==user== \t==week==\t==lines==\n")
	allWeeks := []string{
		"Sunday ",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
	}
	//输出
	for authorName, weekAuthorStat := range weekAuthorStats {
		for _, oneDay := range allWeeks {
			authorStat, ok := weekAuthorStat[oneDay]
			if ok {
				fmt.Printf("%10s\t%5s\t%12d\n", authorName, oneDay, authorStat.AppendLines+authorStat.RemoveLines)
			} else {
				fmt.Printf("%10s\t%5s\t%12d\n", authorName, oneDay, 0)
			}
		}
	}
} /*}}}*/

func showHandle(w http.ResponseWriter, r *http.Request) {
	//filename := r.FormValue("id")
	//imagePath := UPLOAD_DIR + "/" + filename

	//w.Header().Set("Content-Type", "text/html")
	//http.ServeFile(w, r, "src/gostatsvn.html")
	t, err := template.ParseFiles(*chartTemplate)
	if err != nil {
		log.Fatal("not find file: ", err.Error())
	} else {
		locals := make(map[string]interface{})
		xaxis := template.HTML(chartData.XAxis)
		series := template.HTML(chartData.Series)
		locals["xaxis"] = xaxis
		locals["series"] = series
		t.Execute(w, locals)
	}
}

func DrawCharts() {
	http.HandleFunc("/", showHandle)
	log.Println("listen on 8088")
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatal("listen fatal: ", err.Error())
	}
}
