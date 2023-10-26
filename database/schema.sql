CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE pets (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  date_of_birth TIMESTAMP NOT NULL
)
