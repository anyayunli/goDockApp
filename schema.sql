GRANT ALL PRIVILEGES ON DATABASE goDockApp TO api_rw;
--;;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO api_rw;
--;;
GRANT ALL PRIVILEGES ON ALL sequences IN SCHEMA public TO api_rw;
--;;



CREATE TABLE IF NOT EXISTS users(
    id SERIAL NOT NULL PRIMARY KEY,
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp NULL,
    updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP);
--;;
CREATE INDEX user_email ON users USING btree(email);
--;;
