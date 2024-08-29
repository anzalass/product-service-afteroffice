
CREATE TABLE IF NOT EXISTS kategori
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    product_id uuid NOT NULL,
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT kategori_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS product
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    shop_id uuid NOT NULL,
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    description character varying(255) COLLATE pg_catalog."default" NOT NULL,
    harga text COLLATE pg_catalog."default" NOT NULL,
    stok integer NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone,
    penilaian integer,
    merek character varying(255) COLLATE pg_catalog."default",
    CONSTRAINT product_pkey PRIMARY KEY (id)
);



CREATE TABLE IF NOT EXISTS shops
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    description character varying(255) COLLATE pg_catalog."default" NOT NULL,
    terms text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone,
    terjual character varying(255) COLLATE pg_catalog."default",
    CONSTRAINT shops_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS kategori
    ADD CONSTRAINT kategori_product_id_fkey FOREIGN KEY (product_id)
    REFERENCES product (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS product
    ADD CONSTRAINT fk_shop_id FOREIGN KEY (shop_id)
    REFERENCES shops (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;

