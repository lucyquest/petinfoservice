CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE idempotency (
  id UUID NOT NULL,
  function TEXT NOT NULL,
  step TEXT NOT NULL,
  last_used TIMESTAMP NOT NULL
);

CREATE TABLE pets (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  date_of_birth DATE NOT NULL
);
