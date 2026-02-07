module gitlab.com/4uvirik/my-platform/services/notification-service

go 1.24.4

replace gitlab.com/4uvirik/my-platform => ../../

require (
	github.com/segmentio/kafka-go v0.4.50
	gitlab.com/4uvirik/my-platform v0.0.0-00010101000000-000000000000
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)
