CREATE TABLE IF NOT EXISTS schemaName.tag
(
    id       bigserial PRIMARY KEY,
    created  timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    name     varchar(255)
    );

DO
$$
BEGIN
        IF NOT EXISTS (SELECT 1
                       FROM pg_trigger
                       WHERE tgname = 'update_tag_modified_at_schemaName') THEN
CREATE TRIGGER update_tag_modified_at_schemaName
    BEFORE UPDATE
    ON schemaName.tag
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();
END IF;
END
$$;

CREATE TABLE IF NOT EXISTS schemaName.blog
(
    id                bigserial PRIMARY KEY,
    created           timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified          timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    body              text NOT NULL,
    image             text,
    slug              text NOT NULL
);

CREATE TABLE IF NOT EXISTS schemaName.blog_tag
(
    blog_id bigint,
    tag_id  bigint,
    FOREIGN KEY (blog_id) REFERENCES schemaName.blog(id),
    FOREIGN KEY (tag_id) REFERENCES schemaName.tag(id)
    );

DO
$$
BEGIN
        IF NOT EXISTS (SELECT 1
                       FROM pg_trigger
                       WHERE tgname = 'update_blog_modified_at_schemaName') THEN
CREATE TRIGGER update_blog_modified_at_schemaName
    BEFORE UPDATE
    ON schemaName.blog
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();
END IF;
END
$$;
