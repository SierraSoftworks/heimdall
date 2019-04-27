package transports

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"time"

	log "github.com/Sirupsen/logrus"
)

type FileTransport struct {
	f   *os.File
	url *url.URL
}

func NewFileTransport(u *url.URL) (*FileTransport, error) {
	path := filepath.Join(append([]string{u.Host}, strings.Split(u.Path, "/")...)...)
	log.
		WithField("path", path).
		WithField("transport", SafeURLString(u)).
		Debug("Opening file for transport")

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0664)
	if err != nil {
		log.
			WithField("path", path).
			WithField("transport", SafeURLString(u)).
			WithError(err).
			Debug("Failed to open file for transport")
		return nil, err
	}

	log.
		WithField("path", path).
		WithField("transport", SafeURLString(u)).
		Debug("Transport ready")
	return &FileTransport{
		f:   f,
		url: u,
	}, nil
}

func (t *FileTransport) Describe() string {
	return SafeURLString(t.url)
}

func (t *FileTransport) Subscribe(topic string) (Subscription, error) {
	if t.f == nil {
		return nil, fmt.Errorf("not connected")
	}

	log.
		WithField("topic", topic).
		WithField("transport", t.Describe()).
		Debug("Creating new subscriber for transport")

	f, err := os.OpenFile(t.f.Name(), os.O_RDONLY|os.O_SYNC, 0664)
	if err != nil {
		log.
			WithField("path", t.f.Name()).
			Debug("Failed to open file for file transport subscription")
		return nil, err
	}

	c := make(chan []byte)
	go func() {
		r := bufio.NewScanner(newTailReader(f))

		for r.Scan() {
			entry := &fileEntry{}
			if err := json.
				NewDecoder(bytes.NewBuffer(r.Bytes())).
				Decode(&entry); err != nil {
				log.
					WithField("file", f.Name()).
					WithField("transport", t.Describe()).
					WithField("entry", r.Text()).
					WithError(err).
					Warn("Failed to parse entry in transport file")
			} else {
				log.
					WithField("topic", topic).
					WithField("transport", t.Describe()).
					WithField("entry", entry).
					Debug("Read entry in transport file")
				if entry.Topic == topic {
					c <- entry.Data
				}
			}
		}

		if r.Err() != nil {
			log.
				WithField("topic", topic).
				WithField("transport", t.Describe()).
				WithError(r.Err()).
				Warn("File subscription exited with error")
		}

		log.
			WithField("topic", topic).
			WithField("transport", t.Describe()).
			Debug("Closing subscriber")

		close(c)
	}()

	return &fileSubscription{
		f: f,
		c: c,
	}, nil
}

func (t *FileTransport) Publish(topic string, data []byte) error {
	if t.f == nil {
		return fmt.Errorf("not connected")
	}

	log.
		WithField("topic", topic).
		WithField("transport", t.Describe()).
		Debug("Publishing message to transport")

	defer t.f.Sync()

	return json.NewEncoder(t.f).Encode(&fileEntry{
		Topic: topic,
		Data:  data,
	})
}

func (t *FileTransport) Close() error {
	if t.f == nil {
		return nil
	}

	log.
		WithField("transport", t.Describe()).
		Debug("Closing transport")

	f := t.f
	t.f = nil
	return f.Close()
}

type fileEntry struct {
	Topic string `json:"topic"`
	Data  []byte `json:"data"`
}

type fileSubscription struct {
	f *os.File
	c chan []byte
}

func (s *fileSubscription) Channel() <-chan []byte {
	return s.c
}

func (s *fileSubscription) Close() error {
	return s.f.Close()
}

type tailReader struct {
	f io.Reader
}

func newTailReader(r io.Reader) io.Reader {
	return &tailReader{
		f: r,
	}
}

func (r *tailReader) Read(b []byte) (n int, err error) {
	n, err = r.f.Read(b)

	for n == 0 && err == io.EOF {
		time.Sleep(10 * time.Millisecond)
		n, err = r.f.Read(b)
	}

	return n, err
}
