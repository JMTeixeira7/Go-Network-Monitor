DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS schedule;
DROP TABLE IF EXISTS blockedDomains;
DROP TABLE IF EXISTS visitedDomains;

CREATE TABLE visitedDomains (
  id          INT AUTO_INCREMENT NOT NULL,
  domain      VARCHAR(255) NOT NULL,
  time        TIMESTAMP NULL,
  PRIMARY KEY (id),
  UNIQUE(domain)
);

INSERT INTO visitedDomains
  (domain, time)
VALUES
  ('google.com', NULL),
  ('github.com', NULL),
  ('gmail.com', NULL),
  ('lirias2.kuleuven.be', NULL);

CREATE TABLE credentials (
  id          INT AUTO_INCREMENT NOT NULL,
  domain_key  INT NOT NULL,
  username    VARCHAR(255) NOT NULL,
  fingerprint VARCHAR(255) NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (domain_key) REFERENCES visitedDomains(id),
  UNIQUE(domain_key)
);

INSERT INTO credentials
  (domain_key, username, fingerprint)
VALUES
  (1, 'user1', 'ae23c8f573d36c7a1d6cbc8b0bebe52e0666dcf14fa04407328fb944b8f7a82d'), /* password = 12345*/
  (2, 'user2', 'ae23c8f573d36c7a1d6cbc8b0bebe52e0666dcf14fa04407328fb944b8f7a82d');

CREATE TABLE blockedDomains (
  id      INT AUTO_INCREMENT NOT NULL,
  domain  VARCHAR(255) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE(domain)
);

INSERT INTO blockedDomains
  (domain)
VALUES
  ('github.com'),
  ('facebook.com');

CREATE TABLE schedule (
  id                  INT AUTO_INCREMENT NOT NULL,
  blocked_domain_key  INT NOT NULL,
  start_time          TIME NULL,
  end_time            TIME NULL,
  weekday             INT NULL,
  timezone            INT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (blocked_domain_key) REFERENCES blockedDomains(id)
);

INSERT INTO schedule
  (blocked_domain_key, start_time, end_time, weekday, timezone)
VALUES
  (1, NULL, NULL, NULL, NULL),
  (2, NULL, NULL, NULL, NULL);