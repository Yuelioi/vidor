package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0"

func doReqBody(client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func GetClient(proxyURL string, useProxy bool) (client *http.Client, err error) {
	proxyStr := proxyURL
	transport := &http.Transport{}
	baseClient := &http.Client{
		// Timeout: 5 * time.Second,
	}

	if useProxy {
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			fmt.Println("Error parsing proxy URL:", err)
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxyURL)
		baseClient.Transport = transport
		return baseClient, err

	}
	baseClient.Transport = transport
	return baseClient, nil

}
