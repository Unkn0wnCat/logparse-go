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

	batchSize := 100

	for start := 0; start < len(lines); start += batchSize {
		err = db.BatchPersistLines(lines[start:(start+batchSize)], collector)
		if err != nil {
			return nil, err
		}
	}

	err = db.Close()
	if err != nil {
		return nil, err
	}

	log.Println("--- COLLECTOR DUMP ---\n" + collector.DumpString() + "\n--- COLLECTOR END ---")

	return collector, nil
}
