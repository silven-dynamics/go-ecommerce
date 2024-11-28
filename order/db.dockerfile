FROM postgres:17-alpine
COPY ./order/order.sql /docker-entrypoint-initdb.d/1-order.sql
CMD ["postgres"]