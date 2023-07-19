CREATE TABLE IF NOT EXISTS`service` (
  `pdate` date NOT NULL,
  `cluster_uuid` char(32) NOT NULL,
  `uuid` char(32) NOT NULL,
  `created` datetime(6) NOT NULL,
  `name` varchar(255) NOT NULL,
  `summary` varchar(255) DEFAULT NULL,
  `template_uuid` char(32) NOT NULL,
  `flow` text NOT NULL,
  `inputs` text NOT NULL,
  `step_max` int(10) unsigned NOT NULL,
  `priority` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `subscribed_channel` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`pdate`,`cluster_uuid`,`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
 PARTITION BY HASH (to_days(`pdate`))
PARTITIONS 5
;

ALTER TABLE IF EXISTS `service` 
  CONVERT TO CHARACTER SET utf8mb4
;

ALTER TABLE IF EXISTS `service`
  -- add index
  ADD INDEX IF NOT EXISTS `cluster_uuid__uuid` (`cluster_uuid`, `uuid`) 
;

ALTER TABLE IF EXISTS `service`
  -- add index
  ADD INDEX IF NOT EXISTS `cluster_uuid__created` (`cluster_uuid`, `created`) 
;

ALTER TABLE IF EXISTS `service`
  -- add index
  ADD INDEX IF NOT EXISTS `cluster_uuid` (`cluster_uuid`) 
;

ALTER TABLE IF EXISTS `service`
  -- add index
  ADD INDEX IF NOT EXISTS `uuid` (`uuid`) 
;

ALTER TABLE IF EXISTS `service`
  -- add index
  ADD INDEX IF NOT EXISTS `created` (`created`) 
;
