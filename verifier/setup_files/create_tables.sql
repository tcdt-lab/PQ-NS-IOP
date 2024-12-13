DROP TABLE gateways;
CREATE TABLE gateways (
    id INT PRIMARY KEY AUTO_INCREMENT,
    ip VARCHAR(30) NOT NULL,
    port VARCHAR(6) NOT NULL,
    public_key_kem TEXT NOT NULL,
    public_key_sig TEXT NOT NULL,
    kem_scheme VARCHAR(30) NOT NULL,
    sig_scheme VARCHAR(30) NOT NULL,
    ticket TEXT,
    symmetric_key TEXT
);

DROP TABLE verifiers;
CREATE TABLE verifiers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    ip VARCHAR(30) NOT NULL,
    port VARCHAR(6) NOT NULL,
    public_key_sig TEXT NOT NULL,
    sig_scheme VARCHAR(30) NOT NULL,
    symmetric_key TEXT,
    trust_Score double,
    is_in_committee BOOLEAN
);
DROP TABLE verifier_user;
CREATE TABLE verifier_user (
    id INT PRIMARY KEY AUTO_INCREMENT,
    salt VARCHAR(30) NOT NULL,
    password VARCHAR(64) NOT NULL,
    secret_key TEXT NOT NULL,
    public_key TEXT NOT NULL
);