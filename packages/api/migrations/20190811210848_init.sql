-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

/*
this was ripped from https://dba.stackexchange.com/questions/7147/find-highest-level-of-a-hierarchical-field-with-vs-without-ctes
modified to work specifically with the folder table
Also modified to return a max of 50 folders
*/
-- +migrate StatementBegin
DROP FUNCTION IF EXISTS `mentor_doc`.`GetFolderAncestry`;
CREATE FUNCTION `mentor_doc`.`GetFolderAncestry` (GivenID CHAR(36)) RETURNS VARCHAR(1850)
DETERMINISTIC
BEGIN
    DECLARE rv VARCHAR(1850);
    DECLARE cm CHAR(1);
    DECLARE ch CHAR(36);

    SET rv = '';
    SET cm = '';
    SET ch = GivenID;
    WHILE ch != -1 DO
        SELECT IFNULL(parent_folder_id,-1) INTO ch FROM
        (SELECT parent_folder_id FROM folder WHERE id = ch) A;
        IF ch != -1 THEN
            SET rv = CONCAT(rv,cm,ch);
            SET cm = ',';
        END IF;
    END WHILE;
    RETURN rv;
END;
-- +migrate StatementEnd

CREATE TABLE IF NOT EXISTS `organization` (
  `id` CHAR(36) NOT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  `deleted_at` BIGINT NULL DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_organization_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `user` (
  `id` CHAR(36) NOT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  `deleted_at` BIGINT NULL DEFAULT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL UNIQUE,
  PRIMARY KEY (`id`),
  KEY `idx_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `permission` (
  `id` CHAR(36) NOT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  `deleted_at` BIGINT NULL DEFAULT NULL,
  `resource_path` varchar(255) NOT NULL,
  `action` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_perrmission_resource_path_action` (`resource_path`, `action`),
  KEY `idx_permission_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `role` (
  `id` CHAR(36) NOT NULL,
  `name` VARCHAR(36) NOT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  `deleted_at` BIGINT NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_role_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `role_permission` (
  `role_id` CHAR(36) NOT NULL,
  `permission_id` CHAR(36) NOT NULL,
  PRIMARY KEY (`role_id`, `permission_id`),
  FOREIGN KEY (`role_id`) REFERENCES role(`id`),
  FOREIGN KEY (`permission_id`) REFERENCES permission(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `user_role` (
  `user_id` CHAR(36) NOT NULL,
  `role_id` CHAR(36) NOT NULL,
  `resource_id` CHAR(36) NOT NULL,
  PRIMARY KEY (`user_id`, `role_id`, `resource_id`),
  FOREIGN KEY (`user_id`) REFERENCES user(`id`),
  FOREIGN KEY (`role_id`) REFERENCES role(`id`),
  UNIQUE KEY `uk_user_role_user_id_role_id_resource_id` (`user_id`, `role_id`, `resource_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `folder` (
  `id` CHAR(36) NOT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  `deleted_at` BIGINT NULL DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `organization_id` CHAR(36) NOT NULL,
  `parent_folder_id` CHAR(36) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`organization_id`) REFERENCES organization(`id`),
  FOREIGN KEY (`parent_folder_id`) REFERENCES folder(`id`),
  KEY `idx_folder_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `document` (
  `id` CHAR(36) NOT NULL,
  `organization_id` CHAR(36) NOT NULL,
  `folder_id` CHAR(36) NULL DEFAULT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  `deleted_at` BIGINT NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`organization_id`) REFERENCES organization(`id`),
  FOREIGN KEY (`folder_id`) REFERENCES folder(`id`),
  KEY `idx_document_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `document_draft` (
  `id` CHAR(36) NOT NULL,
  `document_id` CHAR(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `creator_id` CHAR(36) NOT NULL,
  `published_at` BIGINT NULL DEFAULT NULL,
  `retracted_at` BIGINT NULL DEFAULT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  `deleted_at` BIGINT NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`document_id`) REFERENCES document(`id`),
  FOREIGN KEY (`creator_id`) REFERENCES user(`id`),
  FULLTEXT (`name`),
  KEY `idx_document_draft_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE IF NOT EXISTS `document_draft_content` (
  `id` CHAR(36) NOT NULL,
  `document_draft_id` CHAR(36) NOT NULL,
  `content` MEDIUMTEXT NOT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  `deleted_at` BIGINT NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`document_draft_id`) REFERENCES document_draft(`id`),
  FULLTEXT (`content`),
  KEY `idx_document_draft_content_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `resource_history` (
  `id` CHAR(36) NOT NULL,
  `resource_id` CHAR(36) NOT NULL,
  `resource_name` varchar(255) NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  `action` varchar(255) NOT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  `deleted_at` BIGINT NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES user(`id`),
  KEY `idx_resource_history_deleted_at` (`deleted_at`),
  KEY `idx_resource_history_identifier` (`resource_id`, `resource_name`, `user_id`, `action`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `resource_history`;
DROP TABLE `document_draft_content`;
DROP TABLE `document_draft`;
DROP TABLE `document`;
DROP TABLE `folder`;
DROP TABLE `user_role`;
DROP TABLE `role_permission`;
DROP TABLE `role`;
DROP TABLE `permission`;
DROP TABLE `organization`;
DROP TABLE `user`;
