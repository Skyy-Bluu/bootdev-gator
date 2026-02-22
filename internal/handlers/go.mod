module github.com/Skyy-Bluu/bootdev-gator/internal/handlers

go 1.25.4

replace github.com/Skyy-Bluu/bootdev-gator/internal/database v0.0.0 => ../database
replace github.com/Skyy-Bluu/bootdev-gator/internal/config v0.0.0 => ../config
replace github.com/Skyy-Bluu/bootdev-gator/internal/rss v0.0.0 => ../rss
require github.com/Skyy-Bluu/bootdev-gator/internal/database v0.0.0
require github.com/Skyy-Bluu/bootdev-gator/internal/config v0.0.0
require github.com/Skyy-Bluu/bootdev-gator/internal/rss v0.0.0

require(
    github.com/google/uuid v1.6.0 // indirect
    github.com/lib/pq v1.11.2 // indirect
)
