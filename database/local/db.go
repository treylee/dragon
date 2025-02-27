package local

import (
	"database/sql"
	"fmt"
	"os"
	 "gdragon/internal/metrics"

	"github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
)

var logger = logrus.New()


func SetupDatabase(testID string) (*sql.DB, error) {

	dir := "test_data"
	dbFile := fmt.Sprintf("%s/%s.db", dir, testID)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logrus.Errorf("Failed to create directory %s: %v", dir, err)
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		logrus.Errorf("Error opening database file: %v", err)
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS test_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		test_id TEXT NOT NULL,
		requests INTEGER,
		failed_requests INTEGER,
		error_rate FLOAT,
		p50_response_time INTEGER,
		p95_response_time INTEGER,
		p99_response_time INTEGER,
		request_per_second INTEGER,
		avg_response_time FLOAT,
		max_response_time INTEGER,
		cpu_usage FLOAT,
		memory_usage FLOAT,
		test_duration INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		logrus.Errorf("Error creating table: %v", err)
		return nil, err
	}

	logrus.Infof("Successfully opened database file=%s testId=%s", dbFile, testID)

	return db, nil
}

func SaveTestResults(db *sql.DB, metrics *metrics.TestMetrics) error {
	insertQuery := `
	INSERT INTO test_results (test_id, requests, failed_requests, error_rate, p50_response_time, p95_response_time, p99_response_time, request_per_second, avg_response_time, max_response_time, cpu_usage, memory_usage, test_duration)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err := db.Exec(insertQuery, metrics.TestID, metrics.Requests, metrics.FailedRequests, metrics.ErrorRate, metrics.P50ResponseTime, metrics.P95ResponseTime, metrics.P99ResponseTime, metrics.RequestPerSecond, metrics.AvgResponseTime, metrics.MaxResponseTime, metrics.CpuUsage, metrics.MemoryUsage, metrics.TestDuration)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"testId": metrics.TestID,
			"error":  err.Error(),
		}).Error("Error inserting test results")
		return err
	}
	logger.WithFields(logrus.Fields{
		"testId":            metrics.TestID,
		"requests":          metrics.Requests,
		"failedRequests":    metrics.FailedRequests,
		"errorRate":         metrics.ErrorRate,
		"p50ResponseTime":   metrics.P50ResponseTime,
		"p95ResponseTime":   metrics.P95ResponseTime,
		"p99ResponseTime":   metrics.P99ResponseTime,
		"requestsPerSecond": metrics.RequestPerSecond,
		"avgResponseTime":   metrics.AvgResponseTime,
		"maxResponseTime":   metrics.MaxResponseTime,
		"cpuUsage":          metrics.CpuUsage,
		"memoryUsage":       metrics.MemoryUsage,
		"testDuration":      metrics.TestDuration,
	}).Info("Test results saved successfully")

	return nil
}

func GetTestResults(testId string) ([]metrics.TestMetrics, error) {
	db, err := SetupDatabase(testId)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"testId": testId,
		}).Error("Error setting up database")
		return nil, err
	}
	defer db.Close()

	selectQuery := `SELECT * FROM test_results WHERE test_id = ?;`
	rows, err := db.Query(selectQuery, testId)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"testId": testId,
			"error":  err.Error(),
		}).Error("Error querying test results")
		return nil, err
	}
	defer rows.Close()

	var results []metrics.TestMetrics
	for rows.Next() {
		var result metrics.TestMetrics
		//err := rows.Scan(&result.ID, &result.TestID, &result.Requests, &result.FailedRequests, &result.ErrorRate, &result.P50ResponseTime, &result.P95ResponseTime, &result.P99ResponseTime, &result.RequestPerSecond, &result.AvgResponseTime, &result.MaxResponseTime, &result.CpuUsage, &result.MemoryUsage, &result.TestDuration, &result.CreatedAt)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"testId": testId,
				"error":  err.Error(),
			}).Error("Error scanning row")
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

