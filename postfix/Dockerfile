# TODO: This docker file is for production but currently it is setup for development, change this

FROM alpine:3.12

RUN apk -U add postfix && echo "postmaster root" > /etc/aliases && newaliases
COPY . /etc/postfix
EXPOSE 25
ENTRYPOINT [ "postfix", "start-fg" ]