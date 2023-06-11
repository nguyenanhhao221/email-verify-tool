package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	log.Println("Welcome to the email check tool")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("domain,hasMX,hasSPF,sprRecord,hasDMARC,dmarRecord")
	for scanner.Scan() {
		checkDomain(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error: could not read from input:", err)
	}

}

func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: failed to look up MX %v", err)
	}
	if len(mxRecords) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error: failed to look up TXT %v", err)
	}
	for _, txtRecord := range txtRecords {
		if strings.HasPrefix(txtRecord, "v=spf1") {
			hasSPF = true
			spfRecord = txtRecord
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: failed to look up dmarc records %v", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}
	log.Printf("%v, %v,%v,%v,%v,%v", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
