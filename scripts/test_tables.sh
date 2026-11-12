#!/usr/bin/env bash
# Run a smoke-test query against every ploy table and report results.
# Exit code is the number of failed tables (0 = all passed).

set -euo pipefail

# Colour codes (disabled if not a terminal)
if [ -t 1 ]; then
  GREEN=$'\033[0;32m'
  RED=$'\033[0;31m'
  YELLOW=$'\033[0;33m'
  RESET=$'\033[0m'
else
  GREEN="" RED="" YELLOW="" RESET=""
fi

PASS=0
FAIL=0
SKIP=0

run_test() {
  local table="$1"
  local query="$2"
  printf "  %-40s" "$table"
  local output
  if output=$(steampipe query "$query" 2>&1); then
    printf "${GREEN}PASS${RESET}\n"
    ((PASS++)) || true
  else
    # A 403 / authorization error is treated as a skip, not a hard failure
    if echo "$output" | grep -qiE "403|unauthorized|forbidden"; then
      printf "${YELLOW}SKIP${RESET} (authorization error)\n"
      ((SKIP++)) || true
    else
      printf "${RED}FAIL${RESET}\n"
      echo "$output" | sed 's/^/    /'
      ((FAIL++)) || true
    fi
  fi
}


echo ""
echo "Ploy Steampipe plugin — table smoke tests"
echo "========================================="
echo ""

run_test "ploy_member"                "select id, email, status from ploy_member limit 1"
run_test "ploy_user"                  "select id, email, role from ploy_user limit 1"
run_test "ploy_resource"              "select id, name, type from ploy_resource limit 1"
run_test "ploy_resource_access"       "select id, resource_id, member_id from ploy_resource_access where member_status = 'inactive' and member_is_managed and active limit 1"
run_test "ploy_access_review_campaign" "select id, name, status from ploy_access_review_campaign limit 1"
run_test "ploy_access_review"         "select id, name, status from ploy_access_review limit 1"
run_test "ploy_access_request"        "select id, status from ploy_access_request limit 1"
run_test "ploy_audit_log"             "select id, type, occurred_at from ploy_audit_log limit 1"
run_test "ploy_integration"           "select id, name, type from ploy_integration limit 1"

echo ""
echo "Results: ${GREEN}${PASS} passed${RESET}  ${YELLOW}${SKIP} skipped${RESET}  ${RED}${FAIL} failed${RESET}"
echo ""

exit "$FAIL"
