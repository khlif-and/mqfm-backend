ALTER TABLE audios ADD COLUMN ogg_path VARCHAR(500) DEFAULT '' AFTER file_path;
ALTER TABLE audios ADD COLUMN file_size BIGINT DEFAULT 0 AFTER duration;
ALTER TABLE audios ADD COLUMN series_id BIGINT UNSIGNED DEFAULT 0 AFTER category_id;

ALTER TABLE playlists ADD COLUMN share_token VARCHAR(64) DEFAULT '' AFTER image_url;
ALTER TABLE playlists ADD COLUMN is_public TINYINT(1) DEFAULT 0 AFTER share_token;
CREATE UNIQUE INDEX idx_playlists_share_token ON playlists(share_token);

CREATE TABLE IF NOT EXISTS bookmarks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    audio_id BIGINT UNSIGNED NOT NULL,
    position_seconds INT NOT NULL DEFAULT 0,
    label VARCHAR(255) DEFAULT '',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_bookmark_user_audio (user_id, audio_id, position_seconds)
);

CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT,
    type VARCHAR(50) DEFAULT '',
    reference_id BIGINT UNSIGNED DEFAULT 0,
    is_read TINYINT(1) DEFAULT 0,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    INDEX idx_notifications_user_id (user_id),
    INDEX idx_notifications_type (type)
);

CREATE TABLE IF NOT EXISTS notification_settings (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    daily_reminder TINYINT(1) DEFAULT 1,
    new_content TINYINT(1) DEFAULT 1,
    event_reminder TINYINT(1) DEFAULT 1,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_notification_settings_user (user_id)
);

CREATE TABLE IF NOT EXISTS audio_progress (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    audio_id BIGINT UNSIGNED NOT NULL,
    last_position INT DEFAULT 0,
    duration INT DEFAULT 0,
    percentage DOUBLE DEFAULT 0,
    completed TINYINT(1) DEFAULT 0,
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_progress_user_audio (user_id, audio_id)
);

CREATE TABLE IF NOT EXISTS downloads (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    audio_id BIGINT UNSIGNED NOT NULL,
    file_size BIGINT DEFAULT 0,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_download_user_audio (user_id, audio_id)
);

CREATE TABLE IF NOT EXISTS listening_stats (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    audio_id BIGINT UNSIGNED NOT NULL,
    duration_seconds INT DEFAULT 0,
    listened_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    INDEX idx_stat_user_date (user_id, listened_at),
    INDEX idx_stat_audio (audio_id)
);

CREATE TABLE IF NOT EXISTS audio_clips (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    audio_id BIGINT UNSIGNED NOT NULL,
    start_time INT NOT NULL,
    end_time INT NOT NULL,
    clip_path VARCHAR(500) DEFAULT '',
    share_token VARCHAR(64) NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_clip_share_token (share_token),
    INDEX idx_clip_user (user_id),
    INDEX idx_clip_audio (audio_id)
);

CREATE TABLE IF NOT EXISTS events (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date DATETIME(3) NOT NULL,
    location VARCHAR(500) DEFAULT '',
    image VARCHAR(500) DEFAULT '',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    INDEX idx_events_date (event_date),
    INDEX idx_events_deleted_at (deleted_at)
);

CREATE TABLE IF NOT EXISTS event_rsvps (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    event_id BIGINT UNSIGNED NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_rsvp_user_event (user_id, event_id)
);

CREATE TABLE IF NOT EXISTS user_preferences (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    playback_speed DOUBLE DEFAULT 1.0,
    sleep_timer_minutes INT DEFAULT 0,
    auto_download_wifi TINYINT(1) DEFAULT 0,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_user_preferences_user (user_id)
);

CREATE TABLE IF NOT EXISTS audio_series (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    artist VARCHAR(255) DEFAULT '',
    image VARCHAR(500) DEFAULT '',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    INDEX idx_audio_series_deleted_at (deleted_at)
);

CREATE TABLE IF NOT EXISTS audio_series_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    series_id BIGINT UNSIGNED NOT NULL,
    audio_id BIGINT UNSIGNED NOT NULL,
    order_num INT NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_series_order (series_id, order_num),
    INDEX idx_series_item_audio (audio_id)
);

CREATE TABLE IF NOT EXISTS audio_votes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    audio_id BIGINT UNSIGNED NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_vote_user_audio (user_id, audio_id),
    INDEX idx_vote_audio (audio_id)
);

CREATE TABLE IF NOT EXISTS audio_rankings (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    audio_id BIGINT UNSIGNED NOT NULL,
    weekly_votes BIGINT DEFAULT 0,
    monthly_votes BIGINT DEFAULT 0,
    total_votes BIGINT DEFAULT 0,
    random_boost DOUBLE DEFAULT 0,
    weekly_rank INT DEFAULT 0,
    monthly_rank INT DEFAULT 0,
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_ranking_audio (audio_id),
    INDEX idx_ranking_weekly (weekly_rank),
    INDEX idx_ranking_monthly (monthly_rank)
);

CREATE TABLE IF NOT EXISTS favorite_artists (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    artist_name VARCHAR(255) NOT NULL,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_fav_user_artist (user_id, artist_name)
);

CREATE TABLE IF NOT EXISTS smart_resumes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    audio_id BIGINT UNSIGNED NOT NULL,
    playlist_id BIGINT UNSIGNED DEFAULT 0,
    position_seconds INT DEFAULT 0,
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_smart_resume_user (user_id)
);

CREATE TABLE IF NOT EXISTS user_locations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    latitude DOUBLE NOT NULL,
    longitude DOUBLE NOT NULL,
    city VARCHAR(255) DEFAULT '',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_user_location_user (user_id)
);

CREATE TABLE IF NOT EXISTS playlist_collaborators (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    playlist_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    role VARCHAR(20) DEFAULT 'contributor',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE INDEX idx_collab_playlist_user (playlist_id, user_id)
);
