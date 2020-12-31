DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` text COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `users` (`name`) VALUES
('まるお'),
('トマト君'),
('マスク侍'),
('腕時計両腕ちゃん'),
('わたしの中の悪魔'),
('眠り姫'),
('べらぼう太郎'),
('ジョン'),
('サンタクロース'),
('先生');