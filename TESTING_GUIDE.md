# OxyWaka Testing Implementation Guide

**Date:** 2025-12-04
**Goal:** Replace Bruno E2E tests with proper Go test suite
**Timeline:** ~5 days
**Status:** Planning phase

---

## 📋 Overview

This guide tracks the implementation of our proper testing strategy, replacing the portfolio bloat (Bruno tests) with fast, reliable, maintainable Go tests.

---

## Phase 1: Repository Tests (~2 days)

### Setup

**Install testcontainers:**
```bash
go get github.com/testcontainers/testcontainers-go
go get github.com/testcontainers/testcontainers-go/modules/postgres
```

**Create test helper** (`repositories/test_helper.go`):
```go
package repositories

import (
    "context"
    "testing"

    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
    ctx := context.Background()

    // Start Postgres container
    postgresContainer, err := postgres.RunContainer(ctx,
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    if err != nil {
        t.Fatal(err)
    }

    // Clean up after test
    t.Cleanup(func() {
        postgresContainer.Terminate(ctx)
    })

    // Get connection string
    connStr, err := postgresContainer.ConnectionString(ctx)
    if err != nil {
        t.Fatal(err)
    }

    // Connect with GORM
    db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
    if err != nil {
        t.Fatal(err)
    }

    // Run migrations
    // (Add your migration logic here)

    return db
}
```

### Repository Tests to Write

#### 1. User Repository Tests (`repositories/user_test.go`)

**Coverage:**
- [ ] Create user
- [ ] Find user by ID
- [ ] Find user by API key
- [ ] Find user by email
- [ ] Update user
- [ ] Delete user
- [ ] Test unique constraints (duplicate API key)
- [ ] Test cascade deletes (user → heartbeats)

#### 2. Heartbeat Repository Tests (`repositories/heartbeat_test.go`)

**Coverage:**
- [ ] Insert single heartbeat
- [ ] Insert bulk heartbeats
- [ ] Find heartbeats by user + time range
- [ ] Find heartbeats by project
- [ ] Test duplicate detection (hash constraint)
- [ ] Test pagination
- [ ] Test ordering (by time)

#### 3. Summary Repository Tests (`repositories/summary_test.go`)

**Coverage:**
- [ ] Create summary
- [ ] Find summary by user + time range
- [ ] Find latest summary for user
- [ ] Update summary
- [ ] Delete old summaries
- [ ] Test summary items (cascade)

#### 4. API Key Repository Tests (`repositories/api_key_test.go`)

**Coverage:**
- [ ] Create API key
- [ ] Find by API key
- [ ] Find all keys for user
- [ ] Delete API key
- [ ] Test read-only flag

#### 5. Alias Repository Tests (`repositories/alias_test.go`)

**Coverage:**
- [ ] Create alias
- [ ] Find alias by type + key
- [ ] Update alias
- [ ] Delete alias
- [ ] Test unique constraints

---

## Phase 2: Service Integration Tests (~2 days)

### Setup

**Use same testcontainer approach** but test services with real DB.

### Service Tests to Write

#### 1. User Service Tests (`services/user_test.go`)

**Coverage:**
- [ ] User signup (success)
- [ ] User signup (duplicate email)
- [ ] User login (success)
- [ ] User login (wrong password)
- [ ] Password reset flow
- [ ] Update user settings
- [ ] Delete user (cascade to heartbeats)

**Mock:** Email service only

#### 2. Heartbeat Service Tests (`services/heartbeat_test.go`)

**Coverage:**
- [ ] Process single heartbeat
- [ ] Process bulk heartbeats
- [ ] Deduplication (hash collision)
- [ ] Validation (invalid data)
- [ ] User association
- [ ] Project/language detection

**Mock:** None (use real DB)

#### 3. Summary Service Tests (`services/summary_test.go`)

**Coverage:**
- [ ] Generate summary for time range
- [ ] Aggregate by project
- [ ] Aggregate by language
- [ ] Aggregate by editor
- [ ] Cache behavior
- [ ] Update existing summary

**Mock:** None (use real DB)

#### 4. Alias Service Tests (`services/alias_test.go`)

**Coverage:**
- [ ] Create alias
- [ ] Apply aliases to heartbeat
- [ ] Update alias
- [ ] Delete alias

