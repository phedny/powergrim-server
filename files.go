package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

type VersionedFile struct {
	lastModified time.Time
	hash         string
	data         []byte
}

func handleFiles[T any](dirName, urlPath, contentType string, fileLoaded func(fileName string, file VersionedFile, decoded T) error) error {
	entries, err := os.ReadDir(dirName)
	if err != nil {
		return fmt.Errorf("failed to read %s directory: %w", dirName, err)
	}
	files := make(map[string]VersionedFile)
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".json" || !entry.Type().IsRegular() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("failed to read info on %s/%s file: %w", dirName, entry.Name(), err)
		}
		data, err := os.ReadFile(path.Join(dirName, entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read %s/%s file: %w", dirName, entry.Name(), err)
		}
		var decoded T
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&decoded)
		if err != nil {
			return fmt.Errorf("failed to parse %s/%s file: %w", dirName, entry.Name(), err)
		}
		hash := sha256.Sum256(data)
		fileName := entry.Name()[:len(entry.Name())-5]
		file := VersionedFile{
			lastModified: info.ModTime(),
			hash:         hex.EncodeToString(hash[:]),
			data:         data,
		}
		files[fileName] = file
		if fileLoaded != nil {
			err = fileLoaded(fileName, file, decoded)
			if err != nil {
				return err
			}
		}
	}
	http.HandleFunc(fmt.Sprintf("/%s/{id}", urlPath), func(resp http.ResponseWriter, req *http.Request) {
		file, ok := files[req.PathValue("id")]
		if !ok {
			header := resp.Header()
			header.Add("Access-Control-Allow-Origin", "*")
			http.Error(resp, "", http.StatusNotFound)
			return
		}
		ifNoneMatch := req.Header["If-None-Match"]
		ifModifiedSince := req.Header["If-Modified-Since"]
		if len(ifNoneMatch) == 1 || len(ifModifiedSince) == 1 {
			sendBody := false
			if len(ifNoneMatch) == 1 {
				sendBody = sendBody || ifNoneMatch[0] != file.hash
			}
			if len(ifModifiedSince) == 1 {
				t, err := time.Parse(http.TimeFormat, ifModifiedSince[0])
				sendBody = sendBody || err != nil || file.lastModified.After(t)
			}
			if !sendBody {
				resp.WriteHeader(http.StatusNotModified)
				return
			}
		}
		header := resp.Header()
		if _, hasOrigin := req.Header["Origin"]; hasOrigin {
			header.Add("Access-Control-Allow-Origin", "*")
		}
		header.Add("Content-Type", contentType)
		header.Add("Last-Modified", file.lastModified.UTC().Format(http.TimeFormat))
		header.Add("ETag", file.hash)
		resp.Write(file.data)
	})
	return nil
}
