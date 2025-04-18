// Copyright 2025 SGNL.ai, Inc.
package awss3

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
)

const FileTypeCSV = "csv"

// TODO: Clean this up by decoupling the attribute value conversion logic from the CSV parsing logic.
// CSVBytesToObject converts a CSV byte array to an array of objects.
func CSVBytesToPage(
	data *[]byte,
	start int64,
	pageSize int64,
	attrConfig []*framework.AttributeConfig,
) ([]map[string]any, bool, error) {
	csvData := csv.NewReader(bytes.NewReader(*data))

	// Read all the CSV data
	records, err := csvData.ReadAll()
	if err != nil {
		return nil, false, fmt.Errorf("failed to read CSV data: %v", err)
	}

	count := len(records)
	if count == 0 {
		return nil, false, fmt.Errorf("no data found in the CSV file")
	}

	objects := make([]map[string]any, 0, pageSize)
	if count == 1 {
		return objects, false, nil
	}

	// Convert CSV data to a slice of maps
	headers := records[0]
	headerToAttributeConfig := headerToAttributeConfig(headers, attrConfig)

	end := min(start+pageSize, int64(count))
	hasNext := end < int64(count)

	for _, record := range records[start:end] {
		row := make(map[string]interface{})

		for i, value := range record {
			attrConfig, found := headerToAttributeConfig[headers[i]]
			if !found {
				// If the value is a complex list of attributes, unmarshal it.
				// "[{\"primary\": true, \"alias\": \"Klein Luis\"},{\"alias\": \"Cline Luis\", \"primary\": false}]"
				if strings.HasPrefix(value, "[{") && strings.HasSuffix(value, "}]") {
					var childObj []map[string]any
					if err := json.Unmarshal([]byte(value), &childObj); err != nil {
						return nil, false, fmt.Errorf(
							`failed to unmarshal the value: "%v" in row: %d, column: %s`,
							value, i, headers[i],
						)
					}

					childArray := make([]any, 0, len(childObj))
					for _, obj := range childObj {
						childArray = append(childArray, obj)
					}

					row[headers[i]] = childArray

					continue
				}

				row[headers[i]] = value

				continue
			}

			// If attributeConfig is present, based on the attribute type, convert the value to a number
			switch attrConfig.Type {
			case framework.AttributeTypeInt64, framework.AttributeTypeDouble:
				floatValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, false, fmt.Errorf(
						`failed to convert the value: "%v" in row: %d, column: %s to a number`,
						value, i, headers[i],
					)
				}

				row[headers[i]] = floatValue
			default:
				row[headers[i]] = value
			}
		}

		objects = append(objects, row)
	}

	return objects, hasNext, nil
}

func headerToAttributeConfig(
	headers []string,
	attrConfig []*framework.AttributeConfig,
) map[string]framework.AttributeConfig {
	attrExternalIDToAttrConfig := make(map[string]framework.AttributeConfig, len(attrConfig))

	for _, attr := range attrConfig {
		if attr != nil {
			attrExternalIDToAttrConfig[attr.ExternalId] = *attr
		}
	}

	headerToAttrType := make(map[string]framework.AttributeConfig, len(headers))

	for _, header := range headers {
		attrConfig, found := attrExternalIDToAttrConfig[header]
		if !found {
			continue
		}

		headerToAttrType[header] = attrConfig
	}

	return headerToAttrType
}
