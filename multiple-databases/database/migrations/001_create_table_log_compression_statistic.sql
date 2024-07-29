-- log_compression_statistic definition

-- UP
CREATE TABLE IF NOT EXISTS "log_compression_statistic" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "job_id" VARCHAR(100) NOT NULL,
    "task_id" VARCHAR(100) NOT NULL,
    "user_id" VARCHAR(100) NOT NULL,
    "compression_id" VARCHAR(100) NOT NULL,
    "machine_id" VARCHAR(10) NOT NULL,
    "compressor_id" VARCHAR(50) NOT NULL,
    "filename" VARCHAR(255) NOT NULL,
    "format_file" VARCHAR(10) NOT NULL,
    "original_size" FLOAT NOT NULL,
    "compressed_size" FLOAT NOT NULL,
    "compressed_duration" FLOAT NOT NULL,
    "space_saving_percentage" FLOAT NOT NULL,
    "started_at" TIMESTAMP NOT NULL,
    "finished_at" TIMESTAMP NOT NULL,
    "timestamp" CURRENT_TIMESTAMP NOT NULL
);