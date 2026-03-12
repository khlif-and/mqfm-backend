CREATE TABLE IF NOT EXISTS radios (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    thumbnail VARCHAR(500),
    dominant_color VARCHAR(7),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    INDEX idx_radio_active (is_active),
    INDEX idx_radio_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS radio_audios (
    radio_id BIGINT UNSIGNED NOT NULL,
    audio_id BIGINT UNSIGNED NOT NULL,
    order_num INT NOT NULL DEFAULT 0,
    PRIMARY KEY (radio_id, audio_id),
    CONSTRAINT fk_radio_audios_radio FOREIGN KEY (radio_id) REFERENCES radios(id) ON DELETE CASCADE,
    CONSTRAINT fk_radio_audios_audio FOREIGN KEY (audio_id) REFERENCES audios(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
