package service

import (
	"regexp"
	"strings"
)

type ParsedCurl struct {
	URL     string
	Version string
}

func ParseCurl(command string) (*ParsedCurl, error) {
	urlRegex := regexp.MustCompile(`https?://[^\s'"]+`)
	url := urlRegex.FindString(command)
	if url == "" {
		return nil, nil
	}

	versionRegex := regexp.MustCompile(`/v([\d.]+)`)
	matches := versionRegex.FindStringSubmatch(url)
	version := ""
	if len(matches) > 1 {
		version = matches[1]
	}

	return &ParsedCurl{
		URL:     url,
		Version: version,
	}, nil
}

type ParsedConfig struct {
	URL         string
	Token       string
	ProjectName string
}

func ParseConfig(command string) (*ParsedConfig, error) {
	urlRegex := regexp.MustCompile(`--url\s+([^\s'"]+)`)
	tokenRegex := regexp.MustCompile(`--token\s+([^\s'"]+)`)

	urlMatches := urlRegex.FindStringSubmatch(command)
	tokenMatches := tokenRegex.FindStringSubmatch(command)

	if len(urlMatches) < 2 || len(tokenMatches) < 2 {
		return nil, nil
	}

	url := urlMatches[1]
	token := tokenMatches[1]

	parts := strings.Split(strings.TrimRight(url, "/"), "/")
	projectName := parts[len(parts)-1]

	return &ParsedConfig{
		URL:         url,
		Token:       token,
		ProjectName: projectName,
	}, nil
}
