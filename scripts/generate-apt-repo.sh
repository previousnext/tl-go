#!/usr/bin/env bash
#
# Generates a signed APT repository structure from the .deb files in ./pool.
#
# Expects to be run from within the apt repo root directory (i.e. the
# directory containing ./pool), and requires a GPG signing key to already be
# imported into the runner's keyring.
set -euo pipefail
IFS=$'\n\t'

main() {
  ARCHS=(amd64 arm64)

  SUITE_DIR="dists/${SUITE:-stable}"
  COMPONENT_DIR="${SUITE_DIR}/${COMPONENTS:-main}"

  echo "Generating Packages files"
  for ARCH in "${ARCHS[@]}"; do
    PACKAGE_DIR="${COMPONENT_DIR}/binary-${ARCH}"
    mkdir -p "${PACKAGE_DIR}"
    dpkg-scanpackages --multiversion --arch "${ARCH}" pool/ > "${PACKAGE_DIR}/Packages"
    gzip -fk "${PACKAGE_DIR}/Packages"
    bzip2 -fk "${PACKAGE_DIR}/Packages"
  done

  pushd "${SUITE_DIR}" >/dev/null
  echo "Making Release file"
  {
    echo "Origin: ${ORIGIN:-tl}"
    echo "Label: ${REPO_OWNER:-previousnext}"
    echo "Suite: ${SUITE:-stable}"
    echo "Codename: ${SUITE:-stable}"
    echo "Version: 1.0"
    echo "Architectures: ${ARCHS[*]}"
    echo "Components: ${COMPONENTS:-main}"
    echo "Description: ${DESCRIPTION:-A repository for packages released by ${REPO_OWNER:-previousnext}}"
    echo "Date: $(date -Ru)"
    generate_hashes MD5Sum md5sum
    generate_hashes SHA1 sha1sum
    generate_hashes SHA256 sha256sum
    generate_hashes SHA512 sha512sum
  } > Release

  echo "Signing Release files"
  SIGN_ARGS=(--batch --yes --armor)
  if [ -n "${GPG_PASSPHRASE:-}" ]; then
    SIGN_ARGS+=(--pinentry-mode loopback --passphrase "${GPG_PASSPHRASE}")
  fi
  gpg "${SIGN_ARGS[@]}" --sign --detach-sign --output Release.gpg Release
  gpg "${SIGN_ARGS[@]}" --sign --detach-sign --clearsign --output InRelease Release

  popd >/dev/null
  echo "Apt repo generated"
}

generate_hashes() {
  HASH_TYPE="$1"
  HASH_COMMAND="$2"
  echo "${HASH_TYPE}:"
  find "${COMPONENTS:-main}" -type f | while read -r file
  do
    echo " $(${HASH_COMMAND} "$file" | cut -d" " -f1) $(wc -c "$file")"
  done
}

main
