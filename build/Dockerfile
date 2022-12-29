FROM harbor.bangcle.net:8029/library/alpine:latest
LABEL maintainer="guowei.cao <guowei.cao@bangcle.com>"
ENV LANG=en_US.UTF-8 LANGAGE=en_US.UTF-8 TZ=Asia/Shanghai

WORKDIR /data

COPY abc .
COPY deploy.ini .
COPY run.sh .
COPY glibc-2.34-r0.apk .
COPY resources ./resources

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add ca-certificates \
    && mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 \
    && apk --update add bash mysql-client tzdata \
    && /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'Asia/Shanghai' > /etc/timezone \
    && apk add automake libtool make gcc pkgconfig g++ zip unzip \
    && apk add --no-cache -U autoconf \
    && apk del automake libtool make g++ gcc pkgconfig autoconf \
    && wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
    && cd /data && apk add --allow-untrusted glibc-2.34-r0.apk \
    && chmod +x run.sh


EXPOSE 9100/tcp
CMD ["./run.sh"]