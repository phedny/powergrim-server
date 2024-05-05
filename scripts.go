package main

import (
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

type Script struct {
	Name          string           `json:"name"`
	Complexity    ScriptComplexity `json:"complexity,omitempty"`
	Tagline       string           `json:"tagline"`
	Url           string           `json:"url,omitempty"`
	Logo          string           `json:"logo,omitempty"`
	Description   string           `json:"description"`
	KeyCharacters []string         `json:"keyCharacters,omitempty"`
	Characters    []string         `json:"characters"`
}

type ScriptComplexity struct {
	Level       string  `json:"level,omitempty"`
	Storyteller float64 `json:"storyteller,omitempty"`
	Player      float64 `json:"player,omitempty"`
}

type ScriptFile struct {
	Author  string   `json:"author,omitempty"`
	Url     string   `json:"url,omitempty"`
	Scripts []Script `json:"scripts"`
}

type VersionedScripts struct {
	lastModified time.Time
	hash         string
	data         []byte
}

var scripts = make(map[string]VersionedScripts)

func loadScripts() error {
	entries, err := os.ReadDir("scripts")
	if err != nil {
		return fmt.Errorf("failed to read scripts directory: %w", err)
	}
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".json" || !entry.Type().IsRegular() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("failed to read info on scripts/%s file: %w", entry.Name(), err)
		}
		d, err := os.ReadFile(path.Join("scripts", entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read scripts/%s file: %w", entry.Name(), err)
		}
		var s ScriptFile
		err = json.Unmarshal(d, &s)
		if err != nil {
			return fmt.Errorf("failed to parse scripts/%s file: %w", entry.Name(), err)
		}
		hash := sha256.Sum256(d)
		scripts[entry.Name()[:len(entry.Name())-5]] = VersionedScripts{
			lastModified: info.ModTime(),
			hash:         hex.EncodeToString(hash[:]),
			data:         d,
		}
	}
	return nil
}

func getScript(resp http.ResponseWriter, req *http.Request) {
	vScripts, ok := scripts[req.PathValue("scriptId")]
	if !ok {
		http.Error(resp, "", http.StatusNotFound)
		return
	}
	ifNoneMatch := req.Header["If-None-Match"]
	ifModifiedSince := req.Header["If-Modified-Since"]
	if len(ifNoneMatch) == 1 || len(ifModifiedSince) == 1 {
		sendBody := false
		if len(ifNoneMatch) == 1 {
			sendBody = sendBody || ifNoneMatch[0] != vScripts.hash
		}
		if len(ifModifiedSince) == 1 {
			t, err := time.Parse(http.TimeFormat, ifModifiedSince[0])
			sendBody = sendBody || err != nil || vScripts.lastModified.After(t)
		}
		if !sendBody {
			resp.WriteHeader(http.StatusNotModified)
			return
		}
	}
	header := resp.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Content-Type", ScriptsContentType)
	header.Add("Last-Modified", vScripts.lastModified.UTC().Format(http.TimeFormat))
	header.Add("ETag", vScripts.hash)
	resp.Write(vScripts.data)
}
