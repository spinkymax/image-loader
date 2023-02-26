CREATE TABLE  users (
                       id serial PRIMARY KEY,
                       name text NOT NULL,
                       description text,
                       login text NOT NULL,
                       password text NOT NULL
);