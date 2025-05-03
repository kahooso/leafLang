DROP TABLE IF EXISTS learning_attempt CASCADE;

DROP TABLE IF EXISTS user_word CASCADE;

DROP TABLE IF EXISTS "user" CASCADE;

DROP TABLE IF EXISTS role CASCADE;

DROP TABLE IF EXISTS word_status CASCADE;

CREATE TABLE
    role (
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR(50) NOT NULL UNIQUE CHECK (name IN ('Student', 'Admin'))
    );

CREATE TABLE
    "user" (
        id BIGSERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        phone VARCHAR(20),
        first_name VARCHAR(100) NOT NULL,
        last_name VARCHAR(100) NOT NULL,
        password VARCHAR(255) NOT NULL,
        role_id BIGINT NOT NULL DEFAULT 1 REFERENCES role (id),
        last_active TIMESTAMP DEFAULT now (),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP,
        image_url VARCHAR(255)
    );

CREATE TABLE
    word_status (
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR(50) NOT NULL UNIQUE CHECK (name in ('To learn', 'Known', 'Learned'))
    );

CREATE TABLE
    user_word (
        id BIGSERIAL PRIMARY KEY,
        user_id BIGINT NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
        original_word VARCHAR(255) NOT NULL,
        translation VARCHAR(255) NOT NULL,
        example TEXT,
        status_id BIGINT NOT NULL REFERENCES word_status (id),
        success_count INTEGER NOT NULL DEFAULT 0,
        fail_count INTEGER NOT NULL DEFAULT 0,
        last_reviewed TIMESTAMP,
        next_review_at TIMESTAMP,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP
    );

CREATE TABLE
    learning_attempt (
        id BIGSERIAL PRIMARY KEY,
        user_word_id BIGINT NOT NULL REFERENCES user_word (id) ON DELETE CASCADE,
        attempt_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        success BOOLEAN NOT NULL
    );

INSERT INTO
    word_status (name)
VALUES
    ('To learn'),
    ('Known'),
    ('Learned');

INSERT INTO
    role (name)
VALUES
    ('Student'),
    ('Admin');