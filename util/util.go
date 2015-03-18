package util

import(
	"fmt"
	"strings"
)

func GetLineDiff( diffBuffer string) {
	lines := strings.Split(diffBuffer, "\r")
	for _,line := range lines {
		if strings.Index(diffBuffer, "+") == 0 {
			fmt.Println("+ one line;")
		}
		if strings.Index(diffBuffer, "-") == 0 {
			fmt.Println("- one line;")
		}
		fmt.Println(line)
	}
}
