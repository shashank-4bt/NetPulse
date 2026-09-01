# NETPULSE Implementation Plan

**Current stage:** 08 — Global internet health map  
**Date:** 2026-09-01  
**Status:** Stages 01–08 implemented. `/map` renders MapLibre with aggregated cells only. Empty stores stay empty; the basemap is not a health heatmap.

This plan is derived from:

1. The NetPulse product contract (Internet Health Intelligence Platform).
2. A full inspection of this repository (greenfield; no application code).

It does **not** invent product features beyond the contract. Items required to ship safely are labeled **PRODUCTION ENGINEERING REQUIREMENT**.

---

## 1. Audit summary

### 1.1 What exists

| Item | Status |
| --- | --- |
| Git repository | Yes (`main` → `origin` `https://github.com/shashank-4bt/NetPulse.git`) |
| Application code | **None** |
| `package.json` / lockfile | **None** |
| Next.js, TypeScript, Tailwind, shadcn | **None** |
| Go API / workers | **None** |
| PostgreSQL / ClickHouse / Redis usage | **None** |
| Auth, API client, domain models | **None** |
| Tests, lint, CI, Docker, Terraform | **None** |
| Design system, fonts, icons, assets | **None** |
| Product docs in-repo | Engineering contract in chat; this plan + `ARCHITECTURE.md` |

### 1.2 Requested analysis (all absent)

| Area | Result |
| --- | --- |
| Duplicate components | N/A — no components |
| Unused dependencies | N/A — no `package.json` |
| Dead code | N/A |
| Inconsistent naming | N/A |
| Poor folder structure | Empty tree; target structure is in `ARCHITECTURE.md` |
| Oversized files | N/A |
| Duplicated logic | N/A |
| Unsafe code | N/A — no runtime code |
| Broken routes | N/A |
| Missing error / loading states | Entire product missing; required states specified for later stages |
| Accessibility problems | N/A — no UI |

### 1.3 Commands run this stage

There is no Node or Go project. Lint, typecheck, tests, and production build **cannot execute**. See Stage 01 report.

### 1.4 Critical foundation fixed this stage

| Issue | Action |
| --- | --- |
| No ignore rules for secrets, `node_modules`, `.next`, Terraform state | Added `.gitignore` (**PRODUCTION ENGINEERING REQUIREMENT**) |

Not done this stage (would be later features or premature scaffolding):

- Next.js app, UI, shadcn, maps, charts
- Go API, workers, databases
- Auth, CI, Docker
- Fake health data or placeholder dashboards

---

## 2. Product invariants (every later stage)

Diagnostic flow:

```text
MEASURE → ISOLATE → COLLECT EVIDENCE → CORRELATE
→ SCORE → EXPLAIN → RECOMMEND → VERIFY
```

Isolation path:

```text
Device → Wi-Fi → DNS → Connectivity → ISP
→ Routing → CDN → TLS → HTTP → Service
```

Hard rules:

- Never present inference as measured fact.
- Never fabricate measurements, incidents, service status, or user/device data.
- Never claim certainty when evidence is incomplete.
- If evidence is insufficient: show **Insufficient evidence**.
- Prefer: "Observed evidence suggests…" — never "AI thinks…"
- No dead buttons, no fake loading, no fake live data.

---

## 3. Stage map

Stages are sequential. Do not skip foundation. Do not implement a later stage's UI with mocked internet health.

### Stage 01 — Repository audit & architecture *(this stage)*

- Inspect repository.
- Record architecture: app, feature, component, API, data, state, auth, testing.
- Record technical debt / gaps.
- Restore `.gitignore`.
- Do not build UI.

### Stage 02 — Web foundation *(next recommended)*

**PRODUCTION ENGINEERING REQUIREMENT**

- Scaffold `apps/web` (Next.js, TypeScript strict, Tailwind, ESLint, Prettier as needed).
- Add `.env.example` (names only).
- App shell: semantic HTML, skip link, heading hierarchy, mobile nav, visible focus.
- Design tokens (calm infrastructure look — not a SaaS template).
- Central `lib/api` + error types; adapter returns `unavailable` when API is absent.
- Placeholder routes with **honest empty/unavailable** copy, not charts.
- Unit test + lint + typecheck + build scripts that actually run.
- Accessibility baseline (keyboard, reduced motion).
- SEO metadata for real pages only.

