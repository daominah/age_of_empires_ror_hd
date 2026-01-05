# Growing Storage Management Strategy

## Problem

Each day, the workflow stores a new `data_lite` zip file (~1 MiB) in the GitHub Pages repository. Over time, this causes the repository to grow:

- After 1 year: ~365 MiB
- After 5 years: ~1.8 GiB
- GitHub repositories have size limits and performance degrades with large repos

## Recommended Solution: Retention Strategy ⭐

### Strategy Overview

**Keep last 100 days full history, before that only keep 1 day per month:**

- **Last 100 days:** Keep all daily files (full granularity for recent data)
- **Older than 100 days:** Keep only 1 file per month
  - Prefer the 1st day of the month if available
  - Fallback to any day we have data for that month if 1st is not available

### Benefits

- **Bounded current commit size:** After initial 100 days, each new commit grows by ~1 file/month instead of ~30 files/month
- **Preserves historical trends:** Monthly snapshots allow long-term trend analysis
- **Maintains recent detail:** Full daily history for last 100 days (3+ months)
- **Automatic fallback:** Handles cases where workflow didn't run on 1st of month
- **No external dependencies:** All data stays in the same repository
- **Simple implementation:** Can be done entirely in the workflow

### Important Limitation: Git History

⚠️ **Note:** This retention strategy only affects the **current commit**. Git history still contains all previous versions of deleted files, so the **total repository size (including history) will continue to grow**.

