-- add column `cluster_client_token_uuid`
ALTER TABLE IF EXISTS `session`
    DROP IF EXISTS `cluster_client_token_uuid`
;