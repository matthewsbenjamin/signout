CREATE TABLE `adults` (
  `email` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `active` tinyint(4) DEFAULT '1',
  `pwd` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`email`)
) ENGINE=InnoDB;

CREATE TABLE `boat_locations` (
  `boat_name` varchar(255) NOT NULL,
  `on_water` tinyint(4) NOT NULL DEFAULT '0',
  `last_on_off` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`boat_name`),
  KEY `boat` (`boat_name`)
) ENGINE=InnoDB;


CREATE TABLE `transactions` (
  `transaction_id` int(11) NOT NULL AUTO_INCREMENT,
  `boat_name` varchar(45) NOT NULL,
  `adult` varchar(255) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `signout` tinyint(4) DEFAULT '0',
  `hazards` varchar(999) DEFAULT NULL,
  `damage` varchar(999) DEFAULT NULL,
  PRIMARY KEY (`transaction_id`),
  KEY `boat` (`boat_name`) /*!80000 INVISIBLE */,
  KEY `individual` (`adult`)
) ENGINE=InnoDB;
