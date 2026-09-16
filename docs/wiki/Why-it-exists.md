# Why it exists

AI assurance evidence is fragmented across repos, CI, dashboards, and PDFs. Inventory standards (SPDX / CycloneDX ML-BOM) and heavy attestation models exist, but developers still lack a **thin, CI-native package** that composes heterogeneous artefacts without claiming legal compliance.

AegisProof aims for a SARIF-like role: one portable file (plus attachments) that tools and humans can verify.

See `docs/RESEARCH.md` in the repository for the honest gap analysis (including why “AI SBOM” is the wrong banner).
