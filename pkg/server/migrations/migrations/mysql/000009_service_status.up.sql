CREATE TABLE IF NOT EXISTS`service_status` (
  `pdate` date NOT NULL,
  `cluster_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `created` datetime(6) NOT NULL,
  `step_max` int(10) unsigned NOT NULL,
  `step_seq` int(10) unsigned NOT NULL,
  `status` tinyint(3) unsigned NOT NULL,
  `started` datetime DEFAULT NULL,
  `ended` datetime DEFAULT NULL,
  `message` text COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`pdate`,`cluster_uuid`,`uuid`,`created`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`pdate`))
PARTITIONS 5
;

ALTER TABLE IF EXISTS `service_status` 
  CONVERT TO CHARACTER SET utf8mb4
;

ALTER TABLE IF EXISTS `service_status`
  -- add index
  ADD INDEX IF NOT EXISTS `cluster_uuid__uuid` (`cluster_uuid`, `uuid`) 
;

ALTER TABLE IF EXISTS `service_status`
  -- add index
  ADD INDEX IF NOT EXISTS `cluster_uuid__created` (`cluster_uuid`, `created`) 
;

ALTER TABLE IF EXISTS `service_status`
  -- add index
  ADD INDEX IF NOT EXISTS `cluster_uuid` (`cluster_uuid`) 
;

ALTER TABLE IF EXISTS `service_status`
  -- add index
  ADD INDEX IF NOT EXISTS `uuid` (`uuid`) 
;

ALTER TABLE IF EXISTS `service_status`
  -- add index
  ADD INDEX IF NOT EXISTS `created` (`created`) 
;
