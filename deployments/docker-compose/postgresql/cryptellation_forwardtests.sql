CREATE USER cryptellation_forwardtests;
ALTER USER cryptellation_forwardtests PASSWORD 'cryptellation_forwardtests';
ALTER USER cryptellation_forwardtests CREATEDB;

CREATE DATABASE cryptellation_forwardtests;
GRANT ALL PRIVILEGES ON DATABASE cryptellation_forwardtests TO cryptellation_forwardtests;
\c cryptellation_forwardtests postgres
GRANT ALL ON SCHEMA public TO cryptellation_forwardtests;