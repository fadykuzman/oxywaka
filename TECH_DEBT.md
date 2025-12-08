# Technical Debt

**Purpose:** Track design issues, anti-patterns, and areas for future refactoring.

**Status:** Living document - add items as discovered during development/refactoring.

---

## Repository Layer

### 1. InsertOrGet Pattern

**Location:** `repositories/user.go:InsertOrGet()`

**Issue:**

- Non-atomic operation (FindOne + Create in two separate queries)
- Race condition vulnerability
- Inefficient (two database roundtrips)
- Unclear semantics (get-or-create vs insert-or-ignore)

**Current Implementation:**

```go
func InsertOrGet(user *User) (*User, bool, error) {
    if existing := FindOne(user.ID); existing != nil {
        return existing, false, nil  // Found existing
    }
    Create(user)
    return user, true, nil  // Created new
}
```

**Problem:**
Between FindOne and Create, another request could insert the same user → duplicate key error or lost race.

**Better Approach:**

```sql
-- PostgreSQL native support
INSERT INTO users (...) VALUES (...)
ON CONFLICT (id) DO NOTHING
RETURNING *, (xmax = 0) AS created;
```

**When to Fix:**

- During GORM → sqlc migration
- Use database-native UPSERT operations
- Single atomic query
- Proper conflict handling

**Priority:** Medium (works in practice due to low concurrency, but not ideal)

---

## Database Layer

### 2. DeleteTx - GORM-Specific Transaction Handling

**Location:** `repositories/user.go:DeleteTx()`

**Issue:**

- Takes `*gorm.DB` as parameter (GORM-specific)
- Not in repository interface (breaks abstraction)
- Makes testing/mocking harder

**Current Implementation:**

```go
func DeleteTx(user *User, tx *gorm.DB) error {
    return tx.Delete(user).Error
}
```

**Better Approach:**

```go
// Define generic transaction interface
type Transaction interface {
    Commit() error
    Rollback() error
}

// Repository method uses generic interface
func DeleteInTx(user *User, tx Transaction) error
```

**When to Fix:**

- During GORM → sqlc migration
- Create transaction abstraction layer
- Works with any database driver

**Priority:** Low (limited usage, can work around)

---

## Testing

### 3. Zero Repository Tests

**Location:** `repositories/*_test.go` (missing)

**Issue:**

- No tests for database layer (highest risk area)
- Makes refactoring dangerous
- No protection against regressions

**Status:** 🚧 **IN PROGRESS** - Writing interface-based tests

**When Complete:**

- All repository CRUD operations tested
- Tests against interfaces (survive GORM → sqlc migration)
- Real Postgres via testcontainers

**Priority:** 🔴 **CRITICAL** (being addressed now)

---

## Architecture

### 4. GORM Dependency Throughout Codebase

**Location:** Multiple files using `gorm.DB`, GORM models, GORM queries

**Issue:**

- ORM abstraction has performance cost
- Less control over SQL

**Current State:**

- ✅ Keeping GORM for now (pragmatic decision)
- ✅ Writing interface-based tests
- ✅ Tests will survive migration

**Future Migration Path:**

1. Write SQL schema (`schema.sql`)
2. Write SQL queries (`queries.sql`)
3. Generate sqlc code
4. Implement `SqlcUserRepository` (same interface)
5. Swap implementations in tests → verify tests pass
6. Update production code
7. Remove GORM

**When to Fix:**

- After product ships and validates
- When performance becomes bottleneck
- Estimated effort: 2-3 weeks

**Priority:** Low (GORM is fine for current scale)

---

## Data Model

### 5. Multi-Database Support (SQLite/MySQL/Postgres)

**Location:** `config/config.go`, `models/shared.go`, migrations

**Issue:**

- ~~Maintains compatibility with 3 databases~~ **MOSTLY RESOLVED**
- ~~Adds complexity (dialect-specific code)~~ **IN PROGRESS**
- ~~More testing surface area~~
- ~~Portfolio bloat~~

**Status:** 🔧 **70% COMPLETE**

**Done:**
- ✅ Removed SQLite, MySQL, MariaDB drivers
- ✅ Removed CockroachDB support (db_opts.go)
- ✅ Git shows active cleanup in repositories/key_value.go, services/duration.go

**Remaining:**
- 🔧 `models/shared.go:78-88` - SQLite string date parsing
- 🔧 `config/config.go` - Multi-dialect configuration handling
- 🔧 `main.go:140` - Remove GetWakapiDBOpts usage
- 🔧 Migrations - Remove dialect conditionals

**Priority:** 🟡 **MEDIUM** (actively being worked on)

---

## Notes

- This document tracks **what exists**, not what we wish existed
- Add items when discovered, not when fixed
- Mark items as resolved when addressed
- Link to relevant GitHub issues when created

**Last Updated:** 2025-12-07
