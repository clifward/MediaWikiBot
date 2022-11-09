#!/bin/bash
set -xe

VERSION=0.1.1
echo "build the application version ${VERSION}"
GOOS=linux CGO_ENABLED=0 go build ./cmd/wikibot

echo "zipping app & config file"
mv wikibot main
cp prod.yaml config.yaml
zip -r "wiki-bot-${VERSION}.zip" main config.yaml

echo "shipping to s3"
aws s3 cp "wiki-bot-${VERSION}.zip" s3://wikibot/

echo "clean up"
rm main -f
rm "wiki-bot-${VERSION}.zip" -f
rm config.yaml

echo "build & deploy finished"