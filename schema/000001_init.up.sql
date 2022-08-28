CREATE TYPE account_status_t AS ENUM ( 'disabled', 'pending', 'active'  );
CREATE TYPE device_status_t AS ENUM ( 'active', 'archived'  );
CREATE TYPE project_status_t AS ENUM ( 'disabled', 'active'  );
CREATE TYPE user_account_role_t AS ENUM ( 'manager', 'user'  );
CREATE TYPE user_account_status_t AS ENUM ( 'disabled', 'invited', 'active'  );


CREATE TABLE users ( 
	id                   serial not null unique,
	email                varchar(200) NOT NULL,
	first_name           varchar(200) DEFAULT ''::character varying NOT NULL,
	last_name            varchar(200) DEFAULT ''::character varying NOT NULL,
	timezone             varchar(128) DEFAULT 'Europe/Kiev'::character varying NOT NULL,
	password_hash        varchar(256) NOT NULL,
	password_salt        varchar(36) NOT NULL,
	password_reset       varchar(36) DEFAULT NULL::character varying,
	created_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	archived_at          timestamptz,
	CONSTRAINT users_pkey PRIMARY KEY ( id ),
	CONSTRAINT email UNIQUE ( email ) 
 );

CREATE TABLE accounts ( 
	id                   serial not null unique,
	signup_user_id       integer NOT NULL,
	billing_user_id      integer NOT NULL,
	u_token              varchar(32),
	name                 varchar(255) DEFAULT ''::character varying NOT NULL,
	address1             varchar(255) DEFAULT ''::character varying NOT NULL,
	address2             varchar(255) DEFAULT ''::character varying NOT NULL,
	city                 varchar(100) DEFAULT ''::character varying NOT NULL,
	region               varchar(255) DEFAULT ''::character varying NOT NULL,
	country              varchar(255) DEFAULT ''::character varying NOT NULL,
	zipcode              varchar(20) DEFAULT ''::character varying NOT NULL,
	status               account_status_t DEFAULT 'active'::account_status_t NOT NULL,
	timezone             varchar(128) DEFAULT 'Europe/Kiev'::character varying NOT NULL,
	created_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	archived_at          timestamptz,
	CONSTRAINT accounts_pkey PRIMARY KEY ( id ),
	CONSTRAINT accounts_billing_user_id_fkey FOREIGN KEY ( billing_user_id ) REFERENCES users( id ) ON DELETE SET NULL,
	CONSTRAINT accounts_signup_user_id_fkey FOREIGN KEY ( signup_user_id ) REFERENCES users( id ) ON DELETE SET NULL  
 );

CREATE TABLE account_preferences ( 
	account_id           integer NOT NULL,
	name                 varchar(200) DEFAULT ''::character varying NOT NULL,
	"value"              varchar(200) DEFAULT ''::character varying NOT NULL,
	created_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	archived_at          timestamptz,
	CONSTRAINT account_preferences_pkey UNIQUE ( account_id, name ),
	CONSTRAINT account_preferences_account_id_fkey FOREIGN KEY ( account_id ) REFERENCES accounts( id )   
 );

CREATE TABLE users_accounts ( 
	id                   serial not null unique,
	user_id              integer NOT NULL,
	account_id           integer NOT NULL,
	roles                _user_account_role_t DEFAULT '{user}'::user_account_role_t[] NOT NULL,
	status               user_account_status_t DEFAULT 'active'::user_account_status_t NOT NULL,
	created_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	archived_at          timestamptz,
	CONSTRAINT users_accounts_pkey PRIMARY KEY ( id ),
	CONSTRAINT user_account UNIQUE ( user_id, account_id ),
	CONSTRAINT users_accounts_account_id_fkey FOREIGN KEY ( account_id ) REFERENCES accounts( id ),
	CONSTRAINT users_accounts_user_id_fkey FOREIGN KEY ( user_id ) REFERENCES users( id )   
 );

CREATE TABLE checklists ( 
	id                   serial not null unique,
	account_id           integer NOT NULL,
	title                varchar(255) DEFAULT ''::character varying NOT NULL,
	description          varchar(255) DEFAULT ''::character varying NOT NULL,
	status               project_status_t DEFAULT 'active'::project_status_t NOT NULL,
	created_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	archived_at          timestamptz,
	CONSTRAINT projects_pkey PRIMARY KEY ( id ),
	CONSTRAINT projects_account_id_fkey FOREIGN KEY ( account_id ) REFERENCES accounts( id ) ON DELETE SET NULL  
 );

CREATE TABLE checklist_items ( 
	id                   serial not null unique,
	checklist_id         integer NOT NULL,
	title                varchar(255) DEFAULT ''::character varying NOT NULL,
	description          varchar(255) DEFAULT ''::character varying NOT NULL,
	done                 boolean DEFAULT false NOT NULL,
	created_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	archived_at          timestamptz,
	CONSTRAINT checklist_items_pkey PRIMARY KEY ( id ),
	CONSTRAINT checklist_items_checklist_id_fkey FOREIGN KEY ( checklist_id ) REFERENCES checklists( id ) ON DELETE SET NULL  
 );

CREATE TABLE aquahubs ( 
	id                   serial not null unique,
	account_id           integer NOT NULL,
	h_token              varchar(32) NOT NULL,
	title                varchar(255),
	description          varchar(255),
	status               device_status_t DEFAULT 'active'::device_status_t,
	created_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at           timestamptz,
	archived_at          timestamptz,
	CONSTRAINT aquahubs_id_key UNIQUE ( id ),
	CONSTRAINT aquahubs_account_id_fkey FOREIGN KEY ( account_id ) REFERENCES accounts( id ) ON DELETE SET NULL  
 );

CREATE TABLE devices ( 
	id                   serial not null unique,
	aquahub_id           integer NOT NULL,
	local_id             integer DEFAULT '-1'::integer NOT NULL,
	title                varchar(255),
	description          varchar(255),
	status               device_status_t DEFAULT 'active'::device_status_t,
	CONSTRAINT devices_id_key UNIQUE ( id ),
	CONSTRAINT devices_aquahub_id_fkey FOREIGN KEY ( aquahub_id ) REFERENCES aquahubs( id )   
 );

CREATE TABLE sensors ( 
	id                   serial not null unique,
	device_id            integer NOT NULL,
	local_id             integer DEFAULT '-1'::integer,
	title                varchar(255),
	description          varchar(255),
	for_engineer         boolean DEFAULT true,
	for_analytics        boolean DEFAULT false,
	CONSTRAINT sensors_id_key UNIQUE ( id ),
	CONSTRAINT sensors_device_id_fkey FOREIGN KEY ( device_id ) REFERENCES devices( id )   
 );

CREATE TABLE sensors_dataset ( 
	id                   serial not null unique,
	account_id           integer NOT NULL,
	aquahub_id           integer NOT NULL,
	device_id            integer NOT NULL,
	sensor_id            integer NOT NULL,
	created_at           timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	user_time            timestamptz,
	local_device_id      integer DEFAULT '-1'::integer NOT NULL,
	local_sensor_id      integer DEFAULT '-1'::integer NOT NULL,
	"value"              varchar(32) NOT NULL,
	CONSTRAINT sensors_dataset_id_key UNIQUE ( id ),
	CONSTRAINT fk_sensors_dataset_accounts FOREIGN KEY ( account_id ) REFERENCES accounts( id ),
	CONSTRAINT sensors_dataset_aquahub_id_fkey FOREIGN KEY ( aquahub_id ) REFERENCES aquahubs( id ),
	CONSTRAINT sensors_dataset_device_id_fkey FOREIGN KEY ( device_id ) REFERENCES devices( id ),
	CONSTRAINT sensors_dataset_sensor_id_fkey FOREIGN KEY ( sensor_id ) REFERENCES sensors( id )   
 );

CREATE TABLE properties ( 
	id                   serial not null unique,
	device_id            integer NOT NULL,
	sensor_id            integer NOT NULL,
	"type"               integer,
	"value"              varchar(255),
	value_json           json,
	CONSTRAINT properties_id_key UNIQUE ( id ),
	CONSTRAINT properties_device_id_fkey FOREIGN KEY ( device_id ) REFERENCES devices( id ),
	CONSTRAINT properties_sensor_id_fkey FOREIGN KEY ( sensor_id ) REFERENCES sensors( id )   
 );



CREATE TABLE countries ( 
	code                 char(2) NOT NULL,
	iso_alpha3           char(3),
	name                 varchar(50),
	capital              varchar(50),
	currency_code        char(3),
	currency_name        char(20),
	phone                varchar(20),
	postal_code_format   varchar(200),
	postal_code_regex    varchar(200),
	CONSTRAINT countries_pkey PRIMARY KEY ( code )
 );

CREATE TABLE country_timezones ( 
	country_code         char(2) NOT NULL,
	timezone_id          varchar(50) NOT NULL,
	CONSTRAINT country_timezones_pkey UNIQUE ( country_code, timezone_id ) 
 );

CREATE TABLE geonames ( 
	country_code         char(2),
	postal_code          varchar(60),
	place_name           varchar(200),
	state_name           varchar(200),
	state_code           varchar(10),
	county_name          varchar(200),
	county_code          varchar(10),
	community_name       varchar(200),
	community_code       varchar(10),
	latitude             double precision,
	longitude            double precision,
	accuracy             integer    
 );

