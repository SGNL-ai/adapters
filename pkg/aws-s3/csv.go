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
	// StreamingChunkSize defines how many bytes to read per S3 range request
	// 1MB chunks provide good balance between memory usage and API calls
	StreamingChunkSize = 1024 * 1024 // 1MB
)

// CSVHeaders extracts just the headers from the first chunk of CSV data
func CSVHeaders(headerChunk *[]byte) ([]string, error) {
	if headerChunk == nil || len(*headerChunk) == 0 {
		return nil, fmt.Errorf("CSV file is empty or could not be read")
	}

	csvData := csv.NewReader(bytes.NewReader(*headerChunk))

	// Read just the first line to get headers
	record, err := csvData.Read()
	if err != nil {
		return nil, fmt.Errorf("CSV file format is invalid or corrupted: %v", err)
	}

	return record, nil
}

// CSVRowCount estimates the total number of rows by reading through the file in chunks
// This is used to determine if we need to continue pagination
func CSVRowCount(handler *S3Handler, ctx context.Context, bucket, key string, fileSize int64) (int64, error) {
	var totalRows int64 = 0
	var currentPos int64 = 0

	for currentPos < fileSize {
		endPos := currentPos + StreamingChunkSize - 1
		if endPos >= fileSize {
			endPos = fileSize - 1
		}

		chunkData, err := handler.GetFileRange(ctx, bucket, key, currentPos, endPos)
		if err != nil {
			return 0, fmt.Errorf("unable to read CSV file data for row counting: %v", err)
		}

		// Count newlines in this chunk
		chunkRows := int64(bytes.Count(*chunkData, []byte("\n")))
		totalRows += chunkRows

		currentPos = endPos + 1
	}

	// Subtract 1 for header row if file has content
	if totalRows > 0 {
		totalRows--
	}

	return totalRows, nil
}

// StreamingCSVToPage processes CSV data by streaming chunks and seeking to the desired start position
func StreamingCSVToPage(
	handler *S3Handler,
	ctx context.Context,
	bucket, key string,
	fileSize int64,
	headers []string,
	start int64,
	pageSize int64,
	attrConfig []*framework.AttributeConfig,
) ([]map[string]any, bool, error) {

	objects := make([]map[string]any, 0, pageSize)
	headerToAttributeConfig := headerToAttributeConfig(headers, attrConfig)

	var currentRow int64 = 0 // Start at 0 to track absolute row position (including header)
	var currentPos int64 = 0 // Current byte position in file
	var collectedRows int64 = 0

	// We need to find our starting position and then collect pageSize rows
	targetEndRow := start + pageSize

	for currentPos < fileSize && collectedRows < pageSize {
		endPos := currentPos + StreamingChunkSize - 1
		if endPos >= fileSize {
			endPos = fileSize - 1
		}

		chunkData, err := handler.GetFileRange(ctx, bucket, key, currentPos, endPos)
		if err != nil {
			return nil, false, fmt.Errorf("unable to read CSV file data: %v", err)
		}

		// Process this chunk and get the actual byte position where processing stopped
		chunkObjects, nextRow, nextBytePos, err := processCSVChunk(
			*chunkData,
			headers,
			headerToAttributeConfig,
			currentRow,
			start,
			targetEndRow,
			currentPos, // Pass current file position
		)
		if err != nil {
			return nil, false, fmt.Errorf("CSV file processing failed: %v", err)
		}

		objects = append(objects, chunkObjects...)
		collectedRows += int64(len(chunkObjects))
		currentRow = nextRow

		// Safety check: if we didn't advance position, prevent infinite loop
		if nextBytePos <= currentPos {
			return nil, false, fmt.Errorf("CSV file contains formatting issues that prevent processing from continuing")
		}

		// CRITICAL FIX: Use actual processed byte position, not chunk boundary
		currentPos = nextBytePos

		// If we've collected enough rows, we can stop
		if collectedRows >= pageSize {
			break
		}
	}

	// Determine if there are more rows after this page
	hasNext := currentRow <= start+pageSize || currentPos < fileSize
	if !hasNext && collectedRows == pageSize {
		// We collected exactly pageSize rows, but there might be more
		// Check if we're at the end of the file
		hasNext = currentPos < fileSize
	}

	return objects, hasNext, nil
}

// processCSVChunk processes a single chunk of CSV data
func processCSVChunk(
	chunkData []byte,
	headers []string,
	headerToAttributeConfig map[string]framework.AttributeConfig,
	startRowNum int64,
	targetStartRow int64,
	targetEndRow int64,
	chunkStartPos int64, // File position where this chunk starts
) ([]map[string]any, int64, int64, error) {

	// Handle the case where chunk doesn't end on a complete line
	// Find the last complete line in the chunk
	lastNewlineIndex := bytes.LastIndex(chunkData, []byte("\n"))
	if lastNewlineIndex == -1 {
		return nil, startRowNum, chunkStartPos, fmt.Errorf("CSV file contains a single row larger than %d MB, which exceeds processing limits", StreamingChunkSize/(1024*1024))
	}

	// Only process up to the last complete line
	completeChunk := chunkData[:lastNewlineIndex+1]

	// Calculate the actual file position where we stopped processing
	// This is crucial for continuing from the right position in the next chunk
	nextBytePos := chunkStartPos + int64(lastNewlineIndex) + 1

	csvReader := csv.NewReader(bytes.NewReader(completeChunk))

	var objects []map[string]any
	currentRowNum := startRowNum

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

		// Check if this row is in our target range
		// Convert currentRowNum to data row number (subtract 1 to account for header)
		dataRowNum := currentRowNum
		if dataRowNum >= targetStartRow && dataRowNum < targetEndRow {
			row := make(map[string]interface{})

			for i, value := range record {
				if i >= len(headers) {
					continue // Skip extra columns
				}

				headerName := headers[i]
				attrConfig, found := headerToAttributeConfig[headerName]

				if !found {
					// Handle complex JSON values
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

				// Convert based on attribute type
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

		currentRowNum++

		// Early exit if we've processed enough rows
		if dataRowNum >= targetEndRow {
			break
		}
	}

	return objects, currentRowNum, nextBytePos, nil
}

// TODO: Clean this up by decoupling the attribute value conversion logic from the CSV parsing logic.
// CSVBytesToPage converts a CSV byte array to an array of objects.
// DEPRECATED: Use StreamingCSVToPage for large files to avoid memory issues.
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
		return nil, false, fmt.Errorf("CSV file format is invalid or corrupted: %v", err)
	}

	count := len(records)
	if count == 0 {
		return nil, false, fmt.Errorf("CSV file contains no data")
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
							`CSV contains invalid JSON data in column "%s" at row %d: %v`,
							headers[i], i, err,
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
						`CSV contains invalid numeric value "%s" in column "%s" at row %d`,
						value, headers[i], i,
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
