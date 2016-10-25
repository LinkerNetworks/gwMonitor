FROM alpine:latest
MAINTAINER Yifa Zhang <zyfdegg@gmail.com>

USER root
WORKDIR /root

# fix library dependencies
# otherwise golang binary may encounter 'not found' error
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY pgw.json /root/pgw.json
COPY sgw.json /root/sgw.json
COPY monitor /root/monitor
COPY monitor.conf /root/monitor.conf

RUN chmod +x /root/monitor

CMD ["/root/monitor"]
