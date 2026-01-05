#!/usr/bin/env python3
"""
Retention strategy for managing growing storage in GitHub Pages repository.

Strategy:
- Keep last RETENTION_DAYS (100) of full daily history
- For files older than RETENTION_DAYS, keep only 1 file per month
  - Prefer the 1st day of the month if available
  - Fallback to any day we have data for that month

Usage:
    python retain_data_lite.py [--execute]

Examples:
    # Dry run (default) - shows what would be removed without deleting
    python retain_data_lite.py

    # Actually execute the cleanup (removes files)
    python retain_data_lite.py --execute
"""

import argparse
import os
import sys
from datetime import datetime, timedelta
from pathlib import Path
from typing import Dict, List, Optional, Tuple


# Retention period: keep last N days of full history
RETENTION_DAYS = 100

# Directory configurations: (directory, prefix, extension)
# Only process directories with large files that need retention management
# Small files (output_charts ~20KB, data_summarized ~small) are commented out
DIRS_TO_PROCESS = [
    ("data_lite", "all_players_", ".zip"),
    # ("data_summarized", "aoe2_rating_percentile_date_", ".csv"),  # Small files, not needed
    # ("output_charts", "chart_", ".html"),  # Small files (~20KB), not needed
]


def extract_date_from_filename(filename: str, prefix: str) -> Optional[str]:
    """
    Extract date string (YYYY-MM-DD) from filename.

    Args:
        filename: The filename to parse
        prefix: The expected prefix before the date

    Returns:
        Date string in YYYY-MM-DD format, or None if not found/invalid
    """
    if not filename.startswith(prefix):
        return None

    # Remove prefix and extension to get date part
    date_part = filename[len(prefix):]
    # Remove extension (find last dot and remove everything after it)
    if '.' in date_part:
        date_part = date_part.rsplit('.', 1)[0]

    # Try to parse as date
    try:
        # Validate date format YYYY-MM-DD
        date_obj = datetime.strptime(date_part, "%Y-%m-%d")
        return date_obj.strftime("%Y-%m-%d")
    except ValueError:
        return None


def get_month_key(date_str: str) -> str:
    """Get month key (YYYY-MM) from date string."""
    return date_str[:7]  # First 7 characters: YYYY-MM


def should_keep_file(
    date_str: str,
    cutoff_date: datetime,
    monthly_kept: Dict[str, Tuple[str, str]],
    filepath: str,
    filename: str
) -> Tuple[bool, Optional[str]]:
    """
    Determine if a file should be kept based on retention strategy.

    Args:
        date_str: Date string in YYYY-MM-DD format
        cutoff_date: Files newer than this are kept (full history)
        monthly_kept: Dict mapping month_key -> (filepath, date_str) of kept files
        filepath: Full path to the file
        filename: Just the filename

    Returns:
        Tuple of (should_keep, reason_for_removal)
        - should_keep: True if file should be kept
        - reason_for_removal: None if keeping, or reason string if removing
    """
    file_date = datetime.strptime(date_str, "%Y-%m-%d")

    # Keep all files within the retention period
    if file_date >= cutoff_date:
        return True, None

    # For older files, apply monthly retention
    month_key = get_month_key(date_str)
    day = date_str.split('-')[2]  # Extract day (DD)

    if month_key not in monthly_kept:
        # First file for this month - keep it
        monthly_kept[month_key] = (filepath, date_str)
        return True, None

    # We already have a file for this month
    existing_filepath, existing_date = monthly_kept[month_key]
    existing_day = existing_date.split('-')[2]

    # Prefer 1st of month
    if day == "01" and existing_day != "01":
        # Replace existing with this one (prefer 1st)
        monthly_kept[month_key] = (filepath, date_str)
        return False, f"Replaced by 1st of month file"
    else:
        # Keep the existing one, remove this
        return False, f"Another file already kept for this month ({existing_date})"


