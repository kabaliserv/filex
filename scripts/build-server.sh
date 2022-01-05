#!/bin/sh

#!/bin/sh

echo "building docker images for ${GOOS}/${GOARCH} ..."

REPO="github.com/kabaliserv/filex"

# compile the server using the cgo
go build -ldflags='-X github.com/kabaliserv/filex/web.Mode=prod' -o release/linux/"${GOARCH}"/filex ${REPO}/cmd/kbs-filex
