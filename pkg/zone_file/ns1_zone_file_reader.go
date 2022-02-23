package zone_file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mfdlabs-grid-development/ns1-to-zone-file/pkg/ns1"
)

func ReadZoneFilesDirectory(dir string) ([]ns1.Zone, error) {
	// First check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %s\n", dir)
		return nil, err
	}

	// Read all the files in the directory
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		fmt.Printf("Error reading directory: %s\n", dir)
		return nil, err
	}

	// Return an error if there are no files in the directory
	if len(files) == 0 {
		fmt.Printf("No files in directory: %s\n", dir)
		return nil, fmt.Errorf("no files in directory: %s", dir)
	}

	// Return an error if there are no json files in the directory
	var jsonFiles []os.FileInfo

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if file.Name()[len(file.Name())-5:] == ".json" {
			jsonFiles = append(jsonFiles, file)
		}
	}

	if len(jsonFiles) == 0 {
		fmt.Printf("No json files in directory: %s\n", dir)
		return nil, fmt.Errorf("no json files in directory: %s", dir)
	}

	// Create a slice to hold all the zones
	zones := make([]ns1.Zone, 0)

	// Iterate over the files
	for _, file := range files {
		// Only go pass if the file is a JSON file
		if file.Name()[len(file.Name())-5:] == ".json" {
			// Read the file
			fileBytes, err := ioutil.ReadFile(dir + "/" + file.Name())

			if err != nil {
				fmt.Printf("Error reading file: %s\n", file.Name())
				return nil, err
			}

			// Create a new zone
			zone := ns1.Zone{}

			// Unmarshal the file into the zone
			err = json.Unmarshal(fileBytes, &zone)

			if err != nil {
				fmt.Printf("Error unmarshalling file: %s\n", file.Name())
				return nil, err
			}

			// Add the zone to the slice
			zones = append(zones, zone)
		}
	}

	return zones, nil
}
