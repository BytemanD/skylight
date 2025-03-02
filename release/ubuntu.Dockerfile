FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ubuntu:22.04

ARG DATE

RUN echo $DATE
COPY build /tmp/build
RUN mkdir -p /var/lib/skylight/gsessions /var/lib/skylight/db /usr/local/bin \
    && cd /tmp \
    && cp build/skylight /usr/local/bin         \
    && chmod u+x /usr/local/bin/skylight        \
    && cp -r build/migrations /var/lib/skylight \
    && cp -r build/config /var/lib/skylight     \
    && cp -r build/dist /var/lib/skylight/WEB

ENTRYPOINT skylight serve
