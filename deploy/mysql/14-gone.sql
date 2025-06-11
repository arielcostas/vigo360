ALTER TABLE publicaciones ADD COLUMN gone_at DATETIME NULL AFTER legally_retired_at;
ALTER TABLE autores ADD COLUMN gone_at DATETIME NULL AFTER biografia;
