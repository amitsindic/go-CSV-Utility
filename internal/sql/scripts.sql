--CREATE TABLE--

CREATE TABLE default.records (
	  first_name varchar(30) NOT NULL,
	  last_name varchar(30) DEFAULT NULL,
	  company_name varchar(30) DEFAULT NULL,
	  address varchar(100) DEFAULT NULL,
	  city varchar(30) DEFAULT NULL,
	  county varchar(30) DEFAULT NULL,
	  postal varchar(60) DEFAULT NULL,
	  phone varchar(30) DEFAULT NULL,
	  email varchar(50) DEFAULT NULL,
	  web varchar(100) DEFAULT null,
	  
	  PRIMARY KEY (first_name)

)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;