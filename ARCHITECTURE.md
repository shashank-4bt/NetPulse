# NETPULSE Architecture

**Product:** NetPulse — Internet Health Intelligence Platform  
**Core question:** "Is it me, my network, my ISP, or the service?"  
**Repository state (Stage 01):** greenfield Git repository. No application source, configs, tests, or runtime infrastructure exist yet.

This document defines the target architecture. It is not a description of implemented code.

Classification note: items required for production quality that were not named as product features are marked **PRODUCTION ENGINEERING REQUIREMENT**.

---

## 1. Current repository facts

Inspected on 2026-08-31.

| Area | Finding |
| --- | --- |
| Application source | None |
| `package.json` | Missing |
| Next.js / TypeScript / Tailwind | Missing |
| Components, routes, API layer | Missing |
| Database, auth, state management | Missing |
| Tests, Docker, CI/CD | Missing |
| Assets, fonts, icons, design system | Missing |
| Git | Initialized; `main` tracks `https://github.com/shashank-4bt/NetPulse.git` |
| Working tree | Empty except Stage 01 docs and `.gitignore` |

There is no working product to preserve. Later stages must **create** the foundation rather than refactor it.

---

## 2. Architectural principles

1. **Facts ≠ hypotheses ≠ recommendations.** UI and APIs must use discriminated types so inferences cannot be rendered as measurements.
2. **Evidence before explanation.** No diagnostic copy without attached evidence references.
3. **Insufficient evidence is a first-class result**, not an error and not a guess.
4. **Do not fabricate** measurements, incidents, service status, or user/device data.
5. **UI must not contain business logic.** Components render; feature modules decide; API clients transport; domain models define truth.
6. **Simplest architecture that can scale.** Introduce Kafka, Neo4j, object storage, and multi-region AWS only when a measured need exists.
7. **Measurement safety.** Outbound probes are SSRF-constrained. No unrestricted requests to loopback, private networks, link-local, or cloud metadata.
8. **Privacy by minimization.** Coarse geography, short retention where appropriate, no full browsing history.

---

## 3. System context

```text
Browser (Next.js)
    │  HTTPS
    ▼
Edge / WAF (future; Cloudflare optional)
    │
    ▼
NetPulse API (Go)
    │
    ├── PostgreSQL     operational data, users, jobs, evidence metadata
    ├── ClickHouse     high-volume measurement time series
    ├── Redis          job queue, rate limits, short-lived cache
    └── Object store   (future) raw artifacts, HAR, packet summaries
           │
           ▼
Measurement workers (Go, async)
    Device → Wi-Fi → DNS → Connectivity → ISP
    → Routing → CDN → TLS → HTTP → Service
```

**Not in the first deployable slice:** Kafka/Redpanda, Neo4j, ECS/Fargate, Terraform, Cloudflare. Those remain compatibility targets.

---

## 4. App architecture

Target layout (create in Stage 02+; do not invent a monorepo until two deployable units exist):

```text
NetPulse/
  apps/
    web/                 Next.js App Router (TypeScript, Tailwind, shadcn/ui)
    api/                 Go HTTP API
    worker/              Go async measurement workers
  packages/
    domain/              shared TypeScript types generated from API contracts
    config/              env schema, feature flags (non-secret)
  docs/                  product + engineering docs (these files live at repo root for Stage 01)
  infra/                 later: Terraform, Compose for local deps
```

**Stage 02 may start with `apps/web` only**, plus typed API contracts and adapters that fail closed when the Go API is absent.

### 4.1 Next.js app (`apps/web`)

- App Router, TypeScript **strict**.
- Server Components by default; Client Components only for interaction, maps, charts, telemetry consent.
- Route groups for marketing, authenticated app, and public diagnostics.
- No business logic in `page.tsx` / presentational components.
- Heavy modules (MapLibre, ECharts) via dynamic import.

Proposed route map (product, not yet implemented):

| Route | Purpose |
| --- | --- |
| `/` | Product home: what NetPulse diagnoses; no fake live metrics |
| `/diagnose` | Start a measurement run; `?run=` shows the current report |
| `/reports/[id]` | Stored report model (process-local until PostgreSQL) |
| `/api/reports` | JSON create/get for the shareable diagnostic document |
| `/status` | Service health (only from real status sources) |
| `/observatory` | Regional / ASN comparison (when data exists) |
| `/account` | Session, devices, telemetry controls |
| `/docs` | Public product/docs later |

Every critical route must support: **loading, success, empty, error, unavailable, insufficient-evidence**.

### 4.2 Go API (`apps/api`)

Responsibilities:

- Authenticate and authorize.
- Accept diagnostic job requests.
- Enqueue work; never run unrestricted probes in the HTTP process.
- Return typed diagnostic results with evidence classes.
- Rate-limit by identity and IP.
- Validate and canonicalize all user-supplied targets (hostnames, URLs) before workers run.

