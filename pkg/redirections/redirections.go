package redirections

import (
	"bytes"
	_ "embed"
	"errors"
	log "github.com/sirupsen/logrus"
	"gloo-redirector/domain"
	"gloo-redirector/pkg/csv"
	"gloo-redirector/pkg/templates"
	"net/url"
	"strconv"
	"strings"
	"text/template"
)

type Redirection struct {
	RouteTableName string
	Host           string
	Matchers       []Matcher
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

func Generate(inputData domain.InputData) (bytes.Buffer, error) {
	var payload bytes.Buffer

	var hostRedirections = make(map[string]Redirection)

	rulesCSV := csv.ReadFile(inputData.File)

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
				RouteTableName: strings.ReplaceAll(fromUrl.Host, ".", "-"),
				Host:           fromUrl.Host,
				Matchers:       []Matcher{},
			}
		}

		red := hostRedirections[fromUrl.Host]
		red.Matchers = append(hostRedirections[fromUrl.Host].Matchers, Matcher{
			ExactPath:    fromUrl.Path,
			HostRewrite:  toUrl.Host,
			PathRewrite:  toUrl.Path,
			RedirectCode: gloomeshRedirectCode,
		})
		hostRedirections[fromUrl.Host] = red
		print("")
	}

	rtTemplate := template.New("routetable.yaml")

	var parse *template.Template
	var err error
	if len(inputData.Template) > 0 {
		parse, err = rtTemplate.Parse(string(inputData.Template))
	} else {
		parse, err = rtTemplate.Parse(templates.RoutetableYaml)
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
