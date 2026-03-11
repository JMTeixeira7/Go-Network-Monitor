DROP TABLE IF EXISTS visitedDomains;
CREATE TABLE visitedDomains (
  id          INT AUTO_INCREMENT NOT NULL,
  domain      VARCHAR(255) NOT NULL,
  time        TIMESTAMP,
  PRIMARY KEY (`id`)
);

INSERT INTO visitedDomains
  (domain, time)
VALUES
  ('google.com', NULL),
  ('github.com', NULL),
  ('gmail.com', NULL),
  ('lirias2.kuleuven.be', NULL);

DROP TABLE IF EXISTS credentials;
CREATE TABLE credentials (
  id          INT AUTO_INCREMENT NOT NULL,
  domain_key  INT NOT NULL,
  username    VARCHAR(255) NOT NULL,
  fingerprint VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
  FOREIGN KEY (`domain_key`) REFERENCES visitedDomains(`id`)
)

INSERT INTO credentials
  (domain_key, username, fingerprint)
VALUES
  (1, 'user1', 'example1'),
  (2, 'user2', 'example2');


DROP TABLE IF EXISTS blockedDomains
CREATE TABLE blockedDomains (
  id            INT AUTO_INCREMENT NOT NULL,
  domain        VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
)

INSERT INTO blockedDomains
  (domain)
VALUES
  ('github.com'),
  ('facebook.com');


DROP TABLE IF EXISTS schedule
CREATE TABLE schedule (
  id                  INT AUTO_INCREMENT NOT NULL,
  blocked_domain_key  INT NOT NULL,
  start_time          TIMESTAMP,
  end_time            TIMESTAMP,
  weekday             INT,
  PRIMARY KEY (`id`)
  FOREIGN KEY (`blocked_domain_key`) REFERENCES blockedDomains(id)
)

INSERT INTO schedule
  (blocked_domain_key, start_time, end_time, weekday)
VALUES
  (1 ,NULL, NULL, TRUE, NULL)


/*TODO: Add constraint that, domain_key is UNIQUE, and overwrite attempt should change last credentials */