package zone_file

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/mfdlabs-grid-development/ns1-to-zone-file/pkg/ns1"
)

func parseRecord(record *ns1.ZoneRecord, zoneName string, defaultTTL int) string {
	recordType := record.Type

	// if the record is NS, SOA or CAA, skip it.
	if recordType == "NS" || recordType == "SOA" || recordType == "CAA" {
		return ""
	}

	// Write the base of the record, if the record matches the zone name, use @ as the record name, else remove the zone name from the record name.
	recordName := record.Domain

	if recordName == zoneName {
		recordName = "@"
	} else {
		recordName = regexp.MustCompile(fmt.Sprintf(".%s", zoneName)).ReplaceAllString(recordName, "")
	}

	// check if recordName is an ipv4
	if regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`).MatchString(recordName) {
		// Warn the user that the record is an ipv4, this is an error and needs to be fixed.
		fmt.Printf("Warning: Record %s is an ipv4, this is an error and needs to be fixed.\n", recordName)
	}

	outRecord := ""

	for _, shortAnswer := range record.ShortAns {

		// The record looks like this `@ OR <record_name>				<ttl (if not default ttl)>	IN	<record_type>	<record_data>\n`

		recordString := ""

		if record.TTL == defaultTTL {
			recordString = fmt.Sprintf(ZoneRecordPrefixNoTTL, recordName, recordType)
		} else {
			recordString = fmt.Sprintf(ZoneRecordPrefix, recordName, record.TTL, record.Type)
		} // We basically only need to append `<data_parsed>\n` to the record string.

		// if the type is TXT, we need to add the quotes around the value and escape the quotes.

		// The record can have multiple "short answers" (A, AAAA, CNAME, MX, NS, PTR, SRV, TXT)
		// We need to iterate through the short answers and append them to the record string.
		// Second case is to search regex for `[a-zA-Z]` in the short answer, and also check if it does not end with a period, OR it is a CNAME and does not end with a period.
		if recordType == "TXT" {
			recordString += fmt.Sprintf("\"%s\"\n", shortAnswer)
		} else if (regexp.MustCompile(`[a-zA-Z]`).MatchString(shortAnswer) && shortAnswer[len(shortAnswer)-1:] != ".") || (recordType == "CNAME" && shortAnswer[len(shortAnswer)-1:] != ".") {
			recordString += fmt.Sprintf("%s.\n", shortAnswer)
		} else {
			recordString += fmt.Sprintf("%s\n", shortAnswer)
		}

		outRecord += recordString
	}

	return outRecord
}

func getBaseZoneFileHeader(zone *ns1.Zone) string {
	zoneName := zone.Zone
	zoneVersion := zone.Serial
	soa := zone.PrimaryMaster
	// The hostmaster is the email address of the person responsible for the zone
	// The SOA will be the hostmaster but replace the @ with a .
	soaAuthority := zone.Hostmaster

	// Replace the @ with a .
	soaAuthority = regexp.MustCompile(`@`).ReplaceAllString(soaAuthority, ".")

	// Header in the format of:
	// ;\n;  Database file <zone_name>.zone for Default Zone scope in zone <zone_name>.\n;      Zone version:  <serial>\n;\n\n
	// SOA record in the format of:
	// @                        IN SOA <primary_master>. <hostmaster>. (\n                        		<serial>         ; serial number\n                        		<refresh>         ; refresh\n                        		<retry>         ; retry\n                        		<expire>         ; expire\n                        		<ttl>         ) ; default ttl\n

	header := fmt.Sprintf(ZoneHeader, zoneName, zoneName, zoneVersion)
	soaRecord := fmt.Sprintf(ZoneSOARecord, soa, soaAuthority, zoneVersion, zone.Refresh, zone.Retry, zone.Expiry, zone.TTL)

	// Append these.
	header += soaRecord

	// Append NS header.
	header += ZoneNSHeader

	// Iterate zone NS records.
	for _, ns := range zone.DNSServers {
		header += fmt.Sprintf(ZoneNSRecord, ns)
	}

	// Append record header.
	header += ZoneRecordHeader

	// Iterate zone records.
	for _, record := range zone.Records {
		header += parseRecord(record, zoneName, zone.TTL)
	}

	return header
}

func WriteZoneFile(zone *ns1.Zone, outputDir string) error {
	zoneFileHeader := getBaseZoneFileHeader(zone)

	// filename is outputDir/<zone_name>.zone
	filename := fmt.Sprintf("%s/%s.zone", outputDir, zone.Zone)

	// Write the zone file.
	return ioutil.WriteFile(filename, []byte(zoneFileHeader), 0644)
}
