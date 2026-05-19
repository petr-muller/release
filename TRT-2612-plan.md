# TRT-2612: Improve agent handling of RHCOS-related payload failures

## Context

The payload agent (Claude running in CI, orchestrated by `ci-operator/step-registry/openshift/claude/payload/`) analyzes rejected OpenShift payloads by invoking `ai-helpers` skills (primarily `/ci:analyze-payload`). It currently performs poorly on RHCOS-related failures because:

1. **No RHCOS data is consumed** — `fetch_payloads.py` only pulls `blockingJobs`/`asyncJobs` from the Release Controller API; the `changeLogJson` (which contains RHCOS component versions, image refs, and — eventually — RPM diffs and feature gates) is ignored.
2. **No multi-RHCOS stream awareness** — payloads now ship both `rhel-coreos` (RHCOS 9) and `rhel-coreos-10` (RHCOS 10), but the analysis treats all jobs uniformly.
3. **No RHCOS-specific failure analysis** — install failure and test failure skills don't know how to trace issues back to RHCOS changes (kernel updates, MCO layering, os-extensions, etc.).
4. **RHCOS changes are invisible in the HTML report** — the report shows candidate PRs and job results but never mentions what RHCOS packages changed between payloads.

### Dependency: OCPCRT-458

TRT-2612 depends on Release Controller API enhancements tracked in OCPCRT-458:

