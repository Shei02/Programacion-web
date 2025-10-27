CREATE TABLE pelicula(
    IDP SERIAL PRIMARY KEY, --EL SERIAL LE DA UN UNICO VALOR DE ID
    titulo varchar(255) NOT NULL,
    duracion varchar(255) NOT NULL,
    director varchar(255) NOT NULL,
    actores varchar(255) NOT NULL,
    edadMin int NOT NULL,
    sinopsis varchar(255) NOT NULL,
    anioEstr int NOT NULL
);

CREATE TABLE usuario(
    IDU SERIAL PRIMARY KEY,
    nomUsu varchar(255) NOT NULL,
    contrasenia varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    fechaNac date NOT NULL
);

CREATE TABLE mira (
    IDP INT NOT NULL,
    IDU INT NOT NULL,
    gustoOno VARCHAR(255) NOT NULL,
    calif INT NOT NULL CHECK (calif BETWEEN 1 AND 5),
    PRIMARY KEY (IDP, IDU),
    FOREIGN KEY (IDP) REFERENCES pelicula(IDP),
    FOREIGN KEY (IDU) REFERENCES usuario(IDU)
);


