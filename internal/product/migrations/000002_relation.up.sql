CREATE TABLE IF NOT EXISTS public.nrb_related_product
(
    nrb_amount_mv integer NOT NULL,
    nrb_product_sku_id uuid,
    nrb_product_mv_id uuid,
    id uuid NOT NULL,
    CONSTRAINT pk_nrb_related_product PRIMARY KEY (id),
    CONSTRAINT fk_nrb_related_product_nrb_product_mv_id_product FOREIGN KEY (nrb_product_mv_id)
        REFERENCES public.product (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT fk_nrb_related_product_nrb_product_sku_id_product FOREIGN KEY (nrb_product_sku_id)
        REFERENCES public.product (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS ix_nrb_related_product_nrb_product_mv_id
    ON public.nrb_related_product USING btree
    (nrb_product_mv_id ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS ix_nrb_related_product_nrb_product_sku_id
    ON public.nrb_related_product USING btree
    (nrb_product_sku_id ASC NULLS LAST)
    TABLESPACE pg_default;