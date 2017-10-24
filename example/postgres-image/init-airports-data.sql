CREATE TABLE airports
(
  rowid integer,
  airport character varying(250),
  city character varying(250),
  country character varying(250),
  iatacode character varying(50),
  icaocode character varying(50),
  latitude character varying(50),
  longitude character varying(50),
  altitude character varying(50),
  timezone character varying(50),
  dst character varying(50),
  tz character varying(50)
);

COPY airports(rowID,Airport,City,Country,IATACode,ICAOCode,Latitude,Longitude,Altitude,TimeZone,DST,TZ) FROM '/airports.csv' DELIMITER ',' CSV HEADER;