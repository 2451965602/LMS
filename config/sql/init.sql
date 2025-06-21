-- 用户表
CREATE TABLE Users (
                           id INT AUTO_INCREMENT PRIMARY KEY,
                           name VARCHAR(50) NOT NULL UNIQUE,
                           password VARCHAR(255) NOT NULL,
                           permission ENUM('admin', 'librarian', 'member') NOT NULL DEFAULT 'member',
                           phone VARCHAR(20),
                           register_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           status ENUM('active', 'suspended', 'inactive') DEFAULT 'active'
) COMMENT '系统用户信息表';

-- 图书类型表（元数据）
CREATE TABLE BookTypes (
                               ISBN VARCHAR(20) PRIMARY KEY,
                               title VARCHAR(100) NOT NULL,
                               author VARCHAR(50) NOT NULL,
                               category VARCHAR(50) NOT NULL,
                               publisher VARCHAR(50) NOT NULL,
                               publish_year INT NOT NULL,
                               description TEXT,
                               total_copies INT NOT NULL DEFAULT 0,
                               available_copies INT NOT NULL DEFAULT 0
) COMMENT '图书元数据信息表';

-- 图书实体表（具体副本）
CREATE TABLE Books (
                           id INT AUTO_INCREMENT PRIMARY KEY,
                           ISBN VARCHAR(20) NOT NULL,
                           location VARCHAR(50) NOT NULL,
                           status ENUM('available', 'checked_out', 'lost', 'damaged') DEFAULT 'available',
                           purchase_date TIMESTAMP NOT NULL,
                           purchase_price DECIMAL(10,2) NOT NULL,
                           last_checkout TIMESTAMP,
                           FOREIGN KEY (ISBN) REFERENCES BookTypes(ISBN) ON UPDATE CASCADE ON DELETE RESTRICT
) COMMENT '图书实体副本表';

-- 借阅记录表
CREATE TABLE BorrowRecords (
                                   id INT AUTO_INCREMENT PRIMARY KEY,
                                   user_id INT NOT NULL,
                                   book_id INT NOT NULL,
                                   title VARCHAR(100) NOT NULL,
                                   checkout_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                   due_date TIMESTAMP NOT NULL,
                                   renewal_count INT DEFAULT 0,
                                   return_date TIMESTAMP,
                                   status ENUM('checked_out', 'returned', 'overdue', 'lost') DEFAULT 'checked_out',
                                   late_fee DECIMAL(10,2) DEFAULT 0.00,
                                   FOREIGN KEY (user_id) REFERENCES Users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
                                   FOREIGN KEY (book_id) REFERENCES Books(id) ON UPDATE CASCADE ON DELETE RESTRICT
) COMMENT '图书借阅记录表';


-- 创建索引
CREATE INDEX idx_users_name ON Users(name);
CREATE INDEX idx_booktypes_title ON BookTypes(title);
CREATE INDEX idx_booktypes_author ON BookTypes(author);
CREATE INDEX idx_books_status ON Books(status);
CREATE INDEX idx_borrowrecords_status ON BorrowRecords(status);
CREATE INDEX idx_borrowrecords_user_book ON BorrowRecords(user_id, book_id);
