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
