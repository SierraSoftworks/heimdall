package transports

import (
	"fmt"
	"net/url"
	"strings"
)

func SafeURLHost(u *url.URL) string {
	return strings.Trim(fmt.Sprintf("%s://%s", u.Scheme, u.Host), "/:")
}

func SafeURLString(u *url.URL) string {
	return strings.Trim(fmt.Sprintf("%s/%s", SafeURLHost(u), strings.TrimLeft(u.Path, "/")), "/:")
}

func GetFullTopic(u *url.URL, topic string) string {
	return strings.TrimLeft(fmt.Sprintf("%s/%s", u.Path, topic), "/")
}

func GetURLOption(u *url.URL, name, defaultValue string) string {
	v := u.Query().Get(name)
	if v == "" {
		return defaultValue
	}

	return v
}
