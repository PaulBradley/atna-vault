package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/oklog/ulid"
)

func main() {
	var av = vault{}
	av.getCommandLineOptions()
	av.openLogFile()
	av.receiveAuditMessage()

	if strings.Contains(av.auditMessage, "AuditMessage") {
		av.generateObjectName()
		av.generateObjectPrefix()
		av.gzipAuditMessage()
		if *av.storeLocally {
			av.writeFileLocally()
		}
	} else {
		log.Println("INFO:" + av.auditMessage)
	}
	av.logFile.Close()
}

func (v *vault) generateObjectName() {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	v.outputFilename = ulid.MustNew(ulid.Timestamp(t), entropy).String() + ".ATNA"
}

func (v *vault) generateObjectPrefix() {
	if *v.storeLocally {
		v.outputFilenamePrefix = v.outputFilename[0:5] + "-"
	} else {
		v.outputFilenamePrefix = v.outputFilename[0:5] + "/"
	}
}

func (v *vault) getCommandLineOptions() {
	v.storeLocally = flag.Bool("store-locally", false, "save files locally")
	flag.Parse()
}

func (v *vault) gzipAuditMessage() {
	zw := gzip.NewWriter(&v.gzBuffer)

	_, v.err = zw.Write([]byte(v.auditMessage))
	if v.err != nil {
		log.Println("WARN:" + v.err.Error())
		v.gzValid = false
		return
	}

	if err := zw.Close(); err != nil {
		log.Println("WARN:" + v.err.Error())
		v.gzValid = false
		return
	}

	v.gzValid = true
}

func (v *vault) openLogFile() {
	v.logFile, v.err = os.OpenFile("atna-vault.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if v.err != nil {
		log.Fatal("could not open log file : atna-vault.log")
	}

	log.SetOutput(v.logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func (v *vault) receiveAuditMessage() {
	v.scanner = bufio.NewScanner(os.Stdin)
	for v.scanner.Scan() {
		v.auditMessage = v.auditMessage + v.scanner.Text()
	}

	if err := v.scanner.Err(); err != nil {
		log.Println(err.Error())
	}
}

func (v *vault) writeFileLocally() {
	v.outputFile, v.err = os.OpenFile(v.outputFilenamePrefix+v.outputFilename,
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if v.err != nil {
		log.Println("WARN:" + v.err.Error())
	}
	defer v.outputFile.Close()

	if v.gzValid {
		v.outputFile.Write(v.gzBuffer.Bytes())
	} else {
		v.outputFile.Write([]byte(v.auditMessage))
	}
}
