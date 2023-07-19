ALTER TABLE IF EXISTS `session`
  -- add column
  ADD IF NOT EXISTS `cluster_client_token_uuid` char(32) NOT NULL AFTER `uuid`,
  -- add index
  ADD INDEX IF NOT EXISTS `cluster_client_token_uuid` (`cluster_client_token_uuid`) 
;
