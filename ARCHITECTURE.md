# Wakapi Architecture Documentation

**Generated:** 2025-12-02
**Purpose:** Understanding system and data flow for OxyWaka fork evaluation

---

## Table of Contents
1. [System Architecture Overview](#1-system-architecture-overview)
2. [Layered Architecture](#2-layered-architecture)
3. [Heartbeat Data Flow](#3-heartbeat-data-flow-critical-path)
4. [Summary Generation Flow](#4-summary-generation-flow)
5. [Service Dependency Graph](#5-service-dependency-graph)
6. [Database Schema](#6-database-schema)
7. [Background Jobs & Scheduling](#7-background-jobs--scheduling)
8. [Authentication Flow](#8-authentication-flow)
9. [WakaTime Compatibility Layer](#9-wakatime-compatibility-layer)
10. [Event-Driven Architecture](#10-event-driven-architecture)

---

## 1. System Architecture Overview

```mermaid
graph TB
    subgraph "Client Layer"
        IDE[IDE/Editor Plugins<br/>VS Code, JetBrains, etc.]
        Browser[Browser Extensions<br/>Chrome, Firefox]
        CLI[WakaTime CLI]
    end

    subgraph "HTTP Layer"
        Router[Chi Router]
        Middleware[Middlewares<br/>Auth, Logging, CORS]
        API[API Routes<br/>/api/heartbeat]
        Compat[Compat Routes<br/>/compat/wakatime/v1]
        MVC[MVC Routes<br/>/summary, /settings]
    end

    subgraph "Business Logic Layer"
        HeartbeatSvc[HeartbeatService]
        SummarySvc[SummaryService]
        DurationSvc[DurationService]
        UserSvc[UserService]
        AggregationSvc[AggregationService]
        LeaderboardSvc[LeaderboardService]
        MailSvc[MailService]
    end

    subgraph "Data Access Layer"
        HeartbeatRepo[HeartbeatRepository]
        SummaryRepo[SummaryRepository]
        DurationRepo[DurationRepository]
        UserRepo[UserRepository]
        OtherRepos[Other Repositories]
    end

    subgraph "Persistence Layer"
        DB[(Database<br/>SQLite/MySQL/Postgres)]
    end

    subgraph "Cross-Cutting Concerns"
        EventBus[Event Bus<br/>Hub]
        JobQueue[Job Queue<br/>Artifex]
        Cache[In-Memory Cache]
    end

    IDE --> CLI
    Browser --> CLI
    CLI --> Router
    Router --> Middleware
    Middleware --> API
    Middleware --> Compat
    Middleware --> MVC

    API --> HeartbeatSvc
    Compat --> HeartbeatSvc
    API --> SummarySvc
    MVC --> SummarySvc

    HeartbeatSvc --> HeartbeatRepo
    SummarySvc --> SummaryRepo
    DurationSvc --> DurationRepo
    UserSvc --> UserRepo

    HeartbeatRepo --> DB
    SummaryRepo --> DB
    DurationRepo --> DB
    UserRepo --> DB
    OtherRepos --> DB

    HeartbeatSvc -.publishes.-> EventBus
    AggregationSvc -.subscribes.-> EventBus
    AggregationSvc --> JobQueue

    HeartbeatSvc --> Cache
    SummarySvc --> Cache

    style HeartbeatSvc fill:#ff9999
    style SummarySvc fill:#99ccff
    style EventBus fill:#ffeb99
    style JobQueue fill:#ffeb99
```

---

## 2. Layered Architecture

```mermaid
graph LR
    subgraph "Presentation Layer"
        Routes[Routes]
        Templates[HTML Templates]
        Static[Static Assets]
    end

    subgraph "Application Layer"
        Handlers[API Handlers]
        Middleware[Middlewares]
    end

    subgraph "Domain Layer"
        Services[Services<br/>Business Logic]
        Models[Domain Models]
        Interfaces[Service Interfaces]
    end

    subgraph "Infrastructure Layer"
        Repositories[Repositories]
        DB[Database<br/>GORM]
        Config[Configuration]
        EventBus[Event Bus]
        Jobs[Job Scheduler]
    end

    Routes --> Handlers
    Handlers --> Services
    Services --> Repositories
    Repositories --> DB
    Services --> EventBus
    Services --> Jobs

    Services --> Interfaces
    Models --> Services
    Config --> Services
    Templates --> Routes
```

**Key Principles:**
- **Dependency Direction**: Always flows inward (outer layers depend on inner layers)
- **Interface-Based**: All services expose interfaces for testability
- **Repository Pattern**: Data access abstracted through repositories
- **Event-Driven**: Services communicate via event bus for loose coupling

---

## 3. Heartbeat Data Flow (Critical Path)

This is the **most important flow** for OxyWaka - tracking every keystroke and activity.

```mermaid
sequenceDiagram
    participant Client as WakaTime Client
    participant Router as Chi Router
    participant Auth as Auth Middleware
    participant Handler as HeartbeatApiHandler
    participant HBSvc as HeartbeatService
    participant LangSvc as LanguageMappingService
    participant Repo as HeartbeatRepository
    participant DB as Database
    participant EventBus as Event Bus
    participant Cache as Cache

    Client->>Router: POST /api/heartbeat<br/>{entity, project, language, time, ...}
    Router->>Auth: Authenticate request
    Auth->>Auth: Validate API Key
    Auth-->>Router: User identified

    Router->>Handler: Handle request
    Handler->>Handler: Parse heartbeat JSON
    Handler->>Handler: Extract user agent<br/>(editor, OS)
    Handler->>Handler: Sanitize data<br/>(category assignment)

    Handler->>HBSvc: InsertBatch([heartbeat])
    HBSvc->>HBSvc: Deduplicate by hash
    HBSvc->>LangSvc: ResolveByUser(userID)
    LangSvc-->>HBSvc: Language mappings
    HBSvc->>HBSvc: Augment language<br/>(apply custom mappings)
    HBSvc->>HBSvc: Sanitize<br/>(set category: coding/browsing)
    HBSvc->>HBSvc: Hash heartbeat<br/>(prevent duplicates)

    HBSvc->>Repo: InsertBatch([heartbeat])
    Repo->>DB: INSERT INTO heartbeats
    DB-->>Repo: Success
    Repo-->>HBSvc: Success

    HBSvc->>EventBus: Publish(EventHeartbeatCreate)
    HBSvc->>Cache: Increment count cache
    HBSvc->>Cache: Invalidate project stats

    HBSvc-->>Handler: Success
    Handler-->>Client: 201 Created

    Note over EventBus: Other services can subscribe<br/>to heartbeat events for<br/>real-time processing
```

**Key Observations for OxyWaka:**
1. **Category field** is assigned during sanitization (coding vs browsing)
2. **Event bus** publishes heartbeat events - perfect hook for insights module
3. **Hash-based deduplication** prevents duplicate entries
4. **Language mapping** allows custom categorization

---

## 4. Summary Generation Flow

Summaries are pre-computed aggregations for fast dashboard rendering.

```mermaid
sequenceDiagram
    participant Cron as Cron Scheduler
    participant AggSvc as AggregationService
    participant SumSvc as SummaryService
    participant HBSvc as HeartbeatService
    participant DurSvc as DurationService
    participant AliasSvc as AliasService
    participant Repo as SummaryRepository
    participant DB as Database

    Cron->>AggSvc: Trigger (daily at 02:15)
    AggSvc->>AggSvc: Get users to process
    AggSvc->>SumSvc: GetLatestByUser()
    SumSvc-->>AggSvc: Last summary times

    loop For each user
        AggSvc->>AggSvc: Calculate date range<br/>(last summary → today)
        AggSvc->>SumSvc: Summarize(from, to, user)

        SumSvc->>HBSvc: GetAllWithin(from, to, user)
        HBSvc-->>SumSvc: Heartbeats[]

        SumSvc->>SumSvc: Group by:<br/>- Project<br/>- Language<br/>- Editor<br/>- OS<br/>- Machine<br/>- Category

        SumSvc->>DurSvc: Get(from, to, user)<br/>(compute durations)
        DurSvc->>DurSvc: Group heartbeats into<br/>continuous sessions
        DurSvc-->>SumSvc: Durations[]

        SumSvc->>SumSvc: Calculate totals per entity
        SumSvc->>AliasSvc: GetByUser(userID)
        AliasSvc-->>SumSvc: Project aliases
        SumSvc->>SumSvc: Merge aliased projects

        SumSvc->>Repo: Insert(summary)
        Repo->>DB: INSERT INTO summaries<br/>INSERT INTO summary_items
        DB-->>Repo: Success
        Repo-->>SumSvc: Success
    end

    AggSvc-->>Cron: Aggregation complete

    Note over DurSvc: Durations = continuous<br/>coding sessions with<br/>max gap of 10 minutes
```

**Summary Structure:**
- **Summary**: Time range container (from, to, user)
- **SummaryItems**: Individual aggregations (project: "wakapi" → 3600 seconds)
- **Types**: Project, Language, Editor, OS, Machine, Category, Branch, Entity

---

## 5. Service Dependency Graph

```mermaid
graph TD
    subgraph "Service Layer Dependencies"
        UserSvc[UserService]
        HeartbeatSvc[HeartbeatService]
        DurationSvc[DurationService]
        SummarySvc[SummaryService]
        AggregationSvc[AggregationService]
        AliasSvc[AliasService]
        LangMapSvc[LanguageMappingService]
        ProjectLabelSvc[ProjectLabelService]
        LeaderboardSvc[LeaderboardService]
        ReportSvc[ReportService]
        MailSvc[MailService]
        KeyValueSvc[KeyValueService]
        ApiKeySvc[ApiKeyService]
    end

    subgraph "Repository Layer"
        UserRepo[UserRepository]
        HeartbeatRepo[HeartbeatRepository]
        DurationRepo[DurationRepository]
        SummaryRepo[SummaryRepository]
        OtherRepos[Other Repositories]
    end

    HeartbeatSvc --> HeartbeatRepo
    HeartbeatSvc --> LangMapSvc

    DurationSvc --> DurationRepo
    DurationSvc --> HeartbeatSvc
    DurationSvc --> UserSvc
    DurationSvc --> LangMapSvc

    SummarySvc --> SummaryRepo
    SummarySvc --> HeartbeatSvc
    SummarySvc --> DurationSvc
    SummarySvc --> AliasSvc
    SummarySvc --> ProjectLabelSvc

    AggregationSvc --> UserSvc
    AggregationSvc --> SummarySvc
    AggregationSvc --> HeartbeatSvc
    AggregationSvc --> DurationSvc

    LeaderboardSvc --> SummarySvc
    LeaderboardSvc --> UserSvc

    ReportSvc --> SummarySvc
    ReportSvc --> UserSvc
    ReportSvc --> MailSvc

    UserSvc --> UserRepo
    UserSvc --> KeyValueSvc
    UserSvc --> MailSvc
    UserSvc --> ApiKeySvc

    style SummarySvc fill:#99ccff
    style HeartbeatSvc fill:#ff9999
    style AggregationSvc fill:#99ff99
```

**Service Initialization Order** (from main.go:179-199):
1. Foundation: MailService, KeyValueService
2. User layer: ApiKeyService, UserService
3. Mapping: LanguageMappingService, ProjectLabelService
4. Core tracking: HeartbeatService
5. Analytics: DurationService, SummaryService
6. Aggregation: AggregationService
7. Features: ReportService, LeaderboardService, ActivityService

---

## 6. Database Schema

```mermaid
erDiagram
    users ||--o{ heartbeats : tracks
    users ||--o{ summaries : owns
    users ||--o{ api_keys : has
    users ||--o{ aliases : configures
    users ||--o{ language_mappings : customizes
    users ||--o{ project_labels : organizes
    users ||--o{ leaderboard_items : appears_in
    users ||--o{ durations : has

    summaries ||--o{ summary_items : contains

    users {
        string id PK
        string api_key UK
        string email
        string password_hash
        datetime created_at
        datetime last_logged_in_at
        int heartbeats_timeout_sec
        string auth_type
        string sub
        bool is_admin
        bool has_data
    }

    heartbeats {
        uint64 id PK
        string user_id FK
        string entity
        string type
        string category
        string project
        string branch
        string language
        bool is_write
        string editor
        string operating_system
        string machine
        datetime time
        string hash UK
        datetime created_at
    }

    summaries {
        uint id PK
        string user_id FK
        datetime from_time
        datetime to_time
        int num_heartbeats
    }

    summary_items {
        uint64 id PK
        uint summary_id FK
        uint8 type
        string key
        duration total
    }

    durations {
        uint64 id PK
        string user_id FK
        datetime time
        duration duration
        string project
        string language
        string editor
        string operating_system
        string machine
        string branch
        string group_hash
    }

    api_keys {
        uint id PK
        string user_id FK
        string key UK
        string name
        datetime created_at
    }

    aliases {
        uint id PK
        string user_id FK
        uint8 type
        string key
        string value
    }

    language_mappings {
        uint id PK
        string user_id FK
        string extension
        string language
    }

    project_labels {
        uint id PK
        string user_id FK
        string project_key
        string label
    }

    leaderboard_items {
        uint id PK
        string user_id FK
        string interval
        uint8 by
        duration total
        string key
    }
```

**Key Tables for OxyWaka:**

1. **heartbeats**: Raw activity data
   - `category`: "coding" vs "browsing" (critical for insights)
   - `is_write`: Active coding indicator
   - `time`: Millisecond precision timestamps
   - Indexed: user_id + time, project, language

2. **durations**: Continuous coding sessions
   - Pre-computed from heartbeats
   - Groups by 10-minute max gap
   - Perfect for "deep work" detection

3. **summaries + summary_items**: Pre-aggregated stats
   - Fast dashboard queries
   - Multiple dimensions (project, language, editor, OS, category)

**Extension Point for Insights:**
Add new tables:
- `insight_scores`: Daily productivity scores
- `activity_patterns`: Deep work sessions, context switches
- `recommendations`: AI-generated suggestions

---

## 7. Background Jobs & Scheduling

```mermaid
graph TB
    subgraph "Scheduled Jobs"
        Cron1[Daily Aggregation<br/>02:15 AM]
        Cron2[Weekly Reports<br/>Friday 6 PM]
        Cron3[Leaderboard Generation<br/>6 AM & 6 PM]
        Cron4[Data Cleanup<br/>Sunday 6 AM]
        Cron5[Database Optimization<br/>1st of month 8 AM]
    end

    subgraph "Job Queues"
        DefaultQ[Default Queue<br/>Artifex]
        ProcessingQ[Processing Queue<br/>Summary Workers]
        Processing2Q[Processing Queue 2<br/>Duration Workers]
    end

    subgraph "Services"
        AggSvc[AggregationService]
        ReportSvc[ReportService]
        LeaderboardSvc[LeaderboardService]
        HousekeepingSvc[HousekeepingService]
        MiscSvc[MiscService]
    end

    Cron1 --> AggSvc
    Cron2 --> ReportSvc
    Cron3 --> LeaderboardSvc
    Cron4 --> HousekeepingSvc
    Cron5 --> MiscSvc

    AggSvc --> ProcessingQ
    AggSvc --> Processing2Q
    ReportSvc --> DefaultQ
    LeaderboardSvc --> DefaultQ
    HousekeepingSvc --> DefaultQ

    ProcessingQ --> |Summary Generation| DB[(Database)]
    Processing2Q --> |Duration Computation| DB

    style Cron1 fill:#ffeb99
    style ProcessingQ fill:#99ff99
    style Processing2Q fill:#99ff99
```

**Background Processing:**
- **Artifex Dispatcher**: Job queue with worker pools
- **Cron Syntax**: Extended cron format (second minute hour day month weekday)
- **Parallel Processing**: Multiple workers for summary/duration generation

**Jobs:**
1. **Aggregation** (02:15): Generate daily summaries for all users
2. **Reports** (Fri 18:00): Email weekly reports to opted-in users
3. **Leaderboard** (06:00 & 18:00): Recompute public rankings
4. **Cleanup** (Sun 06:00): Delete old data based on retention policy
5. **Optimize** (1st @ 08:00): VACUUM (SQLite/Postgres) or OPTIMIZE (MySQL)

---

## 8. Authentication Flow

```mermaid
sequenceDiagram
    participant Client
    participant Router
    participant AuthMW as AuthMiddleware
    participant UserSvc as UserService
    participant ApiKeySvc as ApiKeyService
    participant UserRepo as UserRepository
    participant Handler as API Handler

    Client->>Router: Request with<br/>Authorization: Basic <base64(apiKey)>
    Router->>AuthMW: Process request

    alt API Key Authentication
        AuthMW->>AuthMW: Extract API key from header
        AuthMW->>UserSvc: GetUserByKey(apiKey)
        UserSvc->>ApiKeySvc: GetByApiKey(apiKey)
        ApiKeySvc-->>UserSvc: ApiKey entity
        UserSvc->>UserRepo: GetUserById(userId)
        UserRepo-->>UserSvc: User
        UserSvc-->>AuthMW: User
    else Session Authentication (Web)
        AuthMW->>AuthMW: Check session cookie
        AuthMW->>UserSvc: GetUserById(sessionUserId)
        UserSvc-->>AuthMW: User
    end

    alt User found
        AuthMW->>AuthMW: Store user in request context
        AuthMW->>Handler: Forward request
        Handler->>Handler: Access user from context
        Handler-->>Client: Response
    else User not found
        AuthMW-->>Client: 401 Unauthorized
    end
```

**Authentication Methods:**
1. **API Key** (WakaTime clients): `Authorization: Basic <base64(apiKey)>`
2. **Session Cookie** (Web UI): Gorilla sessions
3. **OpenID Connect** (Optional SSO): OAuth2 with OIDC

**Security Features:**
- Password hashing: Argon2id
- API key storage: Database with unique constraint
- Rate limiting: Chi middleware (per IP)
- CORS: Configurable origins

---

## 9. WakaTime Compatibility Layer

```mermaid
graph TB
    subgraph "Client Tools"
        IDEPlugin[IDE Plugins<br/>wakatime-cli]
    end

    subgraph "Wakapi Routes"
        NativeAPI[Native API<br/>/api/heartbeat]
        CompatAPI[Compat API<br/>/compat/wakatime/v1]
    end

    subgraph "Compatibility Handlers"
        WTHeartbeat[HeartbeatHandler]
        WTSummaries[SummariesHandler]
        WTStats[StatsHandler]
        WTAllTime[AllTimeHandler]
        WTProjects[ProjectsHandler]
        WTLeaders[LeadersHandler]
    end

    subgraph "Native Services"
        HeartbeatSvc[HeartbeatService]
        SummarySvc[SummaryService]
    end

    subgraph "Model Conversion"
        WTModels[WakaTime v1 Models]
        NativeModels[Wakapi Models]
    end

    IDEPlugin --> |Configured to Wakapi URL| CompatAPI
    CompatAPI --> WTHeartbeat
    CompatAPI --> WTSummaries
    CompatAPI --> WTStats
    CompatAPI --> WTAllTime

    WTHeartbeat --> HeartbeatSvc
    WTSummaries --> SummarySvc
    WTStats --> SummarySvc

    HeartbeatSvc --> NativeModels
    WTHeartbeat --> WTModels
    WTModels -.converts to.-> NativeModels

    style CompatAPI fill:#ffeb99
    style WTModels fill:#ff9999
```

**Compatibility Endpoints:**
- `/compat/wakatime/v1/users/current/heartbeats` → HeartbeatHandler
- `/compat/wakatime/v1/users/current/summaries` → SummariesHandler
- `/compat/wakatime/v1/users/current/stats` → StatsHandler
- `/compat/wakatime/v1/users/current/all_time_since_today` → AllTimeHandler

**Model Conversion:**
- WakaTime JSON format → Wakapi models
- Located in: `models/compat/wakatime/v1/`
- Ensures plugins work without modification

---

## 10. Event-Driven Architecture

```mermaid
graph TB
    subgraph "Event Publishers"
        HeartbeatSvc[HeartbeatService]
    end

    subgraph "Event Bus (Hub)"
        EventBus[Hub<br/>In-memory pub/sub]
    end

    subgraph "Event Subscribers"
        CacheUpdater[Cache Updater]
        StatsAggregator[Stats Aggregator]
        Future[Future Subscribers<br/>Insights Module?]
    end

    subgraph "Event Types"
        E1[EventHeartbeatCreate]
    end

    HeartbeatSvc --> |Publish| E1
    E1 --> EventBus

    EventBus --> |Subscribe| CacheUpdater
    EventBus --> |Subscribe| StatsAggregator
    EventBus -.Future.-> Future

    CacheUpdater --> |Increment| Cache[(In-Memory Cache)]
    CacheUpdater --> |Invalidate| ProjectCache[Project Stats Cache]

    style EventBus fill:#ffeb99
    style Future fill:#99ff99,stroke-dasharray: 5 5
```

**Event System:**
- **Hub**: Lightweight pub/sub (github.com/leandro-lugaresi/hub)
- **Events**: `EventHeartbeatCreate` published on every heartbeat
- **Subscribers**: Listen for events and react asynchronously

**Current Uses:**
1. **Cache management**: Increment user heartbeat counts
2. **Stats invalidation**: Clear cached project statistics

**OxyWaka Opportunity:**
```go
// Future: insights module subscribes to heartbeat events
sub := eventBus.Subscribe(0, config.EventHeartbeatCreate)
go func() {
    for msg := range sub.Receiver {
        heartbeat := msg.Fields[config.FieldPayload].(*models.Heartbeat)
        // Real-time insights computation
        insightsService.ProcessHeartbeat(heartbeat)
    }
}()
```

---

## Key Architectural Patterns

### 1. **Repository Pattern**
- Abstracts database operations
- All queries go through repositories
- Enables easy database switching

### 2. **Service Layer Pattern**
- Business logic isolated in services
- Interface-based for testability
- Constructor injection of dependencies

### 3. **Dependency Injection**
- Manual DI (no framework)
- All dependencies injected via constructors
- Main.go orchestrates initialization

### 4. **Event-Driven Communication**
- Loose coupling between components
- Pub/sub pattern via Hub
- Async processing

### 5. **Interface Segregation**
- Each service defines its interface
- Mock implementations for testing
- Clear contracts between layers

---

## Extension Points for OxyWaka

### 1. **Add Insights Service**
```go
// services/insights.go
type InsightsService struct {
    heartbeatSvc IHeartbeatService
    summarySvc   ISummaryService
    durationSvc  IDurationService
}

func (s *InsightsService) ComputeDailyScore(user *User, date time.Time) float64 {
    // Analyze heartbeats for the day
    // Detect deep work sessions
    // Quantify context switches
    // Return productivity score
}
```

### 2. **Subscribe to Heartbeat Events**
```go
// Real-time insights on every heartbeat
sub := eventBus.Subscribe(0, config.EventHeartbeatCreate)
go func() {
    for msg := range sub.Receiver {
        hb := msg.Fields[config.FieldPayload].(*models.Heartbeat)
        insightsService.ProcessRealtimeHeartbeat(hb)
    }
}()
```

### 3. **Add Insights API Routes**
```go
// routes/api/insights.go
func (h *InsightsHandler) GetDailyScore(w http.ResponseWriter, r *http.Request) {
    user := middlewares.GetPrincipal(r)
    score := h.insightsSvc.GetDailyScore(user, time.Now())
    helpers.RespondJSON(w, r, http.StatusOK, score)
}
```

### 4. **Extend Database Schema**
```sql
CREATE TABLE insight_scores (
    id INTEGER PRIMARY KEY,
    user_id TEXT NOT NULL,
    date DATE NOT NULL,
    productivity_score REAL,
    deep_work_minutes INTEGER,
    context_switches INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

---

## Performance Considerations

### 1. **Caching Strategy**
- In-memory cache (patrickmn/go-cache)
- User heartbeat counts cached (24h TTL)
- Project stats cached, invalidated on new heartbeats
- Summary totals cached per user

### 2. **Database Optimizations**
- Indexed queries (user_id + time)
- Batch inserts for heartbeats
- Pre-computed summaries (avoid real-time aggregation)
- Streaming queries for large datasets

### 3. **Background Processing**
- Job queues prevent API blocking
- Worker pools for parallel processing
- Scheduled batch jobs (not real-time)

### 4. **Query Patterns**
- Use durations table for session queries (not raw heartbeats)
- Use summaries table for dashboard stats (not real-time aggregation)
- Stream heartbeats for large date ranges (not in-memory)

---

## Critical Files for Understanding

| File | Purpose | Lines | Key Concepts |
|------|---------|-------|--------------|
| `main.go` | Initialization & DI | 500 | Service wiring, startup |
| `models/heartbeat.go` | Core data model | 177 | Heartbeat structure |
| `services/heartbeat.go` | Heartbeat logic | 400+ | Insert, dedupe, events |
| `services/summary.go` | Aggregation logic | 400+ | Summary generation |
| `services/aggregation.go` | Batch processing | 300+ | Daily jobs |
| `repositories/heartbeat.go` | Data access | 400+ | GORM queries |
| `routes/api/heartbeat.go` | API endpoint | 200+ | HTTP handling |
| `config/config.go` | Configuration | 500+ | Settings management |
| `migrations/migrations.go` | Schema management | 137 | Database evolution |

---

## Questions for OxyWaka?

1. **Real-time vs Batch Insights**: Subscribe to events for real-time, or compute during aggregation?
2. **Storage Strategy**: New tables for insights, or extend existing summaries?
3. **AI Integration**: Embedded Go ML, or external Python service?
4. **UI Extension**: Extend existing dashboard views, or new insights page?

---

**Generated with Claude Code for OxyWaka fork evaluation**
