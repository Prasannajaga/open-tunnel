package utils

import "strconv"

func BuildAddress(host string, port string) string {
	return host + ":" + port
}

func BuildLocalAddress(port int) string {
	return "localhost:" + strconv.Itoa(port)
}
