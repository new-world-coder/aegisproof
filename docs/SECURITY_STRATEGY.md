# AegisProof — Phase-0 Security Strategy & Threat Model

**Status:** Planning artefact (Milestone 0). Not an implemented security control set.  
**Audience:** Founder approval, architecture/spec authors, future maintainers.  
**Merge target:** §11 Security strategy in `docs/IMPLEMENTATION_AND_TESTING_STRATEGY.md`; also feeds `SECURITY.md`, whitepaper §Security/Provenance, and conformance/security test design.

---

## 0. Security thesis (one paragraph)

AegisProof’s job is to make **AI assurance evidence** machine-checkable and **integrity-aware**—not to certify legal compliance, prove model safety, or replace auditors. The verifier must be **boring and hard to lie to about file presence and digests**, while remaining **honest about what a green exit code does and does not mean**. Early cryptography beyond content hashes is optional; **narrative honesty and CI-safe parsing** are mandatory from day one.

**Non-claim:** `aegisproof verify` ≠ “this AI system is safe / compliant / approved.”

---

## 1. Scope & trust boundaries

### In scope (v0.1–v0.3)

| Surface | Role |
|--------|------|
| Spec + schemas | Shape of manifests and evidence descriptors |
| Local CLI (`verify`, report formats, exit codes) | Reads workspace files; emits human/JSON reports |
| CI gate (e.g. GitHub Action wrapping CLI) | Same trust as the CLI binary + pinned version |
| Evidence packages on disk / in git | Manifest + referenced artifacts |

### Out of scope (Phase 0 / early milestones)

- Remote evidence registries, multi-tenant SaaS, or “trust me” cloud attestations
- Runtime agent sandboxing / LLM prompt injection defenses (not AegisProof’s product)
- Determining regulatory compliance or residual risk scores
- Full supply-chain attestation of *models* or *training data* (compose with SBOM/SLSA later)
- Hardware-backed keys, HSM, or enterprise PKI as a v0.1 requirement

### Trust model (who is trusted)

```text
Developer / CI runner
        │  (trusts: OS FS, git checkout, CLI binary integrity)
        ▼
  aegisproof CLI
        │  (trusts: schema; does NOT trust evidence content semantics)
        ▼
  Manifest + artifact paths
        │
        ▼
  Consumers (RiskLine, auditors, agents)
        │  (must re-verify digests if they care; do not trust marketing copy)
```

**Adversary assumptions (early):**

- **Local/CI adversary:** Can craft malicious YAML/JSON and paths in a PR; wants to crash CI, escape path roots, or get a green verify on substituted files.
- **Evidence author adversary:** Controls what evidence *claims*; may “compliance-wash.” Verifier cannot stop semantic lies—only structural/integrity checks.
- **Supply-chain adversary:** Compromises CLI deps or published packages; wants backdoored verify.
- **Later:** Remote substitution / MITM of evidence URLs (when remote refs exist).

---

## 2. Assets

| ID | Asset | Why it matters | Integrity / confidentiality |
|----|--------|----------------|------------------------------|
| **A1** | Evidence integrity | Digests/refs must match on-disk bytes; substitution must fail verify | Integrity primary; confidentiality usually low (evidence often public in OSS) |
| **A2** | Developer trust in verify output | Humans treat ✓ as meaningful; wrong output → wrong decisions | Integrity + clarity of *claims* |
| **A3** | CI gate trustworthiness | Exit codes block merge/release; false pass or DoS is high impact | Integrity + availability |
| **A4** | Spec/schema reputation | If “green verify” is sold as compliance, project credibility collapses | Integrity of *messaging* |
| **A5** | CLI binary & dependency tree | Compromised tool → silent false assurance | Integrity (supply chain) |
| **A6** | Disclosure / patch process | Uncoordinated vulns in YAML parsers or path logic | Process asset |

**Secondary:** Example fixtures and conformance packs (must not ship unsafe defaults that teach bad patterns).

---

## 3. Threat catalogue

STRIDE-flavoured, mapped to AegisProof realities. Severity is for **default OSS + CI use**, not classified environments.

