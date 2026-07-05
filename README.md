# StreamIt API

Backend API for **StreamIt**, a scalable video streaming platform built with Go.

The API handles authentication, video uploads, metadata management, channel management, subscriptions, likes, comments, asynchronous video processing, and an event-driven view counting pipeline.

---

## Tech Stack

- Go
- Gin
- MongoDB
- Redis
- Asynq
- AWS S3
- FFmpeg
- Docker

---

## Features

- JWT authentication
- Video upload with presigned S3 URLs
- Asynchronous video processing using Asynq
- HLS video streaming support
- Channel profiles
- Subscribe / Unsubscribe system
- Like & Unlike videos
- Comment system
- Search API
- Background workers
- Event-driven view counting pipeline (Redis Streams)
- RESTful API
- Dockerized

---

## Architecture

- Gin HTTP API
- MongoDB for persistent storage
- Redis for queues, streams, and caching
- Asynq workers for video processing
- Redis Streams for asynchronous view event ingestion
- AWS S3 for video and asset storage

---

## Status

### Implemented

- Authentication
- Video uploads
- HLS processing pipeline
- Channel management
- Subscription system
- Likes
- Comments
- Search
- Redis Streams producer & consumer for view counting

### In Progress

- 30-second watch validation
- View deduplication
- HyperLogLog analytics
- Hot counters
- Write-behind aggregation
- CDN integration