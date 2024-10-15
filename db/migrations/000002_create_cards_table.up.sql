CREATE TABLE IF NOT EXISTS cards (
	UserLogin   varchar(255) NOT NULL,    
	Title       varchar(255) NOT NULL,
	Bank        varchar(255),
	CardNumber  varchar(16),
	CVV         varchar(3),
	DateExpire  timestamp,
	CardHolder  varchar(255),
	CreatedAt   timestamp DEFAULT (now()), 		 
PRIMARY KEY (UserLogin, Title));
