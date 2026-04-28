-- +goose Up
-- Создание таблицы событий календаря
CREATE TABLE IF NOT EXISTS events (
    id VARCHAR(255) PRIMARY KEY,                        -- Уникальный идентификатор события (UUID)
    title VARCHAR(255) NOT NULL,                        -- Заголовок события (короткий текст)
    start_time TIMESTAMP NOT NULL,                      -- Дата и время начала события
    end_time TIMESTAMP NOT NULL,                        -- Дата и время окончания события (длительность)
    description TEXT,                                   -- Описание события (длинный текст, опционально)
    user_id VARCHAR(255) NOT NULL,                      -- ID пользователя-владельца события
    notify_before BIGINT DEFAULT 0,                     -- За сколько времени высылать уведомление (в наносекундах, опционально)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP      -- Дата создания записи
);

CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_end_time ON events(end_time);

-- +goose Down
DROP TABLE IF EXISTS events;
