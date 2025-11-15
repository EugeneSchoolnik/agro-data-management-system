-- Видаляємо таблиці у зворотному порядку через залежності Foreign Key

DROP TABLE IF EXISTS forecasts;
DROP TABLE IF EXISTS pests;
DROP TABLE IF EXISTS metrics;
DROP TABLE IF EXISTS sensors;
DROP TABLE IF EXISTS fields;
DROP TABLE IF EXISTS crops;