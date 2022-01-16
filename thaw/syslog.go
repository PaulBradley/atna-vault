package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jeromer/syslogparser"
	r "github.com/jeromer/syslogparser/rfc3164"
	t "github.com/jeromer/syslogparser/rfc5424"
)

func (v *vault) sysLogHeaderParse() {
	var err error
	var buffer []byte
	var rfc syslogparser.RFC

	v.containsATNA()

	if strings.Contains(v.auditMessage, `Tiani_Cisco_EHR`) {
		start := strings.Index(v.auditMessage, `<`)
		end := strings.Index(v.auditMessage, `<?xml`)
		buffer = []byte(v.auditMessage[start : end-3])
	} else {
		buffer = []byte(v.auditMessage)
	}

	rfc, err = syslogparser.DetectRFC(buffer)
	if err != nil {
		log.Println("WARN:Could Not Detect RFC of message")
	}

	switch rfc {
	case syslogparser.RFC_UNKNOWN:
		log.Println("WARN:RFC Unknown")
	case syslogparser.RFC_3164:
		p := r.NewParser(buffer)
		err = p.Parse()
		if err != nil {
			log.Println("WARN:" + err.Error())
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
			case `content`:
				if !v.isATNAmessage {
					v.sysLogExtractComment(val.(string))
					fmt.Println(`Comment:` + v.sysLogComment)
				}
			}
		}

	case syslogparser.RFC_5424:
		p := t.NewParser(buffer)
		err = p.Parse()
		if err != nil {
			log.Println("WARN:" + err.Error())
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
			case `content`:
				if !v.isATNAmessage {
					v.sysLogExtractComment(val.(string))
					fmt.Println(`Comment:` + v.sysLogComment)
				}
			}
		}
	}

}

func (v *vault) sysLogExtractComment(content string) {
	v.sysLogComment = ""

	e := len(content)

	start := strings.Index(content, `[Comment = `) + 11
	if start > 0 && start < len(content) {
		content = content[start:e]
	} else {
		return
	}

	end := strings.Index(content, `]`)
	if end > 0 && end < len(content) {
		content = content[0:end]
		v.sysLogComment = content
	} else {
		return
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
