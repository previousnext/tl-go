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
  # dpkg-scanpackages' --arch flag matches *_<arch>.deb filenames, not the
  # package's actual Architecture: field, and nfpm names amd64 packages with a
  # goamd64 suffix (e.g. tl_1.0.0_linux_amd64v3.deb) that doesn't match that
  # pattern. Scan once with no filter (which does read the real Architecture:
  # field) and split the output ourselves.
  ALL_PACKAGES="$(mktemp)"
  dpkg-scanpackages --multiversion pool/ > "${ALL_PACKAGES}"

  for ARCH in "${ARCHS[@]}"; do
    PACKAGE_DIR="${COMPONENT_DIR}/binary-${ARCH}"
    mkdir -p "${PACKAGE_DIR}"
    awk -v RS="" -v ORS="\n\n" -v arch="${ARCH}" '
      $0 ~ ("(^|\n)Architecture: " arch "(\n|$)") { print }
    ' "${ALL_PACKAGES}" > "${PACKAGE_DIR}/Packages"
    gzip -fk "${PACKAGE_DIR}/Packages"
    bzip2 -fk "${PACKAGE_DIR}/Packages"
  done
  rm -f "${ALL_PACKAGES}"

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
