CREATE TABLE IF NOT EXISTS numbers (
 id SERIAL PRIMARY KEY,
 number VARCHAR(255)
);

INSERT INTO numbers (number)
SELECT (floor(random()*10000000000)::bigint)::text
FROM generate_series(1, 10);