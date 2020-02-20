package helpers

import (
	"log"
	"net/http"
	"strings"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func CheckStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
