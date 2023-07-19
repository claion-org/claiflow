CREATE TABLE IF NOT EXISTS `cluster_token` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uuid` char(32) NOT NULL,
  `name` varchar(255) NOT NULL,
  `summary` varchar(255) DEFAULT NULL,
  `cluster_uuid` char(32) NOT NULL,
  `token` varchar(255) NOT NULL,
  `issued_at_time` datetime NOT NULL,
  `expiration_time` datetime NOT NULL,
  `created` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 
;

ALTER TABLE IF EXISTS `cluster_token` 
  CONVERT TO CHARACTER SET utf8mb4
;

ALTER TABLE IF EXISTS `cluster_token` 
  MODIFY IF EXISTS `created` datetime DEFAULT CURRENT_TIMESTAMP,
  MODIFY IF EXISTS `updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
;

ALTER TABLE IF EXISTS `cluster_token`
-- add unique
  ADD CONSTRAINT `uuid` UNIQUE IF NOT EXISTS (`uuid`) 
;

ALTER TABLE IF EXISTS `cluster_token`
  -- add unique
  ADD CONSTRAINT `token` UNIQUE IF NOT EXISTS (`token`) 
;

ALTER TABLE IF EXISTS `cluster_token` 
  -- add index
  DROP INDEX IF EXISTS `IDX_token_cluster_uuid`
;


-- ALTER TABLE IF EXISTS `cluster_token` 
--     ADD CONSTRAINT `FK_cluster`
-- 	FOREIGN KEY (`cluster_uuid`) REFERENCES `cluster` (`uuid`)
-- ; 

ALTER TABLE IF EXISTS `cluster_token` 
  -- add index
  ADD INDEX IF NOT EXISTS `cluster_uuid__token` (`cluster_uuid`, `token`)
;
