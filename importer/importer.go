package importer

import (
	"log"
	"logparse-go/database"
	"logparse-go/logreader"
	"logparse-go/resultcollector"
)

func Import(logFile string, databaseFile string) (*resultcollector.ResultCollector, error) {
	collector := resultcollector.New()

	lines, err := logreader.ParseLogFile(logFile, collector)
	if err != nil {
		return nil, err
	}

	db, err := database.OpenDatabase(databaseFile)
	if err != nil {
		return nil, err
	}

	log.Printf("Starting persistence for %d lines", len(lines))

	batchSize := 200

	for start := 0; start < len(lines); start += batchSize {
		end := start + batchSize
		if end >= len(lines) {
			end = len(lines) - 1
		}

		err = db.BatchPersistLines(lines[start:end], collector)
		if err != nil {
			return nil, err
		}
	}

	err = db.Close()
	if err != nil {
		return nil, err
	}

	//log.Println("--- COLLECTOR DUMP ---\n" + collector.DumpString() + "\n--- COLLECTOR END ---")

	return collector, nil
}
