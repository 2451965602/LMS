-- 用户表
CREATE TABLE Users (
                       ID INT AUTO_INCREMENT PRIMARY KEY,
                       Name VARCHAR(50) NOT NULL UNIQUE,
                       Password VARCHAR(255) NOT NULL,
                       Permissions ENUM('admin', 'librarian', 'member') NOT NULL DEFAULT 'member',
                       Phone VARCHAR(20),
                       RegisterDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       Status ENUM('active', 'suspended', 'inactive') DEFAULT 'active'
) COMMENT '系统用户信息表';

-- 图书类型表（元数据）
CREATE TABLE BookTypes (
                           ISBN VARCHAR(20) PRIMARY KEY,
                           Title VARCHAR(100) NOT NULL,
                           Author VARCHAR(50) NOT NULL,
                           Category VARCHAR(50) NOT NULL,
                           Publisher VARCHAR(50) NOT NULL,
                           PublishYear INT,
                           Description TEXT,
                           TotalCopies INT NOT NULL DEFAULT 0,
                           AvailableCopies INT NOT NULL DEFAULT 0
) COMMENT '图书元数据信息表';

-- 图书实体表（具体副本）
CREATE TABLE Books (
                       ID INT AUTO_INCREMENT PRIMARY KEY,
                       ISBN VARCHAR(20) NOT NULL,
                       Location VARCHAR(50) NOT NULL,
                       Status ENUM('available', 'checked_out', 'reserved', 'lost', 'damaged') DEFAULT 'available',
                       PurchaseDate DATE,
                       PurchasePrice DECIMAL(10,2),
                       LastCheckout TIMESTAMP,
                       FOREIGN KEY (ISBN) REFERENCES BookTypes(ISBN) ON UPDATE CASCADE ON DELETE RESTRICT
) COMMENT '图书实体副本表';

-- 借阅记录表
CREATE TABLE BorrowRecords (
                               ID INT AUTO_INCREMENT PRIMARY KEY,
                               UserID INT NOT NULL,
                               BookID INT NOT NULL,
                               CheckoutDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               DueDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL 14 DAY,
                               ReturnDate TIMESTAMP,
                               Status ENUM('checked_out', 'returned', 'overdue', 'lost') DEFAULT 'checked_out',
                               LateFee DECIMAL(10,2) DEFAULT 0.00,
                               FOREIGN KEY (UserID) REFERENCES Users(ID) ON UPDATE CASCADE ON DELETE RESTRICT,
                               FOREIGN KEY (BookID) REFERENCES Books(ID) ON UPDATE CASCADE ON DELETE RESTRICT,
) COMMENT '图书借阅记录表';

-- 预约记录表
CREATE TABLE Reservations (
                              ID INT AUTO_INCREMENT PRIMARY KEY,
                              UserID INT NOT NULL,
                              BookID INT NOT NULL,
                              ISBN VARCHAR(20) NOT NULL,
                              ReserveDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              ExpiryDate TIMESTAMP NOT NULL,
                              Status ENUM('pending', 'fulfilled', 'canceled', 'expired') DEFAULT 'pending',
                              FOREIGN KEY (UserID) REFERENCES Users(ID) ON UPDATE CASCADE ON DELETE CASCADE,
                              FOREIGN KEY (ISBN) REFERENCES BookTypes(ISBN) ON UPDATE CASCADE ON DELETE CASCADE
) COMMENT '图书预约记录表';

-- 创建索引
CREATE INDEX idx_users_name ON Users(Name);
CREATE INDEX idx_booktypes_title ON BookTypes(Title);
CREATE INDEX idx_booktypes_author ON BookTypes(Author);
CREATE INDEX idx_books_status ON Books(Status);
CREATE INDEX idx_borrowrecords_status ON BorrowRecords(Status);
CREATE INDEX idx_borrowrecords_user_book ON BorrowRecords(UserID, BookID);
CREATE INDEX idx_reservations_user_isbn ON Reservations(UserID, ISBN);
