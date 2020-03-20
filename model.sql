-- MySQL Script generated by MySQL Workbench
-- Sat Mar 21 07:34:43 2020
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema listwo
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema listwo
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `listwo` DEFAULT CHARACTER SET utf8 ;
USE `listwo` ;

-- -----------------------------------------------------
-- Table `listwo`.`User`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `listwo`.`User` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `uuid` CHAR(36) NOT NULL,
  `email` VARCHAR(320) NOT NULL,
  `password` BINARY(60) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  UNIQUE INDEX `uuid_UNIQUE` (`uuid` ASC) VISIBLE,
  UNIQUE INDEX `email_UNIQUE` (`email` ASC) VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `listwo`.`List`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `listwo`.`List` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `uuid` CHAR(36) NOT NULL,
  `title` VARCHAR(64) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) INVISIBLE,
  UNIQUE INDEX `uuid_UNIQUE` (`uuid` ASC) INVISIBLE,
  INDEX `list_owner_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `list_owner`
    FOREIGN KEY (`user_id`)
    REFERENCES `listwo`.`User` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `listwo`.`Column`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `listwo`.`Column` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `uuid` CHAR(36) NOT NULL,
  `title` VARCHAR(64) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `list_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  UNIQUE INDEX `uuid_UNIQUE` (`uuid` ASC) VISIBLE,
  INDEX `column_owner_idx` (`list_id` ASC) VISIBLE,
  CONSTRAINT `column_owner`
    FOREIGN KEY (`list_id`)
    REFERENCES `listwo`.`List` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `listwo`.`Item`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `listwo`.`Item` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `uuid` CHAR(36) NOT NULL,
  `position` INT NOT NULL,
  `title` VARCHAR(64) NOT NULL,
  `description` VARCHAR(64) NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `column_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  UNIQUE INDEX `uuid_UNIQUE` (`uuid` ASC) VISIBLE,
  INDEX `item_owner_idx` (`column_id` ASC) VISIBLE,
  CONSTRAINT `item_owner`
    FOREIGN KEY (`column_id`)
    REFERENCES `listwo`.`Column` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `listwo`.`Session`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `listwo`.`Session` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `access_token` VARCHAR(1024) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  UNIQUE INDEX `access_token_UNIQUE` (`access_token` ASC) VISIBLE,
  INDEX `session_user_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `session_user`
    FOREIGN KEY (`user_id`)
    REFERENCES `listwo`.`User` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
