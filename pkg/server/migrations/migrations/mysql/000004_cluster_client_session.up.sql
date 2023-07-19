CREATE TABLE IF NOT EXISTS `session` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uuid` char(32) NOT NULL,
  `cluster_uuid` char(32) NOT NULL,
  `token` text NOT NULL,
  `issued_at_time` datetime NOT NULL,
  `expiration_time` datetime NOT NULL,
  `created` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
;

ALTER TABLE IF EXISTS `session` 
  CONVERT TO CHARACTER SET utf8mb4
;

ALTER TABLE IF EXISTS `session` 
  MODIFY IF EXISTS `created` datetime DEFAULT CURRENT_TIMESTAMP,
  MODIFY IF EXISTS `updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
;

ALTER TABLE IF EXISTS `session`
  -- add unique 
  ADD CONSTRAINT `uuid` UNIQUE IF NOT EXISTS (`uuid`) 
;

ALTER TABLE IF EXISTS `session`
  -- add index
  -- for check alive
  ADD INDEX IF NOT EXISTS `cluster_expiration_time` (`cluster_uuid`, `expiration_time`) 
;

ALTER TABLE IF EXISTS `session`
  -- add index 
  -- for check alive
  ADD INDEX IF NOT EXISTS `cluster_uuid` (`cluster_uuid`,`uuid`) 
;
