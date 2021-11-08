package main

import "strings"

func Contains(a []string, s string) bool {
	for _, v := range a {
		if s == v {
			return true
		}
	}
	return false
}

func Remove(p []Projectile, i int) []Projectile {
	p[i] = p[len(p)-1]
	return p[:len(p)-1]
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func AbsDiff(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func CamelCase(s string) string {
	var snek bool
	var camel string
	for i := 0; i < len(s); i++ {
		char := string(s[i])
		if char == "_" {
			snek = true
			continue
		}
		if snek {
			snek = false
			camel += strings.ToUpper(char)
		} else {
			camel += char
		}
	}
	return camel
}
