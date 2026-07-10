# StreamIt API

Backend API for **StreamIt**, a distributed video streaming platform built with Go.

The API powers authentication, media uploads, channel management, social interactions, asynchronous video processing, and a scalable event-driven view analytics pipeline.

---

## Tech Stack

- Go
- Gin
- MongoDB
- Redis
- Redis Streams
- Asynq
- AWS S3
- Docker

---

## Features

- JWT authentication
- Video uploads via presigned S3 URLs
- Asynchronous video processing with Asynq
- Adaptive HLS streaming support
- Channel profiles
- Subscribe / Unsubscribe system
- Like & Unlike videos
- Comment system
- Search API
- Distributed view counting pipeline
- Heartbeat-based watch validation
- View deduplication
- HyperLogLog unique viewer analytics
- Hot counters & write-behind aggregation
- RESTful API
- Dockerized

---

## Architecture

- Gin HTTP API
- MongoDB for persistent storage
- Redis for queues, streams, and hot counters
- Redis Streams & Consumer Groups for event processing
- Asynq workers for asynchronous video processing
- AWS S3 for media storage
