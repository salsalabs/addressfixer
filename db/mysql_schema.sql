USE addressfixer;

CREATE TABLE log(at timestamp default current_timestamp,
    id      varchar(32),
    city    varchar(32),
    state   varchar(32),
    zip     varchar(32),
    country varchar(32),
    reason  varchar(128)
);

CREATE TABLE postimage(at timestamp default current_timestamp,
    id      varchar(32),
    city    varchar(32),
    state   varchar(32),
    zip     varchar(32),
    country varchar(32)
);

CREATE TABLE preimage(at timestamp default current_timestamp,
    id      varchar(32),
    city    varchar(32),
    state   varchar(32),
    zip     varchar(32),
    country varchar(32)
);

CREATE UNIQUE INDEX preimage_id
    ON preimage(id);

CREATE UNIQUE INDEX postimage_id
    ON postimage(id); 

