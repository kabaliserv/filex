FROM alpine:3.15.0

VOLUME /data
EXPOSE 3000

ENV XDG_CACHE_HOME /data
ENV FILEX_DATABASE_DRIVER sqlite3
ENV FILEX_DATABASE_DATASOURCE /data/database.sqlite

ADD release/linux/amd64/filex /usr/local/bin/
ENTRYPOINT ["filex"]