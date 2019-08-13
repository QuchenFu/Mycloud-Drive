CREATE TABLE  IF NOT EXISTS `fileinfo`(
 `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '文件编号',
 `file_md5` CHAR(40) NOT NULL DEFAULT'' COMMENT '文件md5',
 `file_name` VARCHAR(256) NOT NULL DEFAULT'' COMMENT '文件名',
 `directory` VARCHAR(1024) NOT NULL DEFAULT'' COMMENT '文件路径',
 `full_path` VARCHAR(1024) NOT NULL DEFAULT'' COMMENT '文件路径',
 `file_size` CHAR(40) NOT NULL DEFAULT'0' COMMENT '文件大小',
 PRIMARY KEY (`id`),
 UNIQUE KEY `full_path` (`full_path`)
)ENGINE = InnoDB
 AUTO_INCREMENT = 1
 DEFAULT CHARACTER SET = utf8
 COLLATE = utf8_general_ci
 COMMENT = '文件表';
 
