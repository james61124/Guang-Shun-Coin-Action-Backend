FROM golang:1.19

RUN apt-get update && \
      apt-get -y install sudo
      
RUN apt-get update \
    && apt-get install -y tmux vim git \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

USER root
CMD /bin/bash
