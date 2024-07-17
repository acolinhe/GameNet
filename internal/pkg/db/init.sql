CREATE DATABASE gamenetdb;

\connect gamenetdb;

CREATE TABLE Games (
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(255) NOT NULL,
                       summary TEXT,
                       release_date VARCHAR(255)
);

CREATE TABLE Developers (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(255) NOT NULL
);

CREATE TABLE Genres (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(255) NOT NULL
);

CREATE TABLE Platforms (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(255) NOT NULL
);

CREATE TABLE GameDevelopers (
                            game_id INTEGER REFERENCES Games(id),
                            developer_id INTEGER REFERENCES Developers(id),
                            PRIMARY KEY (game_id, developer_id)
);

CREATE TABLE GameGenres (
                            game_id INTEGER REFERENCES Games(id),
                            genre_id INTEGER REFERENCES Genres(id),
                            PRIMARY KEY (game_id, genre_id)
);

CREATE TABLE GamePlatforms (
                            game_id INTEGER REFERENCES Games(id),
                            platform_id INTEGER REFERENCES Platforms(id),
                            PRIMARY KEY (game_id, platform_id)
);
