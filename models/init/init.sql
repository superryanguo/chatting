use mysql;
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'yourpassword';
CREATE DATABASE `chat_db` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin */;
use chat_db;

