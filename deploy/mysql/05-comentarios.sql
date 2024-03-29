USE vigo360;

CREATE TABLE comentarios(
	id varchar(13) NOT NULL,
    publicacion_id varchar(40) NOT NULL,
    padre_id varchar(13),
    
    nombre varchar(40) NOT NULL,
    email_hash char(65) NOT NULL, -- SHA256 del correo electrónico usado
    es_autor boolean NOT NULL DEFAULT false, -- Si es autor en Vigo360
    autor_original boolean NOT NULL DEFAULT false, -- Si es el autor original de la publicación

    contenido text(500) NOT NULL,

    fecha_creacion datetime NOT NULL DEFAULT NOW(),
    fecha_moderacion datetime,
    estado enum("pendiente", "aprobado", "rechazado") NOT NULL DEFAULT "pendiente",
    moderador varchar(40) NOT NULL,

    PRIMARY KEY(id),
    FOREIGN KEY (publicacion_id) REFERENCES publicaciones(id),
    FOREIGN KEY (padre_id) REFERENCES comentarios(id),
    FOREIGN KEY (moderador) REFERENCES autores(id),
    CHECK(nombre != ""),
    CHECK(contenido != "")
);