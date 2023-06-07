package exiftool

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var fileExifJson map[string]interface{}

func ExifTool(fileValue string) ([]byte, error) {
	// Gather file raw info with LStat
	fileInfo, err := os.Lstat(fileValue)
	if err != nil {
		return nil, err
	}
	isDirectory := fileInfo.IsDir()
	rawSize := fileInfo.Size()

	// Compute file hash
	hashFile, err := os.Open(fileValue)
	if err != nil {
		return nil, err
	}
	defer hashFile.Close()

	fileHash := sha1.New()
	if _, err := io.Copy(fileHash, hashFile); err != nil {
		return nil, err
	}
	fileHashBytes := fileHash.Sum(nil)
	fileHashValue := hex.EncodeToString(fileHashBytes[:])

	// Compute full file name hash
	fullFileNameHash := sha1.Sum([]byte(fileValue))
	fullFileNameHashValue := hex.EncodeToString(fullFileNameHash[:])

	// Identify file path and compute file path hash
	filePathValue := filepath.Dir(fileValue)
	filePathHash := sha1.Sum([]byte(filePathValue))
	filePathHashValue := hex.EncodeToString(filePathHash[:])

	// Create UUID
	UUID := fullFileNameHashValue
	UUID += ":"
	UUID += fileHashValue

	// Extract file exif data using onboard exiftool
	fileExif, err := exec.Command("exiftool", "-j", fileValue).Output()
	if err != nil {
		return nil, err
	}

	// Unmarshal json
	var exifData []map[string]interface{}

	errUnm := json.Unmarshal(fileExif, &exifData)
	if errUnm != nil {
		return nil, errUnm
	}

	// Add values to JSON
	exifData[0]["IsDirectory"] = isDirectory
	exifData[0]["FileSizeRaw"] = rawSize
	exifData[0]["FileHash"] = fileHashValue
	exifData[0]["SourceFileHash"] = fullFileNameHashValue
	exifData[0]["DirectoryHash"] = filePathHashValue
	exifData[0]["_id"] = UUID

	// Marshal outbound json
	fileJson, err := json.MarshalIndent(&exifData, "", "  ")
	if err != nil {
		return nil, err
	}

	return fileJson, nil
}
