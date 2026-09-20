-- Ativa a extensão espacial do PostGIS (executar apenas uma vez por banco)
CREATE EXTENSION IF NOT EXISTS postgis;

-- 1. Criação do tipo ENUM no Postgres (Opcional, mas garante validação no banco)
CREATE TYPE user_role AS ENUM ('admin', 'user');

-- 2. Tabela de Utilizadores
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,               -- UUID v4
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,     -- Senha criptografada (ex: bcrypt)
    roles user_role[] NOT NULL DEFAULT '{user}', -- Array com as roles do enum
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- 1. AGENCY
CREATE TABLE IF NOT EXISTS agency (
    agency_id TEXT PRIMARY KEY,
    agency_name TEXT NOT NULL,
    agency_url TEXT NOT NULL,
    agency_timezone TEXT NOT NULL,
    agency_lang TEXT,
    agency_phone TEXT,
    agency_fare_url TEXT
);

-- 2. CALENDAR
CREATE TABLE IF NOT EXISTS calendar (
    service_id TEXT PRIMARY KEY,
    monday INT NOT NULL,
    tuesday INT NOT NULL,
    wednesday INT NOT NULL,
    thursday INT NOT NULL,
    friday INT NOT NULL,
    saturday INT NOT NULL,
    sunday INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);

-- 3. CALENDAR_DATES
CREATE TABLE IF NOT EXISTS calendar_dates (
    service_id TEXT NOT NULL,
    date DATE NOT NULL,
    exception_type INT NOT NULL, -- 1 = adicionado, 2 = removido
    PRIMARY KEY (service_id, date)
);

-- 4. FARE_ATTRIBUTES
CREATE TABLE IF NOT EXISTS fare_attributes (
    fare_id TEXT PRIMARY KEY,
    price NUMERIC(10, 2) NOT NULL,
    currency_type TEXT NOT NULL,
    payment_method INT NOT NULL, -- 0 = pago a bordo, 1 = pago antes do embarque
    transfers INT, -- 0 = sem transbordo, 1 = 1 transbordo, 2 = 2 transbordos, NULL = ilimitado
    agency_id TEXT REFERENCES agency(agency_id) ON DELETE CASCADE
);


-- 6. SHAPES
CREATE TABLE IF NOT EXISTS shapes (
    shape_id TEXT NOT NULL,
    shape_pt_lat DOUBLE PRECISION NOT NULL,
    shape_pt_lon DOUBLE PRECISION NOT NULL,
    shape_pt_sequence INT NOT NULL,
    PRIMARY KEY (shape_id, shape_pt_sequence)
);

-- 7. ROUTES
CREATE TABLE IF NOT EXISTS routes (
    route_id TEXT PRIMARY KEY,
    agency_id TEXT REFERENCES agency(agency_id) ON DELETE SET NULL,
    route_short_name TEXT,
    route_long_name TEXT,
    route_desc TEXT,
    route_type INT NOT NULL, -- 3 = Ônibus, 1 = Metrô, 2 = Trem, etc.
    route_url TEXT,
    route_color TEXT,
    route_text_color TEXT
);

-- 8. TRIPS
CREATE TABLE IF NOT EXISTS trips (
    route_id TEXT NOT NULL REFERENCES routes(route_id) ON DELETE CASCADE,
    service_id TEXT NOT NULL,
    trip_id TEXT PRIMARY KEY,
    trip_headsign TEXT,
    trip_short_name TEXT,
    direction_id INT, -- 0 = ida, 1 = volta
    block_id TEXT,
    shape_id TEXT,
    wheelchair_accessible INT
);

-- 9. STOPS
CREATE TABLE IF NOT EXISTS stops (
    stop_id TEXT PRIMARY KEY,
    stop_code TEXT,
    stop_name TEXT NOT NULL,
    stop_desc TEXT,
    stop_lat DOUBLE PRECISION NOT NULL,
    stop_lon DOUBLE PRECISION NOT NULL
);

-- 10. STOP_TIMES
CREATE TABLE IF NOT EXISTS stop_times (
    trip_id TEXT NOT NULL REFERENCES trips(trip_id) ON DELETE CASCADE,
    arrival_time INTERVAL, -- Ex: '08:30:00' (INTERVAL acomoda horários > 24h comuns no GTFS)
    departure_time INTERVAL,
    stop_id TEXT NOT NULL REFERENCES stops(stop_id) ON DELETE CASCADE,
    stop_sequence INT NOT NULL,
    PRIMARY KEY (trip_id, stop_sequence)
);