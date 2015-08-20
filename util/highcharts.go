//图表类,画图  , line-basic
package util

import (
//	"bytes"
	"encoding/json"
//	"io/ioutil"
//	"log"
//	"strconv"
//	"strings"
	"time"
	"sort"
	"statStruct"
)

//line-basic-x
type XAxis struct {
	Categories []string `json:"categories"`
}

//line-basic-data-single
type Serie struct {
	Name string `json:"name"`
	Data []int `json:"data"`
}

//line-basic-data
type Series struct {
	Series []Serie `json:"series"`
}

func GetXAxis (minTimestamp int64, maxTimestamp int64) (xstring string) {/*{{{*/
	var xaxis XAxis
	minTime := time.Unix(minTimestamp, 0)
	minDay := minTime.Format(DATE_DAY)
	minTime, _ = time.Parse(DATE_DAY, minDay)
	minDayTimestamp := minTime.Unix()
	slice := make([]string, 366)
	i := 1
	slice[i] = minDay
	for {
			minDayTimestamp += 86400;
			minTime = time.Unix(minDayTimestamp, 0)
			minDay = minTime.Format(DATE_DAY)
			//todo output minDay
			if (minDayTimestamp > maxTimestamp) {
				break;
			}
			i++
			slice[i] = minDay
	}
	xaxis.Categories = slice[1:i+1]
	xbytes, _ := json.Marshal(xaxis)
	xstring = string(xbytes)
	return
}/*}}}*/

func GetSeries (dayAuthorStats statStruct.AuthorTimeStats) (seriesString string) {/*{{{*/
	var series Series
	tmpSeries := make([]Serie, 50)
	i := 0
	for author, stats := range dayAuthorStats {
		var serie Serie;
		serie.Name = author
		j := 0
		tmpData := make([]int, 365)
		var keys []string
		for k, _ := range stats {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			j++
			tmpData[j] = stats[k].AppendLines + stats[k].RemoveLines
		}
		serie.Data = tmpData[1:j+1]
		i++
		tmpSeries[i] = serie
	}
	series.Series = tmpSeries[1:i+1]
	seriesByte , _ := json.Marshal(series)
	seriesString = string(seriesByte)
	return
}/*}}}*/

