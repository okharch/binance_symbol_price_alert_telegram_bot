create or replace function test_alert_tables_exists() returns text language plpgsql as $$
BEGIN
    -- Connect to database
    PERFORM * FROM pg_catalog.pg_tables WHERE schemaname = 'public' AND tablename = 'alerts';
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Table "alerts" does not exist';
    END IF;
    PERFORM * FROM pg_catalog.pg_tables WHERE schemaname = 'public' AND tablename = 'symbol_prices';
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Table "symbol_prices" does not exist';
    END IF;
    PERFORM * FROM pg_catalog.pg_tables WHERE schemaname = 'public' AND tablename = 'alerts_archive';
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Table "alerts_archive" does not exist';
    END IF;
    return 'OK';
end;
$$;

create or replace function test_alert_populate_prices() returns int
    language plpgsql as $$
DECLARE
    rows_affected INTEGER;
begin
    -- Clear tables
    TRUNCATE TABLE alerts;
    TRUNCATE TABLE symbol_prices;
    TRUNCATE TABLE alerts_archive;

    -- Add test data
    INSERT INTO alerts (user_id, symbol, price, kind, created)
    VALUES
        (1, 'BTCUSD', 22000, 0, now() - interval '1 hour'),
        (2, 'ETHUSD', 1000, 1, now() - interval '1 hour');
    -- this triggers active_since, no alerts yet
    IF (SELECT count(*) FROM update_prices('[{"symbol":"BTCUSD","price":23000},{"symbol":"ETHUSD","price":900}]')) <> 0 then
        RAISE EXCEPTION 'should only set active_since';
    end if;

    if (select count(*) from alerts where active_since is not null) <> 2 then
        raise exception 'active_since has not been set!';
    end if;

    -- now trigger alerts indeed
    drop table if exists new_alerts_test;
    create temporary table new_alerts_test as
    SELECT * FROM update_prices('[{"symbol":"BTCUSD","price":21000},{"symbol":"ETHUSD","price":1020}]');
    GET DIAGNOSTICS rows_affected = ROW_COUNT;
    return rows_affected;
end
$$;

drop function if exists test_alerts();
CREATE OR REPLACE FUNCTION test_alerts()
    RETURNS text
AS $$
BEGIN
    perform test_alert_tables_exists();

    if test_alert_populate_prices() <> 2 then
        RAISE EXCEPTION 'should have triggered 2 alerts!';
    end if;

    -- Check results
    IF (SELECT COUNT(*) FROM alerts) <> 0 THEN
        RAISE EXCEPTION 'alerts table not empty after triggering alerts';
    END IF;

    IF (SELECT COUNT(*) FROM alerts_archive) <> 2 THEN
        RAISE EXCEPTION 'incorrect number of records in alerts_archive';
    END IF;

    IF (SELECT count(*) FROM new_alerts_test a,(
        (SELECT 1 user_id, 'BTCUSD' symbol , 21000 price, 0 kind UNION ALL SELECT 2, 'ETHUSD', 1020, 1)) b
        where a.user_id = b.user_id and a.symbol=b.symbol and a.price=b.price and a.kind=b.kind) <> 2 THEN
        RAISE EXCEPTION 'incorrect alerts triggered';
    END IF;
    return 'ok';
    --raise exception 'tests are ok: exception triggered in order to rollback test transaction!';
END;
$$ LANGUAGE plpgsql;

select * from test_alerts();