### 4.3 Workers (`apps/worker`)

Asynchronous measurement pipeline aligned to the product flow:

```text
MEASURE → ISOLATE → COLLECT EVIDENCE → CORRELATE
→ SCORE → EXPLAIN → RECOMMEND → VERIFY
```

Workers must:

- Bind probes to an allowlisted target after DNS resolution.
- Re-validate after redirects.
- Record **what was attempted**, **what was observed**, and **what was not observed**.
- Time out and mark partial runs as incomplete evidence, not as success.

---

## 5. Feature architecture

Features are vertical slices. Each feature owns UI, feature logic, API calls, and tests. Shared domain types live in `packages/domain`.

| Feature | User outcome | Depends on |
| --- | --- | --- |
| `diagnose` | Run and view a diagnostic | API jobs, workers |
| `evidence` | Inspect measured facts vs inferences | domain evidence model |
| `isolate` | Layer isolation (device → service) | correlation engine |
| `score` | Confidence-bounded health score | scoring policy |
| `explain` | Evidence-backed narrative | explanation templates, not unconstrained LLM output |
| `recommend` | Safe, non-destructive next steps | recommendation catalog |
| `verify` | Re-run after a change | existing run + new run comparison |
| `observatory` | Regional / ASN comparison | aggregated ClickHouse data |
| `status` | Service health | real providers or explicit unavailable |
| `account` | Identity, devices, privacy controls | auth |

**Isolation layers** (ordered, product-defined):

1. Device  
2. Wi-Fi  
3. DNS  
4. Connectivity  
5. ISP  
6. Routing  
7. CDN  
8. TLS  
9. HTTP  
10. Service  

A layer may be `not_measured`, `insufficient_evidence`, `healthy`, `degraded`, or `failed`. Never default a skipped layer to healthy.

---

## 6. Component architecture

```text
apps/web/src/
  app/                      routes only
  components/
    ui/                     shadcn primitives
    layout/                 shell, nav, footer
    evidence/               fact / hypothesis / recommendation display
    charts/                 ECharts wrappers (a11y + empty/error)
    maps/                   MapLibre wrappers (lazy)
  features/<name>/
    components/             feature-specific UI
    hooks/                  UI-facing hooks only
    model.ts                feature view-models (no fetch)
  lib/
    api/                    HTTP client, errors, retries
    auth/                   session helpers (no tokens in localStorage)
    config/                 public env
  styles/
```

Rules:

- Presentational components receive already-classified data.
- No `fetch` inside visual components.
- No duplicated status-color logic; status tokens live in the design system.
- Charts and maps must have text/table fallbacks; color is not the only channel.
- Respect `prefers-reduced-motion`. Framer Motion is optional, not decorative default.

---

## 7. API architecture

### 7.1 Contract style

- Versioned HTTP JSON: `/v1/...`, tenant developer routes under `/v1/dev/...`, and organization routes under `/v1/orgs/...`
- OpenAPI (or equivalent) as source of generated TS types.
- Zod (web) and Go structs (API) validate at the boundary.
- Never trust client-supplied evidence, scores, or service status.

### 7.2 Result envelope

Every diagnostic payload must distinguish classes:

```ts
type EvidenceClass = "measured_fact" | "inferred_hypothesis" | "recommendation";

type Confidence = "none" | "low" | "medium" | "high";

type DiagnosticStatus =
  | "queued"
  | "running"
  | "complete"
  | "partial"
  | "failed"
  | "unavailable"
  | "insufficient_evidence";
```

Partial success is `partial` plus per-layer status. Missing layers are `insufficient_evidence`, not invented values.

### 7.3 Error model

Centralized API errors:

| Code | Meaning | UI |
| --- | --- | --- |
| `validation_error` | Bad input | Inline field errors |
| `unauthorized` | No session | Sign-in |
| `forbidden` | Authenticated but not allowed | 403 page |
| `not_found` | Unknown run/resource | Empty/not found |
| `rate_limited` | Abuse protection | Retry-after |
| `unavailable` | Dependency down | Unavailable state |
| `ssrf_blocked` | Target not allowed | Explicit safety message |
| `internal` | Unexpected | Generic + correlation id |

### 7.4 Client layer (`lib/api`)

Single client:

- Base URL from config (not scattered).
- Timeouts, abort signals, correlation id header.
- Schema parse on responses; discard unknown unsafe fields where needed.
- Maps transport failures to the error model above.

Adapters for missing backends must throw `unavailable` / return `insufficient_evidence` — never mock live internet health.

---

## 8. Data architecture

### 8.1 Stores and roles

| Store | Role | Stage to introduce |
| --- | --- | --- |
| PostgreSQL | Users, sessions, diagnostic runs, evidence metadata, audit | When API exists |
| ClickHouse | Probe time series, regional/ASN aggregates | When volume requires it; Postgres may hold early samples |
| Redis | Job queue, rate limits, ephemeral locks | With workers |
| Object storage | Large artifacts (optional) | When artifacts exist |

