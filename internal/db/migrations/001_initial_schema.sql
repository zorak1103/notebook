-- Meetings table
CREATE TABLE meetings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_by TEXT NOT NULL,        -- Tailscale user (e.g. "user@example.com")
    subject TEXT NOT NULL,            -- Subject (required)
    meeting_date TEXT NOT NULL,       -- Date in YYYY-MM-DD format
    start_time TEXT NOT NULL,         -- Start time in HH:MM format
    end_time TEXT,                    -- End time (optional)
    participants TEXT,                -- Participants (optional, comma-separated or JSON)
    summary TEXT,                     -- Summary (optional)
    keywords TEXT,                    -- Keywords (optional, comma-separated)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for sorting by date/time
CREATE INDEX idx_meetings_date_time ON meetings(meeting_date DESC, start_time DESC);

-- Index for searching in subject
CREATE INDEX idx_meetings_subject ON meetings(subject);

-- Notes table
CREATE TABLE notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    meeting_id INTEGER NOT NULL,
    note_number INTEGER NOT NULL,     -- Sequential number within meeting
    content TEXT NOT NULL,            -- Note content (required)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE CASCADE,
    UNIQUE(meeting_id, note_number)   -- Unique numbers per meeting
);

-- Index for efficient note queries per meeting
CREATE INDEX idx_notes_meeting ON notes(meeting_id, note_number);

-- Config table (Key-Value store)
CREATE TABLE config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Default config keys
INSERT INTO config (key, value) VALUES
    ('llm_provider_url', ''),
    ('llm_api_key', ''),
    ('llm_model', '');

-- Trigger for updated_at (Meetings)
CREATE TRIGGER update_meetings_timestamp
AFTER UPDATE ON meetings
FOR EACH ROW
BEGIN
    UPDATE meetings SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

-- Trigger for updated_at (Notes)
CREATE TRIGGER update_notes_timestamp
AFTER UPDATE ON notes
FOR EACH ROW
BEGIN
    UPDATE notes SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

-- Trigger for updated_at (Config)
CREATE TRIGGER update_config_timestamp
AFTER UPDATE ON config
FOR EACH ROW
BEGIN
    UPDATE config SET updated_at = CURRENT_TIMESTAMP WHERE key = OLD.key;
END;
