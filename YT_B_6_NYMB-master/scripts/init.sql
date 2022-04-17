SET FOREIGN_KEY_CHECKS=0;
DROP TABLE IF EXISTS Transactions;
DROP TABLE IF EXISTS Balances;
DROP TABLE IF EXISTS Currencies;
DROP TABLE IF EXISTS Favorites;
DROP TABLE IF EXISTS Accounts;
DROP TABLE IF EXISTS Permissions;
DROP TABLE IF EXISTS VaultUsers;
DROP TABLE IF EXISTS Vaults;
DROP TABLE IF EXISTS UserRoles;
DROP TABLE IF EXISTS Roles;
DROP TABLE IF EXISTS Users;
SET FOREIGN_KEY_CHECKS=1;

CREATE TABLE Users (
    ID int NOT NULL AUTO_INCREMENT,
    UserName varchar(255) NOT NULL UNIQUE,
    Email varchar(255) NOT NULL UNIQUE,
    Password varchar(255) NOT NULL, 
    FirstName varchar(255),
    LastName varchar(255),
    PRIMARY KEY (ID)
);

CREATE TABLE Roles (
    ID int NOT NULL AUTO_INCREMENT,
    Level varchar(255) NOT NULL UNIQUE,
    PRIMARY KEY (ID)
);

CREATE TABLE UserRoles (
    UserID int NOT NULL UNIQUE,
    RoleID int NOT NULL,
    PRIMARY KEY (UserID, RoleID),
    FOREIGN KEY (UserID) REFERENCES Users(ID),
    FOREIGN KEY (RoleID) REFERENCES Roles(ID)
);

CREATE TABLE Vaults (
    ID int NOT NULL AUTO_INCREMENT,
    Name varchar(255) NOT NULL,
    OwnerID int NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (OwnerID) REFERENCES Users(ID)
);

CREATE TABLE VaultUsers (
    UserID int NOT NULL UNIQUE,
    VaultID int NOT NULL,
    PRIMARY KEY (UserID, VaultID),
    FOREIGN KEY (UserID) REFERENCES Users(ID),
    FOREIGN KEY (VaultID) REFERENCES Vaults(ID)
);

CREATE TABLE Permissions (
    ID int NOT NULL AUTO_INCREMENT,
    UserID int NOT NULL,
    VaultID int NOT NULL,
    RequestTransaction boolean DEFAULT 1,
    ApproveTransaction boolean DEFAULT 0,
    AddUser boolean DEFAULT 0,
    RemoveUser boolean DEFAULT 0,
    AddFunds boolean DEFAULT 1,
    RemoveFunds boolean DEFAULT 0,
    UserName varchar(255) NOT NULL, 
    PRIMARY KEY (ID),
    FOREIGN KEY (UserID) REFERENCES Users(ID),
    FOREIGN KEY (VaultID) REFERENCES Vaults(ID),
    FOREIGN KEY (UserName) REFERENCES Users(UserName)
);

CREATE TABLE Accounts (
    ID int NOT NULL AUTO_INCREMENT,
    UserID int UNIQUE,
    VaultID int UNIQUE,
    CHECK (UserID IS NOT NULL OR VaultID IS NOT NULL),
    PRIMARY KEY (ID),
    FOREIGN KEY (UserID) REFERENCES Users(ID),
    FOREIGN KEY (VaultID) REFERENCES Vaults(ID)
);

CREATE TABLE Favorites (
	ID int NOT NULL AUTO_INCREMENT,
    UserID int NOT NULL,
    AccountID int NOT NULL,
    Name varchar(255) NOT NULL,
    PRIMARY KEY (ID),
    UNIQUE KEY (UserID, AccountID),
    FOREIGN KEY (UserID) REFERENCES Users(ID),
    FOREIGN KEY (AccountID) REFERENCES Accounts(ID)
);

CREATE TABLE Currencies (
    ID int NOT NULL AUTO_INCREMENT,
    Name varchar(255) NOT NULL,
    ShortName varchar(255),
    UnitPrice int NOT NULL,
    PRIMARY KEY (ID)
);

CREATE TABLE Balances (
    ID int NOT NULL AUTO_INCREMENT,
    AccountID int NOT NULL,
    CurrencyID int NOT NULL,
    Amount FLOAT(20, 6) NOT NULL,
    PRIMARY KEY (ID),
    KEY (AccountID, CurrencyID),
    FOREIGN KEY (AccountID) REFERENCES Accounts(ID),
    FOREIGN KEY (CurrencyID) REFERENCES Currencies(ID)
);

CREATE TABLE Transactions (
    ID int NOT NULL AUTO_INCREMENT,
    FromID int,
    ToID int,
    CurrencyID int NOT NULL,
    Amount FLOAT(20, 6) NOT NULL,
    Created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Completed datetime DEFAULT NULL,
    Status varchar(10) DEFAULT 'new' CHECK (Status IN ('new', 'approved', 'denied')),
    CHECK (FromID IS NOT NULL OR ToID IS NOT NULL),
    PRIMARY KEY (ID),
    FOREIGN KEY (FromID) REFERENCES Accounts(ID),
    FOREIGN KEY (ToID) REFERENCES Accounts(ID),
    FOREIGN KEY (CurrencyID) REFERENCES Currencies(ID)
);
