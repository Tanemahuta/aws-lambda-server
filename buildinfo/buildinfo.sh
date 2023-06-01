#!/bin/sh
rm -f buildinfo.txt
if [ -z "${VERSION}" ]; then
  VERSION="local-development"
fi
if [ -z "${COMMIT_SHA}" ]; then
  COMMIT_SHA=$(git rev-parse --short HEAD)
fi
TIMESTAMP=$(date +%Y-%m-%dT%H:%M:%S%z)
cat <<EOF >buildinfo.txt
${VERSION}
${COMMIT_SHA}
${TIMESTAMP}
EOF
