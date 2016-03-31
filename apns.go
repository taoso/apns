package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const api string = "https://api.push.apple.com/3/device/"

var topic string

func makeReq(token string, alert string) (*http.Request, error) {
	msg := "{\"aps\":{\"alert\":\"" + alert + "\",\"sound\":\"default\"}}"

	req, err := http.NewRequest(
		"POST",
		api+token,
		bytes.NewBufferString(msg),
	)

	if err == nil {
		req.Header.Add("apns-topic", topic)
		fmt.Println("topic:" + topic)
	}

	return req, err
}

func makeClient(keyFile string) (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(keyFile, keyFile)
	if err != nil {
		return nil, err
	}

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}
	if len(x509Cert.Subject.CommonName) > 0 {
		commonName := x509Cert.Subject.CommonName
		parts := strings.Split(commonName, ": ")
		topic = parts[1]
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	// tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{
		TLSClientConfig:    tlsConfig,
		DisableCompression: true,
	}
	err = http2.ConfigureTransport(transport)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Transport: transport}

	return client, err
}

func push(client *http.Client, token string, msg string) (string, error) {
	req, err := makeReq(token, msg)
	if err != nil {
		return token, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return token, err
	}

	if resp.StatusCode != 200 {
		contents, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return token, err
		}
		return token, errors.New(string(contents))
	}

	return token, nil
}

func main() {
	var m = flag.String("m", "", "apns alert")
	var c = flag.String("c", "", "cert path")
	flag.Parse()
	msg := *m
	keyFile := *c

	scanner := bufio.NewScanner(os.Stdin)
	client, err := makeClient(keyFile)
	if err != nil {
		os.Exit(1)
	}
	for scanner.Scan() {
		token := scanner.Text()
		token, err := push(client, token, msg)
		f := "+"
		e := "|"
		if err != nil {
			f = "-"
			e = e + err.Error()
		}
		fmt.Println(f + token + e)
	}
}
