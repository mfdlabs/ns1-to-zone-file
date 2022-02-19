package main

import (
	"flag"
	"fmt"
	"os"

	"golang.repositories.vmminfra.local/pkg/ns1-to-zone-file/pkg/zone_file"
)

var (
	zoneFilesPath          = flag.String("zone-files-path", "", "A path to the directory that houses the NS1 exported Zone files (JSON Files).")
	dnsZoneFilesOutputPath = flag.String("dns-zone-files-output-path", "", "A path to the directory that will house the DNS Zone files.")
)

func main() {
	// set flag descriptions
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *zoneFilesPath == "" {
		flag.Usage()
		return
	}

	if *dnsZoneFilesOutputPath == "" {
		flag.Usage()
		return
	}

	// if the dns zone files output path doesn't exist, create it forcefully
	if _, err := os.Stat(*dnsZoneFilesOutputPath); os.IsNotExist(err) {
		os.MkdirAll(*dnsZoneFilesOutputPath, 0755)
	}

	zones, err := zone_file.ReadZoneFilesDirectory(*zoneFilesPath)

	if err != nil {
		panic(err)
	}

	// Print the zone length
	fmt.Printf("%d zones in ns1\n", len(zones))

	for _, zone := range zones {
		zone_file.WriteZoneFile(&zone, *dnsZoneFilesOutputPath)
	}
}
