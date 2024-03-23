# Test meter file read and update

## Description
This is a sample mini application to save data from the meter file with specification specified [here](https://aemo.com.au/-/media/files/electricity/nem/retail_and_metering/market_settlement_and_transfer_solutions/2022/mdff-specification-nem12-nem13-v25.pdf?la=en) into the format sepecified under the [create-tables.sql](https://github.com/AlvinTCH/meter-reader/blob/main/create-tables.sql) file.

### Assumptions
Due to the lack of experience and knowledge in this field, below are some of the assumptions that I have made:

1. Assume that if multiple 200 blocks happen in the same file, consumption under the 300 blocks with the same timestamps are summed up if the 200 block have the same NMI
2. If a file having 200 blocks with the same NMI and with similar timestamps, the consumption will be replaced with the later file

## How to use this application
1. Install docker and docker compose
2. set up the .env file with the following key-value pair
```
POSTGRES_PASSWORD - Postgresql password
POSTGRES_USER - Postgres user name
POSTGRES_DB - Postgres database name

GIN_MODE - debug / release
```
3. Run on bash console:
```
docker compose build --no-cache
docker compose up
```

The application will be available for calling on port **8080** and the database is on port **5432**

You can call **localhost:8080/upload-csv** to upload the csv for each meter file. Sample file given under [scripts](https://github.com/AlvinTCH/meter-reader/tree/main/scripts)

## Testing
Due to time constraints and knowledge constraints, I have written the test cases using Python. It should be migrated to use GORM Sqlmock in the future. Here are the test cases it contains:

1. Normal upload
2. Uploading a bad file (i.e file does not start with a 100 block)
3. Test updating of the values (assumption number 2 above)

Here is now to run it:
1. Install python
2. Go to the [scripts](https://github.com/AlvinTCH/meter-reader/tree/main/scripts) folder
3. Run
```
python3 test.py
```
4. View the test results printed to the console