Do not stand up ClickHouse, Kafka, or Neo4j for the first local slice if PostgreSQL + Redis can store runs and queue jobs.

### 8.2 Core entities (logical)

- `User` (minimal PII)
- `Session`
- `DeviceFingerprint` (coarse, optional, user-visible)
- `DiagnosticRun`
- `Measurement` (layer, probe type, timestamps, raw observation refs)
- `EvidenceItem` (`measured_fact` \| `inferred_hypothesis`)
- `Recommendation`
- `Score` (value + confidence + policy version)
- `VerificationRun` (links to prior run)

Retention: prefer short raw retention; keep aggregates longer. Exact TTLs are a later privacy policy decision.

### 8.3 Evidence rules

- A **measured fact** cites a probe id, timestamp, and observation.
- An **inferred hypothesis** cites one or more evidence ids and a confidence.
- A **recommendation** cites hypotheses/facts and a safety class (`safe_user_action` vs `do_not_auto_execute`).
- Dangerous remediation is never auto-executed.

---

## 9. State architecture

| Kind | Where | Notes |
| --- | --- | --- |
| Server diagnostic data | React Server Components + API | Source of truth is the API |
| Client interactive UI | React local state | Forms, accordions, map viewport |
| Cross-request cache | Next.js `fetch` cache / explicit tags | No fake "live" polling without a real stream |
| Auth session | HTTP-only secure cookies | **PRODUCTION ENGINEERING REQUIREMENT:** never `localStorage` for tokens |
| Telemetry consent | Cookie or account setting | Explicit opt-in |
| Feature flags | Server config | Not in public repo secrets |

No global client store until at least two features share client-only state that cannot live in the URL or server.

---

## 10. Authentication architecture

Target (when auth is implemented):

- Session cookies: `HttpOnly`, `Secure`, `SameSite=Lax` or `Strict` as appropriate.
- CSRF protection for cookie-authenticated mutations.
- Server-side session lookup; revoke on logout.
- RBAC later: `anonymous` (public diagnose if offered), `user`, `operator`.
- IDOR: run ids must be unguessable and authorized.

**Not in Stage 01–02 unless specified:** OAuth providers, SSO, magic links. Start with the smallest real session model when account features begin.

Anonymous diagnostics, if offered, must still be rate-limited and SSRF-safe.

---

## 11. Testing architecture

| Layer | Tooling (target) | What it proves |
| --- | --- | --- |
| Domain | Vitest (web) / Go `testing` | Evidence classification, scoring, SSRF URL policy |
| API | Go tests + httptest | Authz, validation, job enqueue |
| Worker | Go tests with fake resolvers | No private-IP probes; redirect revalidation |
| UI | Vitest + Testing Library | States: loading/empty/error/unavailable/insufficient-evidence |
| a11y | eslint-plugin-jsx-a11y + axe in CI later | Keyboard, names, contrast |
| e2e | Playwright later | Diagnose → results path against stub API |

Policy: no tests that assert fabricated outage or service-status fixtures as if they were production live data. Fixtures must be labeled as fixtures.

---

## 12. Security architecture (measurement)

SSRF policy for workers and any server-side fetch:

**Deny:** localhost, loopback, RFC1918, link-local, IPv6 unique-local, cloud metadata (`169.254.169.254` and equivalents), internal hostnames.

**Require:** DNS resolve → IP allow-check → connect. Re-check after every redirect. Cap redirects and response size. No user-controlled URLs passed to the Next.js server as open proxies.

Other baseline **PRODUCTION ENGINEERING REQUIREMENTS:**

- No secrets in git; `.env` ignored; `.env.example` documents names only.
- XSS: React default escaping; no `dangerouslySetInnerHTML` for diagnostic text.
- SQL: parameterized queries only.
- Rate limits on diagnose endpoints.
- Webhook ingress (future) authenticated and size-limited.

---

## 13. Design system (target)

Feel: precise, technical, trustworthy, calm, modern, fast, professional — infrastructure/reliability software.

Not: crypto, gaming, chatbot, generic SaaS landing.

Tokens to define in Stage 02 (implementation), not as decorative gradients:

- Neutral surfaces, one accent for focus/action.
- Semantic status: info / success / warning / danger / unknown — each with icon + text, not color alone.
- Typography: readable sans; tabular nums for measurements.
- Spacing scale and 44px minimum tap targets on mobile.

MapLibre and ECharts are **data displays**, not dashboard chrome. If there is no data, show empty or insufficient-evidence — never a fake chart.

---

## 14. What Stage 01 does not claim

This architecture is a specification. Until Stage 02 scaffolds the web app and later stages add API/workers, **no diagnostic capability exists**. The UI must not imply otherwise.
