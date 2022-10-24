/*
 * Copyright 2022 Xiongfa Li.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pattern

import "fmt"

//type regexpPatterns struct{}
//
//var Patterns regexpPatterns

func IsNumber() string {
	return "^[0-9]*$"
}

func IsChinese() string {
	return `^[\u4e00-\u9fa5]{0,}$`
}

func IsAlphabet() string {
	return `^[A-Za-z]+$`
}

func IsNumberAndAlphabet() string {
	return `^[A-Za-z0-9]+$`
}

func IsUpper() string {
	return `^[A-Z]+$`
}

func IsLower() string {
	return `^[a-z]+$`
}

func StartWith(s string) string {
	return fmt.Sprintf("^%s.*", s)
}

func EndWith(s string) string {
	return fmt.Sprintf(".*%s$", s)
}

// require github.com/dlclark/regexp2
//func NotStartWith(s string) string {
//	return fmt.Sprintf(`^(?!%s).*`, s)
//}
// require github.com/dlclark/regexp2
//func NotEndWith(s string) string {
//	return fmt.Sprintf(".*(?<!%s)$", s)
//}

func Include(s string) string {
	return fmt.Sprintf(".*%s.*", s)
}

// require github.com/dlclark/regexp2
func Exclude(s string) string {
	return fmt.Sprintf("^((?!%s).)*$", s)
}

func Equal(s string) string {
	return fmt.Sprintf("^%s$", s)
}

// require github.com/dlclark/regexp2
func NotEqual(s string) string {
	return fmt.Sprintf("^(?!%s$)", s)
}
