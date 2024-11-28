FROM postgres:17-alpine
COPY ./account/account.sql /docker-entrypoint-initdb.d/1-account.sql
CMD ["postgres"]