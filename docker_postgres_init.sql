CREATE USER books_user WITH PASSWORD 'password' CREATEDB;
CREATE USER users_user WITH PASSWORD 'password' CREATEDB;
CREATE USER chat_user WITH PASSWORD 'password' CREATEDB;

CREATE DATABASE books WITH OWNER = books_user;
CREATE DATABASE users WITH OWNER = users_user;
CREATE DATABASE chat WITH OWNER = chat_user;
