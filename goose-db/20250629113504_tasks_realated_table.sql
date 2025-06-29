-- +goose Up
-- +goose StatementBegin
create table tasks (
   ID SERIAL PRIMARY KEY,
   Title varchar(200) not null,
   Description text,
   Status status_enum NOT NULL DEFAULT 'Pending',
   Duedate TIMESTAMP,                -- Optional datetime
   CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   UpdatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   CreatedBy varchar(100) not null,
   UpdatedBy varchar(100) not null
);

CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_duedate ON tasks(duedate);

create table users (
ID SERIAL primary key,
Username varchar(200) unique not null,
PasswordHash varchar(200) not null,
CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
UpdatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CreatedBy varchar(100) not null,
UpdatedBy varchar(100) not null
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';

-- +goose StatementEnd
