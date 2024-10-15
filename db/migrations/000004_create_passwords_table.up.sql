CREATE TABLE IF NOT EXISTS passwords (
	UserLogin   varchar(255) NOT NULL,    
	Title       varchar(255) NOT NULL,
	Login       varchar(255),    
	Password    varchar(255),  
	CreatedAt   timestamp DEFAULT (now()), 		 
PRIMARY KEY (UserLogin, Title));
