# FROM hub.oepkgs.net/openeuler/openeuler:22.09
FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ubuntu:22.04

ARG DATE
ARG PACKAGE
COPY ${PACKAGE}.tar.gz /tmp
# COPY openEuler.repo /tmp

# RUN cd /etc/yum.repos.d/ \
#     && mkdir backup && mv *.repo backup \
#     && mv /tmp/openEuler.repo /etc/yum.repos.d/
RUN echo $DATE
RUN cd /tmp && tar -xzf ${PACKAGE}.tar.gz \
    && cd ${PACKAGE} \
    && sh install.sh \
    && cd ../ \
    && rm -rf  ${PACKAGE}.tar.gz ${PACKAGE}

ENTRYPOINT skylight serve