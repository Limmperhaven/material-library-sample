CREATE TABLE IF NOT EXISTS public.material_types
(
    id        bigint GENERATED BY DEFAULT AS IDENTITY
        PRIMARY KEY,
    name      text NOT NULL,
    image_url text NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS public.materials
(
    id                   bigint GENERATED BY DEFAULT AS IDENTITY
        PRIMARY KEY,
    title                text                     NOT NULL,
    subject_id           bigint                   NOT NULL,
    difficultcy_level_id bigint                   NOT NULL,
    type_id              bigint                   NOT NULL REFERENCES public.material_types (id),
    storage_key          text,
    url                  text                     NOT NULL,
    checksum_sha256      bytea,
    size                 bigint,
    created_at           timestamp with time zone NOT NULL default now(),
    updated_at           timestamp with time zone NOT NULL default now(),
    deleted_at           timestamp with time zone
);

INSERT INTO public.material_types (name)
VALUES ('PDF'),
       ('Документ'),
       ('Видеозапись'),
       ('Аудиофайл'),
       ('Ссылка на интернет-ресурс'),
       ('Другое');