DELIMITER //

CREATE FUNCTION fn_tag_name_to_slug (input_string VARCHAR(1000))
RETURNS VARCHAR(1000)
DETERMINISTIC
BEGIN
    DECLARE output_string VARCHAR(1000);

    -- Handle potential NULL input
    IF input_string IS NULL THEN
        RETURN NULL;
    END IF;

    SET output_string = input_string;

    -- Replace lowercase acute accents
    SET output_string = REPLACE(output_string, 'á', 'a');
    SET output_string = REPLACE(output_string, 'é', 'e');
    SET output_string = REPLACE(output_string, 'í', 'i');
    SET output_string = REPLACE(output_string, 'ó', 'o');
    SET output_string = REPLACE(output_string, 'ú', 'u');

    -- Replace uppercase acute accents
    SET output_string = REPLACE(output_string, 'Á', 'A');
    SET output_string = REPLACE(output_string, 'É', 'E');
    SET output_string = REPLACE(output_string, 'Í', 'I');
    SET output_string = REPLACE(output_string, 'Ó', 'O');
    SET output_string = REPLACE(output_string, 'Ú', 'U');

    -- Add more REPLACE lines here if you need to handle other characters (ñ, ç, ü, etc.)
    SET output_string = REPLACE(output_string, 'ñ', 'n');
    SET output_string = REPLACE(output_string, 'Ñ', 'N');
    SET output_string = REPLACE(output_string, 'ç', 'c');
    SET output_string = REPLACE(output_string, 'Ç', 'C');
    SET output_string = REPLACE(output_string, 'ü', 'u');
    SET output_string = REPLACE(output_string, 'Ü', 'U');

    -- Lowercase, hyphenate, and replace all non-alphanumeric characters with empty string
    SET output_string = LOWER(output_string);
    SET output_string = REPLACE(output_string, ' ', '-');

    RETURN output_string;
END //

-- Change the delimiter back to the default
DELIMITER ;

ALTER TABLE tags ADD COLUMN slug VARCHAR(50) NOT NULL AFTER id;
UPDATE tags SET slug = fn_tag_name_to_slug(name);
ALTER TABLE tags ADD UNIQUE INDEX idx_slug (slug);