**Out of scope:** MapLibre/ECharts until a page has real or explicitly fixture-labeled data. Auth. Workers.

### Stage 03 — Domain model & diagnostic contract

- Shared TypeScript types: evidence class, confidence, layer status, run status.
- Zod schemas at the API boundary (even if API is stubbed).
- Feature modules: `diagnose`, `evidence` — logic outside components.
- UI components that can render fact vs hypothesis vs recommendation vs insufficient evidence **from typed props only**.
- Tests for classification rendering (inference cannot look like a fact).

### Stage 04 — Go API skeleton

- `apps/api`: health, version, create/get diagnostic run (queued / unavailable).
- PostgreSQL schema for runs (no fake results).
- Auth session design (HTTP-only cookies) when account is in scope; otherwise documented anonymous + rate limit.
- OpenAPI spec → generated TS client.
- SSRF policy module (unit-tested) even before probes.

### Stage 05 — Measurement workers (minimal real probes)

- Redis (or equivalent) job queue.
- Workers implement a **small, safe** probe set (e.g. DNS + TLS + HTTP to allowlisted public targets only).
- Redirect revalidation; private IP deny.
- Persist measured facts only; hypotheses only from coded rules with confidence.
- Partial runs → `insufficient_evidence` per layer not run.

### Stage 06 — Isolate, score, explain, recommend, verify

- Correlation rules (deterministic, versioned).
- Score with confidence and policy version.
- Explanation templates bound to evidence ids.
- Recommendation catalog with safety class; no auto-destructive actions.
- Verify = new run compared to previous; no invented improvement.

### Stage 07 — Observatory & comparison

- Regional / ASN aggregates from real stored measurements only.
- ClickHouse when volume justifies it (Postgres acceptable until then).
- Maps/charts with empty and insufficient-evidence states; lazy-loaded.

### Stage 08 — Global internet health map *(shipped in this repo as Stage 08)*

- MapLibre GL JS on `/map`, lazy-loaded.
- Aggregates from stored incidents only: clustering, viewport queries, 250-cell cap.
- Hierarchy World → Country → Region → Network/ASN → Service. Layers do not invent health colors.
- Table alternative, mobile summary/filter, no precise individual locations.

The original plan listed accounts here. Accounts remain deferred.

### Stage 09 — Hardening & operations

- CI (lint, typecheck, test, build).
- Rate limits, WAF later if needed.
- Docker Compose for local Postgres/Redis (and ClickHouse when required).
- Structured logging, correlation ids.
- Load/error budgets for diagnose endpoints.

### Stage 10 — Future Internet Observatory (deferred)

- Kafka/Redpanda, object storage, Neo4j, ECS/Fargate, Terraform, Cloudflare — **only with a written need**.

---

## 4. Target folder structure (from Stage 02)

```text
apps/web/                 # Next.js
apps/api/                 # Go (Stage 04)
apps/worker/              # Go (Stage 05)
packages/domain/          # shared types (Stage 03)
packages/config/
ARCHITECTURE.md
NETPULSE_IMPLEMENTATION_PLAN.md
TECHNICAL_DEBT.md
.gitignore
.env.example              # Stage 02
```

Do not create empty `apps/api` until that stage starts.

---

## 5. Quality gates (every implementation stage after 02)

1. Inspect existing code.  
2. Plan.  
3. Implement.  
4. Tests.  
5. Lint.  
6. Typecheck.  
7. Production build.  
8. Responsive audit (320–2560).  
9. Accessibility audit.  
10. Security review (SSRF, XSS, secrets, IDOR).  
11. Fix critical issues.  
12. Remove duplication/dead code.  
13. Document decisions.

Do not mark a stage complete with known critical issues.

---

## 6. Explicit non-goals (until a later stage names them)

- AI chatbot UX or "AI thinks…" copy.
- Fake outage maps, fake ISP scores, fake uptime.
- Malware accusations.
- Automatic dangerous remediation.
- Unrelated SaaS features (billing, social, crypto).
- Kafka, Neo4j, or multi-cloud because they are fashionable.

---

## 7. Next recommended stage

**Stage 09 — Hardening & operations:** CI, Compose for local stores, and operational budgets. Accounts and developer monitoring remain deferred until those surfaces have real data.
