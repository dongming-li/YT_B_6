-- This file contains data to load into the database for testing

-- Roles
INSERT INTO Roles(ID, Level) VALUES(1, "admin"); -- ID: 1
INSERT INTO Roles(ID, Level) VALUES(2, "user"); -- ID: 2

-- Currencies
INSERT INTO Currencies(Name, ShortName, UnitPrice) VALUES("BitCoin", "BTC", 3500); -- ID: 1
INSERT INTO Currencies(Name, ShortName, UnitPrice) VALUES("Ethereum", "LTC", 200); -- ID: 2
INSERT INTO Currencies(Name, ShortName, UnitPrice) VALUES("Litecoin", "ETH", 50); -- ID: 3
INSERT INTO Currencies(Name, ShortName, UnitPrice) VALUES("United States Dollar", "USD", 1); -- ID: 3

-- Default admin user
INSERT INTO Users(UserName, Email, Password, FirstName, LastName)
VALUES ("admin", "admin@admin.com", "admin", "admin", "admin"); -- ID: 1
INSERT INTO UserRoles(UserID, RoleID) VALUES(1, 1);
INSERT INTO Accounts(UserID) VALUES(1); -- ID: 1
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(1, 1, 50);
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(1, 2, 100);
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(1, 3, 150);
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(1, 4, 200);

-- Blake's test user and related artifacts
INSERT INTO Users(UserName, Email, Password, FirstName, LastName)
VALUES ("brob", "brob@iastate.edu", "1234", "Blake", "Roberts"); -- ID: 2
INSERT INTO UserRoles(UserID, RoleID) VALUES(2, 2);
INSERT INTO Accounts(UserID) VALUES(2); -- ID: 2
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(2, 1, 5);
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(2, 2, 10);
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(2, 3, 15);
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(2, 4, 20);
INSERT INTO Vaults(Name, OwnerID) VALUES ("Blake's Vault", 2); -- ID: 1
INSERT INTO Accounts(VaultID) VALUES(1); -- ID: 3
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(3, 2, 1023);
INSERT INTO Favorites(UserID, AccountID, Name) VALUES(2, 1, "2-1");
INSERT INTO Transactions(FromID, ToID, CurrencyID, Amount) VALUES(2, 1, 1, 23);

-- Matt's test user and related artifacts
INSERT INTO Users(UserName, Email, Password, FirstName, LastName)
VALUES ("maschaff", "maschaff@iastate.edu", "4321", "Matt", "Schaffer"); -- ID: 3
INSERT INTO UserRoles(UserID, RoleID) VALUES(3, 2);
INSERT INTO Accounts(UserID) VALUES(3); -- ID: 3
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(3, 1, 42);
INSERT INTO Vaults(Name, OwnerID) VALUES ("Matt's Vault", 3);
INSERT INTO Accounts(VaultID) VALUES(2); -- ID: 4
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(4, 1, 100);
INSERT INTO Favorites(UserID, AccountID, Name) VALUES(3, 1, "3-1");
INSERT INTO Transactions(FromID, ToID, CurrencyID, Amount) VALUES(2, 1, 2, 23);
INSERT INTO Permissions(
    UserID,
    VaultID,
    RequestTransaction,
    ApproveTransaction,
    AddUser,
    RemoveUser,
    AddFunds,
    RemoveFunds,
    UserName
) VALUES (2,1,1,1,1,1,1,1, "brob");

INSERT INTO Favorites(UserID, AccountID, Name) VALUES(1, 2, "1-2");
INSERT INTO Favorites(UserID, AccountID, Name) VALUES(1, 3, "1-3");

-- Lee's test user and related artifacts
INSERT INTO Users(UserName, Email, Password, FirstName, LastName)
VALUES ("lfulbel", "lfulbel@iastate.edu", "leePassword", "Leelabari", "Fulbel"); -- ID: 4
INSERT INTO UserRoles(UserID, RoleID) VALUES(4, 2);
INSERT INTO Accounts(UserID) VALUES(4); -- ID: 5
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(4, 1, 50);
INSERT INTO Vaults(Name, OwnerID) VALUES ("Lee's Vault", 3);
INSERT INTO Accounts(VaultID) VALUES(3); -- ID: 6
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(5, 3, 7000);
INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(5, 2, 500);
INSERT INTO Favorites(UserID, AccountID, Name) VALUES(4, 6, "5-6");
INSERT INTO Transactions(FromID, ToID, CurrencyID, Amount) VALUES(5, 4, 1, 23);
INSERT INTO Permissions(
    UserID,
    VaultID,
    RequestTransaction,
    ApproveTransaction,
    AddUser,
    RemoveUser,
    AddFunds,
    RemoveFunds,
    UserName
) VALUES (3,2,1,1,1,1,1,1, "lfulbel");
