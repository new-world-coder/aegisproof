# AegisProof — Open Questions (Phase 0)

Questions that block or reshape the build. Founder decisions marked ★.

---

## Strategy ★

1. **New schema vs CycloneDX profile?** Dual YAML sugar + CDX export, CDX-native, or independent-only?
2. **Is AegisProof justified** given CycloneDX Attestations + Regula + RiskLine evidence bundles — or should effort go into RiskLine emitters for CDX/Regula?
3. **Positioning:** reject “AI SBOM”? Prefer “machine-readable evidence for AI systems”?
4. **Relationship to RiskLine:** extract format from RiskLine, or define independently with RiskLine as first consumer?
5. **Kill criteria acceptance:** fold/abandon if no third-party emitter in 12–18 months?

## Spec

6. Include `claims[]` in v0.1 or defer?
7. Evidence type enumeration — closed enum + `x-*` extensions, or open strings?
8. Freshness (`expiresAt`): warn vs fail in CI?
9. Remote URI locators in v0.1 or local-only?
10. How to represent model identity without reinventing ML-BOM fields?
11. Backwards compatibility policy before 1.0 (break freely with semver minors?).

## Security

12. YAML subset — require JSON-only for verify to reduce parser risk?
13. Max file size / max evidence count defaults?
14. Signing in v0.2: raw Ed25519 (like RiskLine) vs Sigstore?
15. How to prevent “green verify = compliant” misuse in UX copy?

## Product / adoption

16. Who is the **first wedge user** outside the founder?
17. What is the **default consumer** that rejects releases without the artefact?
18. Licence: Apache-2.0 default OK, or evaluate others?
19. Org hosting: personal GitHub vs new `aegisproof` org day one?
20. Name collision / trademark check for “AegisProof”?

## Interop

21. Regula: compete, ignore, or import their packs as evidence items?
22. OSCAL bridge — ever, or never in scope?
23. XERJ experiments before or after RiskLine import?

## Benchmarks

24. Accept synthetic fixture study as sufficient Phase-1 evidence of value?
25. Minimum time-savings threshold for continuing (e.g. ≥30% audit prep)?
26. Arm A/B hard-cap: 90 vs 150 minutes?
27. n=1 internal Northwind pilot enough past M5, or require one external SME?
28. Checklist length: 15 vs 10 items?
29. Publish blended hourly cost proxy or omit?

## Cloud

30. Any hosted component before 10 external users? (Recommendation: **no**.)
31. Docs hosting: GitHub Pages vs Vercel vs other?

## Governance

32. RFC process starting at which version?
33. When (if ever) consider foundation transfer?
34. Contributor CLA / DCO preference?
