# StreamIt API

![streamit-api](assets/worker.png)

Backend API for StreamIt, built with Go.

Handles authentication, video uploads, metadata management, presigned S3 uploads, background job scheduling, and communication with the video processing worker.

---

## Tech Stack

- Go
- Gin
- Asynq
- Ffmpeg
- Redis 
- AWS S3
- MongoDB

## Features

- JWT authentication
- Presigned S3 upload URLs
- Video metadata management
- Background job scheduling
- Worker callback endpoints
- RESTful API
- Dockerized