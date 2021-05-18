BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS Addresses(
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  Address INTEGER NOT NULL,
  Owner TEXT,
  Remarks TEXT);
CREATE UNIQUE INDEX IF NOT EXISTS Addresses_IDX ON Addresses(Address);
CREATE VIEW IF NOT EXISTS v_Addresses AS
  SELECT ID, Address, Address>>24 AS A, Address>>16&255 AS B, Address>>8&255 AS C, Address&255 AS D, Owner, Remarks
  FROM Addresses;
INSERT INTO Addresses(Address, Owner, Remarks) VALUES(3232235812, 'Jose', 'pi-u'); -- 192.168.1.36
INSERT INTO Addresses(Address, Owner, Remarks) VALUES(3232235811, 'Jose', 'pi-0'); -- 192.168.1.35
INSERT INTO Addresses(Address, Owner, Remarks) VALUES(3232235813, 'Jose', 'pi-1'); -- 192.168.1.37
INSERT INTO Addresses(Address, Owner, Remarks) VALUES(3232235792, 'Jose', 'pi-m'); -- 192.168.1.16
INSERT INTO Addresses(Address, Owner, Remarks) VALUES(3232235803, 'Jose', 'Inspiron 15'); -- 192.168.1.27
INSERT INTO Addresses(Address, Remarks) VALUES(2130706433, 'localhost'); -- 127.0.0.1
CREATE TABLE IF NOT EXISTS Services(
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  Service TEXT NOT NULL,
  Remarks TEXT);
CREATE UNIQUE INDEX IF NOT EXISTS Services_IDX ON Services(Service);
INSERT INTO Services(Service) VALUES('access');
INSERT INTO Services(Service) VALUES('digest');
INSERT INTO Services(Service) VALUES('pubkey');
INSERT INTO Services(Service) VALUES('random');
CREATE TABLE IF NOT EXISTS ACL(
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  Address_ID INTEGER NOT NULL,
  Service_ID INTEGER NOT NULL,
  CanRead INTEGER NOT NULL,
  CanWrite INTEGER NOT NULL,
  Remarks TEXT);
CREATE UNIQUE INDEX IF NOT EXISTS ACL_IDX ON ACL(Address_ID, Service_ID);
CREATE VIEW IF NOT EXISTS v_ACL AS
  SELECT A.ID, B.Address, B.A, B.B, B.C, B.D, B.Owner, (SELECT Service FROM Services C WHERE A.Service_ID = C.ID) Service, A.CanRead, A.CanWrite
  FROM ACL A, v_Addresses B
  WHERE A.Address_ID = B.ID;
INSERT INTO ACL(Address_ID, Service_ID, CanRead, CanWrite) VALUES(6, 1, 1, 1); -- localhost can read & write access
INSERT INTO ACL(Address_ID, Service_ID, CanRead, CanWrite) VALUES(4, 1, 1, 1); -- pi-m can read & write access
INSERT INTO ACL(Address_ID, Service_ID, CanRead, CanWrite) VALUES(5, 1, 1, 1); -- laptop can read & write access
INSERT INTO ACL(Address_ID, Service_ID, CanRead, CanWrite) VALUES(2, 2, 0, 0); -- pi-0 can't access digest
INSERT INTO ACL(Address_ID, Service_ID, CanRead, CanWrite) VALUES(3, 2, 0, 0); -- pi-1 can't access digest
COMMIT;
