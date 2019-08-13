CREATE TABLE  IF NOT EXISTS `filelife`(
 `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '文件编号',
 `file_md5` CHAR(40) NOT NULL DEFAULT'' COMMENT '文件md5',
 `file_life` INT(11) NOT NULL DEFAULT 0 COMMENT '文件life', 
 PRIMARY KEY (`id`),
 UNIQUE KEY `file_md5` (`file_md5`)
)ENGINE = InnoDB
 AUTO_INCREMENT = 1
 DEFAULT CHARACTER SET = utf8
 COLLATE = utf8_general_ci
 COMMENT = '文件寿命表';
 