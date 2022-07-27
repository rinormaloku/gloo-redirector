package mesh

import (
	"bytes"
	_ "embed"
	"errors"
	"github.com/rinormaloku/gloo-redirector/pkg/csv"
	"github.com/rinormaloku/gloo-redirector/pkg/domain"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strconv"
	"strings"
)

type Redirection struct {
	ResourceName string
	Host         string
	Matchers     []Matcher
}

type Matcher struct {
	ExactPath    string
	HostRewrite  string
	PathRewrite  string
	RedirectCode string
}

func convertToCode(redirectCode int) (string, error) {
	switch redirectCode {
	case 301:
		return "MOVED_PERMANENTLY", nil
	case 302:
		return "FOUND", nil
	case 303:
		return "SEE_OTHER", nil
	case 307:
		return "TEMPORARY_REDIRECT", nil
	case 308:
		return "PERMANENT_REDIRECT", nil
	}
	return "", errors.New("the supported redirect codes are 301, 302, 303, 307 and 308, meanwhile the provided code is " + strconv.Itoa(redirectCode))
}

func generate(inputData domain.InputData) (bytes.Buffer, error) {
	var payload bytes.Buffer

	var hostRedirections = make(map[string]Redirection)

	rulesCSV := csv.ReadFile(inputData.CsvFile)

	for _, rule := range rulesCSV {
		fromUrl, err := url.Parse(rule[0])
		if err != nil {
			panic(err)
		}

		toUrl, err := url.Parse(rule[1])
		if err != nil {
			panic(err)
		}

		code, err := strconv.Atoi(rule[2])
		if err != nil {
			panic(err)
		}

		gloomeshRedirectCode, err := convertToCode(code)
		if err != nil {
			panic(err)
		}

		if _, ok := hostRedirections[fromUrl.Host]; !ok {
			hostRedirections[fromUrl.Host] = Redirection{
				ResourceName: strings.ReplaceAll(fromUrl.Host, ".", "-"),
				Host:         fromUrl.Host,
				Matchers:     []Matcher{},
			}
		}

		red := hostRedirections[fromUrl.Host]
		red.Matchers = append(hostRedirections[fromUrl.Host].Matchers, Matcher{
			ExactPath:    fromUrl.Path,
			HostRewrite:  toUrl.Host,
			PathRewrite:  defaultIfEmpty(toUrl.Path, "/"),
			RedirectCode: gloomeshRedirectCode,
		})
		hostRedirections[fromUrl.Host] = red
		print("")
	}

	parse, err := inputData.ParseTemplate()
	if err != nil {
		log.WithError(err).Error("fail to parse template")
		return payload, err
	}

	if err != nil {
		log.WithError(err).Error("fail to parse template")
		return payload, err
	}

	for _, redirection := range hostRedirections {
		err = parse.Execute(&payload, redirection)
		if err != nil {
			log.WithError(err).Error("fail to execute content to template")
			return payload, err
		}
		payload.WriteString("---\n")
	}

	return payload, nil
}

func defaultIfEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