def cleanup_directory(
    directory: str,
    prefix: str,
    extension: str,
    cutoff_date: datetime,
    dry_run: bool = False
) -> Tuple[int, int]:
    """
    Clean up a directory using retention strategy.

    Args:
        directory: Path to directory to clean
        prefix: File prefix to match (e.g., "all_players_")
        extension: File extension (e.g., ".zip")
        cutoff_date: Files newer than this are kept
        dry_run: If True, don't actually delete files

    Returns:
        Tuple of (files_kept_count, files_removed_count)
    """
    if not os.path.isdir(directory):
        print(f"  Directory does not exist: {directory}")
        return 0, 0

    # Find all matching files
    files_by_date: Dict[str, str] = {}  # date_str -> filepath
    dir_path = Path(directory)

    for file_path in dir_path.glob(f"{prefix}*{extension}"):
        if not file_path.is_file():
            continue

        filename = file_path.name
        date_str = extract_date_from_filename(filename, prefix)

        if date_str is None:
            print(f"  Warning: Could not extract date from: {filename}")
            continue

        files_by_date[date_str] = str(file_path)

    if not files_by_date:
        print(f"  No matching files found")
        return 0, 0

    # Apply retention strategy
    monthly_kept: Dict[str, Tuple[str, str]] = {}
    files_to_remove: List[Tuple[str, str, str]] = []  # (filepath, date_str, reason)
    files_kept_list: List[Tuple[str, str]] = []  # (filepath, date_str)

    # Process files sorted by date
    for date_str in sorted(files_by_date.keys()):
        filepath = files_by_date[date_str]
        filename = os.path.basename(filepath)

        should_keep, removal_reason = should_keep_file(
            date_str, cutoff_date, monthly_kept, filepath, filename
        )

        if should_keep:
            files_kept_list.append((filepath, date_str))
        else:
            files_to_remove.append((filepath, date_str, removal_reason or "Monthly retention"))

    # Print files to keep
    kept_count = len(files_kept_list)
    if kept_count > 0:
        print(f"  Files to keep ({kept_count}):")
        # Determine monthly snapshots for all months (including those within retention period)
        # This ensures files like 2026-01-02 (when 2026-01-01 doesn't exist) are labeled correctly
        monthly_snapshot_files: Dict[str, str] = {}  # month_key -> filepath

        # First, identify monthly snapshots from monthly_kept (files older than cutoff)
        for month_key, (filepath, date_str) in monthly_kept.items():
            monthly_snapshot_files[month_key] = filepath

        # Then, for all kept files, determine monthly snapshots (prefer 1st of month)
        for filepath, date_str in files_kept_list:
            month_key = get_month_key(date_str)
            day = date_str.split('-')[2]

            # If no monthly snapshot for this month yet, or this is 1st of month
            if month_key not in monthly_snapshot_files:
                monthly_snapshot_files[month_key] = filepath
            elif day == "01":
                # Prefer 1st of month over existing snapshot
                monthly_snapshot_files[month_key] = filepath

        # Determine reason for each file and print in natural order (chronological)
        for filepath, date_str in sorted(files_kept_list, key=lambda x: x[1]):
            filename = os.path.basename(filepath)
            month_key = get_month_key(date_str)

            # Check monthly snapshot first (prioritized reason)
            # A file is a monthly snapshot if it's the designated snapshot for that month
            if month_key in monthly_snapshot_files and monthly_snapshot_files[month_key] == filepath:
                reason = "monthly snapshot"
            else:
                reason = "within retention period"

            print(f"    {filename} (date {date_str}, reason: {reason})")

    # Remove files
    removed_count = 0
    if files_to_remove:
        print(f"  Files to remove ({len(files_to_remove)}):")
    for filepath, date_str, reason in files_to_remove:
        if dry_run:
            print(f"  [DRY RUN] Would remove: {os.path.basename(filepath)} ({date_str}) - {reason}")
        else:
            try:
                os.remove(filepath)
                print(f"  Removed: {os.path.basename(filepath)} ({date_str}) - {reason}")
                removed_count += 1
            except OSError as e:
                print(f"  Error removing {filepath}: {e}")

    if kept_count > 0 or removed_count > 0:
        print(f"  Summary: Kept {kept_count} file(s), removed {removed_count} file(s)")

    return kept_count, removed_count


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Apply retention strategy to manage repository size",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )
    parser.add_argument(
        "--execute",
        action="store_true",
        help="Actually delete files (default is dry-run mode)"
    )

    args = parser.parse_args()

    # Dry-run is default, only execute if --execute flag is provided
    dry_run = not args.execute

    # Calculate cutoff date
    cutoff_date = datetime.now() - timedelta(days=RETENTION_DAYS)
    cutoff_date_str = cutoff_date.strftime("%Y-%m-%d")

    print(f"Retention strategy: Keep last {RETENTION_DAYS} days full history")
    print(f"Cutoff date: {cutoff_date_str}")
    print(f"Files older than {cutoff_date_str} will be reduced to 1 per month")
    if dry_run:
        print("DRY RUN MODE - No files will be deleted (use --execute to actually delete)")
    else:
        print("EXECUTE MODE - Files will be deleted")
    print()

    # Process default directories (assume current working directory)
    dirs_to_process = [
        (os.path.join(os.getcwd(), dir_name), prefix, ext)
        for dir_name, prefix, ext in DIRS_TO_PROCESS
    ]

    # Process each directory
    total_kept = 0
    total_removed = 0

    for directory, prefix, extension in dirs_to_process:
        print(f"Processing: {directory} (pattern: {prefix}*{extension})")
        kept, removed = cleanup_directory(
            directory, prefix, extension, cutoff_date, dry_run
        )
        total_kept += kept
        total_removed += removed
        print()

    print("=" * 60)
    print(f"Summary: Kept {total_kept} file(s), removed {total_removed} file(s)")
    if dry_run:
        print("(Dry run - no files were actually deleted. Use --execute to apply changes)")

    return 0 if total_removed >= 0 else 1


if __name__ == "__main__":
    sys.exit(main())

