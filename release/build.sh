VERSION=""

function logInfo() {
    echo $(date "+%F %T") "INFO:" "$@" 1>&2
}
function logError() {
    echo $(date "+%F %T") "ERROR:" "$@" 1>&2
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

function prepare() {
    logInfo ">>>>>> 清理构建缓存"
    rm -rf release/build/* && mkdir -p release/build

    logInfo ">>>>>> 安装 tar upx wget which"
    yum install -y tar upx wget which || exit 1

    which go > /dev/null 2>&1
    if [[ $? -ne 0 ]]; then
        logInfo ">>>>>> 安装 go"
        cd /tmp
        wget https://golang.google.cn/dl/go1.24.0.linux-amd64.tar.gz || exit 1
        tar -zxf go1.24.0.linux-amd64.tar.gz  -C /usr/local/     || exit 1
        cd -
    fi
    which node > /dev/null 2>&1
    if [[ $? -ne 0 ]]; then
        logInfo ">>>>>> 安装 nodejs"
        mkdir -p /usr/local/src/
        cd /usr/local/src/

        NODE_VERSION=v23.7.0
        rm -rf node-${NODE_VERSION}.tar.xz
        wget https://mirrors.aliyun.com/nodejs-release/${NODE_VERSION}/node-${NODE_VERSION}-linux-x64.tar.xz || exit
        # wget https://nodejs.org/dist/${NODE_VERSION}/node-${NODE_VERSION}-linux-x64.tar.xz  || exit 1
        tar xf node-${NODE_VERSION}.tar.xz || exit 1
        cd - 
        rm -rf  /usr/bin/node /usr/bin/npm  /usr/bin/npx /usr/bin/corepack
        ln -s /usr/local/src/node-${NODE_VERSION}-linux-x64/bin/node /usr/bin/node
        ln -s /usr/local/src/node-${NODE_VERSION}-linux-x64/bin/npm /usr/bin/npm
        ln -s /usr/local/src/node-${NODE_VERSION}-linux-x64/bin/npx /usr/bin/npx
        ln -s /usr/local/src/node-${NODE_VERSION}-linux-x64/bin/corepack /usr/bin/corepack
    fi

    go env -w GO111MODULE="on"
    go env -w GOPROXY="https://mirrors.aliyun.com/goproxy/,direct"
}

function buildFrontend() {
    # npm config set registry https://npmmirror.com/
    npm config set registry https://registry.npmmirror.com/

    logInfo ">>>>>> npm install --fix-missing"
    npm install --fix-missing || exit 1

    # logInfo ">>>>>> npm audit fix"
    # npm audit fix || exit 1

    rm -rf dist
    logInfo ">>>>>> npm audit build"
    npm run build || exit 1
}

function buildBackend() {
    logInfo ">>>>>> 打包配置资源"
    rm -rf internal/packed/config.go
    gf pack manifest internal/packed/config.go --prefix manifest || exit 1

    logInfo ">>>>>> 执行 go mod download"
    go mod download

    GO_VERSION=$(go version | awk '{print $3}')
    BUILD_DATE=$(date +'%Y-%m-%d %H:%M:%S')
    UNAME=$(uname -si)

    logInfo ">>>>>> 执行 go build"
    rm -rf ./skylight
    go build -ldflags " \
        -X 'main.Version=${VERSION}' \
        -X 'main.GoVersion=${GO_VERSION}' \
        -X 'main.BuildDate=${BUILD_DATE}' \
        -X 'main.BuildPlatform=${UNAME}' -s -w" || exit 1
    # rm -rf internal/packed/resources.go
    rm -rf internal/packed/config.go

    logInfo ">>>>>> 压缩二进制文件"
    upx -q skylight
    ./skylight version
}

function main() {
    docker ps >> /dev/null
    if [[ $? -ne 0 ]]; then
        logError "docker is required" 
        exit 1
    fi
    logInfo "========= 构建前准备 ========="
    prepare

    logInfo "========= 获取项目版本 ========="
    VERSION=$(getVersion)
    logInfo "版本: ${VERSION}"

    RELEASE_PACKAGE="skylight-${VERSION}"

    logInfo "========= 构建前端工程 ========= "
    cd skylight-web && buildFrontend && cd -
    mv skylight-web/dist release/build
    logInfo ">>>>>> 更新项目配置"
    sed -i 's|http://localhost:8081||g' release/build/dist/config.json
    sed -i 's|ws://localhost:8081||g'   release/build/dist/config.json

    logInfo "========= 构建后端工程 ========="
    cd skylight-go && buildBackend && cd -
    mv skylight-go/skylight release/build
    cp -r migrations release/build
    cp -r skylight-go/manifest/config release/build

    logInfo "========= 构建镜像 ========= "
    local imageName="skylight:${VERSION}"
    cd release
    docker build --file ubuntu.Dockerfile --network=host --no-cache \
        --build-arg DATE="$(date)" \
        -t ${imageName} ./   || exit 1
    cd -

    logInfo "========= 推送镜像 ========="
    docker tag ${imageName} registry.cn-hangzhou.aliyuncs.com/fjboy-ec/${imageName} || exit 1
    docker tag ${imageName} registry.cn-hangzhou.aliyuncs.com/fjboy-ec/skylight     || exit 1
    docker push registry.cn-hangzhou.aliyuncs.com/fjboy-ec/${imageName}             || exit 1
    docker push registry.cn-hangzhou.aliyuncs.com/fjboy-ec/skylight                 || exit 1
}

main
