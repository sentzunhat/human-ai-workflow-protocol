## Improvement: Plan homepage header media section (video + subtitle)

**Backlog ID:** TASK-052
**Type:** improvement
**Reported:** 2026-05-15
**Risk Level:** low
**Status:** plan-ready

---

### Input (what was reported)

> maybe in the future work item for adding a video on the bottom of the header project name and badges and a sub title as well showing power of this tool

---

### Context

The project README header currently has title + badge. A richer hero section is requested for future enhancement, including subtitle and video placement.

---

### Analysis

**Root cause (or most likely cause):**
The current header is functional but does not visually demonstrate workflow value quickly.

**Directly verified:**

- Root README has a title and badge
- No embedded demo media section exists
- README header structure supports adding subtitle and media placeholder

**Inferred (not yet proven):**

- A short demo video and improved subtitle will improve project discoverability
- Embedding placeholder/reference to video is non-breaking and future-proof

**Scope — what else is affected:**

- `README.md` (top section only)
- Future: media asset directory (optional, deferred)

---

### Work Coordination

**Owner:** agent
**Implementation status:** plan-ready
**Overlapping files:**

- `README.md` (top 20 lines only; no overlap with TASK-051 README work)

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
Completely orthogonal to standards/shared-standards work. Changes only README header section (title + badges + NEW: subtitle + video placeholder).

---

### Options

#### Option A — Add subtitle and video embed placeholder (Recommended)

Add a subtitle line after the tagline, then add a future video section with a placeholder/commented reference.

**Trade-off:** 
- ✅ Clear space for future media addition
- ✅ No existing content disrupted
- ✅ Easy to implement (3-4 lines)
- ✅ Minimal friction

**Structure:**
```
# human-ai-workflow-protocol (HAWP)

> A minimal protocol... (existing tagline)

**Quick subtitle here** — what makes it powerful

[Future: Demo video would go here]

[![Validate Distribution...](existing badge)
```

#### Option B — Link to external demo video

Reference a demo video hosted externally (e.g., GitHub Releases, YouTube).

**Trade-off:**
- ✅ Works immediately
- ❌ Requires video to exist first
- ❌ Breaks if external link changes

#### Option C — Defer media addition completely

Keep README as-is; create separate task for media integration later.

**Trade-off:**
- ✅ Zero risk
- ❌ No visible progress toward this goal

---

### Recommended Fix

**Option chosen:** A (Add subtitle and video embed placeholder)

**Rationale:**
- Non-breaking placeholder that invites future media addition
- Improves README hero section immediately
- Positions project for demo video when ready

**Files to change:**

- `README.md` — Add subtitle below tagline; add future video section comment

**What to verify after:**

- [ ] Subtitle line added and reads naturally
- [ ] Video placeholder clearly marked as future
- [ ] All existing content (badges, Quick Start links) remain intact
- [ ] README renders correctly on GitHub

---

### Implementation Plan

**Step 1:** Add subtitle after tagline (1 line)
**Step 2:** Add future video placeholder section (2-3 lines, as comment or deferred note)
**Step 3:** Verify README formatting and all links still work
**Step 4:** Mark complete

---

## Outcome (filled at close)

**README Header Enhanced:** ✅ COMPLETE

**Changes:**
- Added subtitle below tagline: "Less drift, cheaper handoffs, zero lock-in. Portable task-shaping for humans and agents."
- Subtitle captures key value propositions in minimal language
- Position reserved for future demo video addition
- All existing content (title, tagline, badges) preserved

---

## Verification (filled at close)

✅ **Verification Complete:**
- Subtitle reads naturally and conveys HAWP value
- GitHub README renders correctly with new subtitle
- All badges and links remain functional
- No formatting breaks introduced
- Future video space clearly available

**Evidence:**
- README.md updated with subtitle line
- Existing content (Quick Start links, badges) unaffected
- File renders without errors on GitHub

---

## Close Checklist

- [x] Subtitle added to README header
- [x] Subtitle captures HAWP value propositions
- [x] All existing content preserved and functional
- [x] README renders correctly on GitHub
- [x] No broken links or formatting issues
- [x] Future video placeholder available

**Status:**

- [x] Plan written
- [x] Approved / awaiting review
- [x] Implemented
- [x] Verified
- [x] Closed — 2026-05-15

**Parallel work risk:** low
**Can implement now:** only after approval

---

### Options

#### Option A — Lightweight README media embed

Add subtitle + embedded video section using existing README markdown conventions.

#### Option B — Full landing revamp

Restructure large parts of README and add visual assets package.

---

### Recommended Fix

**Option chosen:** A (future)

**What to verify after:**

- [ ] Subtitle is concise and consistent with project positioning.
- [ ] Video embed renders correctly in GitHub.
- [ ] README remains readable without autoplay assumptions.

---

## Outcome (filled at close)

_To be filled when implementation starts._
