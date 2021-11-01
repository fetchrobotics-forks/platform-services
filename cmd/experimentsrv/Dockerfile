# Deliver the experiment server as a container
#FROM alpine:latest
FROM ubuntu:18.04

EXPOSE 30001

#RUN apk --no-cache add ca-certificates
# Issues with laptops sleeping can push the clock into strange states, ignore this for our current purposes
RUN apt-get -y -o Acquire::Check-Valid-Until=false -o Acquire::Check-Date=false update  && \
    apt-get -y -o Acquire::Check-Valid-Until=false -o Acquire::Check-Date=false install net-tools iproute2 curl libssl1.0.0 libgflags2.2 iputils-ping telnet wget gnupg gnupg2 lsb-release

RUN adduser --system app --group
USER app

WORKDIR /home/app

COPY run.sh .
COPY bin/experimentsrv .
CMD ["./run.sh"]
