# NETPULSE Technical Debt

**Audit date:** 2026-08-31  
**Repository:** `D:\Cursor\NetPulse` (`https://github.com/shashank-4bt/NetPulse.git`)

This file tracks **actual** gaps and risks. It does not invent code problems in files that do not exist.

Severity:

- **P0** — blocks a trustworthy product or safe next stage  
- **P1** — must be solved before real diagnostics  
- **P2** — needed for production operations  
- **P3** — polish

---

## 1. Stage 01 findings

### 1.1 Greenfield gap (P0)

The working tree had **no application** at Stage 01. Stage 02 added `netpulse-web/` (Next.js design system). Diagnostic product surfaces still do not exist. The public GitHub repo must not be described as a running measurement product.

### 1.2 Missing toolchain (P0) — **partially fixed in Stage 02**

| Tool | Status | Impact |
| --- | --- | --- |
| Node.js / Next.js web app | Present (`netpulse-web`) | Design system only |
| TypeScript strict | Present | Web package |
| ESLint / Vitest | Present | Web package |
| Go module | Absent | No API/workers |
| CI | Absent | No merge protection |

### 1.3 Missing `.gitignore` (P0) — **fixed**

**PRODUCTION ENGINEERING REQUIREMENT.** Without ignore rules, the next `npm install` or `.env` would risk committing secrets and `node_modules`.

**Fix:** added `.gitignore` covering env/secrets, Node/Next, Go artifacts, coverage, Terraform state.

### 1.4 No environment contract (P1) — **web slice fixed in Stage 02**

`netpulse-web/.env.example` documents `NEXT_PUBLIC_SITE_URL` only. API env contract remains for Stage 04.

### 1.5 No SSRF policy implementation (P1)

Architecture requires deny lists for loopback, private networks, link-local, and cloud metadata, plus redirect revalidation.

**When to fix:** Stage 04–05 with unit tests **before** any worker performs network I/O.

### 1.6 No evidence type system (P1)

Without discriminated unions in code, a future UI can accidentally show hypotheses as facts.

**When to fix:** Stage 03, before any diagnose results UI.

### 1.7 No auth session design in code (P1)

Contract forbids `localStorage` for auth tokens. Nothing implements sessions yet.

**When to fix:** when account features start (Stage 08, or earlier if diagnose is authenticated).

### 1.8 No CI/CD or Docker (P2)

Expected later; not debt from abandoned pipelines — they were never created.

### 1.9 No accessibility / SEO / performance baseline (P2)

No UI, so no current a11y bugs. Risk is introducing inaccessible UI in Stage 02. Mitigate with semantic shell and gates in the implementation plan.

### 1.10 Product documentation not in-repo (P2)

The product source of truth was provided in the engineering contract (chat). In-repo product spec is only these Stage 01 docs.

**When to fix:** optional `docs/PRODUCT.md` in a later stage if a durable spec file is needed. Do not duplicate randomly.

---

## 2. Items that are **not** technical debt

| Observation | Why it is not debt |
| --- | --- |
| No Kafka, Neo4j, ECS, Terraform | Intentionally deferred; simplest scalable path |
| No MapLibre/ECharts yet | No data to display; adding them now would invite fake charts |
| No shadcn components | No UI stage yet |
| Duplicate components / dead code / unused deps | No source to contain them |
| Git history contains a removed policy file | Intentional; files were deleted after identity verification |

---

## 3. Risks if later stages skip the plan

| Shortcut | Failure mode |
| --- | --- |
| Mock ISP/outage APIs in the UI | Fabricated incidents — product-contract violation |
| Client-side "diagnosis" without probes | Unsupported claims |
| Charts with random numbers | Decorative metrics |
| `any` / unchecked `fetch` | Inferences leak into fact UI |
| Open URL fetch in Next.js server | SSRF |
| Tokens in `localStorage` | Session theft |

---

## 4. Resolved

| ID | Item | Stage |
| --- | --- | --- |
| TD-001 | Missing `.gitignore` for secrets and build artifacts | 01 |
| TD-002 | No application / web toolchain | 02 |
| TD-003 | Web `.env.example` | 02 |

---

## 5. Open list (carry forward)

| ID | Sev | Item | Target stage |
| --- | --- | --- | --- |
| TD-003b | P1 | API env / config schema | 04 |
| TD-004 | P1 | Full diagnostic evidence model (beyond display taxonomies) | 03 |
| TD-005 | P1 | No API or workers | 04 / 05 |
| TD-006 | P1 | No SSRF implementation | 04 / 05 |
| TD-007 | P1 | No auth (when required) | 08 |
| TD-008 | P2 | No CI | 09 |
| TD-009 | P2 | No local Compose for Postgres/Redis | 09 |
| TD-010 | P2 | Product spec only in chat + Stage 01 docs | later docs |

---

## 8. Stage 03 notes

- Public marketing routes are live. Status, outages, map, and service health use `getPublicHealthSnapshot()` and remain **unavailable**.
- Diagnose form validates hosts (including private/loopback rejection) but does not send probes.
- `/design-system` is noindex.

## 7. Stage 02 decisions

- Web app lives in `netpulse-web/` because npm package names cannot contain capitals (`NetPulse`). Root scripts proxy into that package. This is a **PRODUCTION ENGINEERING REQUIREMENT**, not a product feature.
- shadcn/ui Base Nova + Radix-compatible Base UI primitives.
- Light/dark via `next-themes` class strategy (not `prefers-color-scheme` alone).
- Status/severity/confidence always include icon + text.
- Framer Motion is used only in `LoadingState`, and honors `prefers-reduced-motion`.
- MapLibre and ECharts are not installed yet: `ChartContainer` is a state wrapper with no fabricated series.
- `/` and `/design-system` are foundation/gallery routes, not diagnostic product pages.
- Stage 02 originally noindexed the whole site. Stage 03 allows public pages and noindexes `/design-system`.

## 6. Review rule

When a stage lands:

1. Move fixed items to "Resolved" with the stage id.  
2. Do not close P0/P1 by hiding errors in the UI.  
3. Do not close "no live data" by shipping fake live data.
