package util

import(
	"strings"
	"strconv"
	"os/exec"
	"bytes"
)

func CallSvnDiff(oldVer, newVer int, fileName string) (stdout string, err error) {

	app := "svn"
	param1 := "diff"
	param2 := "--old"
	param3 := fileName + "@" + strconv.Itoa(oldVer)
	param4 := "--new"
	param5 := fileName + "@" + strconv.Itoa(newVer)

	cmd := exec.Command(app, param1, param2, param3, param4, param5)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", err
	} else {
		return out.String(), nil
	}
}

func GetLineDiff( diffBuffer string) (appendLines, removeLines int, err error) {
	//svndiff 结果头部有 --- +++ 标识,从-1开始计数跳过
	appendLines = -1;
	removeLines = -1;
	err = nil
	lines := strings.Split(diffBuffer, "\n")
	for _,line := range lines {
		if strings.Index(line, "+") == 0 {
			appendLines++
		}
		if strings.Index(line, "-") == 0 {
			removeLines++
		}
	}
	if appendLines == -1 || removeLines == -1 {
		appendLines = 0;
		removeLines = 0;
	}
	return
}
