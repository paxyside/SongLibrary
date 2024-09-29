package http_controller

import (
	"context"
	"effictiveMobile/internal/domain/entities"
	"effictiveMobile/internal/domain/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strconv"
)

type SongController struct {
	songService service.SongService
	logger      *slog.Logger
}

func NewSongController(songService service.SongService, logger *slog.Logger) *SongController {
	return &SongController{
		songService: songService,
		logger:      logger.With("controller", "SongController"),
	}
}

// GetSongsHandler
// @Title Get list of songs with filtering and pagination
// @Description Retrieve a list of songs with optional filters and pagination
// @Tag Song
// @Param  limit   query  int  true   "Number of songs to return"   "10"
// @Param  offset  query  int  true   "Offset for pagination"       "0"
// @Param  group   query  string  false "Filter by group"            "Muse"
// @Success  200  object  entities.SongsResponse   "Songs list with pagination"
// @Failure  400  object  entities.ErrorResponse   "Invalid input parameters"
// @Failure  401  object  entities.ErrorResponse   "Unauthorized"
// @Failure  500  object  entities.ErrorResponse   "Internal server error"
// @Route /api/v1/songs [get]
func (c *SongController) GetSongsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.logger.Error("invalid limit parameter", "limit", limitStr, "error", err)
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.logger.Error("invalid offset parameter", "offset", offsetStr, "error", err)
		http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
		return
	}

	filter := map[string]interface{}{}
	group := r.URL.Query().Get("group")
	if group != "" {
		filter["group"] = group
	}

	songs, err := c.songService.GetSongs(ctx, filter, limit, offset)
	if err != nil {
		c.logger.Error("failed to retrieve songs", "error", err)
		http.Error(w, "Failed to retrieve songs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(songs); err != nil {
		c.logger.Error("failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetSongByIDHandler
// @Title Get song details by ID
// @Description Retrieve detailed information about a song by its ID
// @Tag Song
// @Param  id  path  int  true  "ID of the song"  "1"
// @Success  200  object  entities.Song           "Detailed song information"
// @Failure  400  object  entities.ErrorResponse  "Invalid song ID"
// @Failure  401  object  entities.ErrorResponse   "Unauthorized"
// @Failure  404  object  entities.ErrorResponse  "Song not found"
// @Failure  500  object  entities.ErrorResponse  "Internal server error"
// @Route /api/v1/songs/{id} [get]
func (c *SongController) GetSongByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.logger.Error("invalid song ID", "id", idStr, "error", err)
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	song, err := c.songService.GetSongByID(ctx, id)
	if err != nil {
		c.logger.Error("failed to retrieve song", "id", id, "error", err)
		http.Error(w, "Failed to retrieve song: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if song == nil {
		c.logger.Warn("song not found", "id", id)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(song); err != nil {
		c.logger.Error("failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// CreateSongHandler
// @Title Create a new song
// @Description Create a new song using group and song information
// @Tag Song
// @Param song body entities.CreateSongRequest true "Info of the song to create"
// @Success 201 {object} entities.Song "Created song"
// @Failure 400 {object} entities.ErrorResponse "Invalid input"
// @Failure  401  object  entities.ErrorResponse   "Unauthorized"
// @Failure 500 {object} entities.ErrorResponse "Internal server error"
// @Route /api/songs [post]
func (c *SongController) CreateSongHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req entities.CreateSongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем данные о песне через внешний API
	details, err := c.songService.GetSongDetails(ctx, req.Group, req.Song)
	if err != nil {
		http.Error(w, "Failed to get song details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаём объект песни с заполненными данными
	song := entities.Song{
		Group:       req.Group,
		Song:        req.Song,
		ReleaseDate: details.ReleaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}

	// Сохраняем песню в базе данных
	err = c.songService.CreateSong(ctx, &song)
	if err != nil {
		http.Error(w, "Failed to create song: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

// UpdateSongHandler
// @Title Update song details by ID
// @Description Update the information of an existing song by its ID
// @Tag Song
// @Param  id    path  int           true  "ID of the song to update"  "1"
// @Param  song  body  entities.Song  true  "Updated song information"
// @Success  200  object  map[string]string  "Song updated successfully"
// @Failure  400  object  entities.ErrorResponse   "Invalid input data"
// @Failure  401  object  entities.ErrorResponse   "Unauthorized"
// @Failure  404  object  entities.ErrorResponse   "Song not found"
// @Failure  500  object  entities.ErrorResponse   "Internal server error"
// @Route /api/v1/songs/update/{id} [put]
func (c *SongController) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.logger.Error("invalid song ID", "id", idStr, "error", err)
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	var song entities.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		c.logger.Error("invalid request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = c.songService.UpdateSong(ctx, id, &song)
	if err != nil {
		c.logger.Error("failed to update song", "songID", id, "song", song, "error", err)
		http.Error(w, "Failed to update song: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Song updated successfully"}); err != nil {
		c.logger.Error("failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DeleteSongHandler
// @Title Delete song by ID
// @Description Delete an existing song by its ID
// @Tag Song
// @Param  id  path  int  true  "ID of the song to delete"  "1"
// @Success  200  object  map[string]string  "Song deleted successfully"
// @Failure  400  object  entities.ErrorResponse   "Invalid song ID"
// @Failure  401  object  entities.ErrorResponse   "Unauthorized"
// @Failure  404  object  entities.ErrorResponse   "Song not found"
// @Failure  500  object  entities.ErrorResponse   "Internal server error"
// @Route /api/v1/songs/delete/{id} [delete]
func (c *SongController) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.logger.Error("invalid song ID", "id", idStr, "error", err)
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	err = c.songService.DeleteSong(ctx, id)
	if err != nil {
		c.logger.Error("failed to delete song", "songID", id, "error", err)
		http.Error(w, "Failed to delete song: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Song deleted successfully"}); err != nil {
		c.logger.Error("failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
