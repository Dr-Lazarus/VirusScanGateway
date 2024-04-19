CREATE TABLE virustotal_reports (
    file_type TEXT,
    sha256 TEXT PRIMARY KEY NOT NULL,
    sha1 TEXT NOT NULL,
    md5 TEXT NOT NULL,
    meaningful_name TEXT,
    type_description TEXT,
    type_extension TEXT,
    last_modification_date TIMESTAMPTZ NOT NULL,
    first_submission_date TIMESTAMPTZ NOT NULL,
    last_submission_date TIMESTAMPTZ NOT NULL,
    size BIGINT NOT NULL,
    last_analysis_date TIMESTAMPTZ NOT NULL,
    last_analysis_stats JSONB NOT NULL,
    last_analysis_results JSONB NOT NULL
);