### T1 — Malicious evidence (content)

| | |
|--|--|
| **Description** | Attacker crafts evidence that is schema-valid but semantically false (fake eval scores, forged policy text, empty “implemented” controls). |
| **Assets** | A2, A4 |
| **Likelihood / impact** | High / High (misuse + washing) |
| **v0.1 detectability** | **Not fully mitigable by CLI.** Verifier checks structure + digests, not truth. |
| **Mitigation direction** | Honest report language; no “COMPLIANT” result; profiles as *mappings* not verdicts; optional later policy engines live in consumers (e.g. RiskLine). |

### T2 — Path traversal / path escape

| | |
|--|--|
| **Description** | Manifest `location: ../../etc/passwd` or symlink tricks cause CLI to read outside evidence root / workspace. |
| **Assets** | A3 (info leak / unexpected fail), host confidentiality |
| **Likelihood / impact** | Medium / High in CI |
| **Mitigation** | Resolve paths under a declared root; reject `..` escapes; optional symlink policy (fail closed or canonicalize-inside-root); never follow absolute paths unless explicitly opted in. |

### T3 — Untrusted YAML / JSON parsing

| | |
|--|--|
| **Description** | Billion laughs, recursive anchors, exotic tags, prototype pollution (JS), unsafe `load` equivalents, zip bombs in nested structures. |
| **Assets** | A3 (DoS), A5 |
| **Likelihood / impact** | Medium / Medium–High |
| **Mitigation** | Safe loaders only; depth/size limits; disable custom tags; prefer JSON Schema validation after parse; fuzz parsers in security tests. |

### T4 — Oversized inputs / resource exhaustion

| | |
|--|--|
| **Description** | Huge manifests, giant “eval” JSON, pathological nesting → OOM / CI timeout → failed gates or runner abuse. |
| **Assets** | A3 |
| **Likelihood / impact** | Medium / Medium |
| **Mitigation** | Hard caps (file size, total bytes, evidence count, nesting depth); clear errors; configurable limits with safe defaults. |

### T5 — Hash confusion / weak or ambiguous digests

| | |
|--|--|
| **Description** | Algorithm confusion (`sha256:` vs bare hex), truncated hashes, case/encoding tricks, hashing wrong representation (pretty-printed JSON vs canonical bytes), digest of path string instead of content. |
| **Assets** | A1, A2 |
| **Likelihood / impact** | Medium / High |
| **Mitigation** | Single canonical digest form for v0.1 (e.g. `sha256:<hex>`); hash **file bytes as stored**; document non-canonicalisation of JSON; reject unknown algs; never “best effort” match. |

### T6 — Supply-chain attack on CLI dependencies

| | |
|--|--|
| **Description** | Typosquat, compromised transitive dep, malicious postinstall, unsigned release tarball. |
| **Assets** | A5, A2, A3 |
| **Likelihood / impact** | Low–Medium / Critical |
| **Mitigation** | Minimal dependency surface; lockfiles; pinned Action versions; checksums/SLSA for releases **when shipping binaries**; dependabot/renovate; prefer stdlib for crypto hashes. |

### T7 — “Compliance washing” misuse

| | |
|--|--|
| **Description** | Marketing or PRs claim “AegisProof verified = EU AI Act ready.” Green CI badge used as regulatory theatre. |
| **Assets** | A4, A2 |
| **Likelihood / impact** | High / High (reputational + societal) |
| **Mitigation** | Spec + CLI copy: **evidence ≠ compliance**; result vocabulary (`VALID` / `INTEGRITY_OK` / `INCOMPLETE`) never `COMPLIANT`; README/SECURITY FAQ; refuse scoring “compliance %” in core. |

### T8 — Evidence substitution

| | |
|--|--|
| **Description** | Swap artifact after hash recorded; swap manifest without updating digests; CI caches stale artifacts; PR replaces good evidence with lookalike. |
| **Assets** | A1, A3 |
| **Likelihood / impact** | Medium / High |
| **Mitigation** | Digests in manifest must match current files; fail on mismatch/missing; discourage “verify without hashes” mode for CI; later: signed manifests / git-linked commits. |

