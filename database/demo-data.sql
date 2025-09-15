-- Initial roles for account managers and admins
CREATE ROLE account_manager_1 NOBYPASSRLS;
CREATE ROLE account_manager_2 NOBYPASSRLS;
CREATE ROLE administrator;
CREATE ROLE customer NOBYPASSRLS;

-- Set up permissions
GRANT SELECT ON customers TO customer, account_manager_1, account_manager_2, administrator;
GRANT SELECT ON invoices TO customer, account_manager_1, account_manager_2, administrator;
GRANT SELECT ON invoice_positions TO customer, account_manager_1, account_manager_2, administrator;
GRANT SELECT ON account_managers TO account_manager_1, account_manager_2, administrator;
GRANT SELECT ON customer_account_managers TO account_manager_1, account_manager_2, administrator;

-- Initial account managers
INSERT INTO account_managers (
    account_manager_id, manager_name, db_role
) VALUES (
    gen_random_uuid(), 'Account Manager 1', 'account_manager_1'
);

INSERT INTO account_managers (
    account_manager_id, manager_name, db_role
) VALUES (
    gen_random_uuid(),'Account Manager 2','account_manager_2'
);

-- Initial customers
INSERT INTO customers (
    customer_id, customer_name
) VALUES (
    gen_random_uuid(), 'Pied Piper'
);

INSERT INTO customers (
    customer_id, customer_name
) VALUES (
    gen_random_uuid(), 'Hooli'
);

-- Initial Invoice
INSERT INTO invoices (
    invoice_id, customer_id, invoice_date
)
SELECT gen_random_uuid(), customer_id, now()
FROM customers
WHERE customer_name = 'Pied Piper';

INSERT INTO invoice_positions
SELECT gen_random_uuid(),
       invoice_id,
       customer_id,
       'The Box - Signature Edition',
       'A Hooli box, personally signed by our CEO Gavin Belson',
       10,
       999999.99
FROM invoices
LIMIT 1;

-- Assign customers to account managers
INSERT INTO customer_account_managers (
    customer_id, account_manager_id
) VALUES (
    (SELECT customer_id FROM customers WHERE customer_name = 'Pied Piper'),
    (SELECT account_manager_id FROM account_managers WHERE manager_name = 'Account Manager 1')
);

INSERT INTO customer_account_managers (
    customer_id, account_manager_id
) VALUES (
    (SELECT customer_id FROM customers WHERE customer_name = 'Hooli'),
    (SELECT account_manager_id FROM account_managers WHERE manager_name = 'Account Manager 1')
);

INSERT INTO customer_account_managers (
    customer_id, account_manager_id
) VALUES (
    (SELECT customer_id FROM customers WHERE customer_name = 'Hooli'),
    (SELECT account_manager_id FROM account_managers WHERE manager_name = 'Account Manager 2')
);
