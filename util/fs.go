package util

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/henilmalaviya/filic"
)

func WriteJSONToFile(filePath string, data interface{}) error {
	f := filic.NewFile(filePath)

	f.Create()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	if err := f.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write JSON data to file: %w", err)
	}

	return nil
}

// CleanupOldFiles removes the oldest files in a directory if the total count exceeds maxFiles
func CleanupOldFiles(directory string, maxFiles int) error {
	if maxFiles <= 0 {
		return nil // No cleanup needed
	}

	// Use filic to access directory
	dir := filic.NewDirectory(directory)

	// Get list of files in the directory
	files, err := dir.ListFiles()
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", directory, err)
	}

	// If we don't exceed the limit, no cleanup needed
	if len(files) <= maxFiles {
		return nil
	}

	// Collect file info with modification times
	type fileInfo struct {
		file    *filic.File
		modTime time.Time
	}

	var fileInfos []fileInfo
	for _, file := range files {
		// Use os.Stat to get file info since filic doesn't provide it directly
		info, err := os.Stat(file.Path)
		if err != nil {
			continue // Skip files we can't get info for
		}

		fileInfos = append(fileInfos, fileInfo{
			file:    file,
			modTime: info.ModTime(),
		})
	}

	// Sort files by modification time (oldest first)
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].modTime.Before(fileInfos[j].modTime)
	})

	// Calculate how many files to delete
	filesToDelete := len(fileInfos) - maxFiles

	// Delete the oldest files using os.Remove
	for i := 0; i < filesToDelete; i++ {
		if err := os.Remove(fileInfos[i].file.Path); err != nil {
			return fmt.Errorf("failed to remove file %s: %w", fileInfos[i].file.Path, err)
		}
	}

	return nil
}

// ReadJSONFromFile reads JSON data from a file using filic
func ReadJSONFromFile(filePath string, data interface{}) error {
	f := filic.NewFile(filePath)

	jsonData, err := f.Read()
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	if err := json.Unmarshal(jsonData, data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	return nil
}
