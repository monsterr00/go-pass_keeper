CREATE TABLE IF NOT EXISTS files (
	UserLogin   varchar(255) NOT NULL,    
	Title       varchar(255) NOT NULL,
	FileName    varchar(255),
	File        bytea,
	DataType    varchar(50), 
	CreatedAt   timestamp DEFAULT (now()), 		 
PRIMARY KEY (UserLogin, Title));
