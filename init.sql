
CREATE TABLE IF NOT EXISTS cats (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    experience INT DEFAULT 0,
    breed TEXT NOT NULL,
    salary INT DEFAULT 0
);


CREATE TABLE IF NOT EXISTS missions (
    id SERIAL PRIMARY KEY,
    cat_id INT REFERENCES cats(id) ON DELETE SET NULL,
    complete BOOLEAN DEFAULT FALSE
);


CREATE TABLE IF NOT EXISTS targets (
    id SERIAL PRIMARY KEY,
    mission_id INT NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    country TEXT,
    notes TEXT,
    complete BOOLEAN DEFAULT FALSE,
    UNIQUE (mission_id, name)
);

insert into cats (name, experience, breed, salary) values ('Ryzhyk', 4, 'Siamese', 2008);