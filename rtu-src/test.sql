CREATE TABLE substation (
  id    INTEGER PRIMARY KEY AUTOINCREMENT,
  name  TEXT    UNIQUE    NOT NULL,
  desc  TEXT
);

CREATE TABLE voltage_level (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  substation_id  INTEGER               NOT NULL  
                    REFERENCES substation(id),
  name           TEXT                  NOT NULL,
  desc           TEXT
);

CREATE TABLE bay (
  id               INTEGER PRIMARY KEY AUTOINCREMENT,
  voltage_level_id INTEGER               NOT NULL  
                    REFERENCES voltage_level(id),
  name             TEXT                  NOT NULL,
  desc             TEXT
);

CREATE TABLE ied (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  name            TEXT    UNIQUE    NOT NULL,
  type            TEXT,
  manufacturer    TEXT,
  config_version  TEXT
);

CREATE TABLE access_point (
  id     INTEGER PRIMARY KEY AUTOINCREMENT,
  ied_id INTEGER               NOT NULL  
           REFERENCES ied(id),
  name   TEXT                  NOT NULL
);

CREATE TABLE server (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  access_point_id INTEGER               NOT NULL  
                   REFERENCES access_point(id)
);

CREATE TABLE ldevice (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  server_id       INTEGER               NOT NULL  
                   REFERENCES server(id),
  inst            TEXT                  NOT NULL
);

CREATE TABLE ln0 (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  ldevice_id   INTEGER               NOT NULL  
                 REFERENCES ldevice(id),
  ln_class     TEXT                  NOT NULL,
  inst         TEXT                  NOT NULL,
  ln_type      TEXT
);

CREATE TABLE dataset (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  ln0_id         INTEGER               NOT NULL  
                   REFERENCES ln0(id),
  name           TEXT                  NOT NULL,
  desc           TEXT
);

CREATE TABLE fcda (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  dataset_id     INTEGER               NOT NULL  
                   REFERENCES dataset(id),
  ld_inst        TEXT                  NOT NULL,
  prefix         TEXT,
  ln_class       TEXT                  NOT NULL,
  ln_inst        TEXT,
  do_name        TEXT                  NOT NULL,
  fc             TEXT                  NOT NULL
);

