package service

import (
	"context"
	"effictiveMobile/internal/domain/entities"
	"effictiveMobile/internal/infrastrtucture/external_api"
	"effictiveMobile/internal/infrastrtucture/persistence"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
)

type SongService interface {
	GetSongs(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]entities.Song, error)
	GetSongByID(ctx context.Context, id int) (*entities.Song, error)
	CreateSong(ctx context.Context, song *entities.Song) error
	UpdateSong(ctx context.Context, id int, song *entities.Song) error
	DeleteSong(ctx context.Context, id int) error
	GetSongDetails(ctx context.Context, group, song string) (*external_api.SongDetail, error) // Новый метод
}

type SongServiceImpl struct {
	songRepo  persistence.SongRepository
	logger    *slog.Logger
	apiClient *external_api.Client
}

func NewSongService(songRepo persistence.SongRepository, logger *slog.Logger, apiClient *external_api.Client) *SongServiceImpl {
	return &SongServiceImpl{
		songRepo:  songRepo,
		logger:    logger.With("service", "SongService"),
		apiClient: apiClient,
	}
}

// GetSongs валидирует и фильтрует данные перед вызовом репозитория.
func (s *SongServiceImpl) GetSongs(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]entities.Song, error) {
	if limit <= 0 {
		err := errors.New("limit must be greater than 0")
		s.logger.Error("invalid limit", "error", err)
		return nil, err
	}
	if offset < 0 {
		err := errors.New("offset cannot be negative")
		s.logger.Error("invalid offset", "error", err)
		return nil, err
	}

	if group, ok := filter["group"]; ok && group == "" {
		err := errors.New("group filter cannot be empty")
		s.logger.Error("invalid group filter", "error", err)
		return nil, err
	}

	return s.songRepo.GetSongs(ctx, filter, limit, offset)
}

// GetSongByID валидирует ID и вызывает репозиторий для получения песни.
func (s *SongServiceImpl) GetSongByID(ctx context.Context, id int) (*entities.Song, error) {
	if id <= 0 {
		err := errors.New("invalid song ID")
		s.logger.Error("invalid song ID", "error", err)
		return nil, err
	}

	song, err := s.songRepo.GetSongByID(ctx, id)
	if err != nil {
		s.logger.Error("error getting song by ID", "id", id, "error", err)
		return nil, err
	}

	if song == nil {
		err := fmt.Errorf("song with ID %d not found", id)
		s.logger.Warn("song not found", "id", id)
		return nil, err
	}

	return song, nil
}

// CreateSong валидирует входные данные и вызывает репозиторий для создания новой песни.
func (s *SongServiceImpl) CreateSong(ctx context.Context, song *entities.Song) error {
	if err := validateSong(song); err != nil {
		s.logger.Error("validation error while creating song", "error", err)
		return err
	}

	err := s.songRepo.CreateSong(ctx, song)
	if err != nil {
		s.logger.Error("error creating song", "song", song, "error", err)
	}
	return err
}

// UpdateSong валидирует данные и вызывает репозиторий для обновления песни.
func (s *SongServiceImpl) UpdateSong(ctx context.Context, id int, song *entities.Song) error {
	if id <= 0 {
		err := errors.New("invalid song ID")
		s.logger.Error("invalid song ID", "error", err)
		return err
	}

	if err := validateSong(song); err != nil {
		s.logger.Error("validation error while updating song", "error", err)
		return err
	}

	err := s.songRepo.UpdateSong(ctx, id, song)
	if err != nil {
		s.logger.Error("error updating song", "songID", id, "song", song, "error", err)
	}
	return err
}

// DeleteSong валидирует ID перед удалением песни.
func (s *SongServiceImpl) DeleteSong(ctx context.Context, id int) error {
	if id <= 0 {
		err := errors.New("invalid song ID")
		s.logger.Error("invalid song ID", "error", err)
		return err
	}

	err := s.songRepo.DeleteSong(ctx, id)
	if err != nil {
		s.logger.Error("error deleting song", "songID", id, "error", err)
	}
	return err
}

// GetSongDetails получает детали о песне из внешнего API
func (s *SongServiceImpl) GetSongDetails(ctx context.Context, group, song string) (*external_api.SongDetail, error) {
	return s.apiClient.GetSongDetails(ctx, group, song)
}

func validateSong(song *entities.Song) error {
	if song.Group == "" {
		return errors.New("group cannot be empty")
	}
	if song.Song == "" {
		return errors.New("song name cannot be empty")
	}
	if !isValidReleaseDate(song.ReleaseDate) {
		return errors.New("invalid release date format")
	}
	if song.Text == "" {
		return errors.New("song text cannot be empty")
	}
	if !isValidURL(song.Link) {
		return errors.New("invalid song link URL")
	}
	return nil
}

func isValidReleaseDate(date string) bool {
	return len(date) == 10
}

func isValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}
	return true
}
