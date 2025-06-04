// Copyright 2025 SGNL.ai, Inc.
package awss3

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
)

const (
	FileTypeCSV        = "csv"
	MaxHeaderSizeBytes = 8 * 1024        // 8KB
	MaxCSVRowSizeBytes = 1 * 1024 * 1024 // 1MB
)

func CSVHeaders(reader *bufio.Reader) (headers []string, bytesReadForHeader int64, err error) {
	var headerLineBuffer bytes.Buffer
	var currentBytesRead int64
	inQuotes := false

	for {
		if currentBytesRead >= MaxHeaderSizeBytes {
			return nil, currentBytesRead, fmt.Errorf("CSV header line exceeds %dKB size limit", MaxHeaderSizeBytes/1024)
		}

		b, readErr := reader.ReadByte()
		if readErr != nil {
			if readErr == io.EOF {
				if headerLineBuffer.Len() == 0 && currentBytesRead == 0 {
					return nil, 0, fmt.Errorf("CSV header is empty or missing")
				}
				break
			}
			return nil, currentBytesRead, fmt.Errorf("failed to read byte for CSV header: %w", readErr)
		}

		currentBytesRead++
		headerLineBuffer.WriteByte(b)

		if b == '"' {
			inQuotes = !inQuotes
		} else if b == '\n' && !inQuotes {
			break
		}
	}

	headerLineBytes := headerLineBuffer.Bytes()
	if len(headerLineBytes) == 0 {
		return nil, currentBytesRead, fmt.Errorf("CSV header is empty or missing")
	}

	csvReader := csv.NewReader(bytes.NewReader(headerLineBytes))
	parsedHeaders, parseErr := csvReader.Read()
	if parseErr != nil {
		return nil, currentBytesRead, fmt.Errorf("CSV file format is invalid or corrupted: %v", parseErr)
	}

	if len(parsedHeaders) == 0 {
		return nil, currentBytesRead, fmt.Errorf("CSV header is empty or missing")
	}

	return parsedHeaders, currentBytesRead, nil
}

func readNextCSVRow(reader *bufio.Reader, maxRowBytes int64) (rowLineBytes []byte, bytesConsumedThisRow int64, err error) {
	var rowBuffer bytes.Buffer
	var currentBytesRead int64
	inQuotes := false
	var prevByte byte

	for {
		if currentBytesRead >= maxRowBytes {
			return rowBuffer.Bytes(), currentBytesRead, fmt.Errorf("CSV file contains a single row larger than %d MB", maxRowBytes/(1024*1024))
		}

		b, readErr := reader.ReadByte()
		if readErr != nil {
			if readErr == io.EOF {
				if rowBuffer.Len() > 0 {
					return rowBuffer.Bytes(), currentBytesRead, io.EOF
				}
				return nil, currentBytesRead, io.EOF
			}
			return nil, currentBytesRead, fmt.Errorf("failed to read byte for CSV row: %w", readErr)
		}

		currentBytesRead++
		rowBuffer.WriteByte(b)

		if b == '"' {
			if inQuotes && prevByte == '"' {
				prevByte = 0
			} else {
				inQuotes = !inQuotes
				prevByte = b
			}
		} else if b == '\n' && !inQuotes {
			break
		} else {
			prevByte = b
		}
	}
	return rowBuffer.Bytes(), currentBytesRead, nil
}

func StreamingCSVToPage(
	streamReader *bufio.Reader,
	headers []string,
	pageSize int64,
	attrConfig []*framework.AttributeConfig,
	maxProcessingBytesTotal int64,
) (objects []map[string]any, bytesReadFromDataStream int64, hasNext bool, err error) {

	objects = make([]map[string]any, 0, pageSize)
	headerToAttributeConfig := headerToAttributeConfig(headers, attrConfig)
	var totalBytesReadThisCall int64

	hasNext = true

	for int64(len(objects)) < pageSize {
		if maxProcessingBytesTotal > 0 && totalBytesReadThisCall >= maxProcessingBytesTotal {
			break
		}

		rowBytes, bytesForRow, rowReadErr := readNextCSVRow(streamReader, MaxCSVRowSizeBytes)

		if bytesForRow > 0 {
			totalBytesReadThisCall += bytesForRow
		}

		processThisRowData := false
		if rowReadErr == nil {
			if len(rowBytes) > 0 {
				processThisRowData = true
			} else {
				processThisRowData = true
			}
		} else if rowReadErr == io.EOF {
			hasNext = false
			if len(rowBytes) > 0 {
				processThisRowData = true
			} else {
				break
			}
		} else {
			return objects, totalBytesReadThisCall, false, fmt.Errorf("unable to read CSV file data: %w", rowReadErr)
		}

		if !processThisRowData {
			if !hasNext {
				break
			}
			continue
		}

		csvRowReader := csv.NewReader(bytes.NewReader(rowBytes))
		record, recordParseErr := csvRowReader.Read()

		if recordParseErr != nil {
			if recordParseErr == io.EOF {
				if len(record) == 0 {
					if !hasNext {
						break
					}
					continue
				}
			} else {
				return objects, totalBytesReadThisCall, false, fmt.Errorf("CSV file format is invalid or corrupted: %w", recordParseErr)
			}
		}

		if len(record) == 0 {
			if !hasNext && rowReadErr == io.EOF {
				break
			}
			continue
		}

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
					if errUnmarshal := json.Unmarshal([]byte(value), &childObj); errUnmarshal == nil {
						childArray := make([]any, 0, len(childObj))
						for _, obj := range childObj {
							childArray = append(childArray, obj)
						}
						row[headerName] = childArray
					} else {
						return objects, totalBytesReadThisCall, false, fmt.Errorf(
							`failed to unmarshal the value: "%v" in column: %s`,
							value, headerName,
						)
					}
				} else {
					row[headerName] = value
				}
				continue
			}

			switch attrConfig.Type {
			case framework.AttributeTypeInt64, framework.AttributeTypeDouble:
				floatValue, convErr := strconv.ParseFloat(value, 64)
				if convErr != nil {
					return objects, totalBytesReadThisCall, false, fmt.Errorf(
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

		if !hasNext {
			break
		}
	}

	if hasNext && int64(len(objects)) == pageSize {
		_, errPeek := streamReader.Peek(1)
		if errPeek == io.EOF {
			hasNext = false
		}
	}

	return objects, totalBytesReadThisCall, hasNext, nil
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
