#! /bin/sh

VERSION=$1

if [[ "$VERSION" == "" ]]; then
    echo "no version specified"
    exit 1
fi

result=0

PREV_VERSION=$(grep ^VERSION Makefile | tr -d 'VERSION := ')
grep "$VERSION" Makefile 2>&1 > /dev/null || result=$?
if [[ $result -eq 0 && $VERSION == $PREV_VERSION ]]; then
    echo "already updated to $VERSION"
    exit 1
fi

for i in README.md docs/content/_index.md ; do
    sed -e "s#Version-${PREV_VERSION}-information#Version-${VERSION//-/--}-information#g" -e "s#tag/v${PREV_VERSION}#tag/v${VERSION}#g" $i > a
    mv a $i
done

sed "s/VERSION := .*/VERSION := ${VERSION}/g" Makefile > a && mv a Makefile
sed "s/const VERSION = \".*\"/const VERSION = \"${VERSION}\"/g" config.go > a && mv a config.go

echo "Replace version to \"${VERSION}\""
