CREATE TYPE patient_gender AS ENUM('female', 'male');

CREATE TABLE IF NOT EXISTS patients (
  identity_number varchar(16) primary key,
  name varchar(30) not null,
  birth_date date not null,
  phone_number varchar(15) not null,
  gender patient_gender not null,
  card_image varchar(255) not null,
  created_at timestamp
);


CREATE INDEX IF NOT EXISTS patients_identity_number
	ON patients(identity_number);
CREATE INDEX IF NOT EXISTS patients_name
	ON patients(lower(name));
CREATE INDEX IF NOT EXISTS patients_phone_number
	ON patients(phone_number);
CREATE INDEX IF NOT EXISTS patients_created_at_desc
	ON patients(created_at desc);
CREATE INDEX IF NOT EXISTS patients_created_at_asc
	ON patients(created_at asc);