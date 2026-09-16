# Evidence model

A package has:

- `system` identity (`name`, `version`, optional `id`)
- `producer`
- `createdAt`
- `evidence[]` items: `id`, `type`, `locator`, `digest` (`sha256:…`), optional freshness

Verify checks schema structure, path confinement, and digests. Optional `expiresAt` may yield `STALE` (exit 2) unless `--strict-stale`.

Control status matrices and risk classification are **out of core** (belong in consumers/profiles).
