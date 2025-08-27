package internal

import (
	"slices"
	"strings"

	"github.com/jordan-wright/email"
)

type emailHeader struct {
	key, value string
}

func getAllHeaders(email *email.Email) []emailHeader {
	return append(getCanonicalHeaders(email), getMimeHeaders(email)...)
}

func getCanonicalHeaders(email *email.Email) []emailHeader {
	return []emailHeader{
		{"From", email.From},
		{"To", strings.Join(email.To, ", ")},
		{"Cc", strings.Join(email.Cc, ", ")},
		{"Bcc", strings.Join(email.Bcc, ", ")},
		{"Subject", email.Subject},
	}
}

func getMimeHeaders(email *email.Email) []emailHeader {
	selectedHeaders := []string{}
	for header := range email.Headers {
		switch header {
		case "Content-Type", "Content-Transfer-Encoding", "Mime-Version", "Message-Id", "Date":
			// skip header

		default:
			selectedHeaders = append(selectedHeaders, header)
		}
	}

	slices.Sort(selectedHeaders)

	res := []emailHeader{}
	for _, header := range selectedHeaders {
		value := strings.Join(email.Headers.Values(header), ", ")
		res = append(res, emailHeader{header, value})
	}

	return res
}
