# Project agent instructions

## Artifact hygiene

- Do not create `MODIFIED_FILE`, `DIFF_FILE`, `VERIFICATION.txt`,
  `ROLLBACK.sh`, or `change_artifacts/` in this repository.
- Report validation results directly in the final response.
- Write temporary validation and rollback artifacts under `/tmp`, never in the
  repository, and remove them when they are no longer needed.
- Do not generate patch bundles, audit reports, or rollback scripts unless the
  user explicitly requests them.
