-- Table to store symbol prices
CREATE TABLE IF NOT EXISTS symbol_prices (
                                             symbol VARCHAR(10) NOT NULL,
                                             price FLOAT8 NOT NULL
);

-- Table to store active alerts
CREATE TABLE IF NOT EXISTS alerts (
                                      id SERIAL PRIMARY KEY,
                                      user_id INTEGER NOT NULL,
                                      symbol VARCHAR(10) NOT NULL,
                                      price FLOAT8 NOT NULL,
                                      kind INTEGER NOT NULL, -- 0 - trigger when price goes below, 1-when price goes up
                                      created timestamptz, -- when alert has been added
                                      active_since timestamptz null -- not null when it is active, null otherwise
);

-- Table to store triggered alerts
drop table if exists alerts_archive;
CREATE TABLE IF NOT EXISTS alerts_archive (
                                              id int PRIMARY KEY,
                                              user_id INTEGER NOT NULL,
                                              symbol VARCHAR(10) NOT NULL,
                                              price FLOAT8 NOT NULL,
                                              kind INTEGER NOT NULL,
                                              created timestamptz, -- when alert has been added
                                              active_since timestamptz null, -- not null when it is active, null otherwise
                                              last_price FLOAT8 NOT NULL,
                                              triggered_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     telegram_id INTEGER NOT NULL UNIQUE,
                                     first_name TEXT,
                                     last_name TEXT,
                                     username TEXT,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
