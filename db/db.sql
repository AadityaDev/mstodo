create database testdb;

create table testdb.todo (
	Id bigint unsigned auto_increment,
    Description varchar(255),
    Completed smallint,
    primary key (Id)
);

INSERT INTO `testdb`.`todo` (`Description`) VALUES ('fd');
