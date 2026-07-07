# Video Upload & Processing Architecture

![video-processing](../../assets/worker.png)

## Overview

The video processing pipeline is built using an asynchronous, event-driven architecture to keep uploads responsive while offloading CPU-intensive transcoding to background workers.

## Flow

```text
Client
   │
   ▼
Upload Request
   │
   ▼
Generate S3 Presigned URL
   │
   ▼
Upload Raw Video to S3
   │
   ▼
Confirm Upload API
   │
   ▼
Store Metadata (MongoDB)
   │
   ▼
Enqueue Processing Job (Asynq + Redis)
   │
   ▼
Go Worker
   │
   ├── Download Raw Video
   ├── Transcode with FFmpeg
   ├── Generate HLS Segments
   ├── Generate Thumbnail
   └── Upload Processed Assets to S3
   │
   ▼
Worker Callback API
   │
   ▼
Update Video Metadata (MongoDB)
   │
   ▼
Video Ready for Streaming
```

## Components

| Component | Responsibility |
|----------|----------------|
| **Next.js Client** | Requests upload URL and uploads video directly to S3 |
| **Go API (Gin)** | Generates presigned URLs, stores metadata, exposes APIs |
| **AWS S3** | Stores raw videos, HLS segments, and thumbnails |
| **Redis** | Message broker for background jobs |
| **Asynq** | Job queue and scheduling |
| **Go Worker** | Processes videos asynchronously |
| **FFmpeg** | HLS transcoding and thumbnail generation |
| **MongoDB** | Stores video metadata and processing status |

## Key Features

- Direct browser-to-S3 uploads using presigned URLs
- Asynchronous video processing
- HLS adaptive streaming
- Automatic thumbnail generation
- Background job processing with Redis & Asynq
- Processing progress tracking
- Fault isolation between API and workers
- Horizontally scalable worker architecture