CREATE TABLE IF NOT EXISTS texts (
	UserLogin   varchar(255) NOT NULL,    
	Title       varchar(255) NOT NULL,
	Text        text,
	CreatedAt   timestamp DEFAULT (now()), 		 
PRIMARY KEY (UserLogin, Title));