-- Weather parameters / параметри погоди
CREATE TABLE IF NOT EXISTS weather_parameters (
    id SERIAL PRIMARY KEY,
    param_id INT NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    unit VARCHAR(50),
    description TEXT
);

-- Weather stations / метеостанції
CREATE TABLE IF NOT EXISTS weather_stations (
    id SERIAL PRIMARY KEY,
    external_id INT NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    region VARCHAR(255),
    active BOOLEAN DEFAULT TRUE,
    last_seen TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Initial weather station seed
INSERT INTO weather_stations (external_id, name, region, active)
VALUES
    (9, 'Сумська, Сад', 'Сумська', TRUE),
    (20, 'Івано-Франківськ, Котиківка', 'Івано-Франківськ', TRUE),
    (24, 'Львівська, Білий Камінь', 'Львівська', TRUE),
    (25, 'Рівненська (демонтовано)', 'Рівненська', TRUE),
    (26, 'Житомир, Високе', 'Житомир', TRUE),
    (67, 'Чернігівська, Халявин', 'Чернігівська', TRUE),
    (68, 'Чернівецька, Ошихліби', 'Чернівецька', TRUE),
    (69, 'Полтавська, Покровське', 'Полтавська', TRUE),
    (70, 'Хмельницька Тр (демонтовано)', 'Хмельницька', TRUE),
    (233, 'Тернопільська, Плотича', 'Тернопільська', TRUE),
    (234, 'Луганська, Лозовівка', 'Луганська', TRUE),
    (235, 'Закарпатська, Великі Лучки', 'Закарпатська', TRUE),
    (236, 'Волинь, Деревок', 'Волинь', TRUE),
    (237, 'Полтавська, Градизьк', 'Полтавська', TRUE),
    (238, 'Черкаська, Дзензелівка', 'Черкаська', TRUE),
    (239, 'Волинська, Звиняч', 'Волинська', TRUE),
    (240, 'Сумська, Лікарське', 'Сумська', TRUE),
    (241, 'Дніпропетровська, Семенівка', 'Дніпропетровська', TRUE),
    (242, 'Полтавська, Карлівка', 'Полтавська', TRUE),
    (243, 'Кіровоградська, Новоселиця', 'Кіровоградська', TRUE),
    (244, 'Київська, Біла Церква', 'Київська', TRUE),
    (245, 'Харківська, Золочівське', 'Харківська', TRUE),
    (246, 'Вінницька, Голубече', 'Вінницька', TRUE),
    (247, 'Київська, Гостомель', 'Київська', TRUE),
    (335, 'Одеська, Кілія', 'Одеська', TRUE),
    (336, 'Одеська, Новоселівка', 'Одеська', TRUE),
    (337, 'Хмельницька, Іванівка', 'Хмельницька', TRUE),
    (420, 'Рівненська, Верхівськ', 'Рівненська', TRUE),
    (421, 'Хмельницька, Требухівці', 'Хмельницька', TRUE)
ON CONFLICT (external_id) DO NOTHING;

-- Station parameter mapping / відповідність параметрів станції до загальних параметрів
CREATE TABLE IF NOT EXISTS weather_station_parameters (
    id SERIAL PRIMARY KEY,
    station_id INT NOT NULL REFERENCES weather_stations(id) ON DELETE CASCADE,
    weather_parameter_id INT NOT NULL REFERENCES weather_parameters(id) ON DELETE CASCADE,
    station_param INT NOT NULL,
    UNIQUE(station_id, weather_parameter_id)
);

-- Weather stations / метеостанції
CREATE TABLE IF NOT EXISTS weather_stations (
    id SERIAL PRIMARY KEY,
    external_id INT NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    region VARCHAR(255),
    active BOOLEAN DEFAULT TRUE,
    last_seen TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Station parameter mapping / відповідність параметрів станції до загальних параметрів
CREATE TABLE IF NOT EXISTS weather_station_parameters (
    id SERIAL PRIMARY KEY,
    station_id INT NOT NULL REFERENCES weather_stations(id) ON DELETE CASCADE,
    weather_parameter_id INT NOT NULL REFERENCES weather_parameters(id) ON DELETE CASCADE,
    station_param INT NOT NULL,
    UNIQUE(station_id, weather_parameter_id)
);

-- Weather observations / зафіксовані показники
CREATE TABLE IF NOT EXISTS weather_observations (
    id BIGSERIAL PRIMARY KEY,
    station_id INT NOT NULL REFERENCES weather_stations(id) ON DELETE CASCADE,
    weather_parameter_id INT NOT NULL REFERENCES weather_parameters(id) ON DELETE SET NULL,
    station_param INT NOT NULL,
    value FLOAT NOT NULL,
    recorded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Initial parameter seed
INSERT INTO weather_parameters (param_id, name, unit, description)
VALUES
    (1, 't° повітря', '°C', 'Температура повітря'),
    (2, 'Атмосферний тиск', 'мм рт. ст.', 'Атмосферний тиск'),
    (3, 'Швидкість вітру', 'м/с', 'Швидкість вітру'),
    (4, 'Напрямок вітру', '', 'Напрямок вітру'),
    (5, 'Точка роси', '°C', 'Точка роси'),
    (7, 'Опади відносні', 'мм', 'Опади відносні'),
    (8, 'Опади накопичені', 'мм', 'Опади накопичені'),
    (9, 'Відносна вологість', '%', 'Відносна вологість'),
    (10, 't° грунту', '°C', 'Температура ґрунту'),
    (20, 'Сонячна радіація', 'Вт/м2', 'Сонячна радіація')
ON CONFLICT (param_id) DO NOTHING;

