package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	JsonContentType       = "application/json; charset=utf-8"
	ScriptfileContentType = "application/prs.powergrim.scriptfile+json; charset=utf-8"
	LayoutContentType     = "application/prs.powergrim.layout+json; charset=utf-0"
	GameContentType       = "application/prs.powergrim.game+json; charset=utf-8"
	ActionContentType     = "application/prs.powergrim.action+json; charset=utf-8"
	ActionsContentType    = "application/prs.powergrim.actions+json; charset=utf-8"
)

type VersionedGame struct {
	LastModified time.Time
	Version      int
	Game         Game
}

var scriptIdToScriptFileId = make(map[string]string)
var gamesMut sync.Mutex
var games = make(map[string]VersionedGame)

func main() {
	http.HandleFunc("GET /findScript", findScript)
	if err := handleFiles[ScriptFile]("scripts", "script", ScriptfileContentType, collectScriptIds); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := handleFiles[Layout]("layouts", "layout", LayoutContentType, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	http.HandleFunc("POST /game", newGame)
	http.HandleFunc("GET /game/{gameId}", getGame)
	http.HandleFunc("PATCH /game/{gameId}", patchGame)

	err := http.ListenAndServe(":3000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

func newGame(resp http.ResponseWriter, req *http.Request) {
	contentType := req.Header["Content-Type"]
	if len(contentType) != 1 || contentType[0] != GameContentType {
		resp.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	game := VersionedGame{
		LastModified: time.Now().Truncate(time.Second),
		Version:      1,
	}
	err := json.NewDecoder(req.Body).Decode(&game.Game)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	gameId := uuid.NewString()
	gamesMut.Lock()
	games[gameId] = game
	gamesMut.Unlock()
	resp.Header().Add("Location", fmt.Sprintf("/game/%s", gameId))
	resp.WriteHeader(http.StatusCreated)
}

func getGame(resp http.ResponseWriter, req *http.Request) {
	gamesMut.Lock()
	game, ok := games[req.PathValue("gameId")]
	gamesMut.Unlock()
	if !ok {
		http.Error(resp, "", http.StatusNotFound)
		return
	}
	ifNoneMatch := req.Header["If-None-Match"]
	ifModifiedSince := req.Header["If-Modified-Since"]
	if len(ifNoneMatch) == 1 || len(ifModifiedSince) == 1 {
		sendBody := false
		if len(ifNoneMatch) == 1 {
			sendBody = sendBody || ifNoneMatch[0] != fmt.Sprintf("W/%d", game.Version)
		}
		if len(ifModifiedSince) == 1 {
			t, err := time.Parse(http.TimeFormat, ifModifiedSince[0])
			sendBody = sendBody || err != nil || game.LastModified.After(t)
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
	header.Add("Content-Type", GameContentType)
	header.Add("Last-Modified", game.LastModified.UTC().Format(http.TimeFormat))
	header.Add("ETag", fmt.Sprintf("W/%d", game.Version))
	json.NewEncoder(resp).Encode(game.Game)
}

func patchGame(resp http.ResponseWriter, req *http.Request) {
	gameId := req.PathValue("gameId")
	gamesMut.Lock()
	game, ok := games[gameId]
	gamesMut.Unlock()
	if !ok {
		http.Error(resp, "", http.StatusNotFound)
		return
	}
	ifMatch := req.Header["If-Match"]
	if len(ifMatch) == 1 && ifMatch[0] != fmt.Sprintf("W/%d", game.Version) {
		http.Error(resp, "", http.StatusPreconditionFailed)
		return
	}
	ifUnmodifiedSince := req.Header["If-Unmodified-Since"]
	if len(ifUnmodifiedSince) == 1 {
		t, err := time.Parse(http.TimeFormat, ifUnmodifiedSince[0])
		if err != nil || game.LastModified.After(t) {
			http.Error(resp, "", http.StatusPreconditionFailed)
			return
		}
	}
	contentType := req.Header["Content-Type"]
	if len(contentType) != 1 {
		http.Error(resp, "", http.StatusUnsupportedMediaType)
		return
	}
	var actions []WrappedAction
	switch contentType[0] {
	case ActionContentType:
		actions = make([]WrappedAction, 1)
		err := json.NewDecoder(req.Body).Decode(&actions[0])
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
	case ActionsContentType:
		err := json.NewDecoder(req.Body).Decode(&actions)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		http.Error(resp, "", http.StatusUnsupportedMediaType)
		return
	}
	for _, action := range actions {
		newGame, err := game.Game.ApplyAction(action.Action)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		game.Game = newGame
	}
	gamesMut.Lock()
	updated := false
	if games[gameId].Version == game.Version {
		game.LastModified = time.Now().Truncate(time.Second)
		game.Version++
		games[gameId] = game
		updated = true
	}
	gamesMut.Unlock()
	header := resp.Header()
	header.Add("Content-Type", GameContentType)
	header.Add("Last-Modified", game.LastModified.UTC().Format(http.TimeFormat))
	header.Add("ETag", fmt.Sprintf("W/%d", game.Version))
	if !updated {
		resp.WriteHeader(http.StatusConflict)
	}
	json.NewEncoder(resp).Encode(game.Game)
}

func collectScriptIds(scriptFileId string, file VersionedFile, scriptFile ScriptFile) error {
	for _, script := range scriptFile.Scripts {
		if scriptIdToScriptFileId[script.Id] != "" {
			return fmt.Errorf("duplicate script id %q", script.Id)
		}
		scriptIdToScriptFileId[script.Id] = scriptFileId
	}
	return nil
}

func findScript(resp http.ResponseWriter, req *http.Request) {
	if !req.URL.Query().Has("q") {
		http.Error(resp, "", http.StatusBadRequest)
		return
	}
	scriptFileId := scriptIdToScriptFileId[req.URL.Query().Get("q")]
	header := resp.Header()
	if _, hasOrigin := req.Header["Origin"]; hasOrigin {
		header.Add("Access-Control-Allow-Origin", "*")
	}
	header.Add("Content-Type", JsonContentType)
	if scriptFileId == "" {
		json.NewEncoder(resp).Encode(nil)
	} else {
		json.NewEncoder(resp).Encode(scriptFileId)
	}
}
