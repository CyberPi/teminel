package utils

import "net/http"

func RequestShouldSucceed(url string) bool {
	reponse, err := http.Head(url)
	if err != nil {
		return false
	}
	if reponse.StatusCode >= 300 || reponse.StatusCode < 200 {
		return false
	}
	return true
}
