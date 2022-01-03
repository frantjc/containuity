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

func ToMap(a ...string) map[string]string {
	return ArrToMap(a)
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
