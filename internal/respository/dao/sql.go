package dao

var databases = `create database if not exists webbook default charset =utf8mb4;`
var initSqlUser = `CREATE TABLE IF NOT EXISTS webbook.users(
    id bigint primary key auto_increment,
    email varchar(30) not null unique,
    password text not null,
    ctime bigint,
    utime bigint
    );`

var initProfile = `
CREATE TABLE  IF NOT EXISTS webbook.profile(
                                               id bigint(20) primary key auto_increment,
                                               email varchar(30) NOT NULL unique,
                                               username varchar(15),
                                               birthday bigint(20),
                                               personalprofile text
)`
