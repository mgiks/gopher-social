CREATE TABLE IF NOT EXISTS roles (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name varchar(55) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description text
);


/* 
The reason why we didn't delegate the insertion of
these roles to db seeding is due to the fact that
we're not going to seed the production database
 */
INSERT INTO roles (name, description, level)
    VALUES ('user', 'A user can create posts and comments', 1);

INSERT INTO roles (name, description, level)
    VALUES ('moderator', 'A moderator can update other users posts', 2);

INSERT INTO roles (name, description, level)
    VALUES ('admin', 'An admin can update and delete other users posts', 3);
