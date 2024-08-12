VERSION=""

function logInfo() {
    echo $(date "+%F %T") "INFO:" $@ 1>&2
}

function getVersion() {
    latest_tag=$(git describe --tags --abbrev=0)
    commit_count=$(git rev-list ${latest_tag}..HEAD --count)

    is_dirty=$(git status --porcelain)
    if [ -n "$is_dirty" ]; then
        dirty="-dev"
    else
        dirty=""
    fi
    echo "${latest_tag}.${commit_count}${dirty}"
}

function buildFrontend() {
    logInfo "========= build frontend ========= "
    # npm config set registry https://npmmirror.com/
    npm config set registry https://registry.npmmirror.com/

    logInfo ">>>>>> npm install --fix-missing"
    npm install --fix-missing || exit 1

    # logInfo ">>>>>> npm audit fix"
    # npm audit fix || exit 1

    logInfo ">>>>>> npm audit build"
    npm run build || exit 1
}

function buildBackend() {
    logInfo "========= build backend ========= "
    logInfo ">>>>>> pack resources"
    rm -rf internal/packed/resources.go
    rm -rf internal/packed/config.go
    gf pack manifest internal/packed/config.go --prefix manifest || exit 1
    sed -i 's|http://localhost:8081||g' ../skylight-web/dist/config.json
    gf pack ../skylight-web/dist internal/packed/resources.go --prefix resources || exit 1

    # wget -q https://dl.google.com/go/go1.21.4.linux-amd64.tar.gz
    # rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
    # cp /usr/local/go/bin/* /usr/bin/
    # echo 'export PATH=/usr/local/go/bin:$PATH' >> $HOME/.bashrc
    # source $HOME/.bashrc && /usr/local/go/bin/go version

    logInfo ">>>>>> go mod download"
    go mod download

    GO_VERSION=$(go version | awk '{print $3}')
    BUILD_DATE=$(date +'%Y-%m-%d %H:%M:%S')
    UNAME=$(uname -si)

    logInfo ">>>>>> go build"
    rm -rf ./skylight
    go build -ldflags " \
        -X 'main.Version=${VERSION}' \
        -X 'main.GoVersion=${GO_VERSION}' \
        -X 'main.BuildDate=${BUILD_DATE}' \
        -X 'main.BuildPlatform=${UNAME}' -s -w" || exit 1
    rm -rf internal/packed/resources.go
    rm -rf internal/packed/config.go

    logInfo ">>>>>> compress"
    upx -q skylight
    ./skylight version
}

function main() {
    logInfo ">>>>>> make semver"
    VERSION=$(getVersion)
    logInfo "version: ${VERSION}"

    logInfo ">>>>>> install packages"
    yum install -y tar upx wget || exit 1
    go version
    if [[ $? -ne 0 ]]; then
        yum install -y golang || exit 1
    fi
    node --version
    if [[ $? -ne 0 ]]; then
        logInfo ">>>>>> install nodejs"
        mkdir -p /usr/local/src/
        cd /usr/local/src/
        wget https://mirrors.aliyun.com/nodejs-release/v22.5.0/node-v22.5.0-linux-x64.tar.xz || exit
        # wget https://nodejs.org/dist/v22.5.0/node-v22.5.0-linux-x64.tar.xz  || exit 1
        tar xf node-v22.5.0-linux-x64.tar.xz || exit 1
        cd -
        ln -s /usr/local/src/node-v22.5.0-linux-x64/bin/node /usr/bin/node
        ln -s /usr/local/src/node-v22.5.0-linux-x64/bin/npm /usr/bin/npm
        ln -s /usr/local/src/node-v22.5.0-linux-x64/bin/npx /usr/bin/npx
        ln -s /usr/local/src/node-v22.5.0-linux-x64/bin/corepack /usr/bin/corepack
    fi

    go env -w GO111MODULE="on"
    go env -w GOPROXY="https://mirrors.aliyun.com/goproxy/,direct"

    cd skylight-web
    buildFrontend
    cd -

    cd skylight-go
    buildBackend
    cd -

    logInfo ">>>>>> make packages"

    RELEASE_PACKAGE="skylight-${VERSION}"

    rm -rf release/${RELEASE_PACKAGE}
    mkdir release/${RELEASE_PACKAGE}
    mv skylight-go/skylight release/${RELEASE_PACKAGE} || exit 1
    cd release
    cp install.sh ${RELEASE_PACKAGE}
    tar czf ${RELEASE_PACKAGE}.tar.gz ${RELEASE_PACKAGE} || exit 1
    rm -rf ${RELEASE_PACKAGE}
    cd -
    rm -rf skylight-web/dist
}

main
