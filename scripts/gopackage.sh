#!/usr/bin/env bash

set -e

source $(dirname $0)/version.sh

echo "==> Packaging binaries version ${VERSION} ..."

DIST=$(pwd)/dist/artifacts

mkdir -p $DIST/${VERSION}

for i in build/bin/*; do
    if [ ! -e $i ]; then
        continue
    fi

    BASE=build/archive
    DIR=${BASE}/${VERSION}

    rm -rf $BASE
    mkdir -p $BASE $DIR

    EXT=
    if [[ $i =~ .*windows.* ]]; then
        EXT=.exe
    fi

    cp $i ${DIR}/terraform-provider-rke_${VERSION}${EXT}

    (
        cd $DIR
        NAME=$(basename $i | cut -f1 -d_)
        ARCH=$(basename $i | cut -f2,3 -d_ | cut -f1 -d.)
        ARCHIVE=${NAME}_$(echo $VERSION | sed "s/^[v|V]//")_${ARCH}.zip
        echo "Packaging dist/artifacts/${VERSION}/${ARCHIVE} ..."
        zip -q $DIST/${VERSION}/${ARCHIVE} *
    )
done

(
    cd $DIST/${VERSION}/
    shasum -a 256 * > terraform-provider-rke_${VERSION}_SHA256SUMS
)

