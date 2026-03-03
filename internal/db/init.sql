DROP TABLE IF EXISTS visitedDomains;
CREATE TABLE visitedDomains (
  id         INT AUTO_INCREMENT NOT NULL,
  domain      VARCHAR(255) NOT NULL,
  time      TIMESTAMP,
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
  (1, 'user1', 'example1')
  (2, 'user2', 'example2')

/*TODO: Add constraint that, domain_key is UNIQUE, and overwrite attempt should change last credentials */