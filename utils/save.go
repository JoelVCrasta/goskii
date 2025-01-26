package utils

import (
	"fmt"
	"os"
)

func SaveToTextFile(ascii string, path string, fileName string) error {
	file, err := os.Create(fmt.Sprintf("%s/%s.txt", path, fileName))
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(ascii)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}