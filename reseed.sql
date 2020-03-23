USE `mock`;

SET FOREIGN_KEY_CHECKS=0;

DELETE FROM `User`;
DELETE FROM `List`;
DELETE FROM `Column`;
DELETE FROM `Item`;
DELETE FROM `Session`;

INSERT INTO `User` (`uuid`, `email`, `password`) VALUES (UUID(), 'john@example.com', UNHEX('24326124313024532E52516B32645570414E3476427349666536347575393567576931344E7935326D55435346486348756F4D554279754C72683161'));

SET FOREIGN_KEY_CHECKS=1;