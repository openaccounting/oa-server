CREATE INDEX account_orgId_index ON account (orgId);
CREATE INDEX split_accountId_index ON split (accountId);
CREATE INDEX split_transactionId_index ON split (transactionId);
CREATE INDEX split_date_index ON split (date);
CREATE INDEX split_updated_index ON split (updated);
CREATE INDEX budgetitem_orgId_index ON budgetitem (orgId);