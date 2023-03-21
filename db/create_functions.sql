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
    INSERT INTO alerts_archive (id, user_id, symbol, price, kind, created, active_since, last_price)
    SELECT a.* FROM new_alerts a;

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