### T9 — Stale evidence

| | |
|--|--|
| **Description** | Digests match old files that no longer describe the deployed system version; git moves on, evidence does not. |
| **Assets** | A1, A2 |
| **Likelihood / impact** | High / Medium–High |
| **Mitigation** | Optional `producedAt` / `subjectVersion` / git commit ref fields; **warn** on age or version skew in v0.2+; do not invent “freshness = safe.” Consumers decide policy. |

### T10 — Signature theater

| | |
|--|--|
| **Description** | Adding signatures/timestamps that nobody verifies, or signing empty/self-attested claims, creates false cryptographic confidence. |
| **Assets** | A2, A4 |
| **Likelihood / impact** | Medium (later) / High |
| **Mitigation** | Ship hashes first; signing only when verification path is real; document what key identity means; prefer ecosystem standards (Sigstore) over bespoke PKI; never imply signature ⇒ truth of claims. |

### Additional threats (track, lower early priority)

| ID | Threat | Note |
|----|--------|------|
| T11 | Symlink / TOCTOU races | Race between hash and open; mitigate with open-then-hash same fd where practical |
| T12 | Report injection | Control chars / ANSI in system names polluting CI logs; sanitize for human output |
| T13 | Remote URL evidence SSRF | Defer remote fetch; if added later, deny by default |
| T14 | Schema downgrade / version confusion | Reject unknown `specVersion` or unsupported features explicitly |
| T15 | Insider CI privilege | Compromised workflow can skip verify; process/org control, not CLI |

---

## 4. Risk matrix (Phase-0 prioritisation)

| Threat | Early severity | Ship-blocking for v0.1 CLI? |
|--------|----------------|-----------------------------|
| T2 Path traversal | High | **Yes** |
| T3 Unsafe parse | High | **Yes** |
| T4 Oversized inputs | Medium–High | **Yes** (sane defaults) |
| T5 Hash confusion | High | **Yes** if hashes present |
| T8 Substitution | High | Digests required for “integrity verified” claims |
| T1 Malicious content | High | Mitigate via **honesty**, not crypto |
| T7 Compliance washing | High | Docs + result vocabulary |
| T6 Supply chain | Medium→High at release | Process + minimal deps |
| T9 Stale evidence | Medium | Warn later; document now |
| T10 Signature theater | High if premature | **Do not ship fake signing** |

---

## 5. Mitigations by milestone

Aligned with project milestones (hashes early; signing/timestamps later).

### Milestone 0 — Research & architecture (now)

- Publish this threat model; fix vocabulary: evidence ≠ compliance.
- Decide digest algorithm and canonical string form.
- Decide workspace root / path policy.
- Sketch `SECURITY.md` + disclosure contacts (even if placeholder email).
- Security test categories in testing strategy (path, YAML, size, hash mismatch).

### Milestone 1 — Minimal specification

- Schema fields for evidence `id`, `type`, `location`, optional `digest`.
- Explicit non-goals in SPEC: no compliance verdict, no mandatory signing.
- Extension/profile mechanism that cannot smuggle “legal status” into core.

### Milestone 2 — CLI `verify` (security baseline)

| Control | Requirement |
|---------|-------------|
| Path containment | All `location`s resolved under evidence/project root; reject escapes |
| Safe parse | JSON/YAML safe loaders; max depth/size |
| Exit codes | Distinct codes: schema fail, missing evidence, integrity fail, OK-with-warnings |
| Report honesty | Never print “compliant”; integrity section separate from “presence/schema” |
| Hash verify | If digest present → must match; if CI flag `--require-digests` → missing digest fails |
| No network | Default verify is local-only |

### Milestone 3 — Evidence integrity (v0.1 hashes → v0.2 provenance)

| Stage | Capability |
|-------|------------|
| **v0.1** | SHA-256 content digests; fail closed on mismatch; document “bytes as stored” |
| **v0.2** | Optional git commit / tree refs as *subject binding*; `producedAt`; staleness **warnings** |
| **Later** | Signed manifests / Sigstore / in-toto link metadata; optional RFC3161 or Sigstore timestamps |

