FROM postgres:17
COPY order.sql /docker-entrypoint-initdb.d/1.sql
CMD ["postgres"]