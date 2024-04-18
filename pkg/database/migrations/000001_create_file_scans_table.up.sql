CREATE TABLE IF NOT EXISTS file_scans(
    id SERIAL PRIMARY KEY,
    file_name VARCHAR NOT NULL,
    file_size BIGINT,
    upload_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    virustotal_scan_id VARCHAR,
    scan_status VARCHAR,
    scan_result JSONB,
    last_checked TIMESTAMP
);
