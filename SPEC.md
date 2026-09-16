# AegisProof Specification v0.1

**Status:** Draft — implementable  
**specVersion:** `"0.1"`  
**Licence:** Apache-2.0

## 1. Purpose

AegisProof defines a **local-first evidence package**: a small manifest that points at artefacts about an AI system revision, with content digests so tools can check integrity.

It does **not** determine legal compliance, safety, or risk tier.

> Green verify means: the package is well-formed and listed files match their digests (and optional freshness rules). It does **not** mean the system is compliant, certified, or safe.

## 2. Design principles

- Vendor-neutral; no mandatory LLM or network for verify
- Compose SPDX / CycloneDX / SARIF / etc. by **reference**, do not reinvent them
- Regulation-neutral core; profiles are opaque tags in v0.1
- Evidence ≠ compliance

## 3. Package layout

A package is a directory containing:

- `aegisproof.yaml` **or** `aegisproof.json` (exactly one preferred; if both exist, JSON wins)
- Referenced evidence files (typically under `evidence/`)

The **package root** is the directory that contains the manifest.

## 4. Manifest object

### 4.1 Required fields

| Field | Type | Description |
|---|---|---|
| `specVersion` | string | Must be `"0.1"` for this revision |
| `system` | object | `name` (string), `version` (string), optional `id` (string) |
| `producer` | object | `name` (string), optional `version` (string) |
| `createdAt` | string | RFC 3339 timestamp |
| `evidence` | array | ≥1 evidence items |

### 4.2 Optional fields

| Field | Type | Description |
|---|---|---|
| `profiles` | string[] | Opaque profile ids (e.g. `eu-ai-act@draft`) — **no legal meaning** |
| `externalDocuments` | array | `{ type, locator }` pointing at SPDX/CycloneDX/etc. |
| `description` | string | Human note |

**Deferred (not in v0.1):** `claims[]`, control status matrices, risk classification, signatures.

### 4.3 Evidence item

| Field | Required | Description |
|---|---|---|
| `id` | yes | Unique within package |
| `type` | yes | See §5 |
| `locator` | yes | Relative path from package root (POSIX `/` separators) |
| `digest` | yes | `sha256:` + 64 lowercase hex chars |
| `mediaType` | no | IANA media type |
| `producedAt` | no | RFC 3339 |
| `expiresAt` | no | RFC 3339 — freshness hint only |

## 5. Evidence types

Closed set plus extension:

`evaluation` | `policy` | `test-result` | `sbom` | `sarif` | `model-card` | `runtime-receipt` | `other`

Unknown types **fail** schema validation in v0.1 (tight for conformance). Future minors may allow `x-*`.

## 6. Path rules

- Locators MUST be relative
- MUST NOT contain `..` segments
- MUST NOT be absolute
- Resolving a locator MUST stay within the package root (after cleaning)
- Symlinks that escape the package root MUST be rejected

## 7. Verify algorithm

1. Locate manifest; parse YAML or JSON
2. Validate against JSON Schema `schema/aegisproof-0.1.schema.json`
3. Ensure evidence `id` values are unique
4. For each evidence item:
   - Resolve locator under package root
   - Reject path escape / missing file
   - Compute SHA-256 of file bytes; compare to `digest`
5. If any `expiresAt` is before verify time:
   - Default: record **warning**, overall result may be `STALE` if otherwise valid
   - `--strict-stale`: treat as failure (`INVALID`)

### 7.1 Results

| Result | Meaning |
|---|---|
| `VALID` | Schema + digests OK; no stale warnings |
| `STALE` | Schema + digests OK; ≥1 freshness warning |
| `INVALID` | Schema, path, or digest failure |

Forbidden result labels: `COMPLIANT`, `ASSURED`, `CERTIFIED`, `PARTIALLY ASSURED`.

### 7.2 Exit codes (CLI)

| Code | Meaning |
|---|---|
| 0 | `VALID` |
| 1 | `INVALID` |
| 2 | `STALE` (integrity OK) |
| 3 | Usage / IO error |

## 8. Size limits (CLI defaults)

- Manifest ≤ 1 MiB
- Single evidence file ≤ 64 MiB
- Evidence count ≤ 10_000
- Path length ≤ 4096 bytes

## 9. Non-goals (v0.1)

- Risk scoring / regime classification
- Control implementation status
- Remote URI fetch during verify
- Cryptographic signatures
- CycloneDX export (planned later)
- Mandatory cloud or XERJ

## 10. Versioning

Pre-1.0: breaking changes allowed with `specVersion` bump and changelog.  
JSON Schema file name tracks `specVersion`.
