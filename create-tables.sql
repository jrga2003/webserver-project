DROP TABLE IF EXISTS companies;

CREATE TABLE companies (
id	int AUTO_INCREMENT NOT NULL,
name	VARCHAR(128) NOT NULL,
country	VARCHAR(128) NOT NULL,
website	VARCHAR(128) NOT NULL,
phone	VARCHAR(128) NOT NULL,
PRIMARY KEY (`id`)
);

INSERT INTO companies
  (name, country, website, phone)
VALUES
  ('IBM', 'United States', 'ibm.com', '1 4395 35497');

