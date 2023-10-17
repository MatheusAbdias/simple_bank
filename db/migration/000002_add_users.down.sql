ALTER TABLE
    IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";

ALTER TABLE
    IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_onwer_fkey";

DROP TABLE IF EXISTS "users";