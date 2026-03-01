CREATE TABLE IF NOT EXISTS public.product
(
    smr_fnrec character varying(50) COLLATE pg_catalog."default" NOT NULL,
    name character varying(250) COLLATE pg_catalog."default" NOT NULL,
    smr_is_service boolean NOT NULL DEFAULT FALSE,
    is_archive boolean NOT NULL DEFAULT FALSE,
    ask_pack_volume numeric(18,2) NOT NULL DEFAULT 0,
    ask_weight numeric(18,2) NOT NULL DEFAULT 0,
    ask_pack_length numeric(18,2) NOT NULL DEFAULT 0,
    ask_pack_width numeric(18,2) NOT NULL DEFAULT 0,
    ask_pack_height numeric(18,2) NOT NULL DEFAULT 0,
    smr_series_id uuid DEFAULT NULL,
    type_id uuid DEFAULT NULL,
    nrb_status_sku_id uuid DEFAULT NULL,
    nrb_type_production_id uuid DEFAULT NULL,
    category_id uuid DEFAULT NULL,
    nrb_account_product_id uuid DEFAULT NULL,
    ask_non_standart_category_id uuid DEFAULT NULL,
    nrb_integration_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    nrb_category_sku_id uuid DEFAULT NULL,
    nrb_count_mv integer NOT NULL DEFAULT 0,
    id uuid NOT NULL,
    smr_product_group_flag_id uuid DEFAULT NULL,
    CONSTRAINT pk_product PRIMARY KEY (id)
)


TABLESPACE pg_default;
-- Index: ix_product_ask_non_standart_category_id

-- DROP INDEX IF EXISTS public.ix_product_ask_non_standart_category_id;

CREATE INDEX IF NOT EXISTS ix_product_ask_non_standart_category_id
    ON public.product USING btree
    (ask_non_standart_category_id ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: ix_product_category_id

-- DROP INDEX IF EXISTS public.ix_product_category_id;

CREATE INDEX IF NOT EXISTS ix_product_category_id
    ON public.product USING btree
    (category_id ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: ix_product_nrb_account_product_id

-- DROP INDEX IF EXISTS public.ix_product_nrb_account_product_id;

CREATE INDEX IF NOT EXISTS ix_product_nrb_account_product_id
    ON public.product USING btree
    (nrb_account_product_id ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: ix_product_nrb_category_sku_id

-- DROP INDEX IF EXISTS public.ix_product_nrb_category_sku_id;

CREATE INDEX IF NOT EXISTS ix_product_nrb_category_sku_id
    ON public.product USING btree
    (nrb_category_sku_id ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: ix_product_nrb_status_sku_id

-- DROP INDEX IF EXISTS public.ix_product_nrb_status_sku_id;

CREATE INDEX IF NOT EXISTS ix_product_nrb_status_sku_id
    ON public.product USING btree
    (nrb_status_sku_id ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: ix_product_nrb_type_production_id

-- DROP INDEX IF EXISTS public.ix_product_nrb_type_production_id;

CREATE INDEX IF NOT EXISTS ix_product_nrb_type_production_id
    ON public.product USING btree
    (nrb_type_production_id ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: ix_product_smr_fnrec

-- DROP INDEX IF EXISTS public.ix_product_smr_fnrec;

CREATE INDEX IF NOT EXISTS ix_product_smr_fnrec
    ON public.product USING btree
    (smr_fnrec COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: ix_product_smr_series_id

-- DROP INDEX IF EXISTS public.ix_product_smr_series_id;

CREATE INDEX IF NOT EXISTS ix_product_smr_series_id
    ON public.product USING btree
    (smr_series_id ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: ix_product_type_id

-- DROP INDEX IF EXISTS public.ix_product_type_id;

CREATE INDEX IF NOT EXISTS ix_product_type_id
    ON public.product USING btree
    (type_id ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: ux_product_nrb_integration_id

-- DROP INDEX IF EXISTS public.ux_product_nrb_integration_id;

CREATE UNIQUE INDEX IF NOT EXISTS ux_product_nrb_integration_id
    ON public.product USING btree
    (nrb_integration_id COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;