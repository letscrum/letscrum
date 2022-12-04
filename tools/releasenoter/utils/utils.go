package utils

import (
	"os"
	"strings"
)

type OrderMap struct {
	Array []string
	Map   map[string]interface{}
}

func (o *OrderMap) Sort(str string, standard map[string]int) {
	if _, ok := o.Map[str]; ok {
		return
	}
	o.Array = append(o.Array, str)
	for i := len(o.Array) - 1; i > 0; i-- {
		if standard[o.Array[i]] < standard[o.Array[i-1]] {
			o.Array[i], o.Array[i-1] = o.Array[i-1], o.Array[i]
		}
	}
}

func GetDirFiles(folder string) []string {
	var result []string
	files, _ := os.ReadDir(folder)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.IsDir() {
			result = append(result, GetDirFiles(folder+"/"+file.Name())...)
			continue
		}

		if strings.HasSuffix(file.Name(), ".yaml") {
			result = append(result, folder+"/"+file.Name())
		}
	}
	return result
}

func GetVersionName(out string) string {
	str := strings.Split(out, "/")
	if strings.HasSuffix(out, "/") {
		return str[len(str)-2]
	}
	return str[len(str)-1]
}
