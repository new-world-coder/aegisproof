#!/usr/bin/env bash
set -euo pipefail
REPO=new-world-coder/aegisproof

create() {
  local title="$1"; shift
  local body="$1"; shift
  if gh issue list --repo "$REPO" --state all --limit 100 --json title --jq '.[].title' | grep -Fxq "$title"; then
    echo "skip: $title"
    return 0
  fi
  local args=(gh issue create --repo "$REPO" --title "$title" --body "$body")
  local lab
  for lab in "$@"; do args+=(--label "$lab"); done
  "${args[@]}"
}

create "Document YAML authoring rules vs JSON Schema source of truth" "$(cat <<'EOF'
## Context
JSON Schema is source of truth; YAML is authoring sugar.

## Acceptance
- [ ] Authoring rules in SPEC or wiki (dates, digests, path separators)
- [ ] Round-trip expectations YAML↔JSON documented
EOF
)" spec documentation

create "Evidence type registry (v0 draft)" "$(cat <<'EOF'
## Context
Closed enum today: evaluation, policy, test-result, sbom, sarif, model-card, runtime-receipt, other.

## Acceptance
- [ ] Short registry table (meaning + example artefact)
- [ ] Note on future x-* extensions
EOF
)" spec good-first-issue

create "Extension / profile mechanism design" "$(cat <<'EOF'
## Context
Profiles are opaque tags in v0.1. Design how EU/NIST/healthcare profiles attach without polluting core.

## Acceptance
- [ ] ADR or SPEC section
- [ ] Non-goal: legal meaning in core verify
EOF
)" spec architecture

create "Versioning and compatibility policy (pre-1.0)" "$(cat <<'EOF'
## Acceptance
- [ ] Pre-1.0 breaking change policy
- [ ] How specVersion bumps relate to CLI versions
EOF
)" spec discussion

create "Model identity: reference ML-BOM vs duplicate fields" "$(cat <<'EOF'
## Context
Avoid reinventing SPDX/CycloneDX model identity. Prefer externalDocuments / evidence pointers.

## Acceptance
- [ ] Recommendation written into SPEC or ADR
EOF
)" spec research

create "Freshness / expiresAt semantics (warn vs fail)" "$(cat <<'EOF'
## Context
CLI supports STALE (exit 2) and --strict-stale. Confirm defaults for CI docs.

## Acceptance
- [ ] Documented default for GitHub Action consumers
- [ ] Conformance case for stale evidence
EOF
)" spec cli discussion

create "ADR: dual YAML sugar vs CycloneDX-native" "$(cat <<'EOF'
## Context
Phase 0 recommended dual YAML + later CDX export (option A).

## Acceptance
- [ ] Short ADR in docs/
- [ ] Kill-switch criteria to collapse into CDX profile project
EOF
)" architecture discussion

create "Field mapping: CycloneDX Attestations ↔ AegisProof" "$(cat <<'EOF'
## Acceptance
- [ ] Spreadsheet or markdown matrix
- [ ] Gaps called out honestly
EOF
)" research integrations

create "Field mapping: Regula evidence v1 ↔ AegisProof" "$(cat <<'EOF'
## Acceptance
- [ ] Mapping notes
- [ ] Compete vs compose recommendation
EOF
)" research

create "Field mapping: RiskLine EvidenceBundle ↔ AegisProof" "$(cat <<'EOF'
## Acceptance
- [ ] Mapping from AssureResponse / EvidenceBundle
- [ ] Keep RiskLine types out of core schema
EOF
)" research riskline

create "Accept Phase-0 threat model for v0.1" "$(cat <<'EOF'
## Context
docs/SECURITY_STRATEGY.md

## Acceptance
- [ ] Maintainer sign-off comment
- [ ] Any v0.1 gaps filed as follow-ups
EOF
)" security

create "CLI report JSON schema for verify output" "$(cat <<'EOF'
## Acceptance
- [ ] JSON Schema for verify --format json
- [ ] Documented in SPEC or schema/
EOF
)" cli spec

create "Reusable GitHub Action for consumers (aegisproof-verify)" "$(cat <<'EOF'
## Context
Repo CI exists; consumers need a composite/action they can copy.

## Acceptance
- [ ] action.yml or documented workflow snippet
- [ ] Fail on INVALID; document STALE behaviour
EOF
)" integrations cli

create "CycloneDX export spike" "$(cat <<'EOF'
## Acceptance
- [ ] Spike notes: can we emit CDX declarations from a package?
- [ ] Go / no-go for Milestone 3+
EOF
)" integrations research

create "Northwind claims-assistant empirical benchmark fixture" "$(cat <<'EOF'
## Context
docs/TESTING_AND_BENCHMARK.md

## Acceptance
- [ ] Fixture tree under benchmarks/
- [ ] 15-item checklist frozen
- [ ] Facilitator script
EOF
)" testing

create "Add invalid conformance fixture (good first issue)" "$(cat <<'EOF'
## Task
Add one new invalid package under conformance/v0.1/invalid/ (e.g. duplicate evidence id, absolute path, empty evidence).

## Acceptance
- [ ] Fixture + note in conformance/README.md
- [ ] aegisproof verify exits 1
EOF
)" testing good-first-issue

create "Improve example model-card evidence" "$(cat <<'EOF'
## Task
Enrich examples/rag-claims/evidence/model-card.md with clearer synthetic fields (still clearly fake).

## Acceptance
- [ ] Digests updated via aegisproof pack
- [ ] verify still VALID
EOF
)" documentation good-first-issue

create "Wiki stubs: What / Why / Spec / Security / FAQ" "$(cat <<'EOF'
## Acceptance
- [ ] GitHub wiki pages published
- [ ] README links remain correct
EOF
)" documentation

create "RiskLine import experiment design" "$(cat <<'EOF'
## Context
docs/RISKLINE_INTEROPERABILITY.md — no core dependency.

## Acceptance
- [ ] Experiment plan issue/PR in integrations/ (or sibling)
- [ ] Anti-coupling checklist
EOF
)" riskline integrations

create "RiskLine emit experiment design" "$(cat <<'EOF'
## Acceptance
- [ ] How assure output becomes an AegisProof package
- [ ] Opaque evidence type vs generic attachments
EOF
)" riskline integrations

create "XERJ experiment 1: index sample packages" "$(cat <<'EOF'
## Context
docs/XERJ_INTEGRATION_HYPOTHESIS.md — optional only.

## Acceptance
- [ ] Notes from indexing examples/
- [ ] Digest fields filterable? document result (incl. null)
EOF
)" xerj research

create "Discussion: reject AI SBOM positioning?" "$(cat <<'EOF'
Phase 0 recommendation: yes, reject as primary banner. Confirm or challenge.

See docs/RESEARCH.md and IMPLEMENTATION_AND_TESTING_STRATEGY.md.
EOF
)" discussion

create "Discussion: should claims[] exist in v0.1?" "$(cat <<'EOF'
Deferred by default. Re-open only with a concrete CI consumer need.
EOF
)" discussion spec

create "Discussion: accept 12–18 month kill/fold criteria?" "$(cat <<'EOF'
If no third-party emitter or recurring consumer, fold into RiskLine adapters / CycloneDX profile work.

Confirm maintainer commitment to honest kill criteria.
EOF
)" discussion

echo "COUNT=$(gh issue list --repo "$REPO" --state open --limit 100 --json number --jq 'length')"
gh issue list --repo "$REPO" --limit 40
