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

import (
	"github.com/dlclark/regexp2"
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	t.Run("number", func(t *testing.T) {
		reg, err := regexp2.Compile(IsNumber(), regexp2.None)
		if err != nil {
			t.Fatal(err)
		}

		m, err := reg.MatchString("1")
		if err != nil {
			t.Fatal(err)
		}
		if !m {
			t.Fatal("expect true but get ", m)
		}
		m, err = reg.MatchString("1a")
		if err != nil {
			t.Fatal(err)
		}
		if m {
			t.Fatal("expect false but get ", m)
		}

		m, err = reg.MatchString("a1")
		if err != nil {
			t.Fatal(err)
		}
		if m {
			t.Fatal("expect false but get ", m)
		}
	})

	t.Run("startWith", func(t *testing.T) {
		reg, err := regexp.Compile(StartWith("aaa"))
		if err != nil {
			t.Fatal(err)
		}
		m := reg.MatchString("aaabbbccc")
		if !m {
			t.Fatal("expect true but get ", m)
		}

		m = reg.MatchString("cccaaabbb")
		if m {
			t.Fatal("expect false but get ", m)
		}

		m = reg.MatchString("bbbcccaaa")
		if m {
			t.Fatal("expect false but get ", m)
		}
	})

	t.Run("endWith", func(t *testing.T) {
		reg, err := regexp.Compile(EndWith("ccc"))
		if err != nil {
			t.Fatal(err)
		}
		m := reg.MatchString("aaabbbccc")
		if !m {
			t.Fatal("expect true but get ", m)
		}

		m = reg.MatchString("cccaaabbb")
		if m {
			t.Fatal("expect false but get ", m)
		}

		m = reg.MatchString("bbbcccaaa")
		if m {
			t.Fatal("expect false but get ", m)
		}
	})

	t.Run("equal", func(t *testing.T) {
		reg, err := regexp.Compile(Equal("aaabbbccc"))
		if err != nil {
			t.Fatal(err)
		}
		m := reg.MatchString("aaabbbccc")
		if !m {
			t.Fatal("expect true but get ", m)
		}

		m = reg.MatchString("cccaaabbb")
		if m {
			t.Fatal("expect false but get ", m)
		}

		m = reg.MatchString("bbbcccaaa")
		if m {
			t.Fatal("expect false but get ", m)
		}
	})

	t.Run("not equal", func(t *testing.T) {
		reg, err := regexp2.Compile(NotEqual("cccaaabbb"), 0)
		if err != nil {
			t.Fatal(err)
		}
		m, err := reg.MatchString("aaabbbccc")
		if err != nil {
			t.Fatal(err)
		}
		if !m {
			t.Fatal("expect true but get ", m)
		}

		m, err = reg.MatchString("cccaaabbb")
		if err != nil {
			t.Fatal(err)
		}
		if m {
			t.Fatal("expect false but get ", m)
		}

		m, err = reg.MatchString("cccaaabbb")
		if err != nil {
			t.Fatal(err)
		}
		if m {
			t.Fatal("expect false but get ", m)
		}
	})

	t.Run("include", func(t *testing.T) {
		reg, err := regexp.Compile(Include("abbbc"))
		if err != nil {
			t.Fatal(err)
		}
		m := reg.MatchString("abbbc")
		if !m {
			t.Fatal("expect true but get ", m)
		}

		m = reg.MatchString("aaabbbccc")
		if !m {
			t.Fatal("expect true but get ", m)
		}

		m = reg.MatchString("cccaaabbb")
		if m {
			t.Fatal("expect false but get ", m)
		}

		m = reg.MatchString("bbbcccaaa")
		if m {
			t.Fatal("expect false but get ", m)
		}
	})

	t.Run("exclude", func(t *testing.T) {
		reg, err := regexp2.Compile(Exclude("abbbc"), 0)
		if err != nil {
			t.Fatal(err)
		}
		m, err := reg.MatchString("abbbc")
		if m {
			t.Fatal("expect false but get ", m)
		}

		m, err = reg.MatchString("aaabbbccc")
		if m {
			t.Fatal("expect false but get ", m)
		}

		m, err = reg.MatchString("cccaaabbb")
		if !m {
			t.Fatal("expect true but get ", m)
		}

		m, err = reg.MatchString("bbbcccaaa")
		if !m {
			t.Fatal("expect true but get ", m)
		}
	})
}
