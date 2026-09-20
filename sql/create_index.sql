-- Índices Espaciais (GIST) para buscas de proximidade no PostGIS
--CREATE INDEX idx_stops_location ON stops USING GIST (location);
--CREATE INDEX idx_shapes_location ON shapes USING GIST (location);

-- Índices Relacionais (B-Tree) para otimizar JOINs de rotas e viagens
CREATE INDEX idx_trips_route_id ON trips(route_id);
CREATE INDEX idx_trips_service_id ON trips(service_id);
CREATE INDEX idx_trips_shape_id ON trips(shape_id);

CREATE INDEX idx_stop_times_stop_id ON stop_times(stop_id);
CREATE INDEX idx_stop_times_trip_id ON stop_times(trip_id);

CREATE INDEX idx_shapes_shape_id ON shapes(shape_id);

-- Índice para buscas rápidas por e-mail no Login
CREATE INDEX idx_users_email ON users(email);