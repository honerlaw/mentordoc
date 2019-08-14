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

CREATE TABLE IF NOT EXISTS `organization_user` (
  `organization_id` CHAR(36) NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  PRIMARY KEY (`organization_id`, `user_id`),
  CONSTRAINT `fk_organization_user_organization_id` FOREIGN KEY (`organization_id`) REFERENCES organization(`id`),
  CONSTRAINT `fk_organization_user_user_id` FOREIGN KEY (`user_id`) REFERENCES user(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
