package main

import (
	"bufio"
	"bytes"
	"os"
	"time"
)

type vault struct {
	auditMessage         string         // holds the plain text audit message
	err                  error          // holds the last error message
	gzBuffer             bytes.Buffer   // hold the gzip version of the audit message
	gzValid              bool           // true/false that the gzip worked
	logFile              *os.File       // pointer to the application logfile
	outputFile           *os.File       // pointer to the file being generated locally
	outputFilename       string         // hold the unique filename for the message
	outputFilenamePrefix string         // holds the filename prefix
	scanner              *bufio.Scanner // holds the buffered audit message
	testing              *bool          // true/false in developer testing mode
	aws_region           *string
	aws_s3_bucket        *string
	sysLogPriority       int
	sysLogTimestamp      time.Time
	sysLogHostName       string
	sysLogApplication    string
	sysLogFacility       int
	syslogSeverity       int
}
