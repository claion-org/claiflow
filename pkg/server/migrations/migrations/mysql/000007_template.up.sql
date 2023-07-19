CREATE TABLE IF NOT EXISTS `template` (
  `uuid` char(32) NOT NULL,
  `name` varchar(255) NOT NULL,
  `summary` varchar(255) DEFAULT NULL,
  `flow` text NOT NULL,
  `inputs` text DEFAULT NULL,
  `origin` varchar(255) NOT NULL,
  `created` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
;

ALTER TABLE IF EXISTS `template` 
  CONVERT TO CHARACTER SET utf8mb4
;

ALTER TABLE IF EXISTS `template` 
  MODIFY IF EXISTS `created` datetime DEFAULT CURRENT_TIMESTAMP,
  MODIFY IF EXISTS `updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
;