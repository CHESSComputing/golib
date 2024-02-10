FROM golang:latest as go-builder
MAINTAINER Valentin Kuznetsov vkuznet@gmail.com

# args to use
ARG tag
ARG srv

# build procedure
ENV WDIR=/data
WORKDIR $WDIR
RUN mkdir /build
RUN git clone https://github.com/CHESSComputing/$srv
ARG CGO_ENABLED=0
RUN if [ "$tag" = "" ]; then cd $srv && make && cp srv /build; \
else cd $srv && git checkout tags/$tag -b build && make && cp srv /build; fi

# build final image for given image
# FROM alpine as final
# RUN mkdir -p /data
FROM gcr.io/distroless/static as final
COPY --from=go-builder /build/srv /data
