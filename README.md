# hld
Fiber optic network high-level design

# db setup
psql postgres

CREATE DATABASE hld WITH OWNER = hld;

\c hld

CREATE EXTENSION postgis;

CREATE TABLE roads (
    id SERIAL PRIMARY KEY,
    geo_data GEOMETRY
);