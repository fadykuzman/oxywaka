# Wakapi Architecture - Quick Reference

**Purpose:** Fast visual reference for understanding system flow

---

## 📊 High-Level System Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                    WAKAPI SYSTEM OVERVIEW                            │
└─────────────────────────────────────────────────────────────────────┘

   IDE/Editor Plugin                                  User Web Browser
         │                                                   │
         │ API Key Auth                                     │ Session Auth
         ▼                                                   ▼
   ┌──────────────────────────────────────────────────────────────┐
   │                    Chi Router + Middleware                    │
   │  ┌────────────┬─────────────────┬────────────────────────┐   │
   │  │ /api/*     │ /compat/waka/*  │ /summary, /settings    │   │
   │  │ (Native)   │ (Compatible)    │ (MVC Web)              │   │
   │  └────────────┴─────────────────┴────────────────────────┘   │
   └──────────────────────────────────────────────────────────────┘
                              │
         ┌────────────────────┼────────────────────┐
         ▼                    ▼                    ▼
   ┌──────────┐      ┌─────────────┐      ┌──────────────┐
   │Heartbeat │      │  Summary    │      │   User       │
   │Service   │──┬──▶│  Service    │      │   Service    │
   └──────────┘  │   └─────────────┘      └──────────────┘
                 │
                 │   ┌─────────────┐      ┌──────────────┐
                 ├──▶│  Duration   │      │ Aggregation  │
                 │   │  Service    │      │  Service     │
                 │   └─────────────┘      └──────────────┘
                 │
                 ▼
         ┌──────────────┐
         │  Event Bus   │──┐
         │  (Pub/Sub)   │  │  Real-time events
         └──────────────┘  │  for loose coupling
                           │
         ┌─────────────────┴────────────────┐
         ▼                                  ▼
   ┌──────────┐                      ┌──────────┐
   │ Cache    │                      │Job Queue │
   │ Updates  │                      │(Artifex) │
   └──────────┘                      └──────────┘
         │                                  │
         ▼                                  ▼
   ┌──────────────────────────────────────────┐
   │         Repositories (GORM)              │
   └──────────────────────────────────────────┘
                     │
                     ▼
         ┌──────────────────────┐
         │  Database            │
         │  (SQLite/MySQL/      │
         │   PostgreSQL)        │
         └──────────────────────┘
```

---

## 🎯 Critical Data Flow: Heartbeat Ingestion

```
Editor Keystroke → WakaTime CLI → Wakapi API → Service → Database → Dashboard
     (1ms)            (buffered)     (HTTP)     (process)   (persist)   (view)

1. IDE Plugin captures activity
2. WakaTime CLI buffers & batches
3. POST /api/heartbeat (JSON)
4. HeartbeatService:
   - Validate & sanitize
   - Assign category (coding/browsing)
   - Hash for deduplication
   - Apply language mappings
5. Insert into heartbeats table
6. Publish event to EventBus
7. Update caches
8. Return 201 Created
```

---

## 📁 Directory Structure at a Glance

```
oxywaka/
├── main.go                 ⭐ Entry point, DI container
├── config/                 📝 Configuration management
│   ├── config.go
│   ├── eventbus.go         🔔 Event bus setup
│   └── jobqueue.go         ⏰ Job scheduler
│
├── models/                 📊 Data models (GORM)
│   ├── heartbeat.go        ⭐ Core tracking entity
│   ├── summary.go          📈 Aggregated stats
│   ├── duration.go         ⏱️  Coding sessions
│   ├── user.go             👤 User accounts
│   └── compat/             🔌 WakaTime compatibility
│
├── repositories/           💾 Data access layer
│   ├── heartbeat.go
│   ├── summary.go
│   └── base.go
│
├── services/               🧠 Business logic
│   ├── heartbeat.go        ⭐ Heartbeat processing
│   ├── summary.go          📊 Aggregation
│   ├── aggregation.go      🔄 Batch jobs
│   ├── duration.go         ⏱️  Session computation
│   ├── user.go             👤 User management
│   └── services.go         📋 Interface definitions
│
├── routes/                 🌐 HTTP handlers
│   ├── api/                🔌 REST API
│   │   ├── heartbeat.go    ⭐ POST /api/heartbeat
│   │   └── summary.go      📊 GET /api/summary
│   ├── compat/             🔌 WakaTime compat
│   │   └── wakatime/v1/
│   └── *.go                🖥️  MVC web pages
│
├── middlewares/            🛡️  Request processing
│   ├── authenticate.go     🔐 Auth logic
│   └── logging.go          📝 Request logging
│
├── migrations/             🗄️  Database schema
│   └── *.go                📝 Migration files
│
├── views/                  🎨 HTML templates
└── static/                 📦 CSS, JS, images
```

---

## 🔄 Data Flow Patterns

### Pattern 1: Heartbeat → Summary (Batch)
```
Raw Heartbeats ──┐
                 │ Daily Aggregation Job (02:15)
                 ├──▶ Group by entity types
                 ├──▶ Calculate totals
                 ├──▶ Merge aliases
                 └──▶ Store in summaries table

Why: Dashboard needs fast queries, not real-time aggregation
```

### Pattern 2: Heartbeats → Durations (Computed)
```
Heartbeats (time-ordered) ──┐
                            │ Group into sessions
                            ├──▶ Max gap: 10 minutes
                            ├──▶ Same project/language
                            └──▶ Store as Duration record

Why: Detect continuous coding sessions ("deep work")
```

### Pattern 3: Real-time Cache Updates
```
New Heartbeat ──▶ Publish Event ──▶ Subscribers ──┐
                                                   ├──▶ Increment user count
                                                   ├──▶ Invalidate project stats
                                                   └──▶ Clear summary cache

Why: Keep caches fresh without coupling services
```

---

## 🏗️ Layer Responsibilities

| Layer | Responsibility | Example |
|-------|---------------|---------|
| **Routes/Handlers** | HTTP I/O, validation | Parse JSON, return 201 |
| **Services** | Business logic, orchestration | Dedupe, categorize, publish events |
| **Repositories** | Database queries | `INSERT INTO heartbeats` |
| **Models** | Data structures | `Heartbeat`, `Summary` |
| **Middlewares** | Cross-cutting concerns | Auth, logging, CORS |
| **Config** | Infrastructure setup | DB, EventBus, Scheduler |

---

## 🎨 Service Dependency Map (Simplified)

```
UserService ─────────┬─────────────┬──────────────┐
                     │             │              │
HeartbeatService ────┤             │              │
    │                │             │              │
    ├──▶ LanguageMappingService   │              │
    │                              │              │
DurationService ─────┼─────────────┤              │
    │                │             │              │
    └─────────────┬──┘             │              │
                  │                │              │
SummaryService ───┴────────────────┤              │
    │                              │              │
    ├──▶ AliasService              │              │
    ├──▶ ProjectLabelService       │              │
    │                              │              │
AggregationService ────────────────┴──────────────┘
    │
    └──▶ Generates summaries daily
```

---

## 📊 Database Tables (Core 4)

### 1. **heartbeats** (Raw activity)
```sql
id, user_id, entity, type, category, project, branch,
language, is_write, editor, operating_system, machine,
time, hash, created_at
```
**Key Fields:**
- `category`: "coding" | "browsing"
- `is_write`: Active vs passive
- `time`: Millisecond precision

### 2. **durations** (Sessions)
```sql
id, user_id, time, duration, project, language, editor,
operating_system, machine, branch, group_hash
```
**Purpose:** Pre-computed coding sessions

### 3. **summaries** (Aggregations)
```sql
id, user_id, from_time, to_time, num_heartbeats
```
**Purpose:** Container for summary items

### 4. **summary_items** (Breakdowns)
```sql
id, summary_id, type, key, total
```
**Example:** type=0 (project), key="wakapi", total=3600 (seconds)

---

## ⚡ Performance Optimizations

| Technique | Implementation | Benefit |
|-----------|---------------|---------|
| **Caching** | In-memory (go-cache) | Avoid repeated DB queries |
| **Batch Inserts** | InsertBatch() | Reduce DB round-trips |
| **Pre-aggregation** | Daily summary generation | Fast dashboard loads |
| **Streaming** | Channel-based iteration | Handle large datasets |
| **Indexes** | user_id+time, hash | Fast lookups |
| **Connection Pool** | MaxConn = 10 | Reuse connections |

---

## 🔐 Authentication Flow (API Key)

```
1. Client sends: Authorization: Basic <base64(apiKey)>
2. AuthMiddleware extracts API key
3. UserService.GetUserByKey(apiKey)
   ├──▶ ApiKeyService.GetByApiKey(apiKey)
   └──▶ UserRepository.GetUserById(userId)
4. User stored in request context
5. Handler accesses via middlewares.GetPrincipal(r)
```

---

## 🎯 Extension Points for OxyWaka

### 1. **Add Insights Service**
```go
// services/insights.go
type InsightsService struct {
    heartbeatSvc IHeartbeatService
    summarySvc   ISummaryService
    durationSvc  IDurationService
}
```

### 2. **Subscribe to Events**
```go
// Subscribe to real-time heartbeats
eventBus.Subscribe(0, EventHeartbeatCreate)
```

### 3. **Add API Routes**
```go
// routes/api/insights.go
router.Get("/api/insights/score", handler.GetScore)
```

### 4. **Extend Database**
```sql
CREATE TABLE insight_scores (
    user_id, date, productivity_score,
    deep_work_minutes, context_switches
)
```

### 5. **Background Jobs**
```go
// Schedule insights computation
insightsService.Schedule() // Add to main.go:202
```

---

## 🧪 Testing Patterns

```go
// Services use interfaces
type IHeartbeatService interface {
    Insert(*Heartbeat) error
    // ... other methods
}

// Easy to mock
type MockHeartbeatService struct {
    mock.Mock
}

// Test with dependency injection
func TestSummary(t *testing.T) {
    mockHB := new(MockHeartbeatService)
    summarySvc := NewSummaryService(mockHB, ...)
    // Test business logic
}
```

---

## 📈 Monitoring & Observability

| Feature | Implementation | Location |
|---------|---------------|----------|
| **Logging** | slog (structured) | Throughout |
| **Metrics** | Prometheus exports | `/api/metrics` |
| **Health Check** | Database ping | `/api/health` |
| **Diagnostics** | Error reporting | `/api/diagnostics` |
| **Profiling** | pprof (optional) | Config flag |

---

## 🚀 Startup Sequence (main.go)

```
1. Parse flags (--config, --version)
2. Load configuration (config.yml)
3. Connect to database (GORM)
4. Run migrations (Auto-migrate models)
5. Initialize repositories (12 repos)
6. Initialize services (11 services)
7. Schedule background jobs (5 jobs)
8. Initialize routes & handlers
9. Start HTTP server
   - IPv4: 127.0.0.1:3000
   - IPv6: [::1]:3000
   - Unix socket: optional
```

---

## 📝 Key Takeaways

1. **Clean Architecture**: Layered, interface-based, testable
2. **Event-Driven**: Pub/sub for loose coupling
3. **Batch-Oriented**: Pre-compute summaries, not real-time
4. **WakaTime Compatible**: Drop-in replacement
5. **Extensible**: Clear patterns for adding features

---

**Next Steps:**
1. Read `ARCHITECTURE.md` for detailed diagrams
2. Explore `main.go` to see initialization
3. Trace heartbeat flow: routes → services → repos
4. Review `models/heartbeat.go` for data structure
5. Plan insights module integration points

---

**Generated for OxyWaka fork evaluation** 🚀
