CREATE TABLE
  users (
    id character varying(60) NOT NULL,
    username character varying NOT NULL,
    email character varying NOT NULL,
    password character varying NOT NULL
  );

ALTER TABLE
  users
ADD
  CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE TABLE images (
  id character varying(60) NOT NULL,
  user_id character varying(60) NOT NULL,
  filename VARCHAR(255) NOT NULL,
  PRIMARY KEY (id)
);
