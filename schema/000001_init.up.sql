CREATE TABLE users
(
    id serial not null unique,
    username varchar(255) not null,
    password_hash varchar(255) not null
);

CREATE TABLE chats
(
    id  serial not null unique,
    name varchar(255) not null
);

CREATE TABLE messages
(
    id serial not null unique,
    user_id integer references users (id) on delete restrict on update cascade not null,
    chat_id integer references chats (id) on delete cascade on update cascade not null,
    body varchar(255) not null,
    time timestamp not null
);