CREATE INDEX idx_geonames_country_code ON geonames  ( country_code );
CREATE INDEX idx_geonames_postal_code ON geonames  ( postal_code );



INSERT INTO countries( code, iso_alpha3, name, capital, currency_code, currency_name, phone, postal_code_format, postal_code_regex ) VALUES ( 'TT', 'TTO', 'Trinidad and Tobago', 'Port of Spain', 'TTD', 'Dollar              ', '+1-868', '', '');
INSERT INTO countries( code, iso_alpha3, name, capital, currency_code, currency_name, phone, postal_code_format, postal_code_regex ) VALUES ( 'TV', 'TUV', 'Tuvalu', 'Funafuti', 'AUD', 'Dollar              ', '688', '', '');
INSERT INTO countries( code, iso_alpha3, name, capital, currency_code, currency_name, phone, postal_code_format, postal_code_regex ) VALUES ( 'TW', 'TWN', 'Taiwan', 'Taipei', 'TWD', 'Dollar              ', '886', '#####', '^(\\d{5})$');
INSERT INTO countries( code, iso_alpha3, name, capital, currency_code, currency_name, phone, postal_code_format, postal_code_regex ) VALUES ( 'TZ', 'TZA', 'Tanzania', 'Dodoma', 'TZS', 'Shilling            ', '255', '', '');
INSERT INTO countries( code, iso_alpha3, name, capital, currency_code, currency_name, phone, postal_code_format, postal_code_regex ) VALUES ( 'UA', 'UKR', 'Ukraine', 'Kyiv', 'UAH', 'Hryvnia             ', '380', '#####', '^(\\d{5})$');
INSERT INTO country_timezones( country_code, timezone_id ) VALUES ( 'UA', 'Europe/Zaporozhye');
INSERT INTO country_timezones( country_code, timezone_id ) VALUES ( 'UA', 'Europe/Uzhgorod');
INSERT INTO country_timezones( country_code, timezone_id ) VALUES ( 'UA', 'Europe/Kiev');
INSERT INTO country_timezones( country_code, timezone_id ) VALUES ( 'UA', 'Europe/Simferopol');
INSERT INTO users( id, email, first_name, last_name, timezone, password_hash, password_salt, password_reset, created_at, updated_at, archived_at ) VALUES ( 1, 'as@as.as', 'as', 'as', 'Europe/Kiev', '$2a$10$KKx.P3cXd4QOsyYUGz3XpOerq5IAIjPBZeGuZ1GBz6kNuYMXnOEty', 'e964ee0a-eed8-4230-9e86-805acbdc1942', null, '2022-07-11 06:19:52 PM', '2022-07-11 06:37:52 PM', null);
INSERT INTO users( id, email, first_name, last_name, timezone, password_hash, password_salt, password_reset, created_at, updated_at, archived_at ) VALUES ( 2, 'zx@zx.zx', 'zx', 'zx', 'Europe/Kiev', '$2a$10$qMXCxW7sNNZJNMhesIZph.8EGhaUA6Dul8VaRZyIk76AGj08REBRS', '63d918f7-8c4a-4648-a835-1c59daa92014', null, '2022-07-11 06:40:36 PM', '2022-07-30 05:28:57 PM', null);
INSERT INTO users( id, email, first_name, last_name, timezone, password_hash, password_salt, password_reset, created_at, updated_at, archived_at ) VALUES ( 29, 'n30b5@as.as', 'as', 'as', 'Europe/Kiev', '$2a$10$AnDTJVXMGTtLa4YUujz.V.yjda/d3frmgnSDKd3YWpsXHZGpn3PMO', '20960387-aa01-4cf4-9880-3a12696c39d1', null, '2022-08-26 11:30:20 PM', '2022-08-26 11:30:20 PM', null);
INSERT INTO users( id, email, first_name, last_name, timezone, password_hash, password_salt, password_reset, created_at, updated_at, archived_at ) VALUES ( 30, 'ae@ae.ae', 'Andy', 'Sokol', 'Europe/Kiev', '$2a$10$YGoKY9Wi1x/cs/IzuuU8fOnxvS56Mtkt1M.yVTQhKyyWtOXxcIqS.', 'c5bd6420-2478-4e07-9847-079d79cf5b12', null, '2022-08-27 03:43:12 PM', '2022-08-27 03:43:12 PM', null);

INSERT INTO accounts( id, signup_user_id, billing_user_id, name, address1, address2, city, region, country, zipcode, status, timezone, created_at, updated_at, archived_at, u_token ) VALUES ( 0, 1, 1, 'airsoft company template', 'Airsoft 123', '', 'Kyiv', 'Kyivska oblast''', 'UA', '04207', 'active', 'Europe/Kiev', '2022-07-11 06:19:52 PM', '2022-07-11 06:35:32 PM', null, '');
INSERT INTO accounts( id, signup_user_id, billing_user_id, name, address1, address2, city, region, country, zipcode, status, timezone, created_at, updated_at, archived_at, u_token ) VALUES ( 1, 1, 1, 'Airsoft Company', 'Airsoft 123', '', 'Kyiv', 'Kyivska oblast''', 'UA', '04207', 'active', 'Europe/Kiev', '2022-07-11 06:19:52 PM', '2022-07-11 06:35:32 PM', null, 'a39831d103eb4c0d');
INSERT INTO accounts( id, signup_user_id, billing_user_id, name, address1, address2, city, region, country, zipcode, status, timezone, created_at, updated_at, archived_at, u_token ) VALUES ( 2, 2, 2, 'zx', 'zx', 'zx', 'zx', 'zx', 'UA', '12345', 'active', 'America/Anchorage', '2022-07-11 06:40:36 PM', '2022-07-11 06:40:36 PM', null, 'w2u_anMDe7xtwr00');
INSERT INTO accounts( id, signup_user_id, billing_user_id, name, address1, address2, city, region, country, zipcode, status, timezone, created_at, updated_at, archived_at, u_token ) VALUES ( 24, 29, 1, '', '', '', '', '', '', '', 'active', 'Europe/Kiev', '2022-08-26 11:30:21 PM', '2022-08-26 11:30:21 PM', null, 'uGanMDe7HXGD_twr');
INSERT INTO accounts( id, signup_user_id, billing_user_id, name, address1, address2, city, region, country, zipcode, status, timezone, created_at, updated_at, archived_at, u_token ) VALUES ( 25, 30, 1, '', '', '', '', '', '', '', 'active', 'Europe/Kiev', '2022-08-27 03:43:13 PM', '2022-08-27 03:43:13 PM', null, 'MN6LwWPJE4Il586U');

INSERT INTO aquahubs( id, h_token, title, description, status, created_at, updated_at, archived_at, account_id ) VALUES ( 2, 'aqen104Ur2zNX1Ykwv4', 'AquaHub More', 'Морской', 'active', '2022-07-07 05:15:41 PM', '2022-07-07 06:39:58 PM', '2022-07-07 06:40:16 PM', 1);
INSERT INTO aquahubs( id, h_token, title, description, status, created_at, updated_at, archived_at, account_id ) VALUES ( 1, '2323R7vzNX1Yk0uT', 'AquaHub.template', 'template', 'active', '2022-07-07 05:15:41 PM', '2022-07-07 06:39:58 PM', '2022-07-07 06:40:16 PM', 0);

INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 33, 1, 'Title Checklist', 'active', '2022-08-26 01:37:07 PM', '2022-08-27 07:16:17 PM', null, 'Description Checklist');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 35, 1, 'Rocket Launch', 'active', '2022-08-27 07:25:33 PM', '2022-08-27 07:25:33 PM', null, 'Rocket Launch Description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 22, 2, 'zxzxzxzx', 'active', '2022-07-21 04:35:04 PM', '2022-07-21 04:35:04 PM', null, '');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 24, 2, 'xdxdxdx', 'active', '2022-07-21 04:44:40 PM', '2022-07-21 04:44:40 PM', null, '');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 25, 2, 'fdfdfdf', 'active', '2022-07-21 04:46:16 PM', '2022-07-21 04:46:16 PM', null, '');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 26, 2, 'fdfdfdf', 'active', '2022-07-21 04:46:36 PM', '2022-07-21 04:46:36 PM', null, '');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 1, 1, 'UP check list 1', 'active', '2022-07-21 02:53:08 PM', '2022-07-21 02:53:08 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 2, 1, 'UP check list 1', 'active', '2022-07-21 03:00:16 PM', '2022-07-21 03:00:16 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 3, 1, 'UP check list 1', 'active', '2022-07-21 03:01:59 PM', '2022-07-21 03:01:59 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 4, 1, 'UP check list 1', 'active', '2022-07-21 03:02:57 PM', '2022-07-21 03:02:57 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 5, 1, 'UP check list 1', 'active', '2022-07-21 03:21:09 PM', '2022-07-21 03:21:09 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 6, 1, 'UP check list 1', 'active', '2022-07-21 03:24:21 PM', '2022-07-21 03:24:21 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 7, 1, 'UP check list 1', 'active', '2022-07-21 03:30:31 PM', '2022-07-21 03:30:31 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 8, 1, 'UP check list 1', 'active', '2022-07-21 03:33:23 PM', '2022-07-21 03:33:23 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 9, 1, 'UP check list 1', 'active', '2022-07-21 03:36:32 PM', '2022-07-21 03:36:32 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 10, 1, 'UP check list 1', 'active', '2022-07-21 03:47:04 PM', '2022-07-21 03:47:04 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 11, 1, 'UP check list 1', 'active', '2022-07-21 03:59:00 PM', '2022-07-21 03:59:00 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 12, 1, 'UP check list 1', 'active', '2022-07-21 04:01:27 PM', '2022-07-21 04:01:27 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 13, 1, 'UP check list 1', 'active', '2022-07-21 04:04:46 PM', '2022-07-21 04:04:46 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 14, 1, 'UP check list 1', 'active', '2022-07-21 04:08:00 PM', '2022-07-21 04:08:00 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 15, 1, 'UP check list 1', 'active', '2022-07-21 04:08:08 PM', '2022-07-21 04:08:08 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 16, 1, 'UP check list 1', 'active', '2022-07-21 04:05:55 PM', '2022-07-21 04:05:55 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 17, 1, 'UP check list 1', 'active', '2022-07-21 04:08:29 PM', '2022-07-21 04:08:29 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 19, 1, 'UP check list 1', 'active', '2022-07-21 04:10:20 PM', '2022-07-21 04:10:20 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 20, 1, 'UP check list 1', 'active', '2022-07-21 04:10:45 PM', '2022-07-21 04:10:45 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 0, 1, 'UP check list 1', 'active', '2022-07-21 04:31:51 PM', '2022-07-21 04:31:51 PM', null, 'UP description');
INSERT INTO checklists( id, account_id, title, status, created_at, updated_at, archived_at, description ) VALUES ( 34, 1, 'UP check list 1', 'active', '2022-08-26 01:37:53 PM', '2022-08-26 01:37:53 PM', null, 'UP description');
INSERT INTO devices( id, aquahub_id, local_id, title, description, status ) VALUES ( 1, 1, 110, 'Device.template', 'template', 'active');
INSERT INTO devices( id, aquahub_id, local_id, title, description, status ) VALUES ( 3, 2, 1, 'AquaHub', 'Description of the AquaHub', 'active');
INSERT INTO devices( id, aquahub_id, local_id, title, description, status ) VALUES ( 4, 2, 2, 'RayStar', 'Description of the RayStar', 'active');
INSERT INTO devices( id, aquahub_id, local_id, title, description, status ) VALUES ( 5, 2, 3, 'Home_Weather', 'Description of the Home_Weather', 'active');
INSERT INTO devices( id, aquahub_id, local_id, title, description, status ) VALUES ( 6, 2, 4, 'PowerMeter', 'Description of the PowerMeter', 'active');
INSERT INTO devices( id, aquahub_id, local_id, title, description, status ) VALUES ( 7, 2, 40, 'Doser', 'Description of the Doser', 'active');
INSERT INTO devices( id, aquahub_id, local_id, title, description, status ) VALUES ( 8, 2, 110, 'AutoDoliv', 'Description of the AutoDoliv', 'active');
INSERT INTO devices( id, aquahub_id, local_id, title, description, status ) VALUES ( 9, 2, 100, 'Cooler', 'Description of the Cooler', 'active');
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 1, 1, 0, 'Sensor.template', 'template', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 11, 3, 3, 'Sensor 3', 'Description of the sensor 3', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 12, 3, 4, 'Sensor 4', 'Description of the sensor 4', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 13, 3, 5, 'Sensor 5', 'Description of the sensor 5', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 14, 3, 0, 'Sensor 0', 'Description of the sensor 0', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 15, 3, 1, 'Sensor 1', 'Description of the sensor 1', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 16, 3, 2, 'Sensor 2', 'Description of the sensor 2', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 17, 4, 0, 'Sensor 0', 'Description of the sensor 0', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 18, 4, 1, 'Sensor 1', 'Description of the sensor 1', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 19, 4, 2, 'Sensor 2', 'Description of the sensor 2', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 20, 4, 3, 'Sensor 3', 'Description of the sensor 3', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 21, 5, 0, 'Sensor 0', 'Description of the sensor 0', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 22, 5, 1, 'Sensor 1', 'Description of the sensor 1', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 23, 5, 2, 'Sensor 2', 'Description of the sensor 2', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 24, 6, 2, 'Sensor 2', 'Description of the sensor 2', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 25, 6, 3, 'Sensor 3', 'Description of the sensor 3', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 26, 6, 4, 'Sensor 4', 'Description of the sensor 4', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 27, 6, 5, 'Sensor 5', 'Description of the sensor 5', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 28, 6, 0, 'Sensor 0', 'Description of the sensor 0', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 29, 6, 1, 'Sensor 1', 'Description of the sensor 1', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 30, 7, 0, 'Sensor 0', 'Description of the sensor 0', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 31, 7, 1, 'Sensor 1', 'Description of the sensor 1', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 32, 7, 2, 'Sensor 2', 'Description of the sensor 2', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 33, 8, 0, 'Sensor 0', 'Description of the sensor 0', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 34, 9, 1, 'Sensor 1', 'Description of the sensor 1', true, false);
INSERT INTO sensors( id, device_id, local_id, title, description, for_engineer, for_analytics ) VALUES ( 35, 9, 4, 'Sensor 4', 'Description of the sensor 4', true, false);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7922, 2, 6, 24, '2022-05-30 03:00:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7923, 2, 6, 25, '2022-05-30 03:00:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7924, 2, 6, 26, '2022-05-30 03:00:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7925, 2, 6, 27, '2022-05-30 03:00:05 PM', null, 4, 5, '41.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7926, 2, 6, 28, '2022-05-30 03:00:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7927, 2, 6, 29, '2022-05-30 03:00:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7928, 2, 7, 32, '2022-05-30 03:00:05 PM', null, 40, 2, '1.00', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7929, 2, 3, 13, '2022-05-30 03:00:05 PM', null, 1, 5, '1020', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7930, 2, 3, 14, '2022-05-30 03:00:05 PM', null, 1, 0, '1.21', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7931, 2, 3, 15, '2022-05-30 03:00:05 PM', null, 1, 1, '8', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7932, 2, 3, 16, '2022-05-30 03:00:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7933, 2, 3, 11, '2022-05-30 03:00:05 PM', null, 1, 3, '17.00', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7934, 2, 3, 12, '2022-05-30 03:00:05 PM', null, 1, 4, '15.00', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7935, 2, 4, 19, '2022-05-30 03:00:05 PM', null, 2, 2, '671', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7936, 2, 4, 20, '2022-05-30 03:00:05 PM', null, 2, 3, '26.3000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7937, 2, 4, 17, '2022-05-30 03:00:05 PM', null, 2, 0, '17.00', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7938, 2, 4, 18, '2022-05-30 03:00:05 PM', null, 2, 1, '147.22', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7939, 2, 3, 11, '2022-05-30 03:01:05 PM', null, 1, 3, '17.01', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7940, 2, 3, 12, '2022-05-30 03:01:05 PM', null, 1, 4, '15.01', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7941, 2, 3, 13, '2022-05-30 03:01:05 PM', null, 1, 5, '1021', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7942, 2, 3, 14, '2022-05-30 03:01:05 PM', null, 1, 0, '1.22', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7943, 2, 3, 15, '2022-05-30 03:01:05 PM', null, 1, 1, '7', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7944, 2, 3, 16, '2022-05-30 03:01:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7945, 2, 4, 17, '2022-05-30 03:01:05 PM', null, 2, 0, '17.01', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7946, 2, 4, 18, '2022-05-30 03:01:05 PM', null, 2, 1, '146.96', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7947, 2, 4, 19, '2022-05-30 03:01:05 PM', null, 2, 2, '664', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7948, 2, 4, 20, '2022-05-30 03:01:05 PM', null, 2, 3, '26.2000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7949, 2, 6, 29, '2022-05-30 03:01:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7950, 2, 6, 24, '2022-05-30 03:01:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7951, 2, 6, 25, '2022-05-30 03:01:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7952, 2, 6, 26, '2022-05-30 03:01:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7953, 2, 6, 27, '2022-05-30 03:01:05 PM', null, 4, 5, '41.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7954, 2, 6, 28, '2022-05-30 03:01:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7955, 2, 7, 32, '2022-05-30 03:01:05 PM', null, 40, 2, '0.59', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7956, 2, 7, 32, '2022-05-30 03:02:05 PM', null, 40, 2, '0.58', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7957, 2, 3, 15, '2022-05-30 03:02:05 PM', null, 1, 1, '6', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7958, 2, 3, 16, '2022-05-30 03:02:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7959, 2, 3, 11, '2022-05-30 03:02:05 PM', null, 1, 3, '17.02', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7960, 2, 3, 12, '2022-05-30 03:02:05 PM', null, 1, 4, '15.02', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7961, 2, 3, 13, '2022-05-30 03:02:05 PM', null, 1, 5, '1022', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7962, 2, 3, 14, '2022-05-30 03:02:05 PM', null, 1, 0, '1.23', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7963, 2, 4, 17, '2022-05-30 03:02:05 PM', null, 2, 0, '17.02', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7964, 2, 4, 18, '2022-05-30 03:02:05 PM', null, 2, 1, '146.71', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7965, 2, 4, 19, '2022-05-30 03:02:05 PM', null, 2, 2, '657', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7966, 2, 4, 20, '2022-05-30 03:02:05 PM', null, 2, 3, '26.3000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7967, 2, 6, 24, '2022-05-30 03:02:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7968, 2, 6, 25, '2022-05-30 03:02:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7969, 2, 6, 26, '2022-05-30 03:02:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7970, 2, 6, 27, '2022-05-30 03:02:05 PM', null, 4, 5, '40.40', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7971, 2, 6, 28, '2022-05-30 03:02:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7972, 2, 6, 29, '2022-05-30 03:02:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7973, 2, 3, 14, '2022-05-30 03:03:05 PM', null, 1, 0, '1.24', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7974, 2, 3, 15, '2022-05-30 03:03:05 PM', null, 1, 1, '5', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7975, 2, 3, 16, '2022-05-30 03:03:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7976, 2, 3, 11, '2022-05-30 03:03:05 PM', null, 1, 3, '17.03', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7977, 2, 3, 12, '2022-05-30 03:03:05 PM', null, 1, 4, '15.03', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7978, 2, 3, 13, '2022-05-30 03:03:05 PM', null, 1, 5, '1023', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7979, 2, 4, 17, '2022-05-30 03:03:05 PM', null, 2, 0, '17.03', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7980, 2, 4, 18, '2022-05-30 03:03:05 PM', null, 2, 1, '146.46', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7981, 2, 4, 19, '2022-05-30 03:03:05 PM', null, 2, 2, '651', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7982, 2, 4, 20, '2022-05-30 03:03:05 PM', null, 2, 3, '26.3000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7983, 2, 6, 27, '2022-05-30 03:03:05 PM', null, 4, 5, '41.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7984, 2, 6, 28, '2022-05-30 03:03:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7985, 2, 6, 29, '2022-05-30 03:03:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7986, 2, 6, 24, '2022-05-30 03:03:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7987, 2, 6, 25, '2022-05-30 03:03:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7988, 2, 6, 26, '2022-05-30 03:03:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7989, 2, 7, 32, '2022-05-30 03:03:05 PM', null, 40, 2, '0.57', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7990, 2, 7, 32, '2022-05-30 03:04:05 PM', null, 40, 2, '0.56', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7991, 2, 3, 16, '2022-05-30 03:04:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7992, 2, 3, 11, '2022-05-30 03:04:05 PM', null, 1, 3, '17.04', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7993, 2, 3, 12, '2022-05-30 03:04:05 PM', null, 1, 4, '15.04', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7994, 2, 3, 13, '2022-05-30 03:04:05 PM', null, 1, 5, '1024', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7995, 2, 3, 14, '2022-05-30 03:04:05 PM', null, 1, 0, '1.25', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7996, 2, 3, 15, '2022-05-30 03:04:05 PM', null, 1, 1, '4', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7997, 2, 4, 17, '2022-05-30 03:04:05 PM', null, 2, 0, '17.04', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7998, 2, 4, 18, '2022-05-30 03:04:05 PM', null, 2, 1, '146.21', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 7999, 2, 4, 19, '2022-05-30 03:04:05 PM', null, 2, 2, '644', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8000, 2, 4, 20, '2022-05-30 03:04:05 PM', null, 2, 3, '26.3000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8001, 2, 6, 27, '2022-05-30 03:04:05 PM', null, 4, 5, '42.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8002, 2, 6, 28, '2022-05-30 03:04:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8003, 2, 6, 29, '2022-05-30 03:04:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8004, 2, 6, 24, '2022-05-30 03:04:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8005, 2, 6, 25, '2022-05-30 03:04:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8006, 2, 6, 26, '2022-05-30 03:04:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8007, 2, 3, 14, '2022-05-30 03:05:05 PM', null, 1, 0, '1.26', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8008, 2, 3, 15, '2022-05-30 03:05:05 PM', null, 1, 1, '3', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8009, 2, 3, 16, '2022-05-30 03:05:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8010, 2, 3, 11, '2022-05-30 03:05:05 PM', null, 1, 3, '17.05', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8011, 2, 3, 12, '2022-05-30 03:05:05 PM', null, 1, 4, '15.05', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8012, 2, 3, 13, '2022-05-30 03:05:05 PM', null, 1, 5, '1025', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8013, 2, 4, 19, '2022-05-30 03:05:05 PM', null, 2, 2, '637', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8014, 2, 4, 20, '2022-05-30 03:05:05 PM', null, 2, 3, '26.2000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8015, 2, 4, 17, '2022-05-30 03:05:05 PM', null, 2, 0, '17.05', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8016, 2, 4, 18, '2022-05-30 03:05:05 PM', null, 2, 1, '145.95', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8017, 2, 6, 28, '2022-05-30 03:05:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8018, 2, 6, 29, '2022-05-30 03:05:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8019, 2, 6, 24, '2022-05-30 03:05:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8020, 2, 6, 25, '2022-05-30 03:05:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8021, 2, 6, 26, '2022-05-30 03:05:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8022, 2, 6, 27, '2022-05-30 03:05:05 PM', null, 4, 5, '41.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8023, 2, 7, 32, '2022-05-30 03:05:05 PM', null, 40, 2, '0.55', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8024, 2, 3, 14, '2022-05-30 03:06:05 PM', null, 1, 0, '1.27', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8025, 2, 3, 15, '2022-05-30 03:06:05 PM', null, 1, 1, '2', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8026, 2, 3, 16, '2022-05-30 03:06:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8027, 2, 3, 11, '2022-05-30 03:06:05 PM', null, 1, 3, '17.06', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8028, 2, 3, 12, '2022-05-30 03:06:05 PM', null, 1, 4, '15.06', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8029, 2, 3, 13, '2022-05-30 03:06:05 PM', null, 1, 5, '1026', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8030, 2, 4, 17, '2022-05-30 03:06:05 PM', null, 2, 0, '17.06', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8031, 2, 4, 18, '2022-05-30 03:06:05 PM', null, 2, 1, '145.69', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8032, 2, 4, 19, '2022-05-30 03:06:05 PM', null, 2, 2, '630', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8033, 2, 4, 20, '2022-05-30 03:06:05 PM', null, 2, 3, '26.3000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8034, 2, 6, 28, '2022-05-30 03:06:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8035, 2, 6, 29, '2022-05-30 03:06:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8036, 2, 6, 24, '2022-05-30 03:06:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8037, 2, 6, 25, '2022-05-30 03:06:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8038, 2, 6, 26, '2022-05-30 03:06:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8039, 2, 6, 27, '2022-05-30 03:06:05 PM', null, 4, 5, '42.00', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8040, 2, 7, 32, '2022-05-30 03:06:05 PM', null, 40, 2, '0.54', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8041, 2, 3, 11, '2022-05-30 03:07:05 PM', null, 1, 3, '17.07', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8042, 2, 3, 12, '2022-05-30 03:07:05 PM', null, 1, 4, '15.07', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8043, 2, 3, 13, '2022-05-30 03:07:05 PM', null, 1, 5, '1027', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8044, 2, 3, 14, '2022-05-30 03:07:05 PM', null, 1, 0, '1.28', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8045, 2, 3, 15, '2022-05-30 03:07:05 PM', null, 1, 1, '1', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8046, 2, 3, 16, '2022-05-30 03:07:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8047, 2, 4, 20, '2022-05-30 03:07:05 PM', null, 2, 3, '26.1000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8048, 2, 4, 17, '2022-05-30 03:07:05 PM', null, 2, 0, '17.07', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8049, 2, 4, 18, '2022-05-30 03:07:05 PM', null, 2, 1, '145.44', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8050, 2, 4, 19, '2022-05-30 03:07:05 PM', null, 2, 2, '623', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8051, 2, 6, 28, '2022-05-30 03:07:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8052, 2, 6, 29, '2022-05-30 03:07:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8053, 2, 6, 24, '2022-05-30 03:07:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8054, 2, 6, 25, '2022-05-30 03:07:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8055, 2, 6, 26, '2022-05-30 03:07:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8056, 2, 6, 27, '2022-05-30 03:07:05 PM', null, 4, 5, '40.40', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8057, 2, 7, 32, '2022-05-30 03:07:05 PM', null, 40, 2, '0.53', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8058, 2, 4, 18, '2022-05-30 03:08:13 PM', null, 2, 1, '145.20', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8059, 2, 4, 19, '2022-05-30 03:08:13 PM', null, 2, 2, '617', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8060, 2, 4, 20, '2022-05-30 03:08:13 PM', null, 2, 3, '26.2000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8061, 2, 4, 17, '2022-05-30 03:08:13 PM', null, 2, 0, '17.08', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8062, 2, 6, 27, '2022-05-30 03:08:13 PM', null, 4, 5, '42.00', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8063, 2, 6, 28, '2022-05-30 03:08:13 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8064, 2, 6, 29, '2022-05-30 03:08:13 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8065, 2, 6, 24, '2022-05-30 03:08:13 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8066, 2, 6, 25, '2022-05-30 03:08:13 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8067, 2, 6, 26, '2022-05-30 03:08:13 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8068, 2, 7, 32, '2022-05-30 03:08:13 PM', null, 40, 2, '0.52', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8069, 2, 3, 11, '2022-05-30 03:08:13 PM', null, 1, 3, '17.08', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8070, 2, 3, 12, '2022-05-30 03:08:13 PM', null, 1, 4, '15.08', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8071, 2, 3, 13, '2022-05-30 03:08:13 PM', null, 1, 5, '1028', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8072, 2, 3, 14, '2022-05-30 03:08:13 PM', null, 1, 0, '1.29', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8073, 2, 3, 15, '2022-05-30 03:08:13 PM', null, 1, 1, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8074, 2, 3, 16, '2022-05-30 03:08:13 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8075, 2, 3, 12, '2022-05-30 03:09:05 PM', null, 1, 4, '15.09', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8076, 2, 3, 13, '2022-05-30 03:09:05 PM', null, 1, 5, '1029', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8077, 2, 3, 14, '2022-05-30 03:09:05 PM', null, 1, 0, '1.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8078, 2, 3, 15, '2022-05-30 03:09:05 PM', null, 1, 1, '29', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8079, 2, 3, 16, '2022-05-30 03:09:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8080, 2, 3, 11, '2022-05-30 03:09:05 PM', null, 1, 3, '17.09', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8081, 2, 4, 17, '2022-05-30 03:09:05 PM', null, 2, 0, '17.09', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8082, 2, 4, 18, '2022-05-30 03:09:05 PM', null, 2, 1, '144.93', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8083, 2, 4, 19, '2022-05-30 03:09:05 PM', null, 2, 2, '610', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8084, 2, 4, 20, '2022-05-30 03:09:05 PM', null, 2, 3, '26.3000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8085, 2, 6, 25, '2022-05-30 03:09:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8086, 2, 6, 26, '2022-05-30 03:09:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8087, 2, 6, 27, '2022-05-30 03:09:05 PM', null, 4, 5, '40.00', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8088, 2, 6, 28, '2022-05-30 03:09:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8089, 2, 6, 29, '2022-05-30 03:09:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8090, 2, 6, 24, '2022-05-30 03:09:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8091, 2, 7, 32, '2022-05-30 03:09:05 PM', null, 40, 2, '0.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8092, 2, 7, 32, '2022-05-30 03:10:05 PM', null, 40, 2, '0.50', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8093, 2, 3, 13, '2022-05-30 03:10:05 PM', null, 1, 5, '1030', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8094, 2, 3, 14, '2022-05-30 03:10:05 PM', null, 1, 0, '1.31', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8095, 2, 3, 15, '2022-05-30 03:10:05 PM', null, 1, 1, '28', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8096, 2, 3, 16, '2022-05-30 03:10:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8097, 2, 3, 11, '2022-05-30 03:10:05 PM', null, 1, 3, '17.10', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8098, 2, 3, 12, '2022-05-30 03:10:05 PM', null, 1, 4, '15.10', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8099, 2, 4, 17, '2022-05-30 03:10:05 PM', null, 2, 0, '17.10', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8100, 2, 4, 18, '2022-05-30 03:10:05 PM', null, 2, 1, '144.68', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8101, 2, 4, 19, '2022-05-30 03:10:05 PM', null, 2, 2, '603', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8102, 2, 4, 20, '2022-05-30 03:10:05 PM', null, 2, 3, '26.3000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8103, 2, 6, 25, '2022-05-30 03:10:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8104, 2, 6, 26, '2022-05-30 03:10:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8105, 2, 6, 27, '2022-05-30 03:10:05 PM', null, 4, 5, '40.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8106, 2, 6, 28, '2022-05-30 03:10:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8107, 2, 6, 29, '2022-05-30 03:10:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8108, 2, 6, 24, '2022-05-30 03:10:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8109, 2, 3, 14, '2022-05-30 03:11:05 PM', null, 1, 0, '1.32', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8110, 2, 3, 15, '2022-05-30 03:11:05 PM', null, 1, 1, '27', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8111, 2, 3, 16, '2022-05-30 03:11:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8112, 2, 3, 11, '2022-05-30 03:11:05 PM', null, 1, 3, '17.11', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8113, 2, 3, 12, '2022-05-30 03:11:05 PM', null, 1, 4, '15.11', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8114, 2, 3, 13, '2022-05-30 03:11:05 PM', null, 1, 5, '1031', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8115, 2, 4, 18, '2022-05-30 03:11:05 PM', null, 2, 1, '144.44', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8116, 2, 4, 19, '2022-05-30 03:11:05 PM', null, 2, 2, '596', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8117, 2, 4, 20, '2022-05-30 03:11:05 PM', null, 2, 3, '26.3000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8118, 2, 4, 17, '2022-05-30 03:11:05 PM', null, 2, 0, '17.11', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8119, 2, 6, 25, '2022-05-30 03:11:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8120, 2, 6, 26, '2022-05-30 03:11:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8121, 2, 6, 27, '2022-05-30 03:11:05 PM', null, 4, 5, '41.59', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8122, 2, 6, 28, '2022-05-30 03:11:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8123, 2, 6, 29, '2022-05-30 03:11:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8124, 2, 6, 24, '2022-05-30 03:11:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8125, 2, 7, 32, '2022-05-30 03:11:05 PM', null, 40, 2, '0.49', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8126, 2, 6, 28, '2022-05-30 03:12:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8127, 2, 6, 29, '2022-05-30 03:12:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8128, 2, 6, 24, '2022-05-30 03:12:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8129, 2, 6, 25, '2022-05-30 03:12:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8130, 2, 6, 26, '2022-05-30 03:12:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8131, 2, 6, 27, '2022-05-30 03:12:05 PM', null, 4, 5, '40.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8132, 2, 7, 32, '2022-05-30 03:12:05 PM', null, 40, 2, '0.48', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8133, 2, 3, 12, '2022-05-30 03:12:05 PM', null, 1, 4, '15.12', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8134, 2, 3, 13, '2022-05-30 03:12:05 PM', null, 1, 5, '1032', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8135, 2, 3, 14, '2022-05-30 03:12:05 PM', null, 1, 0, '1.33', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8136, 2, 3, 15, '2022-05-30 03:12:05 PM', null, 1, 1, '26', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8137, 2, 3, 16, '2022-05-30 03:12:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8138, 2, 3, 11, '2022-05-30 03:12:05 PM', null, 1, 3, '17.12', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8139, 2, 4, 17, '2022-05-30 03:12:05 PM', null, 2, 0, '17.12', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8140, 2, 4, 18, '2022-05-30 03:12:05 PM', null, 2, 1, '144.18', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8141, 2, 4, 19, '2022-05-30 03:12:05 PM', null, 2, 2, '589', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8142, 2, 4, 20, '2022-05-30 03:12:05 PM', null, 2, 3, '26.3000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8143, 2, 3, 16, '2022-05-30 03:13:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8144, 2, 3, 11, '2022-05-30 03:13:05 PM', null, 1, 3, '17.13', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8145, 2, 3, 12, '2022-05-30 03:13:05 PM', null, 1, 4, '15.13', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8146, 2, 3, 13, '2022-05-30 03:13:05 PM', null, 1, 5, '1033', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8147, 2, 3, 14, '2022-05-30 03:13:05 PM', null, 1, 0, '1.34', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8148, 2, 3, 15, '2022-05-30 03:13:05 PM', null, 1, 1, '25', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8149, 2, 4, 17, '2022-05-30 03:13:05 PM', null, 2, 0, '17.13', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8150, 2, 4, 18, '2022-05-30 03:13:05 PM', null, 2, 1, '143.92', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8151, 2, 4, 19, '2022-05-30 03:13:05 PM', null, 2, 2, '583', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8152, 2, 4, 20, '2022-05-30 03:13:05 PM', null, 2, 3, '26.3000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8153, 2, 6, 28, '2022-05-30 03:13:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8154, 2, 6, 29, '2022-05-30 03:13:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8155, 2, 6, 24, '2022-05-30 03:13:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8156, 2, 6, 25, '2022-05-30 03:13:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8157, 2, 6, 26, '2022-05-30 03:13:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8158, 2, 6, 27, '2022-05-30 03:13:05 PM', null, 4, 5, '40.40', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8159, 2, 7, 32, '2022-05-30 03:13:05 PM', null, 40, 2, '0.47', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8160, 2, 3, 16, '2022-05-30 03:14:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8161, 2, 3, 11, '2022-05-30 03:14:05 PM', null, 1, 3, '17.14', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8162, 2, 3, 12, '2022-05-30 03:14:05 PM', null, 1, 4, '15.14', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8163, 2, 3, 13, '2022-05-30 03:14:05 PM', null, 1, 5, '1034', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8164, 2, 3, 14, '2022-05-30 03:14:05 PM', null, 1, 0, '1.35', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8165, 2, 3, 15, '2022-05-30 03:14:05 PM', null, 1, 1, '24', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8166, 2, 4, 17, '2022-05-30 03:14:05 PM', null, 2, 0, '17.14', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8167, 2, 4, 18, '2022-05-30 03:14:05 PM', null, 2, 1, '143.67', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8168, 2, 4, 19, '2022-05-30 03:14:05 PM', null, 2, 2, '576', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8169, 2, 4, 20, '2022-05-30 03:14:05 PM', null, 2, 3, '26.3000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8170, 2, 6, 26, '2022-05-30 03:14:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8171, 2, 6, 27, '2022-05-30 03:14:05 PM', null, 4, 5, '40.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8172, 2, 6, 28, '2022-05-30 03:14:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8173, 2, 6, 29, '2022-05-30 03:14:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8174, 2, 6, 24, '2022-05-30 03:14:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8175, 2, 6, 25, '2022-05-30 03:14:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8176, 2, 7, 32, '2022-05-30 03:14:05 PM', null, 40, 2, '0.46', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8177, 2, 3, 14, '2022-05-30 03:15:05 PM', null, 1, 0, '1.36', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8178, 2, 3, 15, '2022-05-30 03:15:05 PM', null, 1, 1, '23', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8179, 2, 3, 16, '2022-05-30 03:15:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8180, 2, 3, 11, '2022-05-30 03:15:05 PM', null, 1, 3, '17.15', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8181, 2, 3, 12, '2022-05-30 03:15:05 PM', null, 1, 4, '15.15', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8182, 2, 3, 13, '2022-05-30 03:15:05 PM', null, 1, 5, '1035', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8183, 2, 4, 17, '2022-05-30 03:15:05 PM', null, 2, 0, '17.15', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8184, 2, 4, 18, '2022-05-30 03:15:05 PM', null, 2, 1, '143.42', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8185, 2, 4, 19, '2022-05-30 03:15:05 PM', null, 2, 2, '569', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8186, 2, 4, 20, '2022-05-30 03:15:05 PM', null, 2, 3, '26.3000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8187, 2, 6, 28, '2022-05-30 03:15:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8188, 2, 6, 29, '2022-05-30 03:15:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8189, 2, 6, 24, '2022-05-30 03:15:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8190, 2, 6, 25, '2022-05-30 03:15:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8191, 2, 6, 26, '2022-05-30 03:15:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8192, 2, 6, 27, '2022-05-30 03:15:05 PM', null, 4, 5, '42.00', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8193, 2, 7, 32, '2022-05-30 03:15:05 PM', null, 40, 2, '0.45', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8194, 2, 3, 13, '2022-05-30 03:16:05 PM', null, 1, 5, '1036', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8195, 2, 3, 14, '2022-05-30 03:16:05 PM', null, 1, 0, '1.37', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8196, 2, 3, 15, '2022-05-30 03:16:05 PM', null, 1, 1, '22', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8197, 2, 3, 16, '2022-05-30 03:16:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8198, 2, 3, 11, '2022-05-30 03:16:05 PM', null, 1, 3, '17.16', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8199, 2, 3, 12, '2022-05-30 03:16:05 PM', null, 1, 4, '15.16', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8200, 2, 4, 20, '2022-05-30 03:16:05 PM', null, 2, 3, '26.2000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8201, 2, 4, 17, '2022-05-30 03:16:05 PM', null, 2, 0, '17.16', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8202, 2, 4, 18, '2022-05-30 03:16:05 PM', null, 2, 1, '143.17', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8203, 2, 4, 19, '2022-05-30 03:16:05 PM', null, 2, 2, '562', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8204, 2, 6, 28, '2022-05-30 03:16:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8205, 2, 6, 29, '2022-05-30 03:16:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8206, 2, 6, 24, '2022-05-30 03:16:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8207, 2, 6, 25, '2022-05-30 03:16:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8208, 2, 6, 26, '2022-05-30 03:16:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8209, 2, 6, 27, '2022-05-30 03:16:05 PM', null, 4, 5, '40.00', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8210, 2, 7, 32, '2022-05-30 03:16:05 PM', null, 40, 2, '0.44', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8211, 2, 7, 32, '2022-05-30 03:17:05 PM', null, 40, 2, '0.43', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8212, 2, 3, 13, '2022-05-30 03:17:05 PM', null, 1, 5, '1037', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8213, 2, 3, 14, '2022-05-30 03:17:05 PM', null, 1, 0, '1.38', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8214, 2, 3, 15, '2022-05-30 03:17:05 PM', null, 1, 1, '21', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8215, 2, 3, 16, '2022-05-30 03:17:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8216, 2, 3, 11, '2022-05-30 03:17:05 PM', null, 1, 3, '17.17', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8217, 2, 3, 12, '2022-05-30 03:17:05 PM', null, 1, 4, '15.17', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8218, 2, 4, 17, '2022-05-30 03:17:05 PM', null, 2, 0, '17.17', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8219, 2, 4, 18, '2022-05-30 03:17:05 PM', null, 2, 1, '142.91', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8220, 2, 4, 19, '2022-05-30 03:17:05 PM', null, 2, 2, '555', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8221, 2, 4, 20, '2022-05-30 03:17:05 PM', null, 2, 3, '26.1000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8222, 2, 6, 25, '2022-05-30 03:17:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8223, 2, 6, 26, '2022-05-30 03:17:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8224, 2, 6, 27, '2022-05-30 03:17:05 PM', null, 4, 5, '42.00', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8225, 2, 6, 28, '2022-05-30 03:17:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8226, 2, 6, 29, '2022-05-30 03:17:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8227, 2, 6, 24, '2022-05-30 03:17:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8228, 2, 3, 13, '2022-05-30 03:18:05 PM', null, 1, 5, '1038', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8229, 2, 3, 14, '2022-05-30 03:18:05 PM', null, 1, 0, '1.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8230, 2, 3, 15, '2022-05-30 03:18:05 PM', null, 1, 1, '20', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8231, 2, 3, 16, '2022-05-30 03:18:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8232, 2, 3, 11, '2022-05-30 03:18:05 PM', null, 1, 3, '17.18', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8233, 2, 3, 12, '2022-05-30 03:18:05 PM', null, 1, 4, '15.18', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8234, 2, 4, 17, '2022-05-30 03:18:05 PM', null, 2, 0, '17.18', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8235, 2, 4, 18, '2022-05-30 03:18:05 PM', null, 2, 1, '142.66', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8236, 2, 4, 19, '2022-05-30 03:18:05 PM', null, 2, 2, '549', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8237, 2, 4, 20, '2022-05-30 03:18:05 PM', null, 2, 3, '26.1000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8238, 2, 6, 28, '2022-05-30 03:18:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8239, 2, 6, 29, '2022-05-30 03:18:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8240, 2, 6, 24, '2022-05-30 03:18:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8241, 2, 6, 25, '2022-05-30 03:18:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8242, 2, 6, 26, '2022-05-30 03:18:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8243, 2, 6, 27, '2022-05-30 03:18:05 PM', null, 4, 5, '41.59', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8244, 2, 7, 32, '2022-05-30 03:18:05 PM', null, 40, 2, '0.42', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8245, 2, 3, 14, '2022-05-30 03:19:05 PM', null, 1, 0, '1.40', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8246, 2, 3, 15, '2022-05-30 03:19:05 PM', null, 1, 1, '19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8247, 2, 3, 16, '2022-05-30 03:19:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8248, 2, 3, 11, '2022-05-30 03:19:05 PM', null, 1, 3, '17.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8249, 2, 3, 12, '2022-05-30 03:19:05 PM', null, 1, 4, '15.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8250, 2, 3, 13, '2022-05-30 03:19:05 PM', null, 1, 5, '1039', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8251, 2, 4, 17, '2022-05-30 03:19:05 PM', null, 2, 0, '17.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8252, 2, 4, 18, '2022-05-30 03:19:05 PM', null, 2, 1, '142.41', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8253, 2, 4, 19, '2022-05-30 03:19:05 PM', null, 2, 2, '542', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8254, 2, 4, 20, '2022-05-30 03:19:05 PM', null, 2, 3, '26.1000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8255, 2, 6, 29, '2022-05-30 03:19:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8256, 2, 6, 24, '2022-05-30 03:19:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8257, 2, 6, 25, '2022-05-30 03:19:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8258, 2, 6, 26, '2022-05-30 03:19:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8259, 2, 6, 27, '2022-05-30 03:19:05 PM', null, 4, 5, '41.59', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8260, 2, 6, 28, '2022-05-30 03:19:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8261, 2, 7, 32, '2022-05-30 03:19:05 PM', null, 40, 2, '0.41', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8262, 2, 3, 14, '2022-05-30 03:20:05 PM', null, 1, 0, '1.41', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8263, 2, 3, 15, '2022-05-30 03:20:05 PM', null, 1, 1, '18', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8264, 2, 3, 16, '2022-05-30 03:20:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8265, 2, 3, 11, '2022-05-30 03:20:05 PM', null, 1, 3, '17.20', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8266, 2, 3, 12, '2022-05-30 03:20:05 PM', null, 1, 4, '15.20', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8267, 2, 3, 13, '2022-05-30 03:20:05 PM', null, 1, 5, '1040', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8268, 2, 4, 19, '2022-05-30 03:20:05 PM', null, 2, 2, '535', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8269, 2, 4, 20, '2022-05-30 03:20:05 PM', null, 2, 3, '26.2000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8270, 2, 4, 17, '2022-05-30 03:20:05 PM', null, 2, 0, '17.20', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8271, 2, 4, 18, '2022-05-30 03:20:05 PM', null, 2, 1, '142.16', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8272, 2, 6, 29, '2022-05-30 03:20:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8273, 2, 6, 24, '2022-05-30 03:20:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8274, 2, 6, 25, '2022-05-30 03:20:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8275, 2, 6, 26, '2022-05-30 03:20:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8276, 2, 6, 27, '2022-05-30 03:20:05 PM', null, 4, 5, '41.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8277, 2, 6, 28, '2022-05-30 03:20:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8278, 2, 7, 32, '2022-05-30 03:20:05 PM', null, 40, 2, '0.40', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8279, 2, 3, 12, '2022-05-30 03:21:05 PM', null, 1, 4, '15.21', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8280, 2, 3, 13, '2022-05-30 03:21:05 PM', null, 1, 5, '1041', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8281, 2, 3, 14, '2022-05-30 03:21:05 PM', null, 1, 0, '1.42', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8282, 2, 3, 15, '2022-05-30 03:21:05 PM', null, 1, 1, '17', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8283, 2, 3, 16, '2022-05-30 03:21:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8284, 2, 3, 11, '2022-05-30 03:21:05 PM', null, 1, 3, '17.21', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8285, 2, 4, 17, '2022-05-30 03:21:05 PM', null, 2, 0, '17.21', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8286, 2, 4, 18, '2022-05-30 03:21:05 PM', null, 2, 1, '141.89', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8287, 2, 4, 19, '2022-05-30 03:21:05 PM', null, 2, 2, '528', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8288, 2, 4, 20, '2022-05-30 03:21:05 PM', null, 2, 3, '26.2000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8289, 2, 6, 29, '2022-05-30 03:21:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8290, 2, 6, 24, '2022-05-30 03:21:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8291, 2, 6, 25, '2022-05-30 03:21:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8292, 2, 6, 26, '2022-05-30 03:21:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8293, 2, 6, 27, '2022-05-30 03:21:05 PM', null, 4, 5, '41.59', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8294, 2, 6, 28, '2022-05-30 03:21:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8295, 2, 7, 32, '2022-05-30 03:21:05 PM', null, 40, 2, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8296, 2, 7, 32, '2022-05-30 03:22:05 PM', null, 40, 2, '0.38', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8297, 2, 3, 14, '2022-05-30 03:22:05 PM', null, 1, 0, '1.43', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8298, 2, 3, 15, '2022-05-30 03:22:05 PM', null, 1, 1, '16', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8299, 2, 3, 16, '2022-05-30 03:22:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8300, 2, 3, 11, '2022-05-30 03:22:05 PM', null, 1, 3, '17.22', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8301, 2, 3, 12, '2022-05-30 03:22:05 PM', null, 1, 4, '15.22', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8302, 2, 3, 13, '2022-05-30 03:22:05 PM', null, 1, 5, '1042', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8303, 2, 4, 18, '2022-05-30 03:22:05 PM', null, 2, 1, '141.64', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8304, 2, 4, 19, '2022-05-30 03:22:05 PM', null, 2, 2, '521', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8305, 2, 4, 20, '2022-05-30 03:22:05 PM', null, 2, 3, '26.1000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8306, 2, 4, 17, '2022-05-30 03:22:05 PM', null, 2, 0, '17.22', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8307, 2, 6, 28, '2022-05-30 03:22:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8308, 2, 6, 29, '2022-05-30 03:22:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8309, 2, 6, 24, '2022-05-30 03:22:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8310, 2, 6, 25, '2022-05-30 03:22:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8311, 2, 6, 26, '2022-05-30 03:22:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8312, 2, 6, 27, '2022-05-30 03:22:05 PM', null, 4, 5, '42.80', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8313, 2, 3, 11, '2022-05-30 03:23:05 PM', null, 1, 3, '17.23', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8314, 2, 3, 12, '2022-05-30 03:23:05 PM', null, 1, 4, '15.23', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8315, 2, 3, 13, '2022-05-30 03:23:05 PM', null, 1, 5, '1043', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8316, 2, 3, 14, '2022-05-30 03:23:05 PM', null, 1, 0, '1.44', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8317, 2, 3, 15, '2022-05-30 03:23:05 PM', null, 1, 1, '15', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8318, 2, 3, 16, '2022-05-30 03:23:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8319, 2, 4, 17, '2022-05-30 03:23:05 PM', null, 2, 0, '17.23', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8320, 2, 4, 18, '2022-05-30 03:23:05 PM', null, 2, 1, '141.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8321, 2, 4, 19, '2022-05-30 03:23:05 PM', null, 2, 2, '515', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8322, 2, 4, 20, '2022-05-30 03:23:05 PM', null, 2, 3, '26.1000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8323, 2, 6, 24, '2022-05-30 03:23:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8324, 2, 6, 25, '2022-05-30 03:23:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8325, 2, 6, 26, '2022-05-30 03:23:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8326, 2, 6, 27, '2022-05-30 03:23:05 PM', null, 4, 5, '41.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8327, 2, 6, 28, '2022-05-30 03:23:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8328, 2, 6, 29, '2022-05-30 03:23:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8329, 2, 7, 32, '2022-05-30 03:23:05 PM', null, 40, 2, '0.37', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8330, 2, 7, 32, '2022-05-30 03:24:05 PM', null, 40, 2, '0.36', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8331, 2, 3, 13, '2022-05-30 03:24:05 PM', null, 1, 5, '1044', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8332, 2, 3, 14, '2022-05-30 03:24:05 PM', null, 1, 0, '1.45', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8333, 2, 3, 15, '2022-05-30 03:24:05 PM', null, 1, 1, '14', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8334, 2, 3, 16, '2022-05-30 03:24:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8335, 2, 3, 11, '2022-05-30 03:24:05 PM', null, 1, 3, '17.24', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8336, 2, 3, 12, '2022-05-30 03:24:05 PM', null, 1, 4, '15.24', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8337, 2, 4, 17, '2022-05-30 03:24:05 PM', null, 2, 0, '17.24', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8338, 2, 4, 18, '2022-05-30 03:24:05 PM', null, 2, 1, '141.15', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8339, 2, 4, 19, '2022-05-30 03:24:05 PM', null, 2, 2, '508', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8340, 2, 4, 20, '2022-05-30 03:24:05 PM', null, 2, 3, '26.1000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8341, 2, 6, 29, '2022-05-30 03:24:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8342, 2, 6, 24, '2022-05-30 03:24:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8343, 2, 6, 25, '2022-05-30 03:24:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8344, 2, 6, 26, '2022-05-30 03:24:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8345, 2, 6, 27, '2022-05-30 03:24:05 PM', null, 4, 5, '41.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8346, 2, 6, 28, '2022-05-30 03:24:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8347, 2, 3, 12, '2022-05-30 03:25:05 PM', null, 1, 4, '15.25', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8348, 2, 3, 13, '2022-05-30 03:25:05 PM', null, 1, 5, '1045', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8349, 2, 3, 14, '2022-05-30 03:25:05 PM', null, 1, 0, '1.46', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8350, 2, 3, 15, '2022-05-30 03:25:05 PM', null, 1, 1, '13', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8351, 2, 3, 16, '2022-05-30 03:25:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8352, 2, 3, 11, '2022-05-30 03:25:05 PM', null, 1, 3, '17.25', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8353, 2, 4, 17, '2022-05-30 03:25:05 PM', null, 2, 0, '17.25', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8354, 2, 4, 18, '2022-05-30 03:25:05 PM', null, 2, 1, '140.88', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8355, 2, 4, 19, '2022-05-30 03:25:05 PM', null, 2, 2, '501', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8356, 2, 4, 20, '2022-05-30 03:25:05 PM', null, 2, 3, '26.2000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8357, 2, 6, 28, '2022-05-30 03:25:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8358, 2, 6, 29, '2022-05-30 03:25:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8359, 2, 6, 24, '2022-05-30 03:25:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8360, 2, 6, 25, '2022-05-30 03:25:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8361, 2, 6, 26, '2022-05-30 03:25:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8362, 2, 6, 27, '2022-05-30 03:25:05 PM', null, 4, 5, '41.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8363, 2, 7, 32, '2022-05-30 03:25:05 PM', null, 40, 2, '0.35', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8364, 2, 3, 14, '2022-05-30 03:26:05 PM', null, 1, 0, '1.47', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8365, 2, 3, 15, '2022-05-30 03:26:05 PM', null, 1, 1, '12', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8366, 2, 3, 16, '2022-05-30 03:26:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8367, 2, 3, 11, '2022-05-30 03:26:05 PM', null, 1, 3, '17.26', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8368, 2, 3, 12, '2022-05-30 03:26:05 PM', null, 1, 4, '15.26', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8369, 2, 3, 13, '2022-05-30 03:26:05 PM', null, 1, 5, '1046', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8370, 2, 4, 19, '2022-05-30 03:26:05 PM', null, 2, 2, '494', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8371, 2, 4, 20, '2022-05-30 03:26:05 PM', null, 2, 3, '26.2000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8372, 2, 4, 17, '2022-05-30 03:26:05 PM', null, 2, 0, '17.26', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8373, 2, 4, 18, '2022-05-30 03:26:05 PM', null, 2, 1, '140.63', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8374, 2, 6, 28, '2022-05-30 03:26:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8375, 2, 6, 29, '2022-05-30 03:26:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8376, 2, 6, 24, '2022-05-30 03:26:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8377, 2, 6, 25, '2022-05-30 03:26:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8378, 2, 6, 26, '2022-05-30 03:26:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8379, 2, 6, 27, '2022-05-30 03:26:05 PM', null, 4, 5, '40.40', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8380, 2, 7, 32, '2022-05-30 03:26:05 PM', null, 40, 2, '0.34', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8381, 2, 3, 14, '2022-05-30 03:27:05 PM', null, 1, 0, '1.48', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8382, 2, 3, 15, '2022-05-30 03:27:05 PM', null, 1, 1, '11', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8383, 2, 3, 16, '2022-05-30 03:27:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8384, 2, 3, 11, '2022-05-30 03:27:05 PM', null, 1, 3, '17.27', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8385, 2, 3, 12, '2022-05-30 03:27:05 PM', null, 1, 4, '15.27', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8386, 2, 3, 13, '2022-05-30 03:27:05 PM', null, 1, 5, '1047', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8387, 2, 4, 17, '2022-05-30 03:27:05 PM', null, 2, 0, '17.27', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8388, 2, 4, 18, '2022-05-30 03:27:05 PM', null, 2, 1, '140.38', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8389, 2, 4, 19, '2022-05-30 03:27:05 PM', null, 2, 2, '487', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8390, 2, 4, 20, '2022-05-30 03:27:05 PM', null, 2, 3, '26.2000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8391, 2, 6, 26, '2022-05-30 03:27:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8392, 2, 6, 27, '2022-05-30 03:27:05 PM', null, 4, 5, '42.00', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8393, 2, 6, 28, '2022-05-30 03:27:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8394, 2, 6, 29, '2022-05-30 03:27:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8395, 2, 6, 24, '2022-05-30 03:27:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8396, 2, 6, 25, '2022-05-30 03:27:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8397, 2, 7, 32, '2022-05-30 03:27:05 PM', null, 40, 2, '0.33', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8398, 2, 3, 16, '2022-05-30 03:28:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8399, 2, 3, 11, '2022-05-30 03:28:05 PM', null, 1, 3, '17.28', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8400, 2, 3, 12, '2022-05-30 03:28:05 PM', null, 1, 4, '15.28', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8401, 2, 3, 13, '2022-05-30 03:28:05 PM', null, 1, 5, '1048', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8402, 2, 3, 14, '2022-05-30 03:28:05 PM', null, 1, 0, '1.49', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8403, 2, 3, 15, '2022-05-30 03:28:05 PM', null, 1, 1, '10', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8404, 2, 4, 20, '2022-05-30 03:28:05 PM', null, 2, 3, '26.2000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8405, 2, 4, 17, '2022-05-30 03:28:05 PM', null, 2, 0, '17.28', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8406, 2, 4, 18, '2022-05-30 03:28:05 PM', null, 2, 1, '140.13', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8407, 2, 4, 19, '2022-05-30 03:28:05 PM', null, 2, 2, '481', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8408, 2, 6, 27, '2022-05-30 03:28:05 PM', null, 4, 5, '42.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8409, 2, 6, 28, '2022-05-30 03:28:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8410, 2, 6, 29, '2022-05-30 03:28:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8411, 2, 6, 24, '2022-05-30 03:28:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8412, 2, 6, 25, '2022-05-30 03:28:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8413, 2, 6, 26, '2022-05-30 03:28:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8414, 2, 7, 32, '2022-05-30 03:28:05 PM', null, 40, 2, '0.32', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8415, 2, 4, 17, '2022-05-30 03:29:05 PM', null, 2, 0, '17.29', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8416, 2, 4, 18, '2022-05-30 03:29:05 PM', null, 2, 1, '139.87', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8417, 2, 4, 19, '2022-05-30 03:29:05 PM', null, 2, 2, '474', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8418, 2, 4, 20, '2022-05-30 03:29:05 PM', null, 2, 3, '26.2000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8419, 2, 6, 24, '2022-05-30 03:29:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8420, 2, 6, 25, '2022-05-30 03:29:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8421, 2, 6, 26, '2022-05-30 03:29:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8422, 2, 6, 27, '2022-05-30 03:29:05 PM', null, 4, 5, '42.00', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8423, 2, 6, 28, '2022-05-30 03:29:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8424, 2, 6, 29, '2022-05-30 03:29:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8425, 2, 7, 32, '2022-05-30 03:29:05 PM', null, 40, 2, '0.31', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8426, 2, 3, 16, '2022-05-30 03:29:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8427, 2, 3, 11, '2022-05-30 03:29:05 PM', null, 1, 3, '17.29', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8428, 2, 3, 12, '2022-05-30 03:29:05 PM', null, 1, 4, '15.29', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8429, 2, 3, 13, '2022-05-30 03:29:05 PM', null, 1, 5, '1049', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8430, 2, 3, 14, '2022-05-30 03:29:05 PM', null, 1, 0, '1.50', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8431, 2, 3, 15, '2022-05-30 03:29:05 PM', null, 1, 1, '9', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8432, 2, 7, 32, '2022-05-30 03:30:05 PM', null, 40, 2, '0.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8433, 2, 3, 14, '2022-05-30 03:30:05 PM', null, 1, 0, '1.51', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8434, 2, 3, 15, '2022-05-30 03:30:05 PM', null, 1, 1, '8', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8435, 2, 3, 16, '2022-05-30 03:30:05 PM', null, 1, 2, '0', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8436, 2, 3, 11, '2022-05-30 03:30:05 PM', null, 1, 3, '17.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8437, 2, 3, 12, '2022-05-30 03:30:05 PM', null, 1, 4, '15.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8438, 2, 3, 13, '2022-05-30 03:30:05 PM', null, 1, 5, '1050', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8439, 2, 4, 20, '2022-05-30 03:30:05 PM', null, 2, 3, '26.2000', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8440, 2, 4, 17, '2022-05-30 03:30:05 PM', null, 2, 0, '17.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8441, 2, 4, 18, '2022-05-30 03:30:05 PM', null, 2, 1, '139.62', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8442, 2, 4, 19, '2022-05-30 03:30:05 PM', null, 2, 2, '467', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8443, 2, 6, 25, '2022-05-30 03:30:05 PM', null, 4, 3, '0.79', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8444, 2, 6, 26, '2022-05-30 03:30:05 PM', null, 4, 4, '12.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8445, 2, 6, 27, '2022-05-30 03:30:05 PM', null, 4, 5, '41.19', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8446, 2, 6, 28, '2022-05-30 03:30:05 PM', null, 4, 0, '3.30', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8447, 2, 6, 29, '2022-05-30 03:30:05 PM', null, 4, 1, '0.39', 1);
INSERT INTO sensors_dataset( id, aquahub_id, device_id, sensor_id, created_at, user_time, local_device_id, local_sensor_id, "value", account_id ) VALUES ( 8448, 2, 6, 24, '2022-05-30 03:30:05 PM', null, 4, 2, '4.51', 1);
INSERT INTO users_accounts( id, user_id, account_id, roles, status, created_at, updated_at, archived_at ) VALUES ( 1, 1, 1, '{manager}', 'active', '2022-07-11 06:19:52 PM', '2022-07-11 06:37:52 PM', null);
INSERT INTO users_accounts( id, user_id, account_id, roles, status, created_at, updated_at, archived_at ) VALUES ( 2, 1, 2, '{manager}', 'active', '2022-07-11 06:50:01 PM', '2022-07-11 06:50:01 PM', null);
INSERT INTO users_accounts( id, user_id, account_id, roles, status, created_at, updated_at, archived_at ) VALUES ( 4, 2, 1, '{user}', 'active', '2022-07-11 06:40:36 PM', '2022-07-11 06:40:36 PM', null);
INSERT INTO users_accounts( id, user_id, account_id, roles, status, created_at, updated_at, archived_at ) VALUES ( 3, 2, 2, '{manager}', 'active', '2022-07-11 06:40:36 PM', '2022-07-26 08:16:42 PM', null);
INSERT INTO users_accounts( id, user_id, account_id, roles, status, created_at, updated_at, archived_at ) VALUES ( 9, 29, 24, '{user}', 'active', '2022-08-26 11:30:21 PM', '2022-08-26 11:30:21 PM', null);
INSERT INTO users_accounts( id, user_id, account_id, roles, status, created_at, updated_at, archived_at ) VALUES ( 10, 30, 25, '{user}', 'active', '2022-08-27 03:43:13 PM', '2022-08-27 03:43:13 PM', null);
INSERT INTO account_preferences( account_id, name, "value", created_at, updated_at, archived_at ) VALUES ( 1, 'datetime_format', '02 Jan 06 15:04 -0700', '2022-07-11 06:33:46 PM', '2022-07-11 06:33:46 PM', null);
INSERT INTO account_preferences( account_id, name, "value", created_at, updated_at, archived_at ) VALUES ( 1, 'date_format', '2006-01-02', '2022-07-11 06:33:46 PM', '2022-07-11 06:33:46 PM', null);
INSERT INTO account_preferences( account_id, name, "value", created_at, updated_at, archived_at ) VALUES ( 1, 'time_format', '15:04:05', '2022-07-11 06:33:46 PM', '2022-07-11 06:33:46 PM', null);
INSERT INTO checklist_items( id, checklist_id, title, description, done, created_at, updated_at, archived_at ) VALUES ( 3, 33, 'UP item 2', 'UP description 2', false, '2022-08-26 04:36:45 PM', '2022-08-26 04:36:45 PM', null);
