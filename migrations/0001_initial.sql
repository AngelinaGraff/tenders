CREATE TABLE tenders (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    created_by INT REFERENCES employee(id),
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tender_versions (
    id SERIAL PRIMARY KEY,
    tender_id INT REFERENCES tender(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type VARCHAR(50),
    status VARCHAR(20) NOT NULL,
    organization_id INT,
    created_by UUID,
    version INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bids (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    tender_id INT REFERENCES tenders(id) ON DELETE CASCADE,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    created_by UUID REFERENCES employees(id) ON DELETE SET NULL,
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bid_versions (
    id SERIAL PRIMARY KEY,
    bid_id INT REFERENCES bids(id) ON DELETE CASCADE,
    version INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL,
    created_by UUID REFERENCES employees(id) ON DELETE SET NULL,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    tender_id INT REFERENCES tenders(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    bid_id INT REFERENCES bids(id) ON DELETE CASCADE,
    reviewer_id INT REFERENCES employees(id),
    author_username VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
