CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE idempotency (
  user_id UUID NOT NULL,
  -- user supplied idempotency key
  key     TEXT NOT NULL,

  -- full gRPC method path
  method_path TEXT  NOT NULL,
  -- request in protobuf format
  request     BYTEA NOT NULL,

  PRIMARY KEY (user_id, key, method_path, request),

  -- response in protobuf format
  response    BYTEA,
  status_code INT     NOT NULL,
  status_text TEXT,

  locked_at TIMESTAMPTZ NOT NULL,
  step      TEXT        NOT NULL
);

CREATE TABLE pets (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name          TEXT    NOT NULL,
  date_of_birth DATE    NOT NULL
);
