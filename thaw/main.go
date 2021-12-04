package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	t "github.com/jeromer/syslogparser/rfc5424"
)

func main() {
	var av = vault{}
	av.getCommandLineOptions()
	av.openLogFile()
	//av.getMessageObjectKeysFromS3()

	if *av.testing {
		av.receiveTestMessage()
		av.sysLogHeaderParse()
		av.sysLogHeaderResults()
	}
	av.logFile.Close()
}

func (v *vault) getCommandLineOptions() {
	v.testing = flag.Bool("testing", false, "local developer testing mode")
	v.aws_region = flag.String("region", "", "AWS region i.e. eu-west-2")
	v.aws_s3_bucket = flag.String("bucketname", "", "AWS S3 bucket name")
	flag.Parse()
}

func (v *vault) openLogFile() {
	v.logFile, v.err = os.OpenFile("thaw.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if v.err != nil {
		log.Fatal("could not open log file : thaw.log")
	}

	log.SetOutput(v.logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func (v *vault) receiveTestMessage() {
	v.scanner = bufio.NewScanner(os.Stdin)
	for v.scanner.Scan() {
		v.auditMessage = v.auditMessage + v.scanner.Text()
	}

	if err := v.scanner.Err(); err != nil {
		log.Println(err.Error())
	}
}

func (v *vault) sysLogHeaderParse() {
	// extract the syslog header from the audit message
	start := strings.Index(v.auditMessage, `<`)
	end := strings.Index(v.auditMessage, `<?xml`)
	buffer := []byte(v.auditMessage[start : end-3])

	p := t.NewParser(buffer)
	v.err = p.Parse()
	if v.err != nil {
		log.Println("WARN:" + v.err.Error())
	}

	for key, val := range p.Dump() {
		switch key {
		case `priority`:
			v.sysLogPriority = val.(int)
		case `timestamp`:
			v.sysLogTimestamp = val.(time.Time)
		case `hostname`:
			v.sysLogHostName = val.(string)
		case `app_name`:
			v.sysLogApplication = val.(string)
		case `facility`:
			v.sysLogFacility = val.(int)
		case `severity`:
			v.syslogSeverity = val.(int)
		}
	}
}

func (v *vault) sysLogHeaderResults() {
	fmt.Println(v.sysLogPriority)
	fmt.Println(v.sysLogTimestamp)
	fmt.Println(v.sysLogApplication)
	fmt.Println(v.sysLogHostName)
	fmt.Println(v.sysLogFacility)
	fmt.Println(v.syslogSeverity)
}