### Milestone 4 — CI

- Pin CLI version by digest or release checksum.
- Example workflows: fail on integrity errors; treat warnings as policy (org choice).
- Document that skipping the action is trivial—gate value is organisational.

### Milestone 5+ — Ecosystem

- Conformance suite includes **security negative tests** (traversal, bombs, bad digests).
- Consumers (RiskLine) re-verify digests; never trust imported “status” alone.
- Optional XERJ indexing: store digest + URI; retrieval must not mutate source evidence.

---

## 6. Provenance options (compose, don’t invent)

| Option | What it proves | What it does **not** prove | When |
|--------|----------------|----------------------------|------|
| **Content digest (SHA-256)** | These bytes are unchanged since digest was recorded | Who wrote them; that claims are true; freshness | **v0.1 default** |
| **Git commit / tree ref** | Evidence was associated with a specific repo state (if commit still reachable & trusted) | That CI ran; that deploy matches commit | **v0.2 optional** |
| **in-toto / SLSA attestations** | A defined step ran in a supply-chain layout | Semantic correctness of AI evals | Later, compose |
| **Sigstore (keyless) signatures** | Identity (e.g. GitHub OIDC) signed this artifact at time T | Claim truth; regulatory fitness | Later, when verify path exists |
| **Transparency / rekor** | Signature was publicly logged | Same as above | Optional later |
| **W3C PROV / OpenLineage** | Lineage narrative / job graph | Cryptographic integrity by itself | Interop, not core crypto |

**Recommendation:** Treat digests as the **integrity primitive**; treat git refs as **subject binding**; treat Sigstore/in-toto as **authorship/process binding**. Do not collapse these into one “secure” checkbox.

---

## 7. What NOT to over-engineer early

Defer or reject for v0.1–v0.2 unless a concrete consumer demand appears:

1. **Custom PKI / root-of-trust ceremonies** — ops heavy; Sigstore later if needed.
2. **Homegrown signature formats** — signature theater risk (T10).
3. **Blockchain / “immutable ledger” marketing** — unnecessary for local evidence.
4. **Encrypted evidence envelopes** — confidentiality is rarely the first problem; complicates DX.
5. **Remote evidence fetch in core verify** — SSRF and trust sprawl (T13).
6. **Automated “risk scores” or compliance percentages** — washing accelerator (T7).
7. **Runtime attestation of model weights in CLI** — wrong layer; use SBOM/OCI/SLSA elsewhere.
8. **Mandatory LLM-assisted “semantic verification”** — violates no-mandatory-LLM; unverifiable.
9. **Perfect freshness oracles** — clocks lie; warn, don’t pretend.
10. **Enterprise SSO inside the CLI** — consumer/platform concern.

**Ship instead:** path safety, size limits, safe parsers, SHA-256 digests, honest UX, pinned CI examples, security test fixtures.

---

## 8. Honest risk: false sense of security from green verify

### The failure mode

A green `aegisproof verify` can mean only:

> Manifest parsed; referenced files exist (within root); digests match (if present/required); schema constraints satisfied.

It does **not** mean: controls work, evaluations are scientific, policies are followed, the right model is in production, or any law is satisfied.

### Why this is the #1 project risk

If AegisProof becomes a **badge of compliance**, it will attract washing (T7), then backlash when an incident shows “verified” systems failed. That kills adoption harder than a missing feature.

### Required mitigations (product + process)

| Layer | Control |
|-------|---------|
| Spec | Normative language: verification ≠ assurance of safety/compliance |
| CLI | Result names: e.g. `PASS` / `INTEGRITY_OK` / `INCOMPLETE`—never `ASSURED` without qualifiers; illustrative founder mock “PARTIALLY ASSURED” should be challenged/replaced |
| Reports | Separate sections: Schema · Presence · Integrity · (optional) Profile checks |
| Docs | “Threats we do not stop” FAQ; auditor guidance |
| CI badge | Prefer “evidence integrity” wording over “AI assured” |
| Consumers | RiskLine/governance tools own policy decisions |

