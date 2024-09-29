package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Conn *pgxpool.Pool
}

func (d *DB) Close() error {
	d.Conn.Close()
	return nil
}

func Init(DBURI string) (*DB, error) {
	var result DB
	var err error

	result.Conn, err = openPGConnection(DBURI)
	if err != nil {
		return nil, err
	}
	result.Conn.Config().MaxConns = int32(runtime.NumCPU())

	time.Sleep(5 * time.Second)
	if err = testConnection(result.Conn); err != nil {
		return nil, err
	}

	time.Sleep(5 * time.Second)
	go connectionWorker(result.Conn, DBURI)

	return &result, nil
}

func openPGConnection(dbURI string) (*pgxpool.Pool, error) {
	if err := migrateUp(dbURI); err != nil {
		return nil, fmt.Errorf("migrations error: %w", err)
	}
	conn, err := pgxpool.New(context.Background(), dbURI)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func testConnection(db *pgxpool.Pool) error {
	err := db.Ping(context.Background())
	errPing := errors.New("can't ping postgres")
	if err != nil {
		return errors.Join(err, errPing)
	}
	return nil
}

func connectionWorker(conn *pgxpool.Pool, dbURI string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := testConnection(conn); err != nil {
				fmt.Printf("lost connection to database: %v\n", err)
				var newConn *pgxpool.Pool
				newConn, err = openPGConnection(dbURI)
				if err != nil {
					fmt.Errorf("failed to reconnect to PostgreSQL: %v. Exiting application", err)
					return
				}
				conn = newConn
				fmt.Println("successfully reconnected to PostgreSQL.")
			}
		}
	}
}

func migrateUp(dbURI string) error {

	m, err := migrate.New("file://migrations/", dbURI)

	if err != nil {
		return fmt.Errorf("new error: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("up migrations error: %w", err)
	}

	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		return fmt.Errorf("close source error: %w", sourceErr)
	}

	if dbErr != nil {
		return fmt.Errorf("close db error: %w", dbErr)
	}

	return nil
}
