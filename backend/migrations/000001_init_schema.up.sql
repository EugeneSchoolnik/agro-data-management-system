-- Таблиця культур
CREATE TABLE IF NOT EXISTS crops (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    variety VARCHAR(100), -- сорт
    description TEXT
);

-- Таблиця полів
CREATE TABLE IF NOT EXISTS fields (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    area FLOAT NOT NULL,
    location VARCHAR(255),
    crop_id INT REFERENCES crops(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблиця датчиків
CREATE TABLE IF NOT EXISTS sensors (
    id SERIAL PRIMARY KEY,
    field_id INT REFERENCES fields(id) ON DELETE CASCADE,
    sensor_type VARCHAR(50) NOT NULL, -- 'temperature', 'humidity', 'vision_node'
    status VARCHAR(20) DEFAULT 'active',
    last_sync TIMESTAMP
);

-- Таблиця метрик (сирі дані)
CREATE TABLE IF NOT EXISTS metrics (
    id BIGSERIAL PRIMARY KEY,
    sensor_id INT REFERENCES sensors(id) ON DELETE CASCADE,
    value FLOAT NOT NULL,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблиця шкідників (для бази знань ШІ)
CREATE TABLE IF NOT EXISTS pests (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    scientific_name VARCHAR(100),
    description TEXT
);

-- Таблиця прогнозів (результати роботи Python Engine)
CREATE TABLE IF NOT EXISTS forecasts (
    id SERIAL PRIMARY KEY,
    field_id INT REFERENCES fields(id) ON DELETE CASCADE,
    pest_id INT REFERENCES pests(id),
    probability FLOAT NOT NULL, -- від 0 до 1
    recommendation TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);