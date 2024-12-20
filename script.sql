-- Таблица ролей
DROP TABLE IF EXISTS role CASCADE;
CREATE TABLE role (
    id BIGSERIAL PRIMARY KEY,            -- Уникальный идентификатор роли с автоинкрементом
    access_rights BIGINT                 -- Уровень доступа
);

-- Таблица пользователей
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,            -- Уникальный идентификатор пользователя с автоинкрементом
    username TEXT NOT NULL,              -- Имя пользователя
    password_hash TEXT NOT NULL,         -- Хеш пароля (заменён с BIGINT на TEXT для хранения хэшей)
    email TEXT,                          -- Электронная почта
    role_id BIGINT,                      -- Ссылка на роль
    FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE SET NULL
);

-- Таблица артистов
DROP TABLE IF EXISTS artists CASCADE;
CREATE TABLE artists (
    id BIGSERIAL PRIMARY KEY,            -- Уникальный идентификатор артиста с автоинкрементом
    name TEXT NOT NULL                   -- Имя артиста
);

-- Таблица релизов
DROP TABLE IF EXISTS releases CASCADE;
CREATE TABLE releases (
    id BIGSERIAL PRIMARY KEY,            -- Уникальный идентификатор релиза с автоинкрементом
    name TEXT NOT NULL,                  -- Название релиза
    artist_id BIGINT,                    -- Ссылка на артиста
    FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE SET NULL
);

-- Таблица жанров
DROP TABLE IF EXISTS genre CASCADE;
CREATE TABLE genre (
    id BIGSERIAL PRIMARY KEY,            -- Уникальный идентификатор жанра с автоинкрементом
    genre TEXT                           -- Название жанра (заменён тип с BIGINT на TEXT)
);

-- Таблица связей релизов и жанров
DROP TABLE IF EXISTS release_genres CASCADE;
CREATE TABLE release_genres (
    id BIGSERIAL PRIMARY KEY,            -- Уникальный идентификатор записи с автоинкрементом
    release_id BIGINT,                   -- Ссылка на релиз
    genre_id BIGINT,                     -- Ссылка на жанр
    FOREIGN KEY (release_id) REFERENCES releases(id) ON DELETE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES genre(id) ON DELETE CASCADE
);

-- Таблица оценок релизов
DROP TABLE IF EXISTS release_score CASCADE;
CREATE TABLE release_score (
    id BIGSERIAL PRIMARY KEY,            -- Уникальный идентификатор оценки с автоинкрементом
    score BIGINT NOT NULL,               -- Оценка
    release_id BIGINT,                   -- Ссылка на релиз
    user_id BIGINT,                      -- Ссылка на пользователя
    FOREIGN KEY (release_id) REFERENCES releases(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO role (id, access_rights) VALUES (1, 1), (2, 2);

-- Добавление релизов в таблицу releases, с привязкой к артистам
INSERT INTO releases (name, artist_id)
VALUES
    ('Асфальт 8', 1),  -- Macan
    ('Самая самая', 2),  -- Егор Крид
    ('Судно', 3),  -- Молчат Дома
    ('Новые люди', 4);  -- Сплин

INSERT INTO genre (genre)
VALUES
    ('rap'),
    ('rock'),
    ('pop'),
	('doomer music'),
	('emo'), 
	('electronic');

INSERT INTO release_genres (release_id, genre_id)
VALUES
    (1, 1),
	(1, 3),
    (2, 3),
	(2, 1),
    (3, 4),  
    (4, 2);

-- Таблица логов
DROP TABLE IF EXISTS release_score_log CASCADE;
CREATE TABLE release_score_log (
    id BIGSERIAL PRIMARY KEY,        -- Уникальный идентификатор записи с автоинкрементом
    score BIGINT NOT NULL,           -- Оценка
    release_id BIGINT,               -- Ссылка на релиз
    user_id BIGINT,                  -- Ссылка на пользователя
    created_at TIMESTAMP DEFAULT NOW(), -- Время добавления записи
    FOREIGN KEY (release_id) REFERENCES releases(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Функция для логирования
CREATE OR REPLACE FUNCTION log_release_score()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO release_score_log (score, release_id, user_id)
    VALUES (NEW.score, NEW.release_id, NEW.user_id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер для логирования
CREATE TRIGGER release_score_after_insert
AFTER INSERT ON release_score
FOR EACH ROW
EXECUTE FUNCTION log_release_score();
