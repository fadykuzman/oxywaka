# OxyWaka Stripping & Refactoring Plan

**Date:** 2025-12-04
**Goal:** Strip Wakapi down to lean essentials, remove portfolio bloat, simplify codebase

---

## 🎯 Architecture Decisions

### Database Layer
- **REMOVE:** GORM ORM
- **REMOVE:** SQLite, MySQL, MariaDB support
- **KEEP:** Postgres only
- **USE:** `sqlc` (generates type-safe Go from SQL) + `pgx` driver
- **Rationale:** Leaner, more control, type-safe, better performance

### Web Framework
- **REMOVE:** Custom routing middleware stack
- **USE:** Standard `net/http` with `ServeMux` (Go 1.22+ routing)
- **Rationale:** Zero dependencies for routing, sufficient for our needs

### Frontend Build
- **KEEP:** TailwindCSS + build system
- **REMOVE:** Brotli precompression (`.br`, `.gz` files)
- **Rationale:** Keep styling flexibility, simplify asset serving

---

## ✅ Features to KEEP

### Authentication (All Methods)
- ✅ Cookie authentication
- ✅ API key (header + query param)
- ✅ OpenID Connect (SSO)
- ✅ Trusted header auth (reverse proxy)

### Email Features (All)
- ✅ Weekly email reports
- ✅ Password reset emails
- ✅ SMTP configuration

### Public Features
- ✅ Public badges
- ✅ Public profiles
- ❌ Public leaderboards (but keep the code, just not prominent)

### Import/Export
- ✅ Import from WakaTime
- ✅ Export heartbeats to CSV
- ❌ Import from other Wakapi instances (remove)
- ❌ Forward to WakaTime (remove)

### Monitoring & Observability (All)
- ✅ Prometheus metrics
- ✅ Sentry integration
- ✅ pprof profiling

### API
- ✅ WakaTime compatibility layer
- ✅ Custom Wakapi API
- ✅ Swagger docs

### Security Features (All)
- ✅ Rate limiting (signup/login/password reset)
- ✅ CAPTCHA support
- ✅ Invite codes
- ✅ Data retention policies
- ✅ Inactive account cleanup

### Performance (All)
- ✅ Cache warming
- ✅ In-memory caches
- ✅ Background aggregation jobs
- ✅ Periodic summaries

---

## ❌ Features to REMOVE

### Deployment Files
- ❌ Kubernetes references (`wakapi-helm-chart` mentions in README)
- ❌ GitPod config (`.gitpod.yml`)
- ❌ SystemD service file (`etc/wakapi.service`)
- ✅ KEEP: Docker, Docker Compose, API tests (Bruno)

### Database Support
- ❌ SQLite support (code, config, migrations)
- ❌ MySQL support (code, config, migrations)
- ❌ MariaDB support (code, config, migrations)
- ❌ Database charset config (MySQL-specific)
- ❌ Database socket config (MySQL-specific)
- ❌ All GORM code

### Frontend Assets
- ❌ All `.br` (Brotli) files
- ❌ All `.gz` (gzip) files
- ❌ Precompression scripts (`yarn compress`, `yarn watch:compress`)
- ❌ `gzipped.FileServer` logic in `main.go`

### Import/Export Features
- ❌ Forward heartbeats to WakaTime
- ❌ Import from other Wakapi instances

### Configuration Bloat
- ❌ Multi-database dialect configuration
- ❌ SQLite-specific settings
- ❌ MySQL-specific settings
- ❌ Socket listeners (keep IPv4/IPv6 only)

### Documentation Bloat
- ❌ Kubernetes deployment instructions
- ❌ GitPod setup instructions
- ❌ Multiple database setup examples
- ❌ SystemD service instructions

---

## 📋 Refactoring Steps (Priority Order)

### Phase 1: Clean Up Files (Low Risk)
1. ✅ Remove `.gitpod.yml`
2. ✅ Remove `etc/wakapi.service`
3. ✅ Remove all `.br` and `.gz` files from `static/`
4. ✅ Update `scripts/` to remove compression scripts
5. ✅ Clean up `package.json` (remove compression commands)
6. ✅ Update README (remove K8s, GitPod, multi-DB instructions)

### Phase 2: Database Migration (High Risk - Careful!)
1. ✅ Set up `sqlc` configuration
2. ✅ Write raw SQL schema (from current GORM models)
3. ✅ Write SQL queries for all operations
4. ✅ Generate type-safe Go code with `sqlc`
5. ✅ Replace GORM repositories with `sqlc`-generated code
6. ✅ Remove SQLite/MySQL migration files
7. ✅ Keep only Postgres migrations
8. ✅ Update config to remove multi-DB options

### Phase 3: Web Framework Migration (Medium Risk)
1. ✅ Audit current routes and middleware
2. ✅ Rewrite using `net/http` `ServeMux`
3. ✅ Migrate middleware to standard `http.Handler` pattern
4. ✅ Test all endpoints
5. ✅ Remove custom routing code

### Phase 4: Configuration Cleanup (Low Risk)
1. ✅ Remove SQLite/MySQL/MariaDB from `config.default.yml`
2. ✅ Remove socket listener config
3. ✅ Remove database charset/compression config
4. ✅ Simplify to Postgres-only settings

### Phase 5: Feature Removal (Medium Risk)
1. ✅ Remove WakaTime forwarding code
2. ✅ Remove Wakapi instance import code
3. ✅ Remove precompression serving logic from `main.go`

### Phase 6: Testing & Validation
1. ✅ Run API tests (Bruno)
2. ✅ Test Docker deployment
3. ✅ Test all authentication methods
4. ✅ Verify email functionality
5. ✅ Check monitoring endpoints

---

## 🎯 Expected Benefits

### Code Reduction
- **Database layer:** ~30% reduction (GORM → sqlc)
- **Routing:** ~15% reduction (custom → stdlib)
- **Config:** ~20% reduction (multi-DB → Postgres)
- **Assets:** ~10% reduction (no precompressed files)
- **Overall:** Estimate 20-25% total codebase reduction

### Performance Improvements
- Faster database queries (no ORM overhead)
- Simpler HTTP handling (no middleware stack overhead)
- Clearer code paths (fewer abstractions)

### Maintenance Benefits
- Single database to support
- Standard library routing (less to learn)
- Type-safe database code (fewer runtime errors)
- Smaller Docker images (fewer dependencies)

---

## ⚠️ Risks & Mitigations

### Risk 1: Database Migration Bugs
- **Mitigation:** Keep Bruno API tests, run full test suite
- **Mitigation:** Review all SQL queries carefully
- **Mitigation:** Test with real data before deploying

### Risk 2: Breaking Authentication
- **Mitigation:** Test each auth method individually
- **Mitigation:** Keep existing auth logic structure

### Risk 3: Breaking WakaTime Compatibility
- **Mitigation:** Test with actual WakaTime clients
- **Mitigation:** Keep API compatibility layer intact

---

## 📝 Notes

- Keep git history clean: one logical change per commit
- Test after each phase before moving to next
- Document any breaking changes
- Update Docker setup as we remove dependencies
- Keep PLAN.md updated with progress

---

**Status:** Planning complete, ready to execute
**Next:** Start Phase 1 (file cleanup)
