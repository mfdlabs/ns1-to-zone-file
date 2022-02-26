package zone_file

const ZoneNSHeader string = "\n;\n;  Zone NS records\n;\n\n"
const ZoneRecordHeader string = "\n;\n;  Zone records\n;\n\n"
const ZoneRecordPrefixNoTTL string = "%s					  IN	%s	"
const ZoneRecordPrefix string = "%s				%d	IN	%s	"
const ZoneHeader string = ";\n;  Database file %s.zone for Default Zone scope in zone %s.\n;      Zone version:  %d\n;\n\n"
const ZoneSOARecord string = "@ 					  IN	SOA %s. %s. (\n                        		%d         ; serial number\n                        		%d         ; refresh\n                        		%d         ; retry\n                        		%d         ; expire\n                        		%d         ) ; default ttl\n"
const ZoneNSRecord string = "@ 					  IN	NS	%s.\n"
