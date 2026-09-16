# AegisProof — Testing & Empirical Benchmark Design (Phase 0)

**Status:** Methodology only. No conformance suite or fixtures implemented yet.  
**Sources:** synthesised from Phase 0 research threads, especially [Testing benchmark design](15ff989f-feec-4baf-ac62-efeb06d9dde8).  
**Merge target:** `IMPLEMENTATION_AND_TESTING_STRATEGY.md` §§12–13; whitepaper §§14–15.

---

## 0. Measurement philosophy

AegisProof is an **evidence package format + local verifier**. Tests must prove schema validity, integrity failure detection, hostile-input safety, DX, and CI consumability — and, empirically, that packs assemble faster/more completely than ad-hoc practice.

**Not tested:** legal compliance, risk scores, dashboard UX, LLM quality, stars.

| Class | Use |
|---|---|
| Primary (kill/keep) | Engineer/governance hours; missing evidence count; duplicate formats; false-negative integrity rate |
| Secondary | Install time; TTFV; files required; CI LOC; exit-code correctness |
| Vanity (never as impact) | Stars, “controls passed,” “% assured” |

**Falsification:** green `verify` with incomplete required evidence for a *declared profile* is a product bug.

---

## 1. Conformance / integrity / security (summary)

Full matrices live in the master strategy and `SECURITY_STRATEGY.md`. Gate by milestone:

| Milestone | Gate |
|---|---|
| M1 | Schema valid/invalid fixtures |
| M2 | + integrity single-fault mutants + TTFV pilot |
| M3 | + path/YAML/size security negatives |
| M4 | GHA interop + locked exit codes |
| M5 | Empirical pilot #1 |
| v1.0 | ≥3 pilots **or** documented kill/pivot |

Integrity fault codes (proposed): `EVIDENCE_MISSING`, `DIGEST_MISMATCH`, `DIGEST_INVALID`, `REF_BROKEN`, stale warn/fail via `--strict-stale`.

---

## 2. DX targets (aspirational v0.1)

| Metric | Target |
|---|---|
| Install | &lt; 2 min |
| Time-to-first-verify | &lt; 15 min (tightened from earlier 10–15 band) |
| Files for minimal pack | ≤ 3 |
| Reference GHA | 1 job, 0 secrets, ≤ 20 LOC |

---

## 3. Empirical benchmark — “Northwind Claims Assistant”

### Hypothesis (falsifiable)

For a fixed synthetic AI app and fixed checklist, AegisProof yields lower engineer+governance time **or** fewer missing items than ad-hoc folders — with integrity check available.

### Fixture (design now; files post-approval)

`benchmarks/fixtures/northwind-claims-v1/`:

- Fake RAG claims adjuster assistant; messy “as found” layout.
- Planted: stale eval CSV, draft model card, oversight policy, junit XML, logging config, duplicate PDF+CSV conflict.
- Missing: retention statement, runtime drift note, SBOM pointer, signed release metadata.

### Checklist (15 items — freeze before schema freeze)

System identity · purpose · model/provider · eval results · human oversight · PII/data note · audit/logging · deployment pin · SBOM/lockfile pointer · CI tests · limitations · change/approval · incident pointer · runtime monitoring (or explicit absence) · integrity mechanism.

Score: `Present-clear` | `Present-hard-to-find` | `Insufficient` | `Absent` | `Conflicting`.  
Missing = Insufficient + Absent + Conflicting.

### Protocol (&lt; 2 days, ≈ 6–9 focused hours)

| Block | Cap |
|---|---|
| Setup | 30–45 min |
| Arm A ad-hoc assemble | 90–150 min |
| Arm A review | 45–60 min |
| Reset clean tree | 15 min |
| Arm B AegisProof + verify | 90–150 min |
| Arm B review | 30–45 min |
| Debrief | 30 min |

Always A then B on **fresh** fixture copies. Hard caps; unfinished = missing.

### Primary metrics

M1 engineer hours · M2 review hours · M3 missing items · M4 duplicate formats · M5 time-to-integrity-check · M6 CI wire-up minutes (B only).

### Hard kill

1. `(M1_B+M2_B) ≥ 0.95×(M1_A+M2_A)` **and** `M3_B ≥ M3_A`
2. Integrity false-negative on published single-fault mutants after “stable”
3. Median TTFV &gt; 45 min after one doc iteration
4. Research shows SPDX/CDX/in-toto already cover checklist with ≤ effort → pivot to profile, don’t ego-ship
5. Demand is only “say EU AI Act compliant” → refuse to bend core

### Soft pivot

Time savings only for reviewers → position as auditor exchange format. Completeness↑ but hours↑ → generators before promoting standard.

### Non-kill

Low stars; inconclusive XERJ; no enterprise partners yet.

### Validity threats

Learning effect (acknowledge); facilitator bias (scripted); n=1 = pilot only; faster ≠ safer (disclaimer).

---

## 4. Founder decisions (testing)

1. Arm hard-cap 90 vs 150 min?
2. Stale evidence: warn vs fail by default?
3. n=1 internal pilot enough past M5, or require one external SME?
4. Publish blended hourly cost proxy or omit?
5. Checklist 15 vs 10 items for sharper day runs?

---

## 5. Reporting template

```text
Pilot ID / date / roles
Fixture + CLI/spec version
M1–M6 raw + deltas
Missing item IDs / duplicates
Verify outcome
Threats to validity
Decision: proceed | pivot | kill
```
