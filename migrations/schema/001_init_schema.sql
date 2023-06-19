CREATE TABLE auth_role (
                           id SERIAL PRIMARY KEY,
                           name character varying NOT NULL,
                           enabled bool not null default false,
                           created_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                           created_by bigint NOT NULL,
                           updated_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                           updated_by bigint NOT NULL,
                           deleted_at timestamp with time zone,
                           deleted_by bigint
);

CREATE TABLE auth_scope (
                            id SERIAL PRIMARY KEY,
                            name character varying NOT NULL
);

create table auth_scope_endpoint (
                                    id serial primary key,
                                    name character varying not null,
                                    bit int2 not null,
                                    scope_id bigint NOT NULL REFERENCES auth_scope
);

CREATE TABLE auth_user (
                           id SERIAL PRIMARY KEY,
                           name text,
                           email character varying(255),
                           user_name text NOT NULL,
                           password text NOT NULL,
                           role_id bigint REFERENCES auth_role,
                           created_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                           created_by bigint,
                           updated_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                           updated_by bigint,
                           deleted_at timestamp with time zone,
                           deleted_by bigint
);

CREATE TABLE auth_role_permission (
                                    id SERIAL PRIMARY KEY,
                                    role_id bigint NOT NULL REFERENCES auth_role,
                                    scope_id bigint NOT NULL REFERENCES auth_scope,
                                    bit int2 not null default 4,
                                    created_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                                    created_by bigint NOT NULL REFERENCES auth_user,
                                    updated_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                                    updated_by bigint NOT NULL REFERENCES auth_user,
                                    deleted_at timestamp with time zone,
                                    deleted_by bigint REFERENCES auth_user
);

CREATE TABLE auth_role_grantee (
                              id SERIAL PRIMARY KEY,
                              grantor_role_id bigint NOT NULL REFERENCES auth_role,
                              grantee_role_id bigint NOT NULL REFERENCES auth_role,
                              scope_id bigint NOT NULL REFERENCES auth_scope,
                              bit int2 not null default 4,
                              created_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                              created_by bigint NOT NULL REFERENCES auth_user,
                              updated_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                              updated_by bigint NOT NULL REFERENCES auth_user,
                              deleted_at timestamp with time zone,
                              deleted_by bigint REFERENCES auth_user
);

CREATE TABLE auth_user_session (
                                   id SERIAL PRIMARY KEY,
                                   user_id bigint NOT NULL REFERENCES auth_user,
                                   network_ip inet,
                                   access_token character varying,
                                   created_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                                   expires_at timestamp with time zone
);

CREATE TABLE task_source (
                             id SERIAL PRIMARY KEY,
                             name text NOT NULL,
                             created_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                             created_by bigint NOT NULL REFERENCES auth_user,
                             updated_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                             updated_by bigint NOT NULL REFERENCES auth_user,
                             deleted_at timestamp with time zone,
                             deleted_by bigint REFERENCES auth_user
);

CREATE TABLE task_type (
                           id SERIAL PRIMARY KEY,
                           name text NOT NULL,
                           created_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                           created_by bigint NOT NULL REFERENCES auth_user,
                           updated_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                           updated_by bigint NOT NULL REFERENCES auth_user,
                           deleted_at timestamp with time zone,
                           deleted_by bigint REFERENCES auth_user
);

CREATE TABLE task_customer (
                               id SERIAL PRIMARY KEY,
                               name text NOT NULL,
                               created_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                               created_by bigint NOT NULL REFERENCES auth_user,
                               updated_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                               updated_by bigint NOT NULL REFERENCES auth_user,
                               deleted_at timestamp with time zone,
                               deleted_by bigint REFERENCES auth_user
);

CREATE TABLE task_list (
                           id SERIAL PRIMARY KEY,
                           name text NOT NULL,
                           description text,
                           link text,
                           created_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                           created_by bigint NOT NULL REFERENCES auth_user,
                           updated_at timestamp with time zone DEFAULT LOCALTIMESTAMP NOT NULL,
                           updated_by bigint NOT NULL REFERENCES auth_user,
                           deleted_at timestamp with time zone,
                           deleted_by bigint REFERENCES auth_user,
                           type_id bigint NOT NULL REFERENCES task_type,
                           source_id bigint NOT NULL REFERENCES task_source,
                           customer_id bigint NOT NULL REFERENCES task_customer
);

---- create above / drop below ----

DROP TABLE auth_role;
DROP TABLE auth_scope;
DROP TABLE auth_user;
DROP TABLE auth_role_permission;
DROP TABLE auth_role_grantee;
DROP TABLE auth_user_session;
DROP TABLE task_source;
DROP TABLE task_type;
DROP TABLE task_customer;
DROP TABLE task_list;
