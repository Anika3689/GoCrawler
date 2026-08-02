package main
import (
	"net/url"
	"strings"
	"path"
)

func normalizeURL(originalURL string) (string, error) {
	URL, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}
	
	host := strings.ToLower(URL.Hostname())

	pathNorm := ""
	if URL.Path != "" { 
		pathNorm = path.Clean(URL.Path)
		if !strings.HasPrefix(pathNorm, "/") {
			pathNorm = "/" + pathNorm
		}
	}

	if URL.RawQuery != "" {
		pathNorm += "?" + URL.RawQuery
	}	

	return host + pathNorm, nil  
	
}