**However, Git uses smart compression:**
- ✅ **Delta compression** - Git stores file changes as deltas, not full copies
- ✅ **Object packing** - Git packs similar objects efficiently
- ⚠️ **Binary files** - Zip files compress poorly with deltas (they're already compressed)
- ✅ **Text files** - Compress very well with Git's delta compression

**What this means:**
- ✅ **Current working directory size** - bounded and manageable
- ✅ **Size of new commits** - grows slowly (~1 file/month)
- ⚠️ **Total repository size** - grows over time, but less than linear due to compression
  - After 1 year: ~365 MiB raw, but Git compression reduces this significantly
  - Binary files (zip) compress less effectively than text files

**To truly reduce repository size, you need:**
1. **Git LFS** (recommended) - stores large files outside git history
2. **History rewriting** - use `git filter-branch` or BFG Repo-Cleaner (destructive, requires force push)
3. **Separate archive repository** - move old files to a different repo
4. **Accept the growth** - GitHub allows up to 100GB, so growth is manageable for many years even with compression

### Storage Projection

**Current Commit Size (what this strategy controls):**

| Time Period    | Files Kept         | Approximate Size            |
|----------------|--------------------|-----------------------------|
| First 100 days | 100 files          | ~100 MiB                    |
| After 100 days | 100 + (months × 1) | ~100 MiB + (months × 1 MiB) |
| After 1 year   | ~112 files         | ~112 MiB                    |
| After 5 years  | ~160 files         | ~160 MiB                    |
| After 10 years | ~220 files         | ~220 MiB                    |

**Result:** Current commit size grows very slowly after the initial 100 days.

**Total Repository Size (including Git history):**
- Will continue to grow as all previous versions are preserved in Git history
- Git's delta compression reduces actual size significantly
- After 1 year: ~365 MiB raw files, but Git compression may reduce to ~200-300 MiB
- After 5 years: ~1.8 GiB raw, but Git compression may reduce to ~1-1.5 GiB
- Binary files (zip) compress less effectively than text files
- Still manageable for many years, but consider Git LFS for long-term solution

### Implementation Details

The retention strategy is implemented as a Python script (`retain_data_lite.py`) for better testability and readability.

#### Step 1: Python Script Overview

The script `retain_data_lite.py` handles the retention logic:

- Extracts dates from filenames based on patterns
- Keeps all files within the retention period (last N days)
- For older files, keeps only 1 per month (preferring 1st of month)
- Supports dry-run mode for testing
- Provides detailed logging

**Usage:**

```bash
# Dry run (default) - shows what would be removed without deleting
python retain_data_lite.py

# Actually execute the cleanup (removes files)
python retain_data_lite.py --execute
```

**Note:**
- The script uses a constant `RETENTION_DAYS = 100` (modify in script to change)
- By default, the script runs in dry-run mode (safe, no files deleted)
- Only processes `data_lite` directory (large zip files ~1 MiB each)
- Small directories (`data_summarized`, `output_charts`) are commented out as they don't need retention management

#### Step 2: Add Cleanup Step to Workflow

Insert the cleanup step after "Copy generated files to GitHub Pages" and before "Commit and push":

```yaml
    - name: Copy generated files to GitHub Pages
      run: |
        # ... existing copy logic ...

    - name: Clean up old files using retention strategy
      working-directory: gh-pages
      run: |
        python ../data2_daominah/aoego/z_aoe2_rating_percentile/retain_data_lite.py --execute

    - name: Commit and push to GitHub Pages
      run: |
        # ... existing commit logic ...
```

**Note:** The script processes the `data_lite` directory from the current working directory. Running from `gh-pages` directory ensures it processes the correct files. The `--execute` flag is required to actually delete files (default is dry-run mode).

#### Step 3: Handle Edge Cases

**Case 1: No file on 1st of month**

- Solution: The script automatically falls back to any available day in that month
- The first file encountered for each month is kept, then replaced if a file with day=01 is found

**Case 2: Multiple files per month (older than 100 days)**

- Solution: Script keeps only one file per month, preferring day 01

**Case 3: Workflow runs multiple times per day**

- Solution: Only new files are copied (existing logic), cleanup runs on all files

**Case 4: Manual file restoration needed**

- Solution: Files are in git history, can be restored with `git checkout <commit> -- <file>`

#### Step 4: Testing the Script

Test the Python script before deploying to ensure it works correctly:

```bash
cd /d/game/age_of_empires_ror_hd/data2_daominah/aoego/z_aoe2_rating_percentile

# dry run to see what would be removed (default behavior)
python retain_data_lite.py

# test on actual files (be careful)
python retain_data_lite.py --execute
```

### Monitoring

After implementation, monitor:

- Repository size over time
- Number of files in each directory
- Verify that monthly snapshots are being kept correctly
- Check that recent 100 days remain intact

### Adjusting Retention Period

To change the retention period (e.g., to 60 days or 180 days):

- Modify the `RETENTION_DAYS` constant at the top of `retain_data_lite.py`
- The constant is set to `100` by default

Consider your use case:

- **60 days:** More aggressive cleanup, ~62 files after 1 year
- **100 days:** Balanced (recommended), ~112 files after 1 year
- **180 days:** More recent history, ~192 files after 1 year

---

## Alternative Solutions (Comparison Only)

### Option 1: Git LFS (Large File Storage)

**How it works:** Store large files outside git history using GitHub's LFS service.

**Pros:**

- All history preserved
- Repository stays small
- No code changes needed
- Files remain accessible

**Cons:**

- Requires Git LFS setup
- Free tier: 1 GB storage + 1 GB bandwidth/month
- Paid tier: $5/month for 50 GB
- Slightly slower checkout
- Bandwidth limits may be exceeded with many downloads

**Best for:** When you need ALL history and have budget for LFS

---

### Option 2: Separate Archive Repository

**How it works:** Keep recent files in main repo, archive older files to separate repository.

**Pros:**

- Complete separation of active vs historical
- Main repo stays small and fast
- Archive can use different storage strategy

**Cons:**

- More complex workflow
- Need to modify code to handle archive links
- Two repositories to manage
- More complex user experience

**Best for:** When you want complete separation and don't mind complexity

---

### Option 3: Monthly Aggregation

**How it works:** Keep daily files for recent period, aggregate older data into monthly zip files.

**Pros:**

- Reduces file count significantly
- All data preserved (just aggregated)
- Simpler than separate repo

**Cons:**

- Need to modify Go code to handle monthly archives
- More complex data structure
- Still grows over time (but slower)
- Loss of daily granularity for old data

**Best for:** When you want to reduce file count but keep all data

---

### Option 4: GitHub Releases for Archiving

**How it works:** Create monthly releases with old files as assets, remove from main repo.

**Pros:**

- Uses GitHub's native features
- Releases are permanent and browsable
- Good for versioning

**Cons:**

- Need to modify code to handle release links
- Release assets have download limits
- More complex workflow
- Users need to navigate releases for old data

**Best for:** When you want versioned archives with GitHub's UI

---

### Option 5: External Cloud Storage

**How it works:** Store old files in cloud storage (S3, etc.), keep recent in GitHub Pages.

**Pros:**

- Unlimited storage
- Fast CDN delivery
- Cost-effective for large amounts

**Cons:**

- Requires cloud account setup
- Additional costs (though minimal)
- Need to modify code for cloud URLs
- Dependency on external service
- More complex architecture

**Best for:** When you need unlimited storage and have cloud infrastructure

---

## Recommendation

**Use the Retention Strategy** because:

1. ✅ Bounded growth - repository size stabilizes after 100 days
2. ✅ Simple implementation - can be done entirely in workflow
3. ✅ No external dependencies - everything stays in one repo
4. ✅ Preserves trends - monthly snapshots maintain historical context
5. ✅ Maintains recent detail - full daily history for 3+ months
6. ✅ Zero cost - no additional services needed
7. ✅ Easy to adjust - change `KEEP_DAYS` as needed

The retention strategy provides the best balance of repository size management, implementation simplicity, and data preservation for your use case.

