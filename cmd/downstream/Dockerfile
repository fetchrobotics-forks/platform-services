# Deliver the downstream service as a container
#FROM alpine:latest
FROM ubuntu:18.04

EXPOSE 30001

ENV DEBIAN_FRONTEND=noninteractive
ENV LANG=en_US.UTF-8

# Issues with laptops sleeping can push the clock into strange states, ignore this for our current purposes
RUN apt-get -y -o Acquire::Check-Valid-Until=false -o Acquire::Check-Date=false update  && \
    apt-get -y -o Acquire::Check-Valid-Until=false -o Acquire::Check-Date=false install net-tools iproute2 curl libssl1.0.0 libgflags2.2

RUN adduser --system app --group
USER app

WORKDIR /home/app
COPY bin/downstream .
CMD ["./downstream"]
