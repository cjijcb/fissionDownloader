package main

import (
	"errors"
	"io"
	"net/http"
	"fmt"
)

type (
	NetConn struct {
		Client  *http.Client
		Request *http.Request
	}
)

func buildNetConn(ct *http.Client, req *http.Request) *NetConn {
	nc := &NetConn{
		Client:  ct,
		Request: req,
	}

	return nc
}

func buildReq(method string, rawURL string) *http.Request {
	req, err := http.NewRequest(method, rawURL, nil)
	doHandle(err)
	req.Proto = "http/1.1"
	req.ProtoMajor = 1
	req.ProtoMinor = 1
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "fssn/1.0.0")
	return req
}

func buildClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns: 0,
	}
	ct := &http.Client{Transport: tr}
	return ct
}


func (nc *NetConn) SetScope(br ByteRange) {

	nc.Request.Header.Set("Range", buildRangeHeader(br))

}

func (nc *NetConn) DataCast() io.ReadCloser {

	resp, err := nc.Client.Do(nc.Request)
	doHandle(err)

	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		doHandle(errors.New(resp.Status))
	}

	return resp.Body
}

func buildRangeHeader(br ByteRange) string {

	if br.Start == UnknownSize || br.End == UnknownSize {
		return "none"
	}

	rangeStart := br.Start + br.Offset
	rangeEnd := br.End

	if rangeStart > rangeEnd {
		rangeStart = rangeEnd
	}

	return fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd)

}

func GetHeaders(rawURL string) (http.Header, int64) {
	ct := &http.Client{}
	req, err := http.NewRequest(http.MethodHead, rawURL, nil)
	doHandle(err)
	req.Header.Set("User-Agent", "fssn/1.0.0")
	resp, err := ct.Do(req)
	doHandle(err)
	return resp.Header, resp.ContentLength
}