**Kill criterion (reputation):** If messaging cannot stay disciplined under growth pressure, narrow the CLI to **integrity + schema only** and push “assurance narratives” entirely to profiles/consumers.

---

## 9. Security disclosure process (OSS sketch)

Ship early as `SECURITY.md` (even pre-code):

1. **Contact:** Private channel (security@… or GitHub Security Advisories). Prefer GHAS once repo is public.
2. **Scope:** CLI parsing, path handling, digest verification, published packages/Actions, schema validator bugs that cause incorrect **integrity** results. Out of scope: “my evidence was lying,” social engineering, DoS against third-party CI minutes beyond reasonable limits.
3. **Response targets (aspirational for small maintainer set):** Ack ≤ 3 business days; severity triage; fix/advisory for High/Critical before detailed public PoC.
4. **Coordinated disclosure:** 90-day default; shorter if actively exploited.
5. **Credits:** Thank reporters unless anonymity requested.
6. **No bounty initially** — state honestly; revisit after adoption.
7. **Dependency CVEs:** Document update policy; pin releases; announce breaking security fixes clearly.
8. **Supply-chain:** Signed tags / checksums when binaries exist; document verification steps for `curl | sh` avoidance.

---

## 10. Security requirements for design (normative intent)

These should become MUST/SHOULD in SPEC + CLI once approved:

1. **MUST** confine file reads to an explicit root.
2. **MUST** use safe YAML/JSON parsing with resource limits.
3. **MUST** treat digest mismatch as failure (not warning) when digests are present.
4. **MUST NOT** emit compliance/legal conclusions from core verify.
5. **MUST** document the meaning of exit codes.
6. **SHOULD** support `--require-digests` for CI.
7. **SHOULD NOT** fetch remote evidence by default.
8. **MAY** warn on missing optional freshness/subject bindings.
9. **MUST NOT** claim signature-based trust without a documented verification path.

---

## 11. Mapping to tests (hand-off to testing strategy)

| Threat | Test class |
|--------|------------|
| T2 | Path traversal / symlink fixtures |
| T3 | Malicious YAML/JSON corpus |
| T4 | Oversized / deep nesting fixtures |
| T5 | Algorithm confusion, truncated digests, whitespace |
| T8 | Modified bytes after digest; swapped files |
| T9 | Version skew fixtures (warnings) |
| T6 | Dependency audit in release checklist (process test) |
| T7 | Snapshot tests on report strings (no forbidden vocabulary) |

---

## 12. Open security questions (founder / architecture)

1. **Digest required or optional in v0.1?** Recommendation: optional in schema, **required in CI profile** via flag.
2. **Symlinks:** reject vs allow-if-target-in-root?
3. **Canonical JSON hashing:** avoid in v0.1 (hash raw files only)—confirm.
4. **Result vocabulary:** replace founder-illustrative “ASSURED” language?
5. **Security contact identity** before public launch.
6. **When (if ever) Sigstore:** only after a consumer needs signed evidence exchange.

---

## 13. Merge guide (for master strategy doc)

Suggested §11 Security strategy outline:

1. Security thesis & non-claims  
2. Assets (A1–A6)  
3. Threat catalogue summary table (T1–T10)  
4. Milestone-staged mitigations  
5. Provenance ladder (digest → git → Sigstore/in-toto)  
6. Explicit non-goals / anti-over-engineering  
7. False-sense-of-security risk + messaging controls  
8. Disclosure process pointer → `SECURITY.md`  
9. Open questions needing approval  

Full detail remains in this document.

---

## 14. Phase-0 verdict

AegisProof can be **security-credible** without early cryptography theatre: treat the CLI as a **hostile-input parser + digest checker**, treat green verify as **integrity/schema success**, and defer signing until verification of signatures is a first-class, documented path. The dominant residual risk is **social/semantic** (washing and false assurance)—mitigated by product language and scope discipline more than by more crypto.
