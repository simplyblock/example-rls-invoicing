-- Enable Row Level Security
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;
ALTER TABLE invoices ENABLE ROW LEVEL SECURITY;
ALTER TABLE invoice_positions ENABLE ROW LEVEL SECURITY;

-- Create customer policies (which can only select)
CREATE POLICY customers_tenant_isolation
    ON customers
    FOR SELECT
    TO customer
    USING (customer_id::TEXT = CURRENT_SETTING('app.current_customer_id', FALSE));

CREATE POLICY invoices_tenant_isolation
    ON invoices
    FOR SELECT
    TO customer
    USING (customer_id::TEXT = CURRENT_SETTING('app.current_customer_id', FALSE));

CREATE POLICY invoice_positions_tenant_isolation
    ON invoice_positions
    FOR SELECT
    TO customer
    USING (customer_id::TEXT = CURRENT_SETTING('app.current_customer_id', FALSE));

-- Create accountant policies
CREATE POLICY manager_can_see_assigned_customers
ON customers
TO account_manager_1, account_manager_2
USING (
    customer_id IN (
        SELECT cam.customer_id
        FROM customer_account_managers cam
            JOIN account_managers am
            ON cam.account_manager_id = am.account_manager_id
        WHERE am.db_role = current_user
    )
)
WITH CHECK (
    customer_id IN (
        SELECT cam.customer_id
        FROM customer_account_managers cam
            JOIN account_managers am
            ON cam.account_manager_id = am.account_manager_id
        WHERE am.db_role = current_user
    )
);

CREATE POLICY manager_can_see_assigned_invoices
ON invoices
TO account_manager_1, account_manager_2
USING (
    customer_id IN (
        SELECT cam.customer_id
        FROM customer_account_managers cam
            JOIN account_managers am
            ON cam.account_manager_id = am.account_manager_id
        WHERE am.db_role = current_user
    )
)
WITH CHECK (
    customer_id IN (
        SELECT cam.customer_id
        FROM customer_account_managers cam
            JOIN account_managers am
            ON cam.account_manager_id = am.account_manager_id
        WHERE am.db_role = current_user
    )
);

CREATE POLICY manager_can_see_assigned_invoice_positions
ON invoice_positions
TO account_manager_1, account_manager_2
USING (
    customer_id IN (
        SELECT cam.customer_id
        FROM customer_account_managers cam
            JOIN account_managers am
            ON cam.account_manager_id = am.account_manager_id
        WHERE am.db_role = current_user
    )
)
WITH CHECK (
    customer_id IN (
        SELECT cam.customer_id
        FROM customer_account_managers cam
            JOIN account_managers am
            ON cam.account_manager_id = am.account_manager_id
        WHERE am.db_role = current_user
    )
);

-- Create administrator policies
CREATE POLICY admin_can_see_all_customers
    ON customers
    TO administrator
    USING (TRUE)
    WITH CHECK (TRUE);

CREATE POLICY admin_can_see_all_invoices
    ON invoices
    TO administrator
    USING (TRUE)
    WITH CHECK (TRUE);

CREATE POLICY admin_can_see_all_invoice_positions
    ON invoice_positions
    TO administrator
    USING (TRUE)
    WITH CHECK (TRUE);
