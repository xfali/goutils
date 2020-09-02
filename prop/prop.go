// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package prop

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func CreateFromTempFile(tempFile, targetFile string, props map[string]string) error {
	b, err := ioutil.ReadFile(tempFile)
	if err != nil {
		return err
	}
	return CreateFromTempByte(b, targetFile, props)
}

func CreateFromTempByte(temp []byte, targetFile string, props map[string]string) error {
	file, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	ret := ReplaceTempByte(temp, props)
	_, errW := file.Write(ret)
	return errW
}

func ReplaceTempByte(temp []byte, props map[string]string) []byte {
	tempStr := string(temp)
	for k, v := range props {
		tempStr = strings.Replace(tempStr, fmt.Sprintf("${%s}", k), v, -1)
	}

	return []byte(tempStr)
}

func GetEnvs() map[string]string {
	s := os.Environ()
	ret := map[string]string{}
	for _, env := range s {
		env := strings.TrimSpace(env)
		if env != "" {
			pair := strings.Split(env, "=")
			if len(pair) == 2 {
				ret[pair[0]] = pair[1]
			}
		}
	}

	return ret
}
