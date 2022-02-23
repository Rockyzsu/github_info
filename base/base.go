package base

import (
	"encoding/json"
	"fmt"
	"github.com/gittool/models"
	"io/ioutil"
	"log"
	"os"
)

func Must(err error) {
	if err != nil {
		panic(err)
		//return nil
	}
}
func JsonParse(filename string) (models.GithubToken, error) {
	v := models.GithubToken{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return v, err
	}
	err = json.Unmarshal(data, &v) // 使用这个库得要把 json的定义结构体 的首字母大写
	if err != nil {
		log.Println(err)
		return v, err
	}
	return v, nil
}

func Token() string {
	filename := "token.json"
	token, err := JsonParse(filename)
	if err != nil {
		panic(err)
	}
	return token.Token

}

func RemoveDuplicated(originText []string) []string {
	target := make([]string, 0)
	tmpDict := make(map[string]interface{}, 0)
	for _, value := range originText {
		if _, ok := tmpDict[value]; !ok {
			target = append(target, value)
			tmpDict[value] = nil
		}
	}
	return target
}

func GetMapValue(m map[string]int) int {
	//fmt.Println(m)
	if len(m) < 1 {
		panic("length error")
	}
	for _, v := range m {
		return v
	}

	return 0

}
func GetMapKey(m map[string]int) string {
	//fmt.Println(m)
	if len(m) < 1 {
		panic("length error")
	}
	for k, _ := range m {
		return k
	}

	return ""
}

func Sort(origin map[string]int) []map[string]int {
	length := len(origin)
	//fmt.Println("now length is ", length)
	target := make([]map[string]int, 0)

	for k, v := range origin {
		tmpMap := make(map[string]int)
		tmpMap[k] = v
		target = append(target, tmpMap)
	}

	//fmt.Println(target)

	for i := 0; i < length; i++ {
		for j := i; j < length; j++ {
			if GetMapValue(target[i]) < GetMapValue(target[j]) {
				// swap
				tmpMap := make(map[string]int)
				tmpMap = target[i]
				target[i] = target[j]
				target[j] = tmpMap
			}
		}
	}
	return target
}

func Help() {
	fmt.Println(`
./main user rockyzsu
./main fans rockyzsu
./main repo rockyzsu
`)
	os.Exit(0)
}
