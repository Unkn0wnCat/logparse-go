package parser

import (
	"fmt"
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
}

func ParseLine(line string) (*LogLine, error) {
	regex := regexp.MustCompile("^(?<ip>\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) .* .* \\[(?<timestamp>.*)\\] \"(?<method>[A-Z]+) (?<path>[/A-Za-z0-9$\\-_.+!*'(),;?:@=&]+) (?<http>HTTP/.*)\" (?<status>\\d{3}|0) (?<size>\\d+|-)$")

	result := regex.FindAllSubmatch([]byte(line), -1)

	for _, part := range result {
		for i, part2 := range part {
			fmt.Printf("%d: %s\n", i, string(part2))
		}
		fmt.Println()
	}

	if len(result) != 1 {
		return nil, fmt.Errorf("invalid log line: %s", line)
	}

	resultEntry := result[0]

	statusParsed, err := strconv.ParseInt(string(resultEntry[6]), 10, 32)
	if err != nil {
		return nil, err
	}

	lengthParsed, err := strconv.ParseInt(string(resultEntry[7]), 10, 64)
	if err != nil {
		return nil, err
	}

	lineOut := LogLine{
		IP:         string(resultEntry[1]),
		Timestamp:  string(resultEntry[2]),
		Method:     string(resultEntry[3]),
		Path:       string(resultEntry[4]),
		HttpVesion: string(resultEntry[5]),
		Status:     int32(statusParsed),
		Size:       lengthParsed,
	}

	return &lineOut, nil
}
