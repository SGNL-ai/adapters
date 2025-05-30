// Copyright 2025 SGNL.ai, Inc.
package awss3

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
)

const (
	FileTypeCSV = "csv"
	// StreamingChunkSize defines how many bytes to read per S3 range request.
	StreamingChunkSize = 1024 * 1024 // 1MB
)

func CSVHeaders(headerChunk *[]byte) ([]string, error) {
	if headerChunk == nil || len(*headerChunk) == 0 {
		return nil, fmt.Errorf("CSV header is empty or missing")
	}

	csvData := csv.NewReader(bytes.NewReader(*headerChunk))

	headers, err := csvData.Read()
	if err != nil {
		return nil, fmt.Errorf("CSV file format is invalid or corrupted: %v", err)
	}

	return headers, nil
}

func StreamingCSVToPage(
	ctx context.Context,
	handler *S3Handler,
	bucket, key string,
	fileSize int64,
	headers []string,
	start int64,
	pageSize int64,
	attrConfig []*framework.AttributeConfig,
) ([]map[string]any, bool, error) {
	objects := make([]map[string]any, 0, pageSize)
	headerToAttributeConfig := headerToAttributeConfig(headers, attrConfig)

	var currentRow, currentPos, collectedRows int64

	targetEndRow := start + pageSize

	for (currentPos < fileSize) && (collectedRows < pageSize) {
		endPos := currentPos + StreamingChunkSize - 1
		if endPos >= fileSize {
			endPos = fileSize - 1
		}

		chunkData, err := handler.GetFileRange(ctx, bucket, key, currentPos, endPos)
		if err != nil {
			return nil, false, fmt.Errorf("unable to read CSV file data: %v", err)
		}

		if chunkData == nil {
			return nil, false, fmt.Errorf("received empty response from S3 file range request")
		}

		chunkObjects, nextRow, nextBytePos, err := processCSVChunk(
			*chunkData,
			headers,
			headerToAttributeConfig,
			currentRow,
			start,
			targetEndRow,
			currentPos,
		)
		if err != nil {
			return nil, false, fmt.Errorf("CSV file processing failed: %v", err)
		}

		objects = append(objects, chunkObjects...)
		collectedRows += int64(len(chunkObjects))
		currentRow = nextRow

		if nextBytePos <= currentPos {
			return nil, false, fmt.Errorf("CSV file contains formatting issues that prevent processing from continuing")
		}

		currentPos = nextBytePos

		if collectedRows >= pageSize {
			break
		}
	}

	hasNext := currentRow <= start+pageSize || currentPos < fileSize
	if !hasNext && collectedRows == pageSize {
		hasNext = currentPos < fileSize
	}

	return objects, hasNext, nil
}

func processCSVChunk(
	chunkData []byte,
	headers []string,
	headerToAttributeConfig map[string]framework.AttributeConfig,
	startRowNum int64,
	targetStartRow int64,
	targetEndRow int64,
	chunkStartPos int64,
) ([]map[string]any, int64, int64, error) {
	lastNewlineIndex := bytes.LastIndex(chunkData, []byte("\n"))
	if lastNewlineIndex == -1 {
		return nil, startRowNum, chunkStartPos,
			fmt.Errorf("CSV file contains a single row larger than %d MB", StreamingChunkSize/(1024*1024))
	}

	completeChunk := chunkData[:lastNewlineIndex+1]
	nextBytePos := chunkStartPos + int64(lastNewlineIndex) + 1
	currentRowNum := startRowNum
	csvReader := csv.NewReader(bytes.NewReader(completeChunk))

	var objects []map[string]any

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return nil, currentRowNum, nextBytePos, fmt.Errorf("CSV file format is invalid or corrupted: %v", err)
		}

		// Skip header row (row 0) completely
		if currentRowNum == 0 {
			currentRowNum++

			continue
		}

		// dataRowNum := currentRowNum
		if currentRowNum >= targetStartRow && currentRowNum < targetEndRow {
			row := make(map[string]interface{})

			for i, value := range record {
				if i >= len(headers) {
					continue
				}

				headerName := headers[i]
				attrConfig, found := headerToAttributeConfig[headerName]

				if !found {
					if strings.HasPrefix(value, "[{") && strings.HasSuffix(value, "}]") {
						var childObj []map[string]any
						if err := json.Unmarshal([]byte(value), &childObj); err != nil {
							return nil, currentRowNum, nextBytePos, fmt.Errorf(
								`failed to unmarshal the value: "%v" in row: %d, column: %s`,
								value, currentRowNum, headerName,
							)
						}

						childArray := make([]any, 0, len(childObj))
						for _, obj := range childObj {
							childArray = append(childArray, obj)
						}

						row[headerName] = childArray

						continue
					}

					row[headerName] = value

					continue
				}

				switch attrConfig.Type {
				case framework.AttributeTypeInt64, framework.AttributeTypeDouble:
					floatValue, err := strconv.ParseFloat(value, 64)
					if err != nil {
						return nil, currentRowNum, nextBytePos, fmt.Errorf(
							`CSV contains invalid numeric value "%s" in column "%s" at row %d`,
							value, headerName, currentRowNum,
						)
					}

					row[headerName] = floatValue
				default:
					row[headerName] = value
				}
			}

			objects = append(objects, row)
		}

		if currentRowNum >= targetEndRow {
			break
		}

		currentRowNum++
	}

	return objects, currentRowNum, nextBytePos, nil
}

// TODO: Clean this up by decoupling the attribute value conversion logic from the CSV parsing logic.
// CSVBytesToPage converts a CSV byte array to an array of objects.
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
