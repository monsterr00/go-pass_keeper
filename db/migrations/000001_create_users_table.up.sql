CREATE TABLE IF NOT EXISTS users (
	Login          varchar(255) NOT NULL,    
	Password       varchar(255) NOT NULL,    
	CreatedAt      timestamp DEFAULT (now()), 		 
PRIMARY KEY (Login));
