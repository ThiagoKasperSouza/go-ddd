\copy agency FROM './files/agency.txt' WITH (FORMAT csv, HEADER true, DELIMITER ',');

\copy calendar FROM './files/calendar.txt' WITH (FORMAT csv, HEADER true, DELIMITER ',');

\copy calendar_dates FROM './files/calendar_dates.txt' WITH (FORMAT csv, HEADER true, DELIMITER ',');

\copy fare_attributes FROM './files/fare_attributes.txt' WITH (FORMAT csv, HEADER true, DELIMITER ',');

\copy shapes FROM './files/shapes.txt' WITH (FORMAT csv, HEADER true, DELIMITER ',');

\copy routes FROM './files/routes.txt' WITH (FORMAT csv, HEADER true, DELIMITER ',');

\copy stops FROM './files/stops.txt' WITH (FORMAT csv, HEADER true, DELIMITER ',');

\copy trips FROM './files/trips.txt' WITH (FORMAT csv, HEADER true, DELIMITER ',');

\copy stop_times FROM './files/stop_times.txt' WITH (FORMAT csv, HEADER true, DELIMITER ',');