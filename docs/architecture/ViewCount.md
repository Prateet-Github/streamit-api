# View Count Architecture

```text
[ Client Video Player ]
          │
          │
          ▼
┌──────────────────────────────────────────────┐
│ HEARTBEAT TRACKER                            │
│                                              │
│ - Starts when playback begins                │
│ - Fires every 10 seconds                     │
│ - Stops on pause/end                         │
│ - Resumes on play                            │
│ - Uses navigator.sendBeacon()                │
└──────────────────────────────────────────────┘
          │
          │ POST /api/videos/:id/view
          ▼
┌──────────────────────────────────────────────┐
│ 1. INGESTION LAYER: Gin API                  │
│                                              │
│  - Validate payload                          │
│  - Extract viewer identifier                 │
│    (User ID or Anonymous Session ID)         │
│  - Sliding-window rate limiter               │
│    (Session/User ID + IP)                    │
│  - Return immediately (<2ms)                 │
└──────────────────────────────────────────────┘
          │
          ▼
┌──────────────────────────────────────────────┐
│ 2. EVENT STREAM: Redis Streams               │
│                                              │
│ Key: streamit:view_events                    │
│                                              │
│ - Durable event queue                        │
│ - Decouples HTTP from processing             │
└──────────────────────────────────────────────┘
          │
          │ Consumer Groups (Batch Size = 10)
          ▼
┌──────────────────────────────────────────────┐
│ 3. GO WORKER POOL                            │
│                                              │
│ - XREADGROUP batch consumer                  │
│ - Parallel workers                           │
│ - Pipeline Redis operations                  │
│ - XACK only after successful processing      │
│                                              │
│ Recovery:                                    │
│ XPENDING -> XCLAIM -> Reprocess -> XACK      │
└──────────────────────────────────────────────┘
          │
          ▼
┌──────────────────────────────────────────────┐
│ 4. WATCH SESSION VALIDATION                  │
│                                              │
│ Heartbeats Received:                         │
│                                              │
│ 10s → SADD 10                                │
│ 20s → SADD 20                                │
│ 30s → SADD 30                                │
│                                              │
│ Redis Set                                    │
│ track:<viewer>:<video>                       │
│                                              │
│ EXPIRE 5m                                    │
│                                              │
│ SCARD == 3 ?                                 │
└──────────────────────────────────────────────┘
          │
          ├────────────── No ────────────────► Drop
          │
          ▼ Yes (30 seconds watched)
┌──────────────────────────────────────────────┐
│ 5. EXACT DEDUPLICATION                       │
│                                              │
│ Redis String                                 │
│                                              │
│ SET view:<video>:<viewer>                    │
│     1 EX 4h NX                               │
│                                              │
│ Prevent duplicate views within 4 hours       │
└──────────────────────────────────────────────┘
          │
          ├────────── Key Exists ────────────► Drop
          │
          ▼ New View
          │
          ├──────────────────────────────┐
          │                              │
          ▼                              ▼
┌──────────────────────────────┐  ┌──────────────────────────────┐
│ 6A. UNIQUE ANALYTICS         │  │ 6B. HOT COUNTERS             │
│                              │  │                              │
│ Redis HyperLogLog            │  │ Redis Hash                   │
│                              │  │                              │
│ PFADD                        │  │ HINCRBY                      │
│ unique:<date>:<videoId>      │  │ hot_counters videoId +1      │
│                              │  │                              │
│ Approx. Daily Unique Viewers │  │ Real-Time View Counter       │
└──────────────────────────────┘  └──────────────────────────────┘
                                           │
                                           ▼
───────────────────────────────────────────────────────────────
          Asynq Cron (Every 2–5 Minutes)
───────────────────────────────────────────────────────────────
                     │
                     ▼
┌──────────────────────────────────────────────┐
│ 7. WRITE-BEHIND AGGREGATION                  │
│                                              │
│ HSCAN hot_counters                           │
│                                              │
│ Build:                                       │
│ { videoId -> incrementCount }                │
└──────────────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────┐
│ 8. MONGODB BULK WRITE                        │
│                                              │
│ BulkWrite()                                  │
│                                              │
│ $inc views                                   │
│                                              │
│ Single atomic batch update                   │
└──────────────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────┐
│ 9. SAFE COUNTER RECONCILIATION               │
│                                              │
│ HINCRBY hot_counters                         │
│         videoId -processedCount              │
│                                              │
│ Never HDEL                                   │
│                                              │
│ Prevents race conditions with new traffic    │
└──────────────────────────────────────────────┘
```

## View Definition

A view is counted **once per video per viewer (authenticated user or anonymous visitor)** after **30 cumulative seconds of watch time**, provided another counted view for the same viewer **has not occurred within the last 4 hours**.

Anonymous viewers are identified using a persistent session identifier stored on the client.

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| API Layer | Go + Gin |
| Event Queue | Redis Streams |
| Worker Pool | Go Goroutines + Redis Consumer Groups |
| Worker Recovery | Redis `XPENDING` + `XCLAIM` |
| Watch Validation | Redis Sets |
| Exact Deduplication | Redis `SET NX EX` |
| Unique Viewer Analytics | Redis HyperLogLog |
| Hot Counters | Redis Hash |
| Scheduler | Asynq Cron |
| Persistent Storage | MongoDB `BulkWrite()` |
| Rate Limiting | Sliding Window (Session/User ID + IP) |

---

## Why this Architecture?

- **Redis Streams** decouple request handling from processing.
- **Consumer Groups** enable scalable parallel workers with failure recovery.
- **Redis Sets** validate 30 seconds of cumulative watch time.
- **Redis `SET NX EX`** guarantees exact view deduplication within a configurable time window.
- **HyperLogLog** provides memory-efficient approximate unique viewer analytics.
- **Redis Hashes** serve real-time view counts without hitting MongoDB.
- **Asynq Cron** batches writes, dramatically reducing database load.
- **MongoDB BulkWrite** persists view counts efficiently using atomic `$inc` operations.
