CREATE DATABASE IF NOT EXISTS urlshortner;

CREATE TABLE IF NOT EXISTS urlshortner.urls (
   id         SERIAL,
   key        TEXT UNIQUE PRIMARY KEY NOT NULL,
   url        TEXT UNIQUE
);