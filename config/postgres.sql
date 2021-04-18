
drop table IF EXISTS languages;
drop table IF EXISTS editors;
drop table IF EXISTS authors;
drop index IF EXISTS index_books_on_author_id;
drop index IF EXISTS index_books_on_editor_id;
drop index IF EXISTS index_books_on_language_id;
drop table IF EXISTS books;
drop index IF EXISTS index_sessions_on_user_id;
drop table IF EXISTS sessions;
drop index IF EXISTS index_threads_on_user_id;
drop table IF EXISTS threads;
drop index IF EXISTS index_shadows_on_user_id;
drop table IF EXISTS shadows;
drop table IF EXISTS users;


create TABLE IF NOT EXISTS users (
    id         serial primary key,
    uuid       varchar(64) not null unique,
    cuenta     varchar(16),
    password   varchar(64),
    email      varchar(32),
    level      integer,
    created_at timestamp not null,   
    updated_at timestamp not null   
);

create TABLE IF NOT EXISTS shadows (
    id         serial primary key,
    user_id    integer references users(id) NOT NULL,
    uuid       varchar(64) not null unique,
    password   varchar(64),
    created_at timestamp not null,   
    updated_at timestamp not null   
);

CREATE UNIQUE INDEX index_shadows_on_user_id ON shadows (user_id);

create table IF NOT EXISTS sessions (
  id         serial primary key,
  user_id    integer references users(id) NOT NULL,
  uuid       varchar(64) not null unique,
  created_at timestamp not null,   
  updated_at timestamp not null   
);

CREATE UNIQUE INDEX  index_sessions_on_user_id ON  sessions (user_id);

create table IF NOT EXISTS threads (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  user_id    integer references users(id) NOT NULL,
  created_at timestamp not null,       
  updated_at  timestamp not null   
);

CREATE UNIQUE INDEX  index_threads_on_user_id  ON threads (user_id);

CREATE TABLE IF NOT EXISTS languages (
    id          serial primary key,
    name        varchar(48),
    created_at  timestamp not null,   
    updated_at  timestamp not null   
);


CREATE TABLE IF NOT EXISTS editors (
    id          serial primary key,
    name        varchar(64),
    created_at  timestamp not null,   
    updated_at  timestamp not null   
) ;


CREATE TABLE IF NOT EXISTS authors (
    id          serial primary key,
    name        varchar(64),
    created_at  timestamp not null,   
    updated_at  timestamp not null   
);


CREATE TABLE IF NOT EXISTS books (
    id          serial primary key,
    author_id   int,
    editor_id   int,
    language_id int,
    title       text,
    isbn        varchar(24),
    comment     text,
    year        int,
    created_at  timestamp not null,   
    updated_at  timestamp not null   
);

CREATE INDEX  index_books_on_author_id  ON books (author_id);
CREATE INDEX  index_books_on_editor_id  ON books (editor_id);
CREATE INDEX  index_books_on_language_id  ON books (language_id);


