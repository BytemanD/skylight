# TODO: 动态生成项目版本
VERSION=0.1.0
RELEASE_PACKAGE="skylight-${VERSION}"

function logInfo() {
    echo `date "+%F %T" ` "INFO:" $@ 1>&2
}

function buildFrontend() {
    logInfo "========= build backend ========= "
    apt-get install -y npm nodejs || exit 1
    # npm config set registry https://npmmirror.com/
    npm config set registry https://registry.npmmirror.com/

    logInfo ">>>>>> npm install --fix-missing"
    npm install --fix-missing || exit 1

    # logInfo ">>>>>> npm audit fix"
    # npm audit fix || exit 1

    logInfo ">>>>>> npm audit build"
    npm run build || exit 1
}

function buildBackend(){
    logInfo "========= build backend ========= "
    apt-get install -y upx wget
    # wget -q https://dl.google.com/go/go1.21.4.linux-amd64.tar.gz
    # rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
    # cp /usr/local/go/bin/* /usr/bin/
    # echo 'export PATH=/usr/local/go/bin:$PATH' >> $HOME/.bashrc
    # source $HOME/.bashrc && /usr/local/go/bin/go version
    logInfo ">>>>>> install go"
    apt-get install -y golang || exit 1
    go env -w GO111MODULE="on"
    go env -w GOPROXY="https://mirrors.aliyun.com/goproxy/,direct"

    logInfo ">>>>>> go mod download"
    go mod download

    GO_VERSION=$(go version |awk '{print $3}')
    BUILD_DATE=$(date +'%Y-%m-%d %H:%M:%S')
    UNAME=$(uname -si)

    logInfo ">>>>>> go build"
    rm -rf ./skylight
    go build -ldflags " \
        -X 'main.Version=${VERSION}' \
        -X 'main.GoVersion=${GO_VERSION}' \
        -X 'main.BuildDate=${BUILD_DATE}' \
        -X 'main.BuildPlatform=${UNAME}' -s -w"


    logInfo ">>>>>> compress"
    upx -q skylight
    ./skylight version
}


apt-get install -y tar

cd skylight-web
buildFrontend
cd ..

cd skylight-go
buildBackend
cd ..

releasePath="release/${RELEASE_PACKAGE}"
logInfo ">>>>>> create package: ${releasePath}"
rm -rf ${releasePath}
mkdir -p ${releasePath}

cp -r skylight-web/dist ${releasePath}/web || exit 1
sed -i 's|http://localhost:8081||g' ${releasePath}/web/config.json ||exit 1
cp skylight-go/skylight ${releasePath} || exit 1

cd release
cp install.sh ${RELEASE_PACKAGE}
tar -czf ${RELEASE_PACKAGE}.tar.gz ${RELEASE_PACKAGE} || exit 1
cd ../

rm -rf skylight-web/dist skylight-go/skylight

