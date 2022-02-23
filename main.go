package main

import (
	"fmt"
	"github.com/gittool/base"
	"github.com/gittool/githupapi"
	"os"
	"time"
)

func main() {
	start := time.Now()
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./main help")
		base.Help()
	}
	user := os.Args[2]
	switch os.Args[1] {
	case "user":
		githupapi.GetUserBaseInfos(user)
	case "fans":
		githupapi.GetHotFans(user, 10)
	case "repo":
		githupapi.GetUserRepos(user)
	default:
		base.Help()
	}
	used := time.Since(start)
	fmt.Println("Time used ", used)
}
