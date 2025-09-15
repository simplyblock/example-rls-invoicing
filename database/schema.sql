-- Basic invoicing table schemas
CREATE TABLE customers (
    customer_id   UUID NOT NULL
        PRIMARY KEY,
    customer_name TEXT NOT NULL
);

CREATE TABLE invoices (
    invoice_id   UUID NOT NULL
        PRIMARY KEY,
    customer_id  UUID NOT NULL
        REFERENCES customers,
    invoice_date TIMESTAMPTZ NOT NULL
);

CREATE INDEX invoices_customer_id_index
    ON invoices (customer_id);

CREATE INDEX invoices_invoice_date_index
    ON invoices (invoice_date DESC);

CREATE TABLE invoice_positions (
    invoice_position_id UUID NOT NULL,
    invoice_id          UUID NOT NULL
        REFERENCES invoices,
    customer_id         UUID
        REFERENCES customers,
    name                TEXT NOT NULL,
    description         TEXT,
    amount              INTEGER NOT NULL ,
    price               REAL NOT NULL,
    position            INTEGER NOT NULL,
    PRIMARY KEY (invoice_position_id, invoice_id)
);

-- Additional management table schemas
CREATE TABLE account_managers (
    account_manager_id UUID NOT NULL
        PRIMARY KEY,
    manager_name       TEXT NOT NULL,
    db_role            TEXT NOT NULL UNIQUE
);

CREATE TABLE customer_account_managers (
    customer_id        UUID NOT NULL
        REFERENCES customers(customer_id),
    account_manager_id UUID NOT NULL
        REFERENCES account_managers(account_manager_id),
    PRIMARY KEY (customer_id, account_manager_id)
);
