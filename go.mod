module github.com/Skyy-Bluu/bootdev-gator

go 1.25.4

replace github.com/Skyy-Bluu/bootdev-gator/internal/config v0.0.0 => ./internal/config

replace github.com/Skyy-Bluu/bootdev-gator/internal/database v0.0.0 => ./internal/database

replace github.com/Skyy-Bluu/bootdev-gator/internal/handlers v0.0.0 => ./internal/handlers

replace github.com/Skyy-Bluu/bootdev-gator/internal/rss v0.0.0 => ./internal/rss

require github.com/Skyy-Bluu/bootdev-gator/internal/config v0.0.0

require github.com/Skyy-Bluu/bootdev-gator/internal/database v0.0.0

require (
	github.com/Skyy-Bluu/bootdev-gator/internal/handlers v0.0.0
	github.com/lib/pq v1.11.2
)

require (
	github.com/Skyy-Bluu/bootdev-gator/internal/rss v0.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
)
