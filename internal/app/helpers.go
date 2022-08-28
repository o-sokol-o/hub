package app // TODO: ликвидировать файл после рефакторинга

import (
	"strings"
)

// If base URL is empty, set the default value from the HTTP Host
func baseURLtoDefault(url, host string) string {
	if url == "" {
		baseUrl := host
		if !strings.HasPrefix(baseUrl, "http") {
			if strings.HasPrefix(baseUrl, "0.0.0.0:") {
				pts := strings.Split(baseUrl, ":")
				pts[0] = "127.0.0.1"
				baseUrl = strings.Join(pts, ":")
			} else if strings.HasPrefix(baseUrl, ":") {
				baseUrl = "127.0.0.1" + baseUrl
			}
			baseUrl = "http://" + baseUrl
		}
		return baseUrl
	}
	return url
}
