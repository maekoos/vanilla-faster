CREATE TABLE
  public.categories (
    slug text NOT NULL,
    name text NOT NULL,
    collection_id integer NOT NULL,
    image_url text NULL
  );

ALTER TABLE
  public.categories
ADD
  CONSTRAINT categories_pkey PRIMARY KEY (slug);


CREATE TABLE
  public.collections (
    id serial NOT NULL,
    name text NOT NULL,
    slug text NOT NULL
  );

ALTER TABLE
  public.collections
ADD
  CONSTRAINT collections_pkey PRIMARY KEY (id);


CREATE TABLE
  public.products (
    slug text NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    price numeric NOT NULL,
    subcategory_slug text NOT NULL,
    image_url text NULL
  );

ALTER TABLE
  public.products
ADD
  CONSTRAINT products_pkey PRIMARY KEY (slug);



CREATE TABLE
  public.subcategories (
    slug text NOT NULL,
    name text NOT NULL,
    subcollection_id integer NOT NULL,
    image_url text NULL
  );

ALTER TABLE
  public.subcategories
ADD
  CONSTRAINT subcategories_pkey PRIMARY KEY (slug);


CREATE TABLE
  public.subcollections (
    id serial NOT NULL,
    name text NOT NULL,
    category_slug text NOT NULL
  );

ALTER TABLE
  public.subcollections
ADD
  CONSTRAINT subcollections_pkey PRIMARY KEY (id);