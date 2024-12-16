package parser

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net"
	"regexp"
	"strconv"
)

type LogLine struct {
	IP         string
	Timestamp  string
	Method     string
	Path       string
	HttpVesion string
	Status     int32
	Size       int64
	Filename   string
	LineNumber int
	hashCache  *string
}

func (logLine *LogLine) Hash() string {
	if logLine.hashCache != nil {
		return *logLine.hashCache
	}

	hasher := sha512.New()
	hasher.Write([]byte(logLine.IP))
	hasher.Write([]byte(logLine.Timestamp))
	hasher.Write([]byte(logLine.Method))
	hasher.Write([]byte(logLine.Path))
	hasher.Write([]byte(logLine.HttpVesion))
	hasher.Write([]byte(strconv.FormatInt(int64(logLine.Status), 10)))
	hasher.Write([]byte(strconv.FormatInt(logLine.Size, 10)))

	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	//log.Printf("Hashed line: %s", hash)

	logLine.hashCache = &hash

	return hash
}

func ParseLine(line string, fileName string, lineNumber int) (*LogLine, error) {
	regex := regexp.MustCompile("^(?<ip>\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) .* .* \\[(?<timestamp>.*)\\] \"(?<method>[A-Z]+) (?<path>[/A-Za-z0-9$\\-_.+!*'(),;?:@=&%]+) (?<http>HTTP/.*)\" (?<status>\\d{3}|0) (?<size>\\d+|-)$")

	result := regex.FindAllSubmatch([]byte(line), -1)

	if len(result) != 1 {
		return nil, fmt.Errorf("invalid log line: %s", line)
	}

	resultEntry := result[0]

	var err error

	var statusParsed int64
	if string(resultEntry[6]) != "-" {
		statusParsed, err = strconv.ParseInt(string(resultEntry[6]), 10, 32)
		if err != nil {
			return nil, err
		}
	}

	var lengthParsed int64
	if string(resultEntry[7]) != "-" {
		lengthParsed, err = strconv.ParseInt(string(resultEntry[7]), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	if net.ParseIP(string(resultEntry[1])) == nil {
		return nil, fmt.Errorf("invalid ip address: %s", string(resultEntry[1]))
	}

	lineOut := LogLine{
		IP:         string(resultEntry[1]),
		Timestamp:  string(resultEntry[2]),
		Method:     string(resultEntry[3]),
		Path:       string(resultEntry[4]),
		HttpVesion: string(resultEntry[5]),
		Status:     int32(statusParsed),
		Size:       lengthParsed,
		LineNumber: lineNumber,
		Filename:   fileName,
	}

	return &lineOut, nil
}
