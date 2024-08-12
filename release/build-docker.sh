
TAR=$(ls -1 skylight* |tail -n1)

PACKAGE="${TAR%.tar.gz*}"
VERSION="${PACKAGE##skylight-}"

docker build --network=host --no-cache --build-arg PACKAGE=${PACKAGE} -t skylight:${VERSION} ./
