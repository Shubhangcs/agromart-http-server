-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS businesses (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id UUID REFERENCES users (id) ON DELETE CASCADE,
        business_profile_image TEXT,
        business_name VARCHAR(50) NOT NULL,
        business_email VARCHAR(50) UNIQUE NOT NULL,
        business_phone VARCHAR(50) UNIQUE NOT NULL,
        address VARCHAR(500) NOT NULL,
        city VARCHAR(50) NOT NULL,
        state VARCHAR(50) NOT NULL,
        pincode VARCHAR(50) NOT NULL,
        business_type TEXT NOT NULL,
        is_business_verified BOOLEAN NOT NULL DEFAULT FALSE,
        is_business_trusted BOOLEAN NOT NULL DEFAULT FALSE,
        is_business_approved BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS business_socials (
        business_id UUID PRIMARY KEY REFERENCES businesses (id) ON DELETE CASCADE,
        linkedin TEXT,
        instagram TEXT,
        youtube TEXT,
        telegram TEXT,
        x TEXT,
        facebook TEXT,
        website TEXT,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS business_legals (
        business_id UUID PRIMARY KEY REFERENCES businesses (id) ON DELETE CASCADE,
        aadhaar VARCHAR(12),
        pan VARCHAR(10),
        export_import VARCHAR(10),
        msme VARCHAR(19),
        fassi VARCHAR(14),
        gst VARCHAR(15),
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS business_applications (
        business_id UUID REFERENCES businesses (id) ON DELETE CASCADE,
        status TEXT NOT NULL CHECK (
            status IN ('APPLIED', 'REVIWED', 'ACCEPTED', 'REJECTED')
        ),
        reject_reason VARCHAR(1000),
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS business_ratings (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        business_id UUID REFERENCES businesses (id) ON DELETE CASCADE,
        user_id UUID REFERENCES users (id) ON DELETE CASCADE,
        rating NUMERIC(1, 1) NOT NULL CHECK (
            rating > 0.0
            AND rating <= 5.0
        ),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        UNIQUE (business_id, user_id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- FIX: Drop the correct table names (were previously wrong placeholder names)
DROP TABLE IF EXISTS business_applications;

DROP TABLE IF EXISTS business_legals;

DROP TABLE IF EXISTS business_socials;

DROP TABLE IF EXISTS businesses;

-- +goose StatementEnd