Use letscrum;

CREATE TABLE `epic` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sprint_id` bigint(20) unsigned DEFAULT 0,
  `title` varchar(50) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_epic_sprint_id` FOREIGN KEY (`sprint_id`) REFERENCES `sprint` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `feature` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sprint_id` bigint(20) unsigned DEFAULT 0,
  `epic_id` bigint(20) unsigned DEFAULT 0,
  `title` varchar(50) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_feature_sprint_id` FOREIGN KEY (`sprint_id`) REFERENCES `sprint` (`id`),
  CONSTRAINT `fk_feature_epic_id` FOREIGN KEY (`epic_id`) REFERENCES `epic` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `workitem` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sprint_id` bigint(20) unsigned DEFAULT 0,
  `feature_id` bigint(20) unsigned DEFAULT 0,
  `title` varchar(50) NOT NULL,
  `type` varchar(50) NOT NULL,
  `description` varchar(500) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_workitem_sprint_id` FOREIGN KEY (`sprint_id`) REFERENCES `sprint` (`id`),
  CONSTRAINT `fk_workitem_feature_id` FOREIGN KEY (`feature_id`) REFERENCES `feature` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `task` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sprint_id` bigint(20) unsigned DEFAULT 0,
  `workitem_id` bigint(20) unsigned DEFAULT 0,
  `title` varchar(50) NOT NULL,
  `description` varchar(500) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_task_sprint_id` FOREIGN KEY (`sprint_id`) REFERENCES `sprint` (`id`),
  CONSTRAINT `fk_task_workitem_id` FOREIGN KEY (`workitem_id`) REFERENCES `workitem` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `workitem_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `workitem_id` varchar(50) NOT NULL,
  `log` varchar(500) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_workitem_log_workitem_id` FOREIGN KEY (`workitem_id`) REFERENCES `workitem` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `task_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` varchar(50) NOT NULL,
  `log` varchar(500) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_task_log_task_id` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
