#!/usr/bin/env bash
# check-coverage.sh — per-file statement coverage enforcement
#
# Usage: bash scripts/check-coverage.sh <coverage.out> <threshold>
# Example: bash scripts/check-coverage.sh coverage.out 80
#
# Files containing "// coverage-exempt: <reason>" are skipped.

set -euo pipefail

COVERAGE_FILE="${1:-coverage.out}"
THRESHOLD="${2:-80}"

if [ ! -f "$COVERAGE_FILE" ]; then
  echo "Error: coverage file '$COVERAGE_FILE' not found"
  echo "Run 'go test -coverprofile=$COVERAGE_FILE ./...' first"
  exit 1
fi

echo "Checking per-file coverage (threshold: ${THRESHOLD}%)..."

# Determine module path from go.mod to convert coverage paths to local paths
MODULE=$(grep '^module ' go.mod 2>/dev/null | awk '{print $2}')

# Build associative arrays: file -> covered statements, file -> total statements
declare -A file_covered
declare -A file_total

while IFS= read -r line; do
  # Skip the mode header line
  [[ "$line" == mode:* ]] && continue

  # Format: github.com/owner/repo/pkg/file.go:startline.col,endline.col numstmt count
  file=$(echo "$line" | cut -d: -f1)
  numstmt=$(echo "$line" | awk '{print $2}')
  count=$(echo "$line" | awk '{print $3}')

  file_total["$file"]=$(( ${file_total["$file"]:-0} + numstmt ))
  if [ "$count" -gt 0 ]; then
    file_covered["$file"]=$(( ${file_covered["$file"]:-0} + numstmt ))
  fi
done < "$COVERAGE_FILE"

FAILED=0
FAILED_FILES=()
EXEMPT_COUNT=0

for file in "${!file_total[@]}"; do
  total="${file_total[$file]}"
  covered="${file_covered[$file]:-0}"

  [ "$total" -eq 0 ] && continue

  # Convert module path (github.com/owner/repo/pkg/file.go) to local path (pkg/file.go)
  local_file="$file"
  if [ -n "$MODULE" ]; then
    local_file="${file#${MODULE}/}"
  fi

  # Check for coverage-exempt annotation in the source file
  if [ -f "$local_file" ] && grep -q "// coverage-exempt:" "$local_file" 2>/dev/null; then
    EXEMPT_COUNT=$(( EXEMPT_COUNT + 1 ))
    continue
  fi

  # Calculate percentage (one decimal place)
  pct=$(awk "BEGIN {printf \"%.1f\", ($covered / $total) * 100}")

  if awk "BEGIN {exit !($pct < $THRESHOLD)}"; then
    FAILED_FILES+=("  $file: ${pct}% (${covered}/${total} statements)")
    FAILED=1
  fi
done

if [ "$EXEMPT_COUNT" -gt 0 ]; then
  echo "  Skipped $EXEMPT_COUNT exempt file(s)"
fi

if [ "$FAILED" -ne 0 ]; then
  echo ""
  echo "Coverage check FAILED (threshold: ${THRESHOLD}%):"
  printf '%s\n' "${FAILED_FILES[@]}"
  echo ""
  echo "Fix: increase test coverage, or add '// coverage-exempt: <reason>' to the file."
  exit 1
fi

echo "All files meet the ${THRESHOLD}% coverage threshold."
