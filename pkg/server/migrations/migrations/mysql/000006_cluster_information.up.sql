CREATE TABLE IF NOT EXISTS `cluster_information` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `cluster_uuid` char(32) NOT NULL,
  `polling_count` int(10) DEFAULT NULL,
  `polling_offset` datetime(6) DEFAULT NULL,
  `created` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
;

ALTER TABLE IF EXISTS `cluster_information` 
  CONVERT TO CHARACTER SET utf8mb4
;

ALTER TABLE IF EXISTS `cluster_information` 
  MODIFY IF EXISTS `created` datetime DEFAULT CURRENT_TIMESTAMP,
  MODIFY IF EXISTS `updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
;

ALTER TABLE IF EXISTS `cluster_information`
  -- add unique
  ADD CONSTRAINT `cluster_uuid` UNIQUE IF NOT EXISTS (`cluster_uuid`) 
;
