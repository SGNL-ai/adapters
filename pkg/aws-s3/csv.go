// Copyright 2025 SGNL.ai, Inc.
package awss3

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
)

const (
	FileTypeCSV        = "csv"
	MaxCSVRowSizeBytes = 1 * 1024 * 1024 // 1MB
)

var ErrEmptyOrMissing = errors.New("empty or missing")

func handleQuoteChar(reader *bufio.Reader, lineBuffer *bytes.Buffer, bytesRead *int64, inQuotes *bool) error {
	if !*inQuotes {
		// This is an opening quote
		*inQuotes = true

		return nil
	}

	// We're inside quotes - need to check if this is an escaped quote or closing quote
	nextBytes, peekErr := reader.Peek(1)

	if peekErr != nil {
		if peekErr == io.EOF {
			// End of file - this quote closes the field
			*inQuotes = false

			return nil
		}
		// Some other error occurred
		return fmt.Errorf("failed to peek next byte after quote: %w", peekErr)
	}

	if len(nextBytes) > 0 && nextBytes[0] == '"' {
		// This is "", an escaped quote - consume the second quote
		nextByte, readErr := reader.ReadByte()
		if readErr != nil {
			return fmt.Errorf("failed to read escaped quote: %w", readErr)
		}

		*bytesRead++

		// Write the second quote to buffer
		// Stay in quotes, this was just an escaped quote
		// Note: The first quote was already written to buffer in readCSVLine
		lineBuffer.WriteByte(nextByte)
	} else {
		// This is the closing quote
		*inQuotes = false
	}

	return nil
}

func handleLineEnding(reader *bufio.Reader, lineBuffer *bytes.Buffer, bytesRead *int64) error {
	nextBytes, peekErr := reader.Peek(1)

	if peekErr != nil {
		if peekErr == io.EOF {
			// CR at EOF is a valid line ending
			return nil
		}
		// Other errors should be propagated
		return fmt.Errorf("failed to peek after CR: %w", peekErr)
	}

	if len(nextBytes) > 0 && nextBytes[0] == '\n' {
		// This is CRLF - consume the LF
		nextByte, readErr := reader.ReadByte()

		if readErr != nil {
			// This shouldn't happen - we just peeked successfully
			return fmt.Errorf("failed to read LF after CR: %w", readErr)
		}

		*bytesRead++

		lineBuffer.WriteByte(nextByte)
	}

	return nil
}

func readCSVLine(reader *bufio.Reader) (
	lineBytes []byte, bytesRead int64, err error) {
	var lineBuffer bytes.Buffer

	inQuotes := false

	for bytesRead = 0; ; {
		if bytesRead >= MaxCSVRowSizeBytes {
			return nil, 0, fmt.Errorf("size limit of %d MB exceeded", MaxCSVRowSizeBytes/(1024*1024))
		}

		b, readErr := reader.ReadByte()
		if readErr != nil {
			if readErr == io.EOF {
				if lineBuffer.Len() == 0 && bytesRead == 0 {
					return nil, 0, ErrEmptyOrMissing
				}

				break
			}

			return nil, 0, fmt.Errorf("failed to read byte: %w", readErr)
		}

		bytesRead++

		lineBuffer.WriteByte(b)

		if b == '"' {
			if err := handleQuoteChar(reader, &lineBuffer, &bytesRead, &inQuotes); err != nil {
				return nil, 0, err
			}
		} else if (b == '\n' || b == '\r') && !inQuotes {
			if b == '\r' {
				if err := handleLineEnding(reader, &lineBuffer, &bytesRead); err != nil {
					return nil, 0, err
				}
			}

			break
		}
	}

	lineBytes = lineBuffer.Bytes()

	if len(lineBytes) > 0 && lineBytes[len(lineBytes)-1] == '\r' {
		lineBytes[len(lineBytes)-1] = '\n'
	}

	return lineBytes, bytesRead, nil
}

func CSVHeaders(reader *bufio.Reader) (headers []string, bytesReadForHeader int64, err error) {
	headerLineBytes, bytesRead, err := readCSVLine(reader)

	if err != nil {
		return nil, 0, fmt.Errorf("CSV header error: %w", err)
	}

	csvReader := csv.NewReader(bytes.NewReader(headerLineBytes))
	parsedHeaders, parseErr := csvReader.Read()

	if parseErr != nil {
		return nil, 0, fmt.Errorf("CSV file format is invalid or corrupted: %v", parseErr)
	}

	if len(parsedHeaders) == 0 {
		return nil, 0, fmt.Errorf("CSV header error: empty or missing")
	}

	return parsedHeaders, bytesRead, nil
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

	var totalBytesRead int64

	hasNext = true

	for int64(len(objects)) < pageSize {
		rowBytes, bytesRead, rowReadErr := readCSVLine(streamReader)

		if bytesRead > 0 {
			if (totalBytesRead + bytesRead) > maxProcessingBytesTotal {
				break
			}

			totalBytesRead += bytesRead
		}

		if rowReadErr != nil && rowReadErr != io.EOF && !errors.Is(rowReadErr, ErrEmptyOrMissing) {
			return nil, 0, false, fmt.Errorf("CSV row error: %w", rowReadErr)
		}

		if len(rowBytes) == 0 {
			if rowReadErr == io.EOF || errors.Is(rowReadErr, ErrEmptyOrMissing) {
				hasNext = false

				return objects, totalBytesRead, hasNext, nil
			}

			continue
		}

		if rowReadErr == io.EOF || errors.Is(rowReadErr, ErrEmptyOrMissing) {
			hasNext = false
		}

		csvRowReader := csv.NewReader(bytes.NewReader(rowBytes))
		record, recordParseErr := csvRowReader.Read()

		if recordParseErr != nil && recordParseErr != io.EOF {
			return nil, 0, false, fmt.Errorf("CSV file format is invalid or corrupted: %w", recordParseErr)
		}

		if len(record) == 0 {
			if !hasNext {
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
						return nil, 0, false, fmt.Errorf(
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
					return nil, 0, false, fmt.Errorf(
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

	return objects, totalBytesRead, hasNext, nil
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
