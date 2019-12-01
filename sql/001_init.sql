-- +migrate Up
CREATE TABLE site (
    id serial PRIMARY KEY,
    name varchar(50) NOT NULL
);

CREATE TABLE location (
    id serial PRIMARY KEY,
    name varchar(50) NOT NULL,
    site integer NOT NULL,
    CONSTRAINT location_site_fkey
        FOREIGN KEY (site)
        REFERENCES site (id)
);

CREATE TABLE device (
    id serial PRIMARY KEY,
    mac macaddr NOT NULL,
    location integer NOT NULL,
    position integer NOT NULL,
    CONSTRAINT device_location_fkey
        FOREIGN KEY (location)
        REFERENCES location(id)
);

CREATE TABLE content (
    id serial PRIMARY KEY,
    url varchar(2000) NOT NULL,
    active boolean NOT NULL,
    duration interval NOT NULL DEFAULT '1 minute'
);

CREATE TYPE header_type AS ENUM ('RAW', 'FILE', 'ENV');

CREATE TABLE header (
    id serial PRIMARY KEY,
    type header_type NOT NULL,
    key varchar(255) NOT NULL,
    value varchar(255) NOT NULL,
    content integer NOT NULL,
    CONSTRAINT header_content_fkey
        FOREIGN KEY (content)
        REFERENCES content (id)
);

CREATE TABLE display (
    id serial PRIMARY KEY,
    device integer NOT NULL,
    content integer NOT NULL,
    CONSTRAINT display_device_fkey
        FOREIGN KEY (device)
        REFERENCES device (id),
    CONSTRAINT display_content_fkey
        FOREIGN KEY (content)
        REFERENCES content (id)
);

CREATE TABLE list (
    id serial PRIMARY KEY,
    name varchar(50) NOT NULL
);

CREATE TABLE list_device (
    list integer NOT NULL,
    device integer NOT NULL,
    CONSTRAINT list_device_list_fkey
        FOREIGN KEY (list)
        REFERENCES list (id),
    CONSTRAINT list_device_device_fkey
        FOREIGN KEY (device)
        REFERENCES device (id)
);

-- +migrate Down
DROP TABLE IF EXISTS list_device;
DROP TABLE IF EXISTS list;
DROP TABLE IF EXISTS display;
DROP TABLE IF EXISTS header;
DROP TYPE IF EXISTS header_type;
DROP TABLE IF EXISTS content;
DROP TABLE IF EXISTS device;
DROP TABLE IF EXISTS location;
DROP TABLE IF EXISTS site;
