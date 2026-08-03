#!/usr/bin/env bash
#
# Prunes older .deb packages from an apt repo pool directory, keeping only
# the N most recent versions (by filename version, sorted with `sort -V`).
#
# Usage: prune-apt-pool.sh <pool-dir> <keep-count>
set -euo pipefail
IFS=$'\n\t'

main() {
  POOL_DIR="${1:?usage: prune-apt-pool.sh <pool-dir> <keep-count>}"
  KEEP="${2:?usage: prune-apt-pool.sh <pool-dir> <keep-count>}"

  [ -d "${POOL_DIR}" ] || { echo "Pool dir ${POOL_DIR} does not exist, nothing to prune"; return 0; }

  VERSIONS="$(
    find "${POOL_DIR}" -name '*.deb' -printf '%f\n' \
      | sed -n 's/^tl_\([^_]*\)_.*\.deb$/\1/p' \
      | sort -uV
  )"

  [ -n "${VERSIONS}" ] || { echo "No .deb packages found in ${POOL_DIR}"; return 0; }

  KEEP_VERSIONS="$(echo "${VERSIONS}" | tail -n "${KEEP}")"
  PRUNE_VERSIONS="$(comm -23 <(echo "${VERSIONS}") <(echo "${KEEP_VERSIONS}"))"

  if [ -z "${PRUNE_VERSIONS}" ]; then
    echo "Nothing to prune, ${POOL_DIR} has $(echo "${VERSIONS}" | wc -l) version(s), keeping up to ${KEEP}"
    return 0
  fi

  echo "Pruning versions:"
  echo "${PRUNE_VERSIONS}" | sed 's/^/  - /'

  while IFS= read -r VERSION; do
    find "${POOL_DIR}" -name "tl_${VERSION}_*.deb" -print -delete
  done <<< "${PRUNE_VERSIONS}"
}

main "$@"
