package sql_service

import (
	"goskeleton/app/model"
	"strconv"
)

func UserInit() {
	for i := 1; i <= 100; i++ {
		db := model.GetDBUser()
		if i == 1 {
			db.Exec("SET SESSION sql_require_primary_key=0;")
			db.Exec("CREATE TABLE `users` (\n  `id` int UNSIGNED NOT NULL, \n  `shard_key` int DEFAULT 0 NOT NULL, \n `username` varchar(255) COLLATE utf8_general_ci NOT NULL \n) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;")
			//db.Exec("ALTER TABLE `users`\n  ADD PRIMARY KEY (`id`),\n  ADD UNIQUE KEY `users_mobile_no_unique` (`mobile_no`);")
			db.Exec("ALTER TABLE `users`\n  ADD PRIMARY KEY (`id`),\n ADD UNIQUE KEY `users_username_unique` (`username`);")
		}
		count := strconv.Itoa(i)
		tableName := "users_"+count
		db.Exec("CREATE TABLE `"+tableName+"` (\n  `id` int UNSIGNED NOT NULL,\n  `username` varchar(255) COLLATE utf8_general_ci NOT NULL,\n  `email` varchar(255) COLLATE utf8_general_ci DEFAULT NULL,\n `mobile_no` varchar(255) COLLATE utf8_general_ci DEFAULT NULL,\n `password` varchar(255) COLLATE utf8_general_ci NOT NULL,\n `secondary_password` varchar(255) COLLATE utf8_general_ci NOT NULL,\n  `created_at` timestamp NULL DEFAULT NULL DEFAULT CURRENT_TIMESTAMP,\n  `updated_at` timestamp NULL DEFAULT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP\n) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;")
		db.Exec("ALTER TABLE `"+tableName+"`\n  ADD PRIMARY KEY (`id`),\n  ADD UNIQUE KEY `users_username_unique` (`username`);")
	}
}

func MenuInit() {
	for i := 1; i <= 100; i++ {
		db := model.GetDBStore()
		count := strconv.Itoa(i)
		tableName := "store_menu_"+count
		db.Exec("CREATE TABLE `"+tableName+"` (`id` bigint(20) UNSIGNED NOT NULL PRIMARY KEY,`serial_no` varchar(255) CHARACTER SET utf8 NOT NULL,`category` varchar(255) CHARACTER SET utf8 NOT NULL,`product_name` varchar(255) CHARACTER SET utf8 NOT NULL,`product_desc` varchar(255) CHARACTER SET utf8 DEFAULT NULL,`file_url` varchar(255) CHARACTER SET utf8 DEFAULT NULL,`unit_price` decimal(13,2) NOT NULL DEFAULT '0.00',`check_qty` tinyint(4) NOT NULL DEFAULT '0',`qty` varchar(255) CHARACTER SET utf8 NOT NULL,`status` varchar(255) CHARACTER SET utf8 NOT NULL,`created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,`updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;")
		db.Exec("ALTER TABLE `"+tableName+"` MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;")
		db.Exec("ALTER TABLE `"+tableName+"` ADD `b_choices` tinyint(4) NOT NULL DEFAULT '0' AFTER `qty`;")
	}
}

func MenuVarietyInit() {
	for i := 1; i <= 100; i++ {
		db := model.GetDBStore()
		count := strconv.Itoa(i)
		tableName := "store_menu_variety_"+count
		db.Exec("SET SESSION sql_require_primary_key=0;")
		db.Exec("CREATE TABLE `"+tableName+"` (`menu_id` int(11) NOT NULL PRIMARY KEY,`variety` text,`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,`updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	}
}