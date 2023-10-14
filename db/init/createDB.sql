DROP TABLE IF EXISTS orders CASCADE;

CREATE TABLE orders (
    data json
);
COPY orders(data) FROM '/docker-entrypoint-initdb.d/model.json';
