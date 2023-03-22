-- Stored procedure to trigger alerts based on updated prices
CREATE OR REPLACE FUNCTION trigger_alerts()
    RETURNS TABLE (
                      user_id INTEGER,
                      symbol VARCHAR(10),
                      price FLOAT8,
                      kind INTEGER
                  ) AS $$
DECLARE
    rows_affected INTEGER;
BEGIN
    -- Drop temporary table if it already exists
    DROP TABLE IF EXISTS new_alerts;

    -- Check whether price moves trigger any alerts
    CREATE TEMPORARY TABLE new_alerts AS
    SELECT a.*, p.price AS last_price FROM alerts a, symbol_prices p
    WHERE a.symbol = p.symbol AND a.active_since IS NOT NULL AND (
            (a.kind = 0 AND a.price >= p.price) OR
            (a.kind = 1 AND a.price <= p.price)
        );

    -- Get number of rows affected
    GET DIAGNOSTICS rows_affected = ROW_COUNT;

    -- If no alerts were triggered, return message and exit
    IF rows_affected = 0 THEN
        RAISE NOTICE 'No alerts triggered';
        DROP TABLE new_alerts;
        RETURN;
    END IF;

    -- Move triggered alerts to archive
    INSERT INTO alerts_archive (user_id, symbol, price, kind, created_at, active_since, last_price)
    SELECT a.user_id, a.symbol, a.price, a.kind, a.created_at, a.active_since, a.last_price
    FROM new_alerts a;

    -- Delete triggered alerts from alerts table
    DELETE FROM alerts WHERE id IN (SELECT a.id FROM new_alerts a);

    -- Return triggered alerts as result set
    RETURN QUERY SELECT a.user_id, a.symbol, a.last_price, a.kind FROM new_alerts a;
END;
$$ LANGUAGE plpgsql;


-- Stored procedure to update symbol prices and trigger alerts
CREATE OR REPLACE FUNCTION update_prices(prices JSON)
    RETURNS TABLE (
                      user_id INTEGER,
                      symbol VARCHAR(10),
                      price FLOAT8,
                      kind INTEGER
                  ) AS $$
BEGIN
    -- Truncate symbol_prices table
    TRUNCATE TABLE symbol_prices;

    -- Insert new prices from JSON argument
    INSERT INTO symbol_prices (symbol, price)
    SELECT p.symbol, p.price FROM json_populate_recordset(null::symbol_prices, prices) p;

    -- Set active_since when price goes up or down based on kind (0: -p.price > a.price 1: opposite)
    UPDATE alerts a
    SET active_since = NOW()
    FROM symbol_prices p
    WHERE a.symbol = p.symbol AND (
            (a.kind = 0 AND a.price < p.price)
            OR
            (a.kind = 1 AND a.price > p.price)
        );

    -- Trigger alerts based on updated prices and return triggered alerts
    RETURN QUERY SELECT a.* FROM trigger_alerts() a;
END;
$$ LANGUAGE plpgsql;

-- add_alert checks the current price for the given symbol,
-- calculates an_active_since based on whether the condition for setting it is met,
-- and then performs  INSERT into the alerts table.
CREATE OR REPLACE FUNCTION add_alert(
    IN a_user_id INTEGER,
    IN a_symbol VARCHAR(10),
    IN a_kind INTEGER,
    IN a_price FLOAT8
) RETURNS VOID AS $$
DECLARE
    an_active_since TIMESTAMP;
    current_price FLOAT8;
BEGIN
    SELECT price INTO current_price FROM symbol_prices WHERE symbol = a_symbol;

    IF (a_kind = 1 AND current_price < a_price) OR (a_kind = 0 AND current_price > a_price) THEN
        an_active_since := NOW();
    END IF;

    INSERT INTO alerts (user_id, symbol, kind, price, active_since)
    VALUES (a_user_id, a_symbol, a_kind, a_price, an_active_since);
END;
$$ LANGUAGE plpgsql;
