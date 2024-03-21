FROM golang:1.21.3-bullseye

WORKDIR /usr/src/app

COPY . .

RUN apt-get -y update \
    && apt-get -y install netcat \
        build-essential \
        openssl \
        tzdata git \
    && make install \
    && make build

RUN apt-get -y install supervisor   

ENV TZ=Asia/Jakarta   

ARG APP_VERSION
ENV APP_VERSION=$APP_VERSION

# Expose port
EXPOSE 9000

RUN chmod 775 /temp
COPY build/image/supervisord.conf /etc/supervisord.conf
RUN chmod 775 /etc/supervisord.conf
RUN touch /usr/src/app/supervisord.log && chmod 775 /usr/src/app/supervisord.log

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
