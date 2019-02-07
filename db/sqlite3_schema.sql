CREATE TABLE log(
    at timestamp default current_timestamp,
    id      integer primary key autoincrement,
    city    varchar(32),
    state   varchar(32),
    zip     varchar(32),
    country varchar(32),
    reason  varchar(128)
);

CREATE TABLE postimage(
    at timestamp default current_timestamp,
    id      integer primary key autoincrement,
    city    varchar(32),
    state   varchar(32),
    zip     varchar(32),
    country varchar(32)
);

CREATE TABLE preimage(
    at timestamp default current_timestamp,
    id      integer primary key autoincrement,
    city    varchar(32),
    state   varchar(32),
    zip     varchar(32),
    country varchar(32)
);
