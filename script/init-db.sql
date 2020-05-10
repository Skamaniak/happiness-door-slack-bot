DROP TABLE IF EXISTS happiness_door_feedback;
DROP TABLE IF EXISTS happiness_door;

CREATE TABLE happiness_door
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(256) NOT NULL,
    happy   INT NOT NULL DEFAULT 0,
    neutral INT NOT NULL DEFAULT 0,
    sad     INT NOT NULL DEFAULT 0
);

CREATE TABLE happiness_door_feedback
(
    id                SERIAL PRIMARY KEY,
    happiness_door_id INT REFERENCES happiness_door (id),
    feedback          TEXT
);