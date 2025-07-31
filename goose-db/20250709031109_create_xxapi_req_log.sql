-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE xxapi_req_log (
    id SERIAL PRIMARY KEY,
    request_id TEXT,
    respbody TEXT,
    responseStatus BIGINT,
    reqDateTime TIMESTAMP,
    realip TEXT,
    forwardedip TEXT,
    method TEXT,
    path TEXT,
    host TEXT,
    remoteaddr TEXT,
    header TEXT,
    endpoint TEXT,
    respDateTime TIMESTAMP,
    reqbody TEXT,
    requesttime BIGINT,
    responsetime BIGINT
);
-- +goose Down
-- +goose StatementBegin
Drop if exists xxapi_req_log;
-- +goose StatementEnd
