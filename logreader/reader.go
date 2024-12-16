package logreader

import (
	"bufio"
	"fmt"
	"logparse-go/parser"
	"logparse-go/resultcollector"
	"os"
)

func ParseLogFile(fileName string, collector *resultcollector.ResultCollector) ([]parser.LogLine, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		collector.Add(resultcollector.Result{
			Scope:    resultcollector.ScopeFile,
			Filename: fileName,
			Line:     0,
			Message:  err.Error(),
		})

		return nil, err
	}

	var logLines []parser.LogLine

	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() {
		line++

		logLine, err := parser.ParseLine(scanner.Text(), fileName, line)
		if err != nil {
			collector.Add(resultcollector.Result{
				Scope:    resultcollector.ScopeLine,
				Filename: fileName,
				Line:     line,
				Message:  fmt.Sprintf("Could not parse line: %s", err.Error()),
			})
			continue
		}

		logLines = append(logLines, *logLine)
	}

	return logLines, nil
}
