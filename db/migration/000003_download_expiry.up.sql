ALTER TABLE downloads ADD COLUMN playlist_id BIGINT UNSIGNED DEFAULT 0 AFTER audio_id;
ALTER TABLE downloads ADD COLUMN expires_at DATETIME(3) NOT NULL DEFAULT (CURRENT_TIMESTAMP(3) + INTERVAL 30 DAY) AFTER file_size;
CREATE INDEX idx_download_expires ON downloads(expires_at);
