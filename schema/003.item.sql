Use letscrum;

CREATE TABLE `epic` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` bigint(20) unsigned NOT NULL,
  `sprint_id` bigint(20) unsigned DEFAULT NULL,
  `title` varchar(50) NOT NULL,
  `assign_to` bigint(20) unsigned DEFAULT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_epic_project_id` FOREIGN KEY (`project_id`) REFERENCES `project` (`id`),
  CONSTRAINT `fk_epic_sprint_id` FOREIGN KEY (`sprint_id`) REFERENCES `sprint` (`id`),
  CONSTRAINT `fk_epic_assign_to` FOREIGN KEY (`assign_to`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_epic_created_by` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `feature` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` bigint(20) unsigned NOT NULL,
  `sprint_id` bigint(20) unsigned DEFAULT NULL,
  `epic_id` bigint(20) unsigned DEFAULT NULL,
  `title` varchar(50) NOT NULL,
  `assign_to` bigint(20) unsigned DEFAULT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_feature_project_id` FOREIGN KEY (`project_id`) REFERENCES `project` (`id`),
  CONSTRAINT `fk_feature_sprint_id` FOREIGN KEY (`sprint_id`) REFERENCES `sprint` (`id`),
  CONSTRAINT `fk_feature_epic_id` FOREIGN KEY (`epic_id`) REFERENCES `epic` (`id`),
  CONSTRAINT `fk_feature_assign_to` FOREIGN KEY (`assign_to`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_feature_created_by` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `work_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` bigint(20) unsigned NOT NULL,
  `sprint_id` bigint(20) unsigned DEFAULT NULL,
  `feature_id` bigint(20) unsigned DEFAULT NULL,
  `title` varchar(50) NOT NULL,
  `type` ENUM('Backlog', 'Bug') NOT NULL,
  `description` varchar(500) DEFAULT NULL,
  `status` ENUM('New', 'Approved', 'Committed', 'Done', 'Removed') NOT NULL,
  `assign_to` bigint(20) unsigned DEFAULT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_work_item_project_id` FOREIGN KEY (`project_id`) REFERENCES `project` (`id`),
  CONSTRAINT `fk_work_item_sprint_id` FOREIGN KEY (`sprint_id`) REFERENCES `sprint` (`id`),
  CONSTRAINT `fk_work_item_feature_id` FOREIGN KEY (`feature_id`) REFERENCES `feature` (`id`),
  CONSTRAINT `fk_work_item_assign_to` FOREIGN KEY (`assign_to`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_work_item_created_by` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `task` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` bigint(20) unsigned NOT NULL,
  `sprint_id` bigint(20) unsigned DEFAULT NULL,
  `work_item_id` bigint(20) unsigned DEFAULT NULL,
  `title` varchar(50) NOT NULL,
  `description` varchar(500) DEFAULT NULL,
  `status` ENUM('To Do', 'In Progress', 'Done', 'Removed') NOT NULL,
  `assign_to` bigint(20) unsigned DEFAULT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `remaining` float unsigned DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_task_project_id` FOREIGN KEY (`project_id`) REFERENCES `project` (`id`),
  CONSTRAINT `fk_task_sprint_id` FOREIGN KEY (`sprint_id`) REFERENCES `sprint` (`id`),
  CONSTRAINT `fk_task_work_item_id` FOREIGN KEY (`work_item_id`) REFERENCES `work_item` (`id`),
  CONSTRAINT `fk_task_assign_to` FOREIGN KEY (`assign_to`) REFERENCES `user` (`id`),
  CONSTRAINT `fk_task_created_by` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `work_item_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `work_item_id` varchar(50) NOT NULL,
  `log` varchar(500) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_work_item_log_workitem_id` FOREIGN KEY (`work_item_id`) REFERENCES `work_item` (`id`)
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
