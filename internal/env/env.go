package env

import (
	"fmt"
	"strings"
)

func ArrToMap(a []string) map[string]string {
	m := map[string]string{}
	for _, s := range a {
		kv := strings.Split(s, "=")
		if len(kv) >= 2 && kv[0] != "" && kv[1] != "" {
			m[kv[0]] = kv[1]
		}
	}
	return m
}

func ToMap(ss ...string) map[string]string {
	return ArrToMap(ss)
}

func MapToArr(m map[string]string) []string {
	a := []string{}
	for k, v := range m {
		if k != "" && v != "" {
			a = append(a, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return a
}

func ToArr(ss ...string) []string {
	a := []string{}
	for i, s := range ss {
		if i%2 == 1 {
			sm1 := ss[i-1]
			if s != "" && sm1 != "" {
				a = append(a, fmt.Sprintf("%s=%s", sm1, s))
			}
		}
	}
	return a
}
