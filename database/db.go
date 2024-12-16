package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"logparse-go/parser"
	"logparse-go/resultcollector"
	"strconv"
	"strings"
)

type Database struct {
	db *sql.DB
}

func OpenDatabase(fileName string) (*Database, error) {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	database := &Database{db: db}

	isBootstrapped, err := database.isBootstrapped()
	if err != nil {
		return nil, err
	}

	if !isBootstrapped {
		err = database.bootstrap()
		if err != nil {
			return nil, err
		}
	}

	return database, nil
}

func (inst *Database) Close() error {
	return inst.db.Close()
}

func (inst *Database) isBootstrapped() (bool, error) {
	rows, err := inst.db.Query(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='logs';`)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return false, err
		}
	}

	return count > 0, nil
}

func (inst *Database) bootstrap() error {
	sqlStmt := `
		create table logs (
			hash text not null primary key, 
			ip text not null,
			timestamp datetime not null,
			method text not null,
			path text not null,
			httpVersion text not null,
			status numeric not null,
			size numeric not null
		);
	`

	_, err := inst.db.Exec(sqlStmt)
	return err
}

func (inst *Database) BatchPersistLines(lines []parser.LogLine, collector *resultcollector.ResultCollector) error {
	log.Printf("BatchPersistLines starting for %d lines", len(lines))

	tx, err := inst.db.Begin()
	if err != nil {
		return err
	}

	log.Printf("Collecting hashes...")

	var hashes []string
	for _, line := range lines {
		hashes = append(hashes, line.Hash())
	}

	log.Printf("Preparing dedupe...")

	dedupeCheck, err := tx.Prepare(`SELECT COUNT(*) FROM logs WHERE hash = ?`)
	if err != nil {
		return err
	}

	log.Printf("Querying dedupe...")

	var dupeHashes []string

	for _, hash := range hashes {
		dedupeRow, err := dedupeCheck.Query(hash)
		if err != nil {
			return err
		}

		var count int
		for dedupeRow.Next() {
			if err := dedupeRow.Scan(&count); err != nil {
				return err
			}
		}

		if count > 0 {
			dupeHashes = append(dupeHashes, hash)
		}

		dedupeRow.Close()
	}

	log.Printf("Deduplicating...")

	var dedupedLines []parser.LogLine

	for _, line := range lines {
		hash := line.Hash()

		isDupe := false
		for _, dupeHash := range dupeHashes {
			if hash == dupeHash {
				isDupe = true
				break
			}
		}

		if !isDupe {
			dedupedLines = append(dedupedLines, line)
			continue
		}

		// Is Duplicate.
		log.Printf("Found duplicate line: %v", line)
		collector.Add(resultcollector.Result{
			Scope:    resultcollector.ScopeLine,
			Filename: line.Filename,
			Line:     line.LineNumber,
			Message:  "Already in database",
		})
	}

	insertStmt, err := tx.Prepare(`
		INSERT INTO logs (hash, ip, timestamp, method, path, httpVersion, status, size) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		return err
	}

	for _, line := range dedupedLines {
		_, err := insertStmt.Exec(line.Hash(), line.IP, line.Timestamp, line.Method, line.Path, line.HttpVesion, line.Status, line.Size)
		if err != nil {
			collector.Add(resultcollector.Result{
				Scope:    resultcollector.ScopeLine,
				Filename: line.Filename,
				Line:     line.LineNumber,
				Message:  fmt.Sprintf("Error inserting line: %v", err),
			})
			continue
		}

		collector.Add(resultcollector.Result{
			Scope:    resultcollector.ScopeLine,
			Filename: line.Filename,
			Line:     line.LineNumber,
			Success:  true,
			Message:  "Line inserted",
		})
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (inst *Database) Query(ipFilter string, timestampFilterFrom string, timestampFilterTo string) ([]parser.LogLine, error) {
	sql := `SELECT hash, ip, timestamp, method, path, httpVersion, status, size FROM logs WHERE`
	var params []any

	if ipFilter != "" {
		sql += ` ip = ?`
		params = append(params, ipFilter)
	}
	if timestampFilterFrom != "" {
		if len(params) > 0 {
			sql += ` AND`
		}
		sql += ` timestamp > ?`
		params = append(params, timestampFilterFrom)
	}
	if timestampFilterTo != "" {
		if len(params) > 0 {
			sql += ` AND`
		}
		sql += ` timestamp < ?`
		params = append(params, timestampFilterTo)
	}

	if len(params) == 0 {
		sql += ` 1 = 1`
	}

	sql += ` ORDER BY timestamp DESC`

	log.Printf("%s", sql)

	query, err := inst.db.Prepare(sql)
	if err != nil {
		return nil, fmt.Errorf("Error preparing query: %v", err)
	}

	defer query.Close()

	log.Printf("Query prepared")

	rows, err := query.Query(params...)
	if err != nil {
		return nil, fmt.Errorf("Error querying logs: %v", err)
	}

	log.Printf("Query ok")

	defer rows.Close()

	var lines []parser.LogLine
	for rows.Next() {
		log.Printf("Processing row %d", len(lines)+1)
		var hash string
		var ip string
		var timestamp string
		var method string
		var path string
		var httpVersion string
		var status int
		var size int

		err = rows.Scan(&hash, &ip, &timestamp, &method, &path, &httpVersion, &status, &size)
		if err != nil {
			return nil, fmt.Errorf("Reading row: %v", err)
		}

		lines = append(lines, parser.LogLine{
			IP:         ip,
			Timestamp:  timestamp,
			Method:     method,
			Path:       path,
			HttpVesion: httpVersion,
			Status:     int32(status),
			Size:       int64(size),
			Filename:   "DATABASE",
			LineNumber: 0,
		})

		log.Printf("Line %d read", len(lines))
	}

	return lines, nil
}

type IPCount struct {
	IP    string
	Count int
}

func (inst *Database) QueryIPCounts(ipFilter string, timestampFilterFrom string, timestampFilterTo string) ([]IPCount, error) {
	sql := `SELECT ip, COUNT(*) AS count FROM logs WHERE`
	var params []any

	if ipFilter != "" {
		sql += ` ip = ?`
		params = append(params, ipFilter)
	}
	if timestampFilterFrom != "" {
		if len(params) > 0 {
			sql += ` AND`
		}
		sql += ` timestamp > ?`
		params = append(params, timestampFilterFrom)
	}
	if timestampFilterTo != "" {
		if len(params) > 0 {
			sql += ` AND`
		}
		sql += ` timestamp < ?`
		params = append(params, timestampFilterTo)
	}

	if len(params) == 0 {
		sql += ` 1 = 1`
	}

	sql += ` GROUP BY ip ORDER BY count DESC`

	log.Printf("%s", sql)

	query, err := inst.db.Prepare(sql)
	if err != nil {
		return nil, fmt.Errorf("Error preparing query: %v", err)
	}

	defer query.Close()

	log.Printf("Query prepared")

	rows, err := query.Query(params...)
	if err != nil {
		return nil, fmt.Errorf("Error querying logs: %v", err)
	}

	log.Printf("Query ok")

	defer rows.Close()

	var lines []IPCount
	for rows.Next() {
		log.Printf("Processing row %d", len(lines)+1)
		var ip string
		var count int

		err = rows.Scan(&ip, &count)
		if err != nil {
			return nil, fmt.Errorf("Reading row: %v", err)
		}

		lines = append(lines, IPCount{
			IP:    ip,
			Count: count,
		})

		log.Printf("Line %d read", len(lines))
	}

	return lines, nil
}

type StatusCount struct {
	Status int
	Count  int
}

func (inst *Database) QueryStatusCounts(ipFilter string, timestampFilterFrom string, timestampFilterTo string, codes string) ([]StatusCount, error) {
	sql := `SELECT status, COUNT(*) AS count FROM logs WHERE`
	var params []any

	if ipFilter != "" {
		sql += ` ip = ?`
		params = append(params, ipFilter)
	}
	if timestampFilterFrom != "" {
		if len(params) > 0 {
			sql += ` AND`
		}
		sql += ` timestamp > ?`
		params = append(params, timestampFilterFrom)
	}
	if timestampFilterTo != "" {
		if len(params) > 0 {
			sql += ` AND`
		}
		sql += ` timestamp < ?`
		params = append(params, timestampFilterTo)
	}

	if codes != "" {
		codeSql := ""

		if len(params) > 0 {
			sql += ` AND`
		}

		codeList := strings.Split(codes, ",")
		for _, code := range codeList {
			code, err := strconv.Atoi(code)
			if err != nil {
				return nil, fmt.Errorf("Error converting code to int: %v", err)
			}

			if codeSql != "" {
				codeSql += " OR "
			}
			codeSql += "status = ?"
			params = append(params, code)
		}

		sql += " (" + codeSql + ")"
	}

	if len(params) == 0 {
		sql += ` 1 = 1`
	}

	sql += ` GROUP BY status ORDER BY count DESC`

	log.Printf("%s", sql)

	query, err := inst.db.Prepare(sql)
	if err != nil {
		return nil, fmt.Errorf("Error preparing query: %v", err)
	}

	defer query.Close()

	log.Printf("Query prepared")

	rows, err := query.Query(params...)
	if err != nil {
		return nil, fmt.Errorf("Error querying logs: %v", err)
	}

	log.Printf("Query ok")

	defer rows.Close()

	var lines []StatusCount
	for rows.Next() {
		log.Printf("Processing row %d", len(lines)+1)
		var status int
		var count int

		err = rows.Scan(&status, &count)
		if err != nil {
			return nil, fmt.Errorf("Reading row: %v", err)
		}

		lines = append(lines, StatusCount{
			Status: status,
			Count:  count,
		})

		log.Printf("Line %d read", len(lines))
	}

	return lines, nil
}
