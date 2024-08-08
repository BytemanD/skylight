
TAR=$(ls -1 skylight* |tail -n1)

PACKAGE="${TAR%.tar.gz}"

docker build --network=host --no-cache --build-arg PACKAGE=${PACKAGE} -t ${PACKAGE} ./
