FROM hub.oepkgs.net/openeuler/openeuler:22.09

ARG PACKAGE
COPY ${PACKAGE}.tar.gz /tmp
COPY openEuler.repo /tmp

RUN cd /etc/yum.repos.d/ \
    && mkdir backup && mv *.repo backup \
    && mv /tmp/openEuler.repo /etc/yum.repos.d/

RUN rpm -q glibc
RUN yum install glibc-devel -y
RUN cd /tmp && tar -xzf ${PACKAGE}.tar.gz \
    && cd ${PACKAGE} \
    && sh install.sh