| Deliverable | Status | What it provides |
|---|---|---|
| `oc` PR [openshift/oc#2247](https://github.com/openshift/oc/pull/2247) | **Merged** | `featureGates` field in `changeLogJson` — but NOT yet appearing in live API (release-controller not rebuilt with new `oc`) |
| RC PR [openshift/release-controller#750](https://github.com/openshift/release-controller/pull/750) | **Open** | `nodeImageStreams` field with per-stream RPM diffs (`changed`/`added`/`removed` packages) |
| Possibly more | Unknown | OCPCRT-458 AC4 ("oc updated or alternative for multi-RHCOS stream awareness") may need additional work |

**What the live API returns today** (verified on `5.0.0-0.nightly`, 2026-05-20):
- `changeLogJson.components` includes `"Red Hat Enterprise Linux CoreOS 10.2"` with a version string
- `changeLogJson.rebuiltImages` / `updatedImages` include `rhel-coreos`, `rhel-coreos-10`, `rhel-coreos-extensions`, `rhel-coreos-10-extensions`, `machine-os-images` entries with commit hashes
- `changeLogJson.featureGates` — **absent** (oc#2247 merged but RC not rebuilt)
- `nodeImageStreams` — **absent** (PR 750 not merged)
- No RPM diff data available anywhere in the JSON API

---

## Acceptance Criteria — Individual Deliverables

Each AC is a separate piece of work, sequenced by what can be done now vs. what is blocked.

### AC3: Auto-flag failures isolated to rhcos10 jobs
**Status: Can be done now** — no OCPCRT-458 dependency

**Problem:** When payload failures are isolated to jobs with `rhcos10` in the name, the agent doesn't distinguish them from general regressions. RHCOS 10 (RHEL 10 based) is GA in 5.0 but still a distinct OS variant — failures isolated to one RHCOS variant strongly indicate an OS-specific root cause (kernel, systemd, SELinux, package differences) rather than a product-wide regression. This distinction is critical for correct root cause attribution and actionable recommendations.

**Changes (ai-helpers):**

1. **`plugins/ci/skills/analyze-payload/SKILL.md`** — Extend Step 5's "Cross-Platform and Cross-Job Failure Pattern Recognition" section:
   - After collecting subagent results, explicitly check whether any failures are **isolated to rhcos10 jobs** (job name contains `rhcos10`). "Isolated" means: the failure's root cause / error pattern does NOT appear in any non-rhcos10 job.
   - Similarly check the converse: failures isolated to non-rhcos10 (i.e. RHCOS 9) jobs that don't appear in rhcos10 jobs.
   - When a failure is variant-isolated, mark it with `failure_scope: "rhcos10-only"` or `failure_scope: "rhcos9-only"` in the analysis output and add a prominent badge in the HTML report.
   - RHCOS 10 is GA in 5.0 — treat rhcos10 failures with the **same revert severity** as any other failure. The variant isolation is diagnostic context (narrows root cause to OS-specific changes), not a reason to reduce severity.

2. **`plugins/ci/skills/analyze-payload/SKILL.md`** — Extend Step 7.2 (Blocking Jobs Summary Table):
   - Add visual indicator (badge or icon) for rhcos10 jobs in the summary table.
   - When failures are variant-isolated, group them visually or add a note.

3. **`plugins/ci/skills/analyze-payload/SKILL.md`** — Extend Step 7.3 (Failed Job Details):
   - When a failure is variant-isolated, add a callout: e.g. "This failure is isolated to RHCOS 10 jobs and does not appear in RHCOS 9 jobs, indicating an OS-variant-specific root cause."

4. **`plugins/ci/docs/jobs.md`** — Update RHCOS 10 documentation: RHCOS 10 is GA in 5.0 (currently says "4.22" and "TechPreview only").

**Verification:** Run the analyze-payload skill against a recent 5.0 rejected payload that has rhcos10 failures, confirm the report correctly flags variant-isolated failures.

---

### AC4: RHCOS-specific analysis skill identifies RHCOS root causes
**Status: Can be partially done now** — full RPM correlation needs PR 750

**Problem:** When an install or test failure is caused by an RHCOS change (kernel update, package regression, MCO layering issue, Ignition config problem), the subagent doesn't know to look for RHCOS-specific indicators and produces generic or wrong root cause analysis.

**Changes (ai-helpers):**

1. **`plugins/ci/skills/prow-job-analyze-install-failure/SKILL.md`** — Extend with RHCOS-specific diagnostic section:
   - When analyzing install failures, check for RHCOS-specific failure patterns:
     - `machine-config-operator` / MCO errors related to OS image pivot/rebase/layering
     - `rpm-ostree` errors (package conflicts, rebase failures)
     - Kernel boot failures or panics (match against `kernel:` log patterns, dracut errors)
     - Ignition config failures specific to OS provisioning
     - `machine-os-content` image pull failures
     - CoreOS live ISO / PXE boot failures (metal/bare-metal jobs)
   - When RHCOS-specific root cause is identified, flag it as `affected_components: "rhcos"` in the `ANALYSIS_RESULT` block.
   - For rhcos10 jobs specifically: note that RHEL 10 base means different package versions, kernel version, and possibly different systemd/SELinux behavior.

2. **Consider whether to create a new `prow-job-analyze-rhcos-failure` skill** or keep it as an extension:
   - Recommendation: extend the existing install failure skill rather than creating a new one, to avoid routing complexity. The subagent already loads the install failure skill; adding a section is simpler than making it conditionally load a second skill.
   - If the RHCOS section grows large enough to justify separation, it can be extracted later.

3. **(When PR 750 merges)** — Extend the install failure skill to accept and use RHCOS RPM diff data:
   - When the parent (analyze-payload) passes RHCOS package changes to the subagent prompt, the install failure skill should cross-reference the failing component (e.g., kernel panic) with the RPM diff (e.g., kernel package was updated).
   - This is a follow-up enhancement within the same AC; the skill should be designed with a placeholder for this data now.

4. **RHCOS-caused failure recommendations** — When a failure is traced to an RHCOS package change (not a component PR), the report should recommend:
   - Filing a bug with the RHCOS/CoreOS team
   - Requesting a respin with the offending package reverted
   - This is different from PR reverts — the mechanism is an RHCOS team action, not a GitHub revert

**Verification:** Identify a recent payload failure where the root cause was an RHCOS change (e.g., kernel update causing boot failure). Re-run analysis with the extended skill and confirm it correctly identifies the RHCOS root cause and recommends RHCOS team engagement.

---

### AC1: Agent consumes RHCOS build data and changelog from the Release Controller API
**Status: Blocked on OCPCRT-458** — needs PR 750 for RPM diffs; feature gates need RC rebuild

**Problem:** `fetch_payloads.py` discards all changelog data. The analyze-payload skill has no RHCOS context to work with.

**Changes (ai-helpers):**

1. **`plugins/ci/skills/fetch-payloads/fetch_payloads.py`** — Extend to fetch and include changelog data:
   - For each payload, also fetch the detailed release info **with `?from=<previous_tag>`** to get the changelog diff between consecutive payloads. Currently `fetch_release_details()` fetches the release info without `?from=`, so it gets no changelog.
   - From the response, extract and include:
     - `changeLogJson.components` — RHCOS version info (already available)
     - `changeLogJson.rebuiltImages` / `updatedImages` — filtered to RHCOS-related images (`rhel-coreos`, `rhel-coreos-10`, `machine-os-images`, etc.)
     - `changeLogJson.featureGates` — when available (after RC rebuild)
     - `nodeImageStreams` — when available (after PR 750 merges), with full RPM diffs
   - Add this data to each payload in the output JSON under a new `rhcos` key.

2. **`plugins/ci/skills/fetch-payloads/SKILL.md`** — Update documentation to describe the new `rhcos` output field.

3. **`plugins/ci/skills/analyze-payload/SKILL.md`** — Extend Step 2 to consume the new RHCOS data:
   - After fetching payloads, extract RHCOS change data from the `rhcos` field.
   - Pass RHCOS context (which packages changed, kernel version delta, etc.) to subagents in Step 5 so they can correlate failures with RHCOS changes.
   - Use RHCOS data in Step 6.1 scoring: if the failure correlates with an RHCOS package change, it should score highly even though it's not a "PR" in the traditional sense.

**Verification:** Run `fetch_payloads.py` against a 5.0 nightly pair and confirm the output includes RHCOS component versions and image entries. Once PR 750 merges, re-verify that `nodeImageStreams` RPM diffs appear.

---

### AC2: Multi-RHCOS streams detected and handled correctly
**Status: Blocked on OCPCRT-458** — needs PR 750 for `nodeImageStreams`

**Problem:** Payloads contain both `rhel-coreos` (RHCOS 9) and `rhel-coreos-10` (RHCOS 10) with potentially different RPM changes. The agent needs to handle each stream separately and correlate failures with the correct stream.

**Changes (ai-helpers):**

1. **`plugins/ci/skills/analyze-payload/SKILL.md`** — Add multi-stream handling:
   - When `nodeImageStreams` contains multiple entries (e.g., `rhel-coreos` and `rhel-coreos-10`), track each stream's RPM changes separately.
   - When correlating failures with RHCOS changes (Step 6.1), match the job's OS stream (`rhcos10` jobs → `rhel-coreos-10` stream, all others → `rhel-coreos` stream) to the correct RPM diff.
   - The `OSSTREAM` / `OS_IMAGE_STREAM` env vars in the step-registry determine which stream a job uses. Job names containing `rhcos10` use `rhel-10` / `rhel-coreos-10`.

2. **`plugins/ci/skills/fetch-payloads/fetch_payloads.py`** — Ensure the `rhcos` output (from AC1) preserves per-stream data:
   - When `nodeImageStreams` has multiple entries, output all of them keyed by tag name (e.g., `rhel-coreos`, `rhel-coreos-10`).

3. **Consider edge cases:**
   - Payload where only one stream has changes (e.g., kernel update in RHCOS 10 but not RHCOS 9)
   - Payload where both streams changed but to different versions
   - Fallback behavior when `nodeImageStreams` is absent (pre-PR 750)

**Verification:** Find a 5.0 payload where `rhel-coreos` and `rhel-coreos-10` have different RPM contents. Confirm the analysis correctly attributes rhcos10-specific job failures to the rhcos10 stream's changes.

---

### AC5: RHCOS changes surfaced in the analysis HTML report
**Status: Blocked on AC1/AC2** (needs RHCOS data to surface)

**Problem:** The HTML report shows PR candidates and job results but doesn't mention what changed in RHCOS between payloads.

**Changes (ai-helpers):**

1. **`plugins/ci/skills/analyze-payload/SKILL.md`** — Add new HTML report section (between Executive Summary and Recommended Reverts):
   - **"RHCOS Changes"** section showing:
     - RHCOS version(s) in this payload vs. the previous payload
     - Per-stream RPM diff table (package name, old version, new version) when available
     - Kernel version changes highlighted prominently (kernel updates are a common cause of regressions)
     - Feature gate changes (when `featureGates` data is available)
   - When a failed job's root cause is traced to an RHCOS change, cross-link from the job detail to the relevant RHCOS change.

2. **Styling:**
   - Use existing CSS variable palette (`var(--orange)` for RHCOS-related items, consistent with infra badge styling)
   - Collapsible RPM diff table for large diffs (many packages)
   - Visual indicator when RHCOS 10 stream has different changes than RHCOS 9

3. **`plugins/ci/skills/payload-autodl-json/SKILL.md`** — Extend the JSON data schema to include RHCOS change metadata (if this skill exists and is used for database ingestion).

**Verification:** Generate a report for a payload with RHCOS changes and confirm the RHCOS section renders correctly with version info, RPM diffs (when available), and cross-links to failures.

---

## Sequencing Summary

```
AC3 (rhcos10 auto-flag)          ← Can start NOW
  │
AC4 (RHCOS analysis skill)      ← Can start NOW (partial), extends after PR 750
  │
  ├── OCPCRT-458 dependency ──── PR 750 merges + RC rebuilt with new oc
  │
AC1 (consume RHCOS data)         ← Blocked on OCPCRT-458
  │
AC2 (multi-stream handling)      ← Blocked on OCPCRT-458, builds on AC1
  │
AC5 (HTML report)                ← Builds on AC1 + AC2
```

AC3 and AC4 (partial) can proceed immediately. AC1, AC2, AC5 are blocked on OCPCRT-458 deliverables.

## Key Files

| File | Repo | Role |
|---|---|---|
| `plugins/ci/skills/analyze-payload/SKILL.md` | ai-helpers | Main payload analysis skill — most changes go here |
| `plugins/ci/skills/fetch-payloads/fetch_payloads.py` | ai-helpers | Fetches payload data from RC API — needs RHCOS data extraction |
| `plugins/ci/skills/fetch-payloads/SKILL.md` | ai-helpers | Documentation for fetch-payloads |
| `plugins/ci/skills/prow-job-analyze-install-failure/SKILL.md` | ai-helpers | Install failure analysis — needs RHCOS section |
| `plugins/ci/skills/payload-autodl-json/SKILL.md` | ai-helpers | JSON output schema — may need RHCOS fields |
| `plugins/ci/skills/payload-results-yaml/SKILL.md` | ai-helpers | Results YAML schema — may need RHCOS fields |
| `plugins/ci/docs/jobs.md` | ai-helpers | Job pattern documentation — update for 5.0 |
| `ci-operator/step-registry/openshift/claude/payload/agent/openshift-claude-payload-agent-commands.sh` | openshift/release | Agent script — likely no changes needed |

## Resolved Questions

1. **RHCOS 10 GA status in 5.0:** RHCOS 10 is **GA in 5.0**. No reduced revert confidence for rhcos10 failures — treat with same severity. Variant isolation is diagnostic context only.

2. **RHCOS-caused failure action:** When root cause is an RHCOS package change (not a PR), recommend **filing a bug with the RHCOS/CoreOS team** and requesting a respin.

3. **OCPCRT-458 remaining work:** Scope beyond oc#2247 and RC#750 is currently unknown. Feature gates are not appearing in the live API despite oc#2247 being merged — likely the release-controller needs to be rebuilt. Track and revisit when starting AC1.
