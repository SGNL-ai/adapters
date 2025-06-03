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
	startBytePos int64,
	pageSize int64,
	attrConfig []*framework.AttributeConfig,
) ([]map[string]any, bool, int64, error) {
	const maxProcessingLimit = 200 * StreamingChunkSize

	objects := make([]map[string]any, 0, pageSize)
	headerToAttributeConfig := headerToAttributeConfig(headers, attrConfig)

	var currentPos, totalProcessedBytes, actualNextBytePos int64 = startBytePos, 0, startBytePos

	for currentPos < fileSize && int64(len(objects)) < pageSize {
		endPos := currentPos + StreamingChunkSize - 1
		if endPos >= fileSize {
			endPos = fileSize - 1
		}

		chunkData, err := handler.GetFileRange(ctx, bucket, key, currentPos, endPos)
		if err != nil {
			return nil, false, 0, fmt.Errorf("unable to read CSV file data: %v", err)
		}

		if chunkData == nil {
			return nil, false, 0, fmt.Errorf("received empty response from S3 file range request")
		}

		chunkSize := int64(len(*chunkData))

		chunkObjects, chunkNextBytePos, objectBytePositions, err := processCSVChunk(
			*chunkData,
			headers,
			headerToAttributeConfig,
			currentPos,
		)
		if err != nil {
			return nil, false, 0, fmt.Errorf("CSV file processing failed: %v", err)
		}

		totalProcessedBytes += chunkSize

		// Add objects but respect pageSize limit
		remainingSlots := pageSize - int64(len(objects))
		objectsToAdd := int64(len(chunkObjects))

		if objectsToAdd <= remainingSlots {
			// Add all chunk objects
			objects = append(objects, chunkObjects...)
			actualNextBytePos = chunkNextBytePos
		} else {
			// Add only what we need to reach pageSize
			objects = append(objects, chunkObjects[:remainingSlots]...)
			// Calculate cursor position based on the last object we're including
			if int(remainingSlots) > 0 && len(objectBytePositions) > int(remainingSlots) {
				actualNextBytePos = objectBytePositions[remainingSlots-1]
			} else {
				actualNextBytePos = chunkNextBytePos
			}
		}

		if totalProcessedBytes >= maxProcessingLimit && len(objects) > 0 {
			hasNext := true
			return objects, hasNext, actualNextBytePos, nil
		}

		if chunkNextBytePos <= currentPos {
			return nil, false, 0, fmt.Errorf("CSV file contains formatting issues that prevent processing from continuing")
		}

		currentPos = chunkNextBytePos

		// Break if we've reached pageSize
		if int64(len(objects)) >= pageSize {
			break
		}
	}

	hasNext := currentPos < fileSize || actualNextBytePos < fileSize

	return objects, hasNext, actualNextBytePos, nil
}

// findLastCompleteRowEnd finds the last complete row in the chunk (simplified version)
func findLastCompleteRowEnd(data []byte) int {
	var inQuotes bool
	lastValidNewline := -1

	for i, b := range data {
		if b == '"' {
			if i+1 < len(data) && data[i+1] == '"' {
				i++ // Skip escaped quote
				continue
			}
			inQuotes = !inQuotes
		} else if b == '\n' && !inQuotes {
			lastValidNewline = i
		}
	}

	return lastValidNewline
}

// processCSVChunkWithPositions processes a single chunk of CSV data and tracks byte positions
func processCSVChunk(
	chunkData []byte,
	headers []string,
	headerToAttributeConfig map[string]framework.AttributeConfig,
	chunkStartPos int64,
) ([]map[string]any, int64, []int64, error) {
	lastNewlineIndex := findLastCompleteRowEnd(chunkData)
	if lastNewlineIndex == -1 {
		return nil, chunkStartPos, nil, fmt.Errorf("CSV file contains a single row larger than %d MB", StreamingChunkSize/(1024*1024))
	}

	completeChunk := chunkData[:lastNewlineIndex+1]
	nextBytePos := chunkStartPos + int64(lastNewlineIndex) + 1

	// Track byte positions for each row
	var objectBytePositions []int64
	var inQuotes bool

	// Find the start position of each row
	for i, b := range completeChunk {
		if b == '"' {
			if i+1 < len(completeChunk) && completeChunk[i+1] == '"' {
				i++ // Skip escaped quote
				continue
			}
			inQuotes = !inQuotes
		} else if b == '\n' && !inQuotes {
			// End of current row, next row starts after this newline
			nextRowStart := chunkStartPos + int64(i) + 1
			objectBytePositions = append(objectBytePositions, nextRowStart)
		}
	}

	csvReader := csv.NewReader(bytes.NewReader(completeChunk))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nextBytePos, nil, fmt.Errorf("CSV file format is invalid or corrupted: %v", err)
	}

	var objects []map[string]any

	// Process all records in this chunk - no header skipping since we're past headers
	for _, record := range records {
		row := make(map[string]interface{})

		for j, value := range record {
			if j >= len(headers) {
				continue
			}

			headerName := headers[j]
			attrConfig, found := headerToAttributeConfig[headerName]

			if !found {
				if strings.HasPrefix(value, "[{") && strings.HasSuffix(value, "}]") {
					var childObj []map[string]any
					if err := json.Unmarshal([]byte(value), &childObj); err != nil {
						return nil, nextBytePos, nil, fmt.Errorf(
							`failed to unmarshal the value: "%v" in column: %s`,
							value, headerName,
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
					return nil, nextBytePos, nil, fmt.Errorf(
						`CSV contains invalid numeric value "%s" in column "%s"`,
						value, headerName,
					)
				}

				row[headerName] = floatValue
			default:
				row[headerName] = value
			}
		}

		objects = append(objects, row)
	}

	// Ensure we have position information for all objects
	// If we have fewer positions than objects, extend with nextBytePos
	for len(objectBytePositions) < len(objects) {
		objectBytePositions = append(objectBytePositions, nextBytePos)
	}

	return objects, nextBytePos, objectBytePositions, nil
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
