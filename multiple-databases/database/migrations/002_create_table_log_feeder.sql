CREATE TABLE IF NOT EXISTS "log_feeder" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "status" VARCHAR(100) NOT NULL,
    "timestamp" CURRENT_TIMESTAMP NOT NULL
);