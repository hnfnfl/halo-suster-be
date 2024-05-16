CREATE TYPE user_roles AS ENUM('it', 'nurse');

CREATE TABLE IF NOT EXISTS users (
  user_id varchar(255) primary key,
  nip varchar(13) unique not null,
  name varchar(50) not null,
  password_hash BYTEA not null,
  role user_roles,
  card_image varchar(255),
  created_at timestamp DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS users_user_id
    ON users(user_id);
CREATE INDEX IF NOT EXISTS users_name
    ON users(lower(name));
CREATE INDEX IF NOT EXISTS users_nip
    ON users(nip);
CREATE INDEX IF NOT EXISTS users_role
    ON users(role);
CREATE INDEX IF NOT EXISTS users_created_at_desc
    ON users(created_at desc);
CREATE INDEX IF NOT EXISTS users_created_at_asc
    ON users(created_at asc);