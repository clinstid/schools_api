FROM alpine:3.7

COPY bin/schools_api /usr/bin

COPY dist /dist

COPY spec/schools_api_spec.yaml /dist

WORKDIR /

ENTRYPOINT /usr/bin/schools_api
