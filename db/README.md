# Price Alert System in PostgreSQL

This is a simple implementation of price alert system using PostgreSQL. 
The system allows users to create price alerts for specific symbols, 
and receive notifications when the price of that symbol goes above or below a certain value.

The system is implemented using three tables: 
`symbol_prices`, `alerts`, and `alerts_archive`. 
`symbol_prices` stores the current prices for each symbol, 
`alerts` stores the active alerts for each user, and 
`alerts_archive` stores a record of triggered alerts.

The system is implemented using PostgreSQL stored procedures. 
The `update_prices` procedure is called whenever new prices are available, 
returns triggered alerts if the prices cross a threshold. 
The `trigger_alerts` procedure is responsible for checking 
which alerts have been triggered and moving them to the archive table.

The system also includes unit tests implemented as PostgreSQL functions. 
These tests verify that the database is functioning correctly, and that alerts are triggered when expected.

## Requirements

- PostgreSQL 9.5 or later
- Go 1.16 or later

## Installation

To install the system, execute the SQL scripts in the following order:

1. `create_tables.sql`: This script creates the necessary tables for the system.
2. `create_functions.sql`: This script creates the stored procedures for the system.
3. `unit_tests.sql`: This script runs the unit tests for the system.

After executing these scripts, the system will be ready to use.

## Usage

### Creating Alerts

To create an alert, insert a record into the `alerts` table with the following fields:

- `user_id`: The ID of the user who created the alert.
- `symbol`: The symbol to watch for price changes.
- `price`: The price at which to trigger the alert.
- `kind`: The kind of alert to trigger: 0 for when the price goes below the trigger value, and 1 for when the price goes above the trigger value.

### Updating Prices && Retrieving Triggered Alerts

To update the prices and retrieve the alerts, call the `update_prices` function with a JSON array of price updates. 
The JSON array should have the following format:

[
{
"symbol": "BTCUSD",
"price": 23000
},
{
"symbol": "ETHUSD",
"price": 900
}
]


The `symbol` field should match the symbol used in the `alerts` table. 
The `price` field should be the current price for that symbol.

it returns a result set with the following fields:

- `user_id`: The ID of the user who created the alert.
- `symbol`: The symbol that triggered the alert.
- `price`: The current price of the symbol.
- `kind`: The kind of alert that was triggered.

chat bot should send messages about price alerts to specific users for given symbols. 

### Unit Tests

To run the unit tests, execute the `unit_tests.sql` script. This script will run several tests to verify that the database is functioning correctly.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
