
drop table gateways;
CREATE table gateways (
    id int primary key AUTO_INCREMENT,
    ip varchar(30) not null,
    port varchar (6) not null,
    public_key text not null,
    ticket varchar(50) ,
    symmetric_key text
);

drop table verifiers;
create table verifiers (
    id int primary key auto_increment,
    ip varchar(30) not null,
    port varchar (6) not null,
    public_key text not null,
    symmetric_key text
);

drop table gateway_user;
create table gateway_user (
    id int primary key auto_increment,
    salt varchar(30) not null,
    password varchar(64) not null,
    public_key_dsa text not null,
    secret_key_dsa text not null,
    public_key_kem text not null,
    secret_key_kem text not null,
    dsa_scheme varchar(30) not null,
    kem_scheme varchar(30) not null
);

