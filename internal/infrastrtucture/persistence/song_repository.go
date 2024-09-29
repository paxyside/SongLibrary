package persistence

import (
	"context"
	"effictiveMobile/internal/domain/entities"
	"effictiveMobile/pkg/database"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type SongRepository interface {
	GetSongs(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]entities.Song, error)
	GetSongByID(ctx context.Context, id int) (*entities.Song, error)
	CreateSong(ctx context.Context, song *entities.Song) error
	UpdateSong(ctx context.Context, id int, song *entities.Song) error
	DeleteSong(ctx context.Context, id int) error
}

type SongRepositoryImpl struct {
	db     *database.DB
	logger *slog.Logger
}

func NewSongRepository(db *database.DB, logger *slog.Logger) *SongRepositoryImpl {
	return &SongRepositoryImpl{
		db:     db,
		logger: logger.With(slog.String("repository", "SongRepository")),
	}
}

// GetSongs возвращает список песен с возможностью фильтрации и пагинации.
func (r *SongRepositoryImpl) GetSongs(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]entities.Song, error) {
	var songs []entities.Song

	query := "SELECT id, \"group\", song, release_date, text, link FROM songs WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	for field, value := range filter {
		query += " AND " + field + " = $" + string(argIndex)
		args = append(args, value)
		argIndex++
	}

	query += " LIMIT $" + string(argIndex) + " OFFSET $" + string(argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Conn.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("error querying songs", "error", err, "query", query)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song entities.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			r.logger.Error("error scanning song row", "error", err)
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

// GetSongByID возвращает одну песню по ID.
func (r *SongRepositoryImpl) GetSongByID(ctx context.Context, id int) (*entities.Song, error) {
	query := "SELECT id, \"group\", song, release_date, text, link FROM songs WHERE id = $1"
	row := r.db.Conn.QueryRow(ctx, query, id)

	var song entities.Song
	if err := row.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("error querying song by ID", "error", err, "id", id)
		return nil, err
	}

	return &song, nil
}

// CreateSong добавляет новую песню в базу данных.
func (r *SongRepositoryImpl) CreateSong(ctx context.Context, song *entities.Song) error {
	query := `INSERT INTO songs ("group", song, release_date, text, link) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Conn.Exec(ctx, query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		r.logger.Error("error creating song", "error", err, "song", song)
	}
	return err
}

// UpdateSong обновляет данные о песне по ID.
func (r *SongRepositoryImpl) UpdateSong(ctx context.Context, id int, song *entities.Song) error {
	query := `
		UPDATE songs
		SET "group" = $1, song = $2, release_date = $3, text = $4, link = $5
		WHERE id = $6
	`

	_, err := r.db.Conn.Exec(ctx, query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, id)
	if err != nil {
		r.logger.Error("error updating song", "error", err, "songID", id, "song", song)
	}
	return err
}

// DeleteSong удаляет песню по ID.
func (r *SongRepositoryImpl) DeleteSong(ctx context.Context, id int) error {
	query := "DELETE FROM songs WHERE id = $1"
	_, err := r.db.Conn.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("error deleting song", "error", err, "songID", id)
	}
	return err
}
