INSERT INTO schemaName.permissions (id, created, modified, name, description)
VALUES (1, '2021-01-01 00:00:00', '2021-01-01 00:00:00', 'blog_write', 'Write blog posts');

INSERT INTO schemaName.permissions (id, created, modified, name, description)
VALUES (2, '2021-01-01 00:00:00', '2021-01-01 00:00:00', 'blog_read', 'Read blog posts');

INSERT INTO schemaName.roles (id, created, modified, name, description)
VALUES (1, '2021-01-01 00:00:00', '2021-01-01 00:00:00', 'admin', 'Admin role');

INSERT INTO schemaName.roles (id, created, modified, name, description)
VALUES (2, '2021-01-01 00:00:00', '2021-01-01 00:00:00', 'user', 'User role');

INSERT INTO schemaName.role_permissions (role_id, permission_id)
VALUES (1, 1);

INSERT INTO schemaName.role_permissions (role_id, permission_id)
VALUES (1, 2);

INSERT INTO schemaName.role_permissions (role_id, permission_id)
VALUES (2, 2);
