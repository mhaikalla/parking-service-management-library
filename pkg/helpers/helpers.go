package helpers

import "time"

const dateLayout = "2006-01-02"

func GetString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func SplitJSONObjects(data []byte) [][]byte {
	var objects [][]byte
	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] == '}' {
			end := i + 1
			objects = append(objects, data[start:end])
			start = end
			// Move to the start of the next object
			for start < len(data) && (data[start] == '\n' || data[start] == '\r') {
				start++
			}
		}
	}
	return objects
}

func StringToDate(s string) (time.Time, error) {
	result, err := time.Parse(dateLayout, s)
	if err != nil {
		return result, err
	}
	return result, nil
}

func StringToDatetime(s string) (time.Time, error) {
	result, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return result, err
	}
	return result, nil
}
