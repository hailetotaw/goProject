use goBotDB;

CREATE TABLE CommandList (
    id int AUTO_INCREMENT,
    command varchar(30),
    CONSTRAINT CommandList_pk PRIMARY KEY (id)
);

CREATE TABLE AmharicWords (
    id int AUTO_INCREMENT,
    commandId int,
    word varchar(500),
    CONSTRAINT AmharicWords_pk PRIMARY KEY (id)
);

INSERT INTO CommandList(command) values ('@learn_amharic numbers'),
('@learn_amharic days'),
('@learn_amharic months'),
('@learn_amharic greetings');

INSERT INTO AmharicWords(commandId, word) VALUES
(1, 'Ethipian Numbers'),
(2, 'Ethiopian Days'),
(3, 'Ethiopian Months'),
(4, 'Ethiopian Greetings');