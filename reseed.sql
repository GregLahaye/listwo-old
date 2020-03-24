USE `mock`;

SET FOREIGN_KEY_CHECKS=0;

DELETE FROM `User`;
DELETE FROM `List`;
DELETE FROM `Column`;
DELETE FROM `Item`;
DELETE FROM `Session`;

-- John
INSERT INTO `User` (`uuid`, `email`, `password`) VALUES (UUID(), 'john@example.com', UNHEX('24326124313024532E52516B32645570414E3476427349666536347575393567576931344E7935326D55435346486348756F4D554279754C72683161'));
INSERT INTO `Session` (`access_token`, `user_id`) VALUES ('@uth0r1z@t10n_t0k3n', (SELECT `id` FROM `User` WHERE `email` = 'john@example.com'));
INSERT INTO `List` (`uuid`, `title`, `user_id`) VALUES ('schoolwork-uuid', 'Schoolwork', (SELECT `id` FROM `User` WHERE `email` = 'john@example.com'));
INSERT INTO `List` (`uuid`, `title`, `user_id`) VALUES ('housework-uuid', 'Housework', (SELECT `id` FROM `User` WHERE `email` = 'john@example.com'));
INSERT INTO `Column` (`uuid`, `title`, `list_id`) VALUES ('kitchen-uuid', 'Kitchen', (SELECT `id` FROM `List` WHERE `uuid` = 'housework-uuid'));
INSERT INTO `Column` (`uuid`, `title`, `list_id`) VALUES ('bathroom-uuid', 'Bathroom', (SELECT `id` FROM `List` WHERE `uuid` = 'housework-uuid'));
INSERT INTO `Column` (`uuid`, `title`, `list_id`) VALUES ('laundry-uuid', 'Laundry', (SELECT `id` FROM `List` WHERE `uuid` = 'housework-uuid'));

-- Emily
INSERT INTO `User` (`uuid`, `email`, `password`) VALUES (UUID(), 'emily@example.com', UNHEX('24326124313024532E52516B32645570414E3476427349666536347575393567576931344E7935326D55435346486348756F4D554279754C72683161'));
INSERT INTO `List` (`uuid`, `title`, `user_id`) VALUES ('gardening-uuid', 'Gardening', (SELECT `id` FROM `User` WHERE `email` = 'emily@example.com'));


SET FOREIGN_KEY_CHECKS=1;
