CREATE TABLE IF NOT EXISTS medical_records(
  unique_id varchar(255) primary key,
  identity_number varchar(16) not null,
  creator_id varchar(255),
  symptoms text not null,
  medication text not null,
  created_at timestamp
);

ALTER TABLE medical_records
    ADD CONSTRAINT fk_creator_id FOREIGN KEY (creator_id) REFERENCES users(user_id) ON DELETE CASCADE;
ALTER TABLE medical_records
    ADD CONSTRAINT fk_patient_id FOREIGN KEY (identity_number) REFERENCES patients(identity_number) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS medical_records_unique_id
    ON medical_records(unique_id);
CREATE INDEX IF NOT EXISTS medical_records_identity_number
    ON medical_records(identity_number);
CREATE INDEX IF NOT EXISTS medical_records_creator_id
    ON medical_records(lower(creator_id));
CREATE INDEX IF NOT EXISTS medical_records_created_at_desc
    ON medical_records(created_at desc);
CREATE INDEX IF NOT EXISTS medical_records_created_at_asc
    ON medical_records(created_at asc);