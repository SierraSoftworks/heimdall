package transports

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type TransportURL struct {
	Protocol    string
	User        string
	Password    string
	Host        string
	TopicPrefix string

	Options map[string]string

	url *url.URL
}

func ParseURL(u string) (*TransportURL, error) {
	p, err := url.Parse(os.ExpandEnv(u))
	if err != nil {
		return nil, err
	}

	if p.Host == "" {
		return nil, fmt.Errorf("no host specified")
	}

	user, pass := "", ""
	if p.User != nil {
		user = p.User.Username()
		pass, _ = p.User.Password()
	}

	options := map[string]string{}
	for k, v := range p.Query() {
		options[k] = v[0]
	}

	return &TransportURL{
		Protocol:    strings.TrimRight(p.Scheme, ":"),
		User:        user,
		Password:    pass,
		Host:        p.Host,
		TopicPrefix: strings.TrimLeft(p.Path, "/"),
		Options:     options,

		url: p,
	}, nil
}

func (u *TransportURL) GetFullTopic(topic string) string {
	return strings.TrimLeft(fmt.Sprintf("%s/%s", u.TopicPrefix, topic), "/")
}

func (u *TransportURL) GetOption(option, def string) string {
	val, ok := u.Options[option]
	if !ok {
		return def
	}
	return val
}

func (u *TransportURL) SafeString() string {
	return strings.TrimRight(fmt.Sprintf("%s://%s/%s", u.Protocol, u.Host, u.TopicPrefix), "/")
}

func (u *TransportURL) String() string {
	return u.url.String()
}
