
-- +goose Up
CREATE SCHEMA shiftr;

CREATE TABLE IF NOT EXISTS shiftr.postit_subscribers
(
    subscriber_id uuid UNIQUE NOT NULL,
    subscriber_email character varying(300) UNIQUE NOT NULL,
    subscriber_phone_number character varying(200) UNIQUE NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (subscriber_id)
);

CREATE TABLE IF NOT EXISTS shiftr.visitor_details
(
    visitor_id uuid UNIQUE NOT NULL,
    visitor_name character varying(300),
    visitor_email character varying(300) UNIQUE,
    visitor_phone_number character varying(300) UNIQUE,
    visitor_message text,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (visitor_id)
);

-- SQL in section 'Up' is executed when this migration is applied

-- +goose Down
DROP TABLE IF EXISTS shiftr.postit_subscribers;
DROP TABLE IF EXISTS shiftr.visitor_details;
DROP SCHEMA IF EXISTS shiftr CASCADE;
-- SQL section 'Down' is executed when this migration is rolled back