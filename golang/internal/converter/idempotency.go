package converter

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"
)

func IsProcessed(db *sql.DB, videoId int) bool {
	var IsProcessed bool
	query := "SELECT EXISTS(SELECT 1 FROM processed_videos WHERE video_id = $1 AND status='SUCCESS')"
	err := db.QueryRow(query, videoId).Scan(&IsProcessed)
	if err != nil {
		slog.Error("Error checking if video is processed", slog.Int("video_id", videoId))
		return false
	}
	return IsProcessed
}

func MarkProcessed(db *sql.DB, videoId int) error {
	query := "INSERT INTO processed_videos (video_id, status, processed_at) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, videoId, "SUCCESS", time.Now())
	if err != nil {
		slog.Error("Error marking video as processed", slog.Int("video_id", videoId))
		return err
	}
	return nil
}

func RegisterError(db *sql.DB, errorData map[string]any, err error) {
	serializedError, _ := json.Marshal(errorData)
	query := "INSERT INTO process_errors_log (error_details, created_at) VALUES ($1, $2)"
	_, dbErr := db.Exec(query, serializedError, time.Now())
	if dbErr != nil {
		slog.Error("Error registering error", slog.String("error_details", string(serializedError)))
	}
}
