package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
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

func (v *vault) containsATNA() {
	if strings.Contains(v.auditMessage, `<AuditMessage>`) {
		v.isATNAmessage = true
	} else {
		v.isATNAmessage = false
	}
}
