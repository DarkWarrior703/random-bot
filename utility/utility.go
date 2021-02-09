package utility

import (
	"io/ioutil"
	"net/http"
)

// GetData gets the data from a website and returns bytes
func GetData(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
