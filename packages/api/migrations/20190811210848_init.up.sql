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
  KEY `idx_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `role` (
  `id` CHAR(36) NOT NULL,
  `name` VARCHAR(36) NOT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  `deleted_at` BIGINT NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_deleted_at` (`deleted_at`)
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
  PRIMARY KEY (`user_id`, `role_id`),
  FOREIGN KEY (`user_id`) REFERENCES user(`id`),
  FOREIGN KEY (`role_id`) REFERENCES role(`id`),
  UNIQUE KEY `uk_user_role_user_id_role_id_resource_id` (`user_id`, `role_id`, `resource_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

