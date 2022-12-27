package web1

import "net/http"

func IsBadRequest(response *http.Response) bool {
	return response.StatusCode == 400
}
