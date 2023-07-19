CREATE TABLE IF NOT EXISTS `webhook` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uuid` char(32) NOT NULL,
  `name` varchar(255) NOT NULL,
  `summary` varchar(255) DEFAULT NULL,
  `url` varchar(255) NOT NULL,
  `method` varchar(255) NOT NULL,
  `headers` text DEFAULT NULL,
  `timeout` int DEFAULT NULL,
  `condition_validator` int DEFAULT NULL,
  `condition_filter` text DEFAULT NULL,
  `created` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
;

ALTER TABLE IF EXISTS `webhook` 
  CONVERT TO CHARACTER SET utf8mb4
;

ALTER TABLE IF EXISTS `webhook` 
  MODIFY IF EXISTS `created` datetime DEFAULT CURRENT_TIMESTAMP,
  MODIFY IF EXISTS `updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
;

ALTER TABLE IF EXISTS `webhook`
  -- add unique
  ADD CONSTRAINT `uuid` UNIQUE IF NOT EXISTS (`uuid`) 
;