**Mock:** None (use real DB)

---

## Phase 3: Critical API Tests (~1 day)

### Setup

**Use Go stdlib `httptest`:**
```go
import (
    "net/http/httptest"
    "testing"
)

func TestHeartbeatEndpoint(t *testing.T) {
    // Setup
    db := SetupTestDB(t)
    router := setupRouter(db) // Your router setup

    // Create test request
    req := httptest.NewRequest("POST", "/api/heartbeat", body)
    req.Header.Set("Authorization", "Bearer "+testAPIKey)

    // Record response
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, 201, w.Code)
}
```

### API Tests to Write

#### 1. Authentication Tests (`routes/api/auth_test.go`)

**Coverage:**
- [ ] Login with username/password
- [ ] Login with API key (header)
- [ ] Login with API key (query param)
- [ ] Login failure (invalid credentials)
- [ ] Session cookie validation
- [ ] OIDC login (if applicable)

#### 2. Heartbeat API Tests (`routes/api/heartbeat_integration_test.go`)

**Coverage:**
- [ ] POST /api/heartbeat (success)
- [ ] POST /api/heartbeat (invalid data)
- [ ] POST /api/heartbeat (unauthorized)
- [ ] Bulk heartbeat ingestion

#### 3. Summary API Tests (`routes/api/summary_integration_test.go`)

**Coverage:**
- [ ] GET /api/summary (today)
- [ ] GET /api/summary (custom range)
- [ ] GET /api/summary (unauthorized)
- [ ] Summary filters (project, language)

#### 4. Badge API Tests (`routes/api/badge_integration_test.go`)

**Coverage:**
- [ ] GET /api/badge (public user)
- [ ] GET /api/badge (private user → 403)
- [ ] Badge customization (colors, labels)

---

## Phase 4: Delete Bruno (~5 minutes)

### Files to Remove

- [ ] `testing/wakapi_api_tests/` (entire directory)
- [ ] `testing/run_api_tests.sh`
- [ ] `testing/run_mail_tests.sh`
- [ ] `testing/compose.yml`
- [ ] `testing/config.sqlite.yml`
- [ ] `testing/config.mysql.yml`
- [ ] `testing/config.postgres.yml` (if only used for Bruno)
- [ ] `testing/schema.sql`
- [ ] `testing/data.sql`
- [ ] `testing/wakapi_testing.db*`

### README Updates

- [ ] Remove Bruno CLI installation instructions
- [ ] Remove "API tests" section referencing Bruno
- [ ] Add new "Testing" section explaining Go tests
- [ ] Update test running instructions

---

## 📊 Progress Tracking

### Phase 1: Repository Tests
- [ ] Setup testcontainers helper
- [ ] User repository tests (0/8)
- [ ] Heartbeat repository tests (0/7)
- [ ] Summary repository tests (0/6)
- [ ] API key repository tests (0/5)
- [ ] Alias repository tests (0/5)

**Total:** 0/31 repository tests

### Phase 2: Service Integration Tests
- [ ] User service tests (0/7)
- [ ] Heartbeat service tests (0/6)
- [ ] Summary service tests (0/6)
- [ ] Alias service tests (0/4)

**Total:** 0/23 service tests

### Phase 3: Critical API Tests
- [ ] Authentication tests (0/6)
- [ ] Heartbeat API tests (0/4)
- [ ] Summary API tests (0/4)
- [ ] Badge API tests (0/3)

**Total:** 0/17 API tests

### Phase 4: Cleanup
- [ ] Delete Bruno test files
- [ ] Update README

---

## 🎯 Success Criteria

**Done when:**
1. All repository CRUD operations have tests
2. All critical services have integration tests
3. All critical API endpoints have tests
4. All tests pass with `gotestsum`
5. Test suite runs in < 2 minutes
6. Bruno completely removed
7. README updated

---

## 📝 Notes

- Use `t.Parallel()` where tests don't share state
- Each test should clean up its own data
- Use descriptive test names: `TestUserRepository_Create_DuplicateEmail`
- Keep test data minimal (don't over-seed)
- Prefer table-driven tests for multiple cases

---

**Last Updated:** 2025-12-04
**Next Session:** Start Phase 1 - Repository test helper setup
