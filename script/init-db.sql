DROP TABLE IF EXISTS happiness_door_user_action;
DROP TABLE IF EXISTS happiness_door;

CREATE TABLE happiness_door
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL
);

CREATE TABLE happiness_door_user_action
(
    id                SERIAL PRIMARY KEY,
    happiness_door_id INT REFERENCES happiness_door (id),
    user_id           VARCHAR(32) NOT NULL,
    user_name         VARCHAR(64) NOT NULL,
    action_id         VARCHAR(32) NOT NULL,
    CONSTRAINT unique_user_vote UNIQUE (happiness_door_id, user_id) -- 1 Vote per person per meeting
);