
CREATE TABLE public.airports
(
  rowid integer,
  airport character varying(250) COLLATE pg_catalog."default",
  city character varying(250) COLLATE pg_catalog."default",
  country character varying(250) COLLATE pg_catalog."default",
  iatacode character varying(50) COLLATE pg_catalog."default",
  icaocode character varying(50) COLLATE pg_catalog."default",
  latitude character varying(50) COLLATE pg_catalog."default",
  longitude character varying(50) COLLATE pg_catalog."default",
  altitude character varying(50) COLLATE pg_catalog."default",
  timezone character varying(50) COLLATE pg_catalog."default",
  dst character varying(50) COLLATE pg_catalog."default",
  tz character varying(50) COLLATE pg_catalog."default"
)
WITH (
  OIDS = FALSE
)
TABLESPACE pg_default;

COPY airports(rowID,Airport,City,Country,IATACode,ICAOCode,Latitude,Longitude,Altitude,TimeZone,DST,TZ)
FROM '/airports.csv' DELIMITER ',' CSV HEADER;