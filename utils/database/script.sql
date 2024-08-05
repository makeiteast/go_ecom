--
-- PostgreSQL database dump
--

-- Dumped from database version 14.4 (Debian 14.4-1.pgdg110+1)
-- Dumped by pg_dump version 14.4 (Debian 14.4-1.pgdg110+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
--
-- Name: statet; Type: TYPE; Schema: public; Owner: market
--

CREATE TYPE public.statet AS ENUM (
    'enabled',
    'disabled',
    'deleted'
);


ALTER TYPE public.statet OWNER TO market;

--
-- Name: userrole; Type: TYPE; Schema: public; Owner: market
--

CREATE TYPE public.userrole AS ENUM (
    'SUPERADMIN',
    'ADMIN',
    'USER'
);


ALTER TYPE public.userrole OWNER TO market;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: brand; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.brand (
    id integer NOT NULL,
    brand_name character varying(50) NOT NULL,
    brand_type character varying(50) NOT NULL,
    brand_icon character varying(100),
    create_ts timestamp without time zone NOT NULL,
    update_ts timestamp without time zone NOT NULL,
    state public.statet NOT NULL,
    version integer
);


ALTER TABLE public.brand OWNER TO market;

--
-- Name: brand_id_seq; Type: SEQUENCE; Schema: public; Owner: market
--

CREATE SEQUENCE public.brand_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.brand_id_seq OWNER TO market;

--
-- Name: brand_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: market
--

ALTER SEQUENCE public.brand_id_seq OWNED BY public.brand.id;


--
-- Name: cart; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.cart (
    user_id integer NOT NULL,
    sku_id bigint NOT NULL,
    quantity integer NOT NULL,
    create_ts timestamp without time zone NOT NULL,
    update_ts timestamp without time zone NOT NULL,
    state public.statet NOT NULL,
    version integer
);


ALTER TABLE public.cart OWNER TO market;

--
-- Name: category; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.category (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    parent integer NOT NULL,
    image character varying(100),
    icon character varying(100),
    create_ts timestamp without time zone NOT NULL,
    update_ts timestamp without time zone NOT NULL,
    state public.statet NOT NULL,
    version integer
);


ALTER TABLE public.category OWNER TO market;

--
-- Name: category_id_seq; Type: SEQUENCE; Schema: public; Owner: market
--

CREATE SEQUENCE public.category_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.category_id_seq OWNER TO market;

--
-- Name: category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: market
--

ALTER SEQUENCE public.category_id_seq OWNED BY public.category.id;


--
-- Name: currency; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.currency (
    id integer NOT NULL,
    value double precision,
    local_value double precision,
    start_at timestamp without time zone NOT NULL,
    end_at timestamp without time zone NOT NULL,
    create_ts timestamp without time zone NOT NULL,
    update_ts timestamp without time zone NOT NULL,
    state public.statet NOT NULL,
    version integer
);


ALTER TABLE public.currency OWNER TO market;

--
-- Name: currency_id_seq; Type: SEQUENCE; Schema: public; Owner: market
--

CREATE SEQUENCE public.currency_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.currency_id_seq OWNER TO market;

--
-- Name: currency_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: market
--

ALTER SEQUENCE public.currency_id_seq OWNED BY public.currency.id;


--
-- Name: discount; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.discount (
    id integer NOT NULL,
    product_id bigint NOT NULL,
    user_id integer NOT NULL,
    percent double precision,
    status boolean NOT NULL,
    start_at timestamp without time zone NOT NULL,
    end_at timestamp without time zone NOT NULL,
    create_ts timestamp without time zone NOT NULL,
    update_ts timestamp without time zone NOT NULL,
    state public.statet NOT NULL,
    version integer
);


ALTER TABLE public.discount OWNER TO market;

--
-- Name: option; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.option (
    id bigint NOT NULL,
    category_id bigint NOT NULL,
    name character varying(50) NOT NULL,
    create_ts timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_ts timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    state public.statet DEFAULT 'enabled'::public.statet NOT NULL,
    version integer DEFAULT 0
);


ALTER TABLE public.option OWNER TO market;

--
-- Name: option_id_seq; Type: SEQUENCE; Schema: public; Owner: market
--

CREATE SEQUENCE public.option_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.option_id_seq OWNER TO market;

--
-- Name: option_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: market
--

ALTER SEQUENCE public.option_id_seq OWNED BY public.option.id;


--
-- Name: option_value; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.option_value (
    id bigint NOT NULL,
    option_id bigint NOT NULL,
    name character varying(50) NOT NULL,
    create_ts timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_ts timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    state public.statet DEFAULT 'enabled'::public.statet NOT NULL,
    version integer DEFAULT 0
);


ALTER TABLE public.option_value OWNER TO market;

--
-- Name: option_value_id_seq; Type: SEQUENCE; Schema: public; Owner: market
--

CREATE SEQUENCE public.option_value_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.option_value_id_seq OWNER TO market;

--
-- Name: option_value_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: market
--

ALTER SEQUENCE public.option_value_id_seq OWNED BY public.option_value.id;


--
-- Name: order_item; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.order_item (
    id integer NOT NULL,
    order_id uuid NOT NULL,
    sku_id integer NOT NULL,
    quantity integer,
    price double precision,
    create_ts timestamp without time zone NOT NULL,
    update_ts timestamp without time zone NOT NULL,
    state public.statet NOT NULL,
    version integer
);


ALTER TABLE public.order_item OWNER TO market;

--
-- Name: order_item_id_seq; Type: SEQUENCE; Schema: public; Owner: market
--

CREATE SEQUENCE public.order_item_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.order_item_id_seq OWNER TO market;

--
-- Name: order_item_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: market
--

ALTER SEQUENCE public.order_item_id_seq OWNED BY public.order_item.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.orders (
    id uuid NOT NULL,
    user_id integer NOT NULL,
    address character varying(100),
    phone character varying(100),
    comment text,
    notes text,
    status character varying(50) NOT NULL,
    create_ts timestamp without time zone NOT NULL,
    update_ts timestamp without time zone NOT NULL,
    state public.statet NOT NULL,
    version integer
);


ALTER TABLE public.orders OWNER TO market;

--
-- Name: product; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.product (
    id bigint NOT NULL,
    product_name character varying(100) NOT NULL,
    description text,
    category_id integer NOT NULL,
    brand_id integer NOT NULL,
    region_id integer NOT NULL,
    create_ts timestamp without time zone NOT NULL,
    update_ts timestamp without time zone NOT NULL,
    state public.statet NOT NULL,
    version integer
);


ALTER TABLE public.product OWNER TO market;

--
-- Name: product_id_seq; Type: SEQUENCE; Schema: public; Owner: market
--

CREATE SEQUENCE public.product_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.product_id_seq OWNER TO market;

--
-- Name: product_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: market
--

ALTER SEQUENCE public.product_id_seq OWNED BY public.product.id;


--
-- Name: region; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.region (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    main_region integer NOT NULL,
    create_ts timestamp without time zone NOT NULL,
    update_ts timestamp without time zone NOT NULL,
    state public.statet NOT NULL,
    version integer
);


ALTER TABLE public.region OWNER TO market;

--
-- Name: region_id_seq; Type: SEQUENCE; Schema: public; Owner: market
--

CREATE SEQUENCE public.region_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.region_id_seq OWNER TO market;

--
-- Name: region_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: market
--

ALTER SEQUENCE public.region_id_seq OWNED BY public.region.id;


--
-- Name: sku; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.sku (
    id bigint NOT NULL,
    product_id bigint NOT NULL,
    sku character varying(30) NOT NULL,
    price double precision,
    quantity integer,
    large_name character varying(300),
    small_name character varying(100),
    thumb_name character varying(100),
    count_viewed integer,
    create_ts timestamp without time zone NOT NULL,
    update_ts timestamp without time zone NOT NULL,
    state public.statet NOT NULL,
    version integer
);


ALTER TABLE public.sku OWNER TO market;

--
-- Name: sku_id_seq; Type: SEQUENCE; Schema: public; Owner: market
--

CREATE SEQUENCE public.sku_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sku_id_seq OWNER TO market;

--
-- Name: sku_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: market
--

ALTER SEQUENCE public.sku_id_seq OWNED BY public.sku.id;


--
-- Name: sku_value; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.sku_value (
    id bigint NOT NULL,
    sku_id integer NOT NULL,
    option_id integer NOT NULL,
    option_value_id integer NOT NULL,
    create_ts timestamp without time zone NOT NULL,
    update_ts timestamp without time zone NOT NULL,
    state public.statet NOT NULL,
    version integer
);


ALTER TABLE public.sku_value OWNER TO market;

--
-- Name: sku_value_id_seq; Type: SEQUENCE; Schema: public; Owner: market
--

CREATE SEQUENCE public.sku_value_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sku_value_id_seq OWNER TO market;

--
-- Name: sku_value_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: market
--

ALTER SEQUENCE public.sku_value_id_seq OWNED BY public.sku_value.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: market
--

CREATE TABLE public.users (
    id integer NOT NULL,
    public_id character varying(50),
    username character varying(20) NOT NULL,
    password character varying(60) NOT NULL,
    email character varying(70),
    phone_number character varying(70),
    address character varying(70),
    photo character varying(70),
    role public.userrole NOT NULL,
    region_id integer,
    parent integer NOT NULL,
    create_ts timestamp without time zone NOT NULL,
    update_ts timestamp without time zone NOT NULL,
    state public.statet NOT NULL,
    version integer
);


ALTER TABLE public.users OWNER TO market;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: market
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO market;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: market
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: brand id; Type: DEFAULT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.brand ALTER COLUMN id SET DEFAULT nextval('public.brand_id_seq'::regclass);


--
-- Name: category id; Type: DEFAULT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.category ALTER COLUMN id SET DEFAULT nextval('public.category_id_seq'::regclass);


--
-- Name: currency id; Type: DEFAULT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.currency ALTER COLUMN id SET DEFAULT nextval('public.currency_id_seq'::regclass);


--
-- Name: option id; Type: DEFAULT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.option ALTER COLUMN id SET DEFAULT nextval('public.option_id_seq'::regclass);


--
-- Name: option_value id; Type: DEFAULT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.option_value ALTER COLUMN id SET DEFAULT nextval('public.option_value_id_seq'::regclass);


--
-- Name: order_item id; Type: DEFAULT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.order_item ALTER COLUMN id SET DEFAULT nextval('public.order_item_id_seq'::regclass);


--
-- Name: product id; Type: DEFAULT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.product ALTER COLUMN id SET DEFAULT nextval('public.product_id_seq'::regclass);


--
-- Name: region id; Type: DEFAULT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.region ALTER COLUMN id SET DEFAULT nextval('public.region_id_seq'::regclass);


--
-- Name: sku id; Type: DEFAULT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.sku ALTER COLUMN id SET DEFAULT nextval('public.sku_id_seq'::regclass);


--
-- Name: sku_value id; Type: DEFAULT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.sku_value ALTER COLUMN id SET DEFAULT nextval('public.sku_value_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: brand; Type: TABLE DATA; Schema: public; Owner: market
--

INSERT INTO public.brand VALUES (1, 'LG', 'Electronic', 'static/brand_image/lg.png', '2022-10-25 10:25:25.3141', '2022-10-25 10:25:25.314106', 'enabled', 0);
INSERT INTO public.brand VALUES (2, 'beko', 'Electronic', 'static/brand_image/beko.png', '2022-10-25 10:25:25.316643', '2022-10-25 10:25:25.316648', 'enabled', 0);
INSERT INTO public.brand VALUES (3, 'Xiaomi', 'Electronic', 'static/brand_image/xiaomi.png', '2022-10-25 10:25:25.318388', '2022-10-25 10:25:25.318392', 'enabled', 0);
INSERT INTO public.brand VALUES (4, 'Samsung', 'Electronic', 'static/brand_image/samsung.png', '2022-10-25 10:25:25.319757', '2022-10-25 10:25:25.319759', 'enabled', 0);
INSERT INTO public.brand VALUES (5, 'Apple', 'Electronic', 'static/brand_image/apple.png', '2022-10-25 10:25:25.321025', '2022-10-25 10:25:25.321027', 'enabled', 0);
INSERT INTO public.brand VALUES (6, 'Asus', 'Electronic', 'static/brand_image/asus.png', '2022-10-25 10:25:25.322477', '2022-10-25 10:25:25.322478', 'enabled', 0);
INSERT INTO public.brand VALUES (7, 'Acer', 'Computer', 'static/brand_image/acer.png', '2022-10-25 10:25:25.324149', '2022-10-25 10:25:25.324154', 'enabled', 0);
INSERT INTO public.brand VALUES (8, 'Intel', 'Computer', 'static/brand_image/intel.png', '2022-10-25 10:25:25.326429', '2022-10-25 10:25:25.326434', 'enabled', 0);
INSERT INTO public.brand VALUES (9, 'Nvidia', 'Computer', 'static/brand_image/nvidia.png', '2022-10-25 10:25:25.328408', '2022-10-25 10:25:25.328413', 'enabled', 0);
INSERT INTO public.brand VALUES (10, 'AMD', 'Computer', 'static/brand_image/amd.png', '2022-10-25 10:25:25.330172', '2022-10-25 10:25:25.330176', 'enabled', 0);
INSERT INTO public.brand VALUES (11, 'Lenovo', 'Electronic', 'static/brand_image/lenovo.png', '2022-10-25 10:25:25.331905', '2022-10-25 10:25:25.331909', 'enabled', 0);
INSERT INTO public.brand VALUES (12, 'Philips', 'Electronic', 'static/brand_image/philips.png', '2022-10-25 10:25:25.333658', '2022-10-25 10:25:25.333662', 'enabled', 0);
INSERT INTO public.brand VALUES (13, 'Panasonic', 'Electronic', 'static/brand_image/panasonic.png', '2022-10-25 10:25:25.33533', '2022-10-25 10:25:25.335334', 'enabled', 0);
INSERT INTO public.brand VALUES (14, 'Sony', 'Electronic', 'static/brand_image/sony.png', '2022-10-25 10:25:25.337013', '2022-10-25 10:25:25.337017', 'enabled', 0);
INSERT INTO public.brand VALUES (15, 'LC Waikiki', 'Clothes', 'static/brand_image/lcwaikiki.png', '2022-10-25 10:25:25.338701', '2022-10-25 10:25:25.338704', 'enabled', 0);
INSERT INTO public.brand VALUES (16, 'koton', 'Clothes', 'static/brand_image/koton.png', '2022-10-25 10:25:25.340372', '2022-10-25 10:25:25.340375', 'enabled', 0);
INSERT INTO public.brand VALUES (17, 'Pierre Cardin', 'Clothes', 'static/brand_image/pierrecardin.png', '2022-10-25 10:25:25.341833', '2022-10-25 10:25:25.341835', 'enabled', 0);


--
-- Data for Name: cart; Type: TABLE DATA; Schema: public; Owner: market
--

INSERT INTO public.cart VALUES (4, 2, 1, '2022-10-25 10:25:25.383956', '2022-10-25 10:25:25.383959', 'enabled', 0);
INSERT INTO public.cart VALUES (4, 3, 1, '2022-10-25 10:25:25.385741', '2022-10-25 10:25:25.385743', 'enabled', 0);
INSERT INTO public.cart VALUES (4, 4, 1, '2022-10-25 10:25:25.387053', '2022-10-25 10:25:25.387056', 'enabled', 0);
INSERT INTO public.cart VALUES (4, 9, 3, '2022-10-25 10:25:25.388367', '2022-10-25 10:25:25.388369', 'enabled', 0);
INSERT INTO public.cart VALUES (4, 5, 1, '2022-10-25 10:25:25.38968', '2022-10-25 10:25:25.389683', 'enabled', 0);
INSERT INTO public.cart VALUES (4, 8, 2, '2022-10-25 10:25:25.391125', '2022-10-25 10:25:25.391128', 'enabled', 0);


--
-- Data for Name: category; Type: TABLE DATA; Schema: public; Owner: market
--

INSERT INTO public.category VALUES (1, 'Egin eşik', 0, 'static/cat_image/1.png', 'static/cat_image/1icon.png', '2022-10-25 10:25:25.301587', '2022-10-25 10:25:25.30159', 'enabled', 0);
INSERT INTO public.category VALUES (2, 'Aýakgap', 1, 'static/cat_image/2.png', 'static/cat_image/2icon.png', '2022-10-25 10:25:25.303844', '2022-10-25 10:25:25.303846', 'enabled', 0);
INSERT INTO public.category VALUES (3, 'Balak', 1, 'static/cat_image/3.png', 'static/cat_image/3icon.png', '2022-10-25 10:25:25.306175', '2022-10-25 10:25:25.30618', 'enabled', 0);
INSERT INTO public.category VALUES (4, 'Elektronika enjamlary', 0, 'static/cat_image/4.png', 'static/cat_image/4icon.png', '2022-10-25 10:25:25.308054', '2022-10-25 10:25:25.308058', 'enabled', 0);
INSERT INTO public.category VALUES (5, 'Noutbuklar', 4, 'static/cat_image/5.png', 'static/cat_image/5icon.png', '2022-10-25 10:25:25.309587', '2022-10-25 10:25:25.309589', 'enabled', 0);
INSERT INTO public.category VALUES (6, 'Ykjam telefonlar', 4, 'static/cat_image/6.png', 'static/cat_image/6icon.png', '2022-10-25 10:25:25.310974', '2022-10-25 10:25:25.310976', 'enabled', 0);


--
-- Data for Name: currency; Type: TABLE DATA; Schema: public; Owner: market
--



--
-- Data for Name: discount; Type: TABLE DATA; Schema: public; Owner: market
--



--
-- Data for Name: option; Type: TABLE DATA; Schema: public; Owner: market
--

INSERT INTO public.option VALUES (7, 3, 'aaaaa', '2022-10-25 11:12:15.223125', '2022-10-25 11:12:15.223125', 'enabled', 0);
INSERT INTO public.option VALUES (9, 3, 'bbbbb', '2022-10-25 11:42:20.943496', '2022-10-25 11:42:20.943496', 'enabled', 0);
INSERT INTO public.option VALUES (6, 3, 'cat1opt', '2022-10-25 11:10:05.0859', '2022-10-25 11:10:05.0859', 'enabled', 0);


--
-- Data for Name: option_value; Type: TABLE DATA; Schema: public; Owner: market
--

INSERT INTO public.option_value VALUES (3, 6, 'optval1', '2022-10-25 11:10:05.08923', '2022-10-25 11:10:05.08923', 'enabled', NULL);
INSERT INTO public.option_value VALUES (4, 7, 'sfgddfge', '2022-10-25 11:12:15.226092', '2022-10-25 11:12:15.226092', 'enabled', 0);
INSERT INTO public.option_value VALUES (6, 9, 'gggg', '2022-10-25 11:43:41.660302', '2022-10-25 11:43:41.660302', 'enabled', 0);


--
-- Data for Name: order_item; Type: TABLE DATA; Schema: public; Owner: market
--



--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: market
--



--
-- Data for Name: product; Type: TABLE DATA; Schema: public; Owner: market
--

INSERT INTO public.product VALUES (1, 'LC Waikiki', 'Balak LC Waiki goni kese jubili', 3, 15, 42, '2022-10-25 10:25:25.344128', '2022-10-25 10:25:25.344131', 'enabled', 0);
INSERT INTO public.product VALUES (2, 'LC Waikiki', 'Balak LC Waiki goni pratoy jubili', 3, 15, 42, '2022-10-25 10:25:25.347583', '2022-10-25 10:25:25.34759', 'enabled', 0);
INSERT INTO public.product VALUES (3, 'Pierre Cardin', 'Balak Pierre Cardin goni prastoy jubili', 3, 17, 42, '2022-10-25 10:25:25.350043', '2022-10-25 10:25:25.350045', 'enabled', 0);
INSERT INTO public.product VALUES (4, 'Pierre Cardin', 'Balak Pierre Cardin goni pratoy jubili', 3, 17, 42, '2022-10-25 10:25:25.351909', '2022-10-25 10:25:25.351912', 'enabled', 0);
INSERT INTO public.product VALUES (5, 'kotton', 'Balak Cotton goni pratoy jubili', 3, 16, 42, '2022-10-25 10:25:25.353872', '2022-10-25 10:25:25.353876', 'enabled', 0);
INSERT INTO public.product VALUES (6, 'kotton', 'Balak Cotton goni pratoy jubili', 3, 16, 42, '2022-10-25 10:25:25.356032', '2022-10-25 10:25:25.356034', 'enabled', 0);
INSERT INTO public.product VALUES (7, 'kotton', 'Balak Cotton goni pratoy jubili', 3, 16, 42, '2022-10-25 10:25:25.357998', '2022-10-25 10:25:25.358001', 'enabled', 0);
INSERT INTO public.product VALUES (8, 'LC Waikiki', 'Balak LC Waiki kese pratoy jubili', 3, 15, 42, '2022-10-25 10:25:25.359875', '2022-10-25 10:25:25.359878', 'enabled', 0);
INSERT INTO public.product VALUES (9, 'LC Waikiki', 'Balak LC Waiki goni pratoy jubili', 3, 15, 42, '2022-10-25 10:25:25.361919', '2022-10-25 10:25:25.361923', 'enabled', 0);
INSERT INTO public.product VALUES (10, 'LC Waikiki', 'Balak LC Waiki goni pratoy jubili', 3, 15, 42, '2022-10-25 10:25:25.364276', '2022-10-25 10:25:25.36428', 'enabled', 0);
INSERT INTO public.product VALUES (11, 'LC Waikiki', 'Balak LC Waiki goni pratoy jubili', 3, 15, 42, '2022-10-25 10:25:25.366696', '2022-10-25 10:25:25.3667', 'enabled', 0);
INSERT INTO public.product VALUES (12, 'LC Waikiki', 'Balak LC Waiki goni pratoy jubili', 3, 15, 42, '2022-10-25 10:25:25.368874', '2022-10-25 10:25:25.368878', 'enabled', 0);
INSERT INTO public.product VALUES (13, 'Xiaomi Mi 11 Lite 6GB/128GB', 'Android, экран 6.55'' AMOLED (1080x2400), Qualcomm Snapdragon 732G, ОЗУ 6 ГБ, флэш-память 128 ГБ, карты памяти, камера 64 Мп, аккумулятор 4250 мАч, 2 SIM', 6, 3, 42, '2022-10-25 10:25:25.371686', '2022-10-25 10:25:25.371691', 'enabled', 0);
INSERT INTO public.product VALUES (14, 'Samsung Galaxy A52 SM-A525F/DS 4GB/128GB', 'Android, экран 6.5'' AMOLED (1080x2400), Qualcomm Snapdragon 720G, ОЗУ 4 ГБ, флэш-память 128 ГБ, карты памяти, камера 64 Мп, аккумулятор 4500 мАч, 2 SIM', 6, 4, 42, '2022-10-25 10:25:25.373817', '2022-10-25 10:25:25.373819', 'enabled', 0);
INSERT INTO public.product VALUES (15, 'Xiaomi Redmi Note 10 4GB/64GB', 'Android, экран 6.43'' AMOLED (1080x2400), Qualcomm Snapdragon 678, ОЗУ 4 ГБ, флэш-память 64 ГБ, карты памяти, камера 48 Мп, аккумулятор 5000 мАч, 2 SIM', 6, 3, 42, '2022-10-25 10:25:25.375694', '2022-10-25 10:25:25.375697', 'enabled', 0);
INSERT INTO public.product VALUES (16, 'Samsung Galaxy M31 SM-M315F/DSN 6GB/128GB', 'Android, экран 6.4'' AMOLED (1080x2340), Exynos 9611, ОЗУ 6 ГБ, флэш-память 128 ГБ, карты памяти, камера 64 Мп, аккумулятор 6000 мАч, 2 SIM', 6, 4, 42, '2022-10-25 10:25:25.377754', '2022-10-25 10:25:25.377756', 'enabled', 0);
INSERT INTO public.product VALUES (17, 'Samsung Galaxy A31 SM-A315F/DS 4GB/64GB', 'Android, экран 6.4'' AMOLED (1080x2400), Mediatek MT6768 Helio P65, ОЗУ 4 ГБ, флэш-память 64 ГБ, карты памяти, камера 48 Мп, аккумулятор 5000 мАч, 2 SIM', 6, 4, 42, '2022-10-25 10:25:25.379531', '2022-10-25 10:25:25.379533', 'enabled', 0);
INSERT INTO public.product VALUES (18, 'Lenovo V130-15IKB', '15.6'' 1920 x 1080 TN+Film, 60 Гц, несенсорный, Intel Core i3 7020U 2300 МГц, 8 ГБ, SSD 128 ГБ, видеокарта встроенная, без ОС, DVD, цвет крышки темно-серый', 5, 11, 42, '2022-10-25 10:25:25.381392', '2022-10-25 10:25:25.381395', 'enabled', 0);


--
-- Data for Name: region; Type: TABLE DATA; Schema: public; Owner: market
--

INSERT INTO public.region VALUES (1, 'Balkan', 0, '2022-10-25 10:25:24.569748', '2022-10-25 10:25:24.569753', 'enabled', 0);
INSERT INTO public.region VALUES (2, 'Balkanabat', 1, '2022-10-25 10:25:24.571415', '2022-10-25 10:25:24.571417', 'enabled', 0);
INSERT INTO public.region VALUES (3, 'Türkmenbaşy', 1, '2022-10-25 10:25:24.571419', '2022-10-25 10:25:24.571419', 'enabled', 0);
INSERT INTO public.region VALUES (4, 'Serdar', 1, '2022-10-25 10:25:24.57142', '2022-10-25 10:25:24.571421', 'enabled', 0);
INSERT INTO public.region VALUES (5, 'Magtymguly', 1, '2022-10-25 10:25:24.571422', '2022-10-25 10:25:24.571423', 'enabled', 0);
INSERT INTO public.region VALUES (6, 'Bereket', 1, '2022-10-25 10:25:24.571424', '2022-10-25 10:25:24.571424', 'enabled', 0);
INSERT INTO public.region VALUES (7, 'Esenguly', 1, '2022-10-25 10:25:24.571425', '2022-10-25 10:25:24.571426', 'enabled', 0);
INSERT INTO public.region VALUES (8, 'Etrek', 1, '2022-10-25 10:25:24.571427', '2022-10-25 10:25:24.571427', 'enabled', 0);
INSERT INTO public.region VALUES (9, 'Ahal', 0, '2022-10-25 10:25:24.571428', '2022-10-25 10:25:24.571429', 'enabled', 0);
INSERT INTO public.region VALUES (10, 'Ak-bugday', 9, '2022-10-25 10:25:24.572945', '2022-10-25 10:25:24.572947', 'enabled', 0);
INSERT INTO public.region VALUES (11, 'Babadaýhan', 9, '2022-10-25 10:25:24.572949', '2022-10-25 10:25:24.57295', 'enabled', 0);
INSERT INTO public.region VALUES (12, 'Gökdepe', 9, '2022-10-25 10:25:24.572951', '2022-10-25 10:25:24.572952', 'enabled', 0);
INSERT INTO public.region VALUES (13, 'Kaka', 9, '2022-10-25 10:25:24.572952', '2022-10-25 10:25:24.572953', 'enabled', 0);
INSERT INTO public.region VALUES (14, 'Derweze', 9, '2022-10-25 10:25:24.572954', '2022-10-25 10:25:24.572955', 'enabled', 0);
INSERT INTO public.region VALUES (15, 'Sarahs', 9, '2022-10-25 10:25:24.572956', '2022-10-25 10:25:24.572956', 'enabled', 0);
INSERT INTO public.region VALUES (16, 'Tejen', 9, '2022-10-25 10:25:24.572957', '2022-10-25 10:25:24.572958', 'enabled', 0);
INSERT INTO public.region VALUES (17, 'Mary', 0, '2022-10-25 10:25:24.572959', '2022-10-25 10:25:24.57296', 'enabled', 0);
INSERT INTO public.region VALUES (18, 'Altyn sähra', 17, '2022-10-25 10:25:24.574304', '2022-10-25 10:25:24.574309', 'enabled', 0);
INSERT INTO public.region VALUES (19, 'Baýramaly', 17, '2022-10-25 10:25:24.574312', '2022-10-25 10:25:24.574314', 'enabled', 0);
INSERT INTO public.region VALUES (20, 'Garagum', 17, '2022-10-25 10:25:24.574316', '2022-10-25 10:25:24.574317', 'enabled', 0);
INSERT INTO public.region VALUES (21, 'Mary', 17, '2022-10-25 10:25:24.574319', '2022-10-25 10:25:24.57432', 'enabled', 0);
INSERT INTO public.region VALUES (22, 'Murgap', 17, '2022-10-25 10:25:24.574322', '2022-10-25 10:25:24.574323', 'enabled', 0);
INSERT INTO public.region VALUES (23, 'Serhetabat', 17, '2022-10-25 10:25:24.574325', '2022-10-25 10:25:24.574326', 'enabled', 0);
INSERT INTO public.region VALUES (24, 'Tagtabazar', 17, '2022-10-25 10:25:24.574328', '2022-10-25 10:25:24.574329', 'enabled', 0);
INSERT INTO public.region VALUES (25, 'Lebap', 0, '2022-10-25 10:25:24.574331', '2022-10-25 10:25:24.574333', 'enabled', 0);
INSERT INTO public.region VALUES (26, 'Atamyrat', 25, '2022-10-25 10:25:24.5756', '2022-10-25 10:25:24.575603', 'enabled', 0);
INSERT INTO public.region VALUES (27, 'Beýik Türkmenbaşy', 25, '2022-10-25 10:25:24.575606', '2022-10-25 10:25:24.575607', 'enabled', 0);
INSERT INTO public.region VALUES (28, 'Birata', 25, '2022-10-25 10:25:24.575609', '2022-10-25 10:25:24.57561', 'enabled', 0);
INSERT INTO public.region VALUES (29, 'Döwletli ', 25, '2022-10-25 10:25:24.575612', '2022-10-25 10:25:24.575613', 'enabled', 0);
INSERT INTO public.region VALUES (30, 'Farap', 25, '2022-10-25 10:25:24.575615', '2022-10-25 10:25:24.575616', 'enabled', 0);
INSERT INTO public.region VALUES (31, 'Galkynyş', 25, '2022-10-25 10:25:24.575618', '2022-10-25 10:25:24.575619', 'enabled', 0);
INSERT INTO public.region VALUES (32, 'Garabekewül', 25, '2022-10-25 10:25:24.575621', '2022-10-25 10:25:24.575623', 'enabled', 0);
INSERT INTO public.region VALUES (33, 'Daşoguz', 0, '2022-10-25 10:25:24.575625', '2022-10-25 10:25:24.575626', 'enabled', 0);
INSERT INTO public.region VALUES (34, 'Akdepe', 33, '2022-10-25 10:25:24.576833', '2022-10-25 10:25:24.576838', 'enabled', 0);
INSERT INTO public.region VALUES (35, 'Boldumsaz', 33, '2022-10-25 10:25:24.57684', '2022-10-25 10:25:24.576841', 'enabled', 0);
INSERT INTO public.region VALUES (36, 'Görogly', 33, '2022-10-25 10:25:24.576843', '2022-10-25 10:25:24.576845', 'enabled', 0);
INSERT INTO public.region VALUES (37, 'Gubadag ', 33, '2022-10-25 10:25:24.576847', '2022-10-25 10:25:24.576848', 'enabled', 0);
INSERT INTO public.region VALUES (38, 'Gurbansoltan Eje', 33, '2022-10-25 10:25:24.57685', '2022-10-25 10:25:24.576851', 'enabled', 0);
INSERT INTO public.region VALUES (39, 'Köneürgenç', 33, '2022-10-25 10:25:24.576853', '2022-10-25 10:25:24.576854', 'enabled', 0);
INSERT INTO public.region VALUES (40, 'Ruhybelent', 33, '2022-10-25 10:25:24.576856', '2022-10-25 10:25:24.576857', 'enabled', 0);
INSERT INTO public.region VALUES (41, 'Saparmyrat Nyýazow', 33, '2022-10-25 10:25:24.576859', '2022-10-25 10:25:24.57686', 'enabled', 0);
INSERT INTO public.region VALUES (42, 'Aşgabat', 0, '2022-10-25 10:25:24.576862', '2022-10-25 10:25:24.576863', 'enabled', 0);
INSERT INTO public.region VALUES (43, 'Arçabil etraby', 42, '2022-10-25 10:25:24.577971', '2022-10-25 10:25:24.577974', 'enabled', 0);
INSERT INTO public.region VALUES (44, 'Bagtyýarlyk etraby', 42, '2022-10-25 10:25:24.577976', '2022-10-25 10:25:24.577977', 'enabled', 0);
INSERT INTO public.region VALUES (45, 'Berkararlyk etraby', 42, '2022-10-25 10:25:24.577979', '2022-10-25 10:25:24.57798', 'enabled', 0);
INSERT INTO public.region VALUES (46, 'Çandybil etraby', 42, '2022-10-25 10:25:24.577982', '2022-10-25 10:25:24.577983', 'enabled', 0);
INSERT INTO public.region VALUES (47, 'Köpetdag etraby', 42, '2022-10-25 10:25:24.577984', '2022-10-25 10:25:24.577985', 'enabled', 0);


--
-- Data for Name: sku; Type: TABLE DATA; Schema: public; Owner: market
--

INSERT INTO public.sku VALUES (1, 1, 'sku-16666935253449433', 23, 15, 'static/product_photo/1L.jpg', 'static/product_photo/1S.jpg', 'static/product_photo/1S.jpg', 0, '2022-10-25 10:25:25.345417', '2022-10-25 10:25:25.34542', 'enabled', 0);
INSERT INTO public.sku VALUES (2, 2, 'sku-16666935253483262', 13, 15, 'static/product_photo/2L.jpg', 'static/product_photo/2S.jpg', 'static/product_photo/2S.jpg', 0, '2022-10-25 10:25:25.348538', '2022-10-25 10:25:25.348541', 'enabled', 0);
INSERT INTO public.sku VALUES (3, 3, 'sku-1666693525350414', 43, 15, 'static/product_photo/3L.jpg', 'static/product_photo/3S.jpg', 'static/product_photo/3S.jpg', 0, '2022-10-25 10:25:25.350552', '2022-10-25 10:25:25.350554', 'enabled', 0);
INSERT INTO public.sku VALUES (4, 4, 'sku-1666693525352291', 33, 15, 'static/product_photo/4L.jpg', 'static/product_photo/4S.jpg', 'static/product_photo/4S.jpg', 0, '2022-10-25 10:25:25.352427', '2022-10-25 10:25:25.352429', 'enabled', 0);
INSERT INTO public.sku VALUES (5, 5, 'sku-16666935253543475', 33.3, 15, 'static/product_photo/5L.jpg', 'static/product_photo/5S.jpg', 'static/product_photo/5S.jpg', 0, '2022-10-25 10:25:25.35453', '2022-10-25 10:25:25.354532', 'enabled', 0);
INSERT INTO public.sku VALUES (6, 6, 'sku-16666935253564272', 23.6, 15, 'static/product_photo/6L.jpg', 'static/product_photo/6S.jpg', 'static/product_photo/6S.jpg', 0, '2022-10-25 10:25:25.35656', '2022-10-25 10:25:25.356562', 'enabled', 0);
INSERT INTO public.sku VALUES (7, 7, 'sku-16666935253583484', 65, 15, 'static/product_photo/7L.jpg', 'static/product_photo/7S.jpg', 'static/product_photo/7S.jpg', 0, '2022-10-25 10:25:25.358509', '2022-10-25 10:25:25.358512', 'enabled', 0);
INSERT INTO public.sku VALUES (8, 8, 'sku-16666935253602104', 125, 15, 'static/product_photo/8L.jpg', 'static/product_photo/8S.jpg', 'static/product_photo/8S.jpg', 0, '2022-10-25 10:25:25.360347', '2022-10-25 10:25:25.36035', 'enabled', 0);
INSERT INTO public.sku VALUES (9, 9, 'sku-1666693525362415', 423, 15, 'static/product_photo/9L.jpg', 'static/product_photo/9S.jpg', 'static/product_photo/9S.jpg', 0, '2022-10-25 10:25:25.36259', '2022-10-25 10:25:25.362593', 'enabled', 0);
INSERT INTO public.sku VALUES (10, 10, 'sku-16666935253646905', 263, 15, 'static/product_photo/10L.jpg', 'static/product_photo/10S.jpg', 'static/product_photo/10S.jpg', 0, '2022-10-25 10:25:25.364871', '2022-10-25 10:25:25.364874', 'enabled', 0);
INSERT INTO public.sku VALUES (11, 11, 'sku-16666935253671086', 223, 15, 'static/product_photo/11L.jpg', 'static/product_photo/11S.jpg', 'static/product_photo/11S.jpg', 0, '2022-10-25 10:25:25.367266', '2022-10-25 10:25:25.367269', 'enabled', 0);
INSERT INTO public.sku VALUES (12, 12, 'sku-16666935253695138', 263, 15, 'static/product_photo/12L.jpg', 'static/product_photo/12S.jpg', 'static/product_photo/12S.jpg', 0, '2022-10-25 10:25:25.369781', '2022-10-25 10:25:25.369785', 'enabled', 0);
INSERT INTO public.sku VALUES (13, 13, 'sku-1666693525372156', 456, 15, '{static/product_photo/13L1.jpeg,static/product_photo/13L2.jpeg,static/product_photo/13L3.jpeg}', 'static/product_photo/13S.jpg', 'static/product_photo/13S.jpg', 0, '2022-10-25 10:25:25.372333', '2022-10-25 10:25:25.372336', 'enabled', 0);
INSERT INTO public.sku VALUES (14, 14, 'sku-16666935253741367', 466, 15, '{static/product_photo/14L1.jpeg,static/product_photo/14L2.jpeg,static/product_photo/14L3.jpeg}', 'static/product_photo/14S.jpg', 'static/product_photo/14S.jpg', 0, '2022-10-25 10:25:25.374269', '2022-10-25 10:25:25.374271', 'enabled', 0);
INSERT INTO public.sku VALUES (15, 15, 'sku-1666693525376127', 564, 15, '{static/product_photo/15L1.jpeg,static/product_photo/15L2.jpeg,static/product_photo/15L3.jpeg}', 'static/product_photo/15S.jpg', 'static/product_photo/15S.jpg', 0, '2022-10-25 10:25:25.376295', '2022-10-25 10:25:25.376297', 'enabled', 0);
INSERT INTO public.sku VALUES (16, 16, 'sku-16666935253780735', 654, 15, '{static/product_photo/16L1.jpeg,static/product_photo/16L2.jpeg,static/product_photo/16L3.jpeg}', 'static/product_photo/16S.jpg', 'static/product_photo/16S.jpg', 0, '2022-10-25 10:25:25.378206', '2022-10-25 10:25:25.378208', 'enabled', 0);
INSERT INTO public.sku VALUES (17, 17, 'sku-1666693525379845', 555, 15, '{static/product_photo/17L1.jpeg,static/product_photo/17L2.jpeg,static/product_photo/17L3.jpeg}', 'static/product_photo/17S.jpg', 'static/product_photo/17S.jpg', 0, '2022-10-25 10:25:25.379986', '2022-10-25 10:25:25.379988', 'enabled', 0);
INSERT INTO public.sku VALUES (18, 18, 'sku-16666935253817098', 557, 15, '{static/product_photo/18L1.jpeg,static/product_photo/18L2.jpeg,static/product_photo/18L3.jpeg}', 'static/product_photo/18S.jpg', 'static/product_photo/18S.jpg', 0, '2022-10-25 10:25:25.38184', '2022-10-25 10:25:25.381842', 'enabled', 0);


--
-- Data for Name: sku_value; Type: TABLE DATA; Schema: public; Owner: market
--



--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: market
--

INSERT INTO public.users VALUES (1, 'c420e27f-5411-4dce-8b59-efbd8ae261e9', 'Super', '$2b$12$VUHnyCv0de1RhlwdOy5t2.xPauQvKUekl9SwadUJegCwzd97B53fu', 'superadmin@plan.com', '+99365777777', NULL, 'static/user_photo/user1.jpeg', 'SUPERADMIN', 5, 0, '2022-10-25 10:25:24.762513', '2022-10-25 10:25:24.762519', 'enabled', 0);
INSERT INTO public.users VALUES (2, '8297e0ee-cc71-4d72-b297-11a53d488b54', 'Admin', '$2b$12$OkZa4SgJId89Q9mpT7OvZeNsZ8TnH2yVDcX9lNEMutI1EltAfI4re', 'admin@plan.com', '+99365788777', NULL, 'static/user_photo/user2.jpeg', 'ADMIN', 7, 1, '2022-10-25 10:25:24.941413', '2022-10-25 10:25:24.941418', 'enabled', 0);
INSERT INTO public.users VALUES (3, '214bfbf8-470f-40b1-9df4-373ace171506', 'Admin2', '$2b$12$ck8FsQkw09JkETD6skPZIuBkcDswN04aoWFjBdrYzJf5/H216KQQK', 'admin2@plan.com', '+99365997777', NULL, 'static/user_photo/user3.jpeg', 'ADMIN', 5, 1, '2022-10-25 10:25:25.118843', '2022-10-25 10:25:25.118846', 'enabled', 0);
INSERT INTO public.users VALUES (4, '98aaf108-a272-4d99-bd61-01136e1ef040', 'user', '$2b$12$elzuw1gzQLmuCylr3jh.P.Nsy/DrejbSlprJ.EfBJhPsM.J2j0fvm', 'user@plan.com', '+99365777755', NULL, 'static/user_photo/user4.jpeg', 'USER', 5, 2, '2022-10-25 10:25:25.296309', '2022-10-25 10:25:25.296316', 'enabled', 0);


--
-- Name: brand_id_seq; Type: SEQUENCE SET; Schema: public; Owner: market
--

SELECT pg_catalog.setval('public.brand_id_seq', 17, true);

SELECT pg_catalog.setval('public.cart_id_seq', 9, true);

--
-- Name: category_id_seq; Type: SEQUENCE SET; Schema: public; Owner: market
--

SELECT pg_catalog.setval('public.category_id_seq', 6, true);


--
-- Name: currency_id_seq; Type: SEQUENCE SET; Schema: public; Owner: market
--

SELECT pg_catalog.setval('public.currency_id_seq', 1, false);


--
-- Name: option_id_seq; Type: SEQUENCE SET; Schema: public; Owner: market
--

SELECT pg_catalog.setval('public.option_id_seq', 16, true);


--
-- Name: option_value_id_seq; Type: SEQUENCE SET; Schema: public; Owner: market
--

SELECT pg_catalog.setval('public.option_value_id_seq', 6, true);


--
-- Name: order_item_id_seq; Type: SEQUENCE SET; Schema: public; Owner: market
--

SELECT pg_catalog.setval('public.order_item_id_seq', 1, false);


--
-- Name: product_id_seq; Type: SEQUENCE SET; Schema: public; Owner: market
--

SELECT pg_catalog.setval('public.product_id_seq', 18, true);


--
-- Name: region_id_seq; Type: SEQUENCE SET; Schema: public; Owner: market
--

SELECT pg_catalog.setval('public.region_id_seq', 47, true);


--
-- Name: sku_id_seq; Type: SEQUENCE SET; Schema: public; Owner: market
--

SELECT pg_catalog.setval('public.sku_id_seq', 18, true);


--
-- Name: sku_value_id_seq; Type: SEQUENCE SET; Schema: public; Owner: market
--

SELECT pg_catalog.setval('public.sku_value_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: market
--

SELECT pg_catalog.setval('public.users_id_seq', 4, true);


--
-- Name: brand brand_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.brand
    ADD CONSTRAINT brand_pkey PRIMARY KEY (id);


--
-- Name: cart cart_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.cart
    ADD CONSTRAINT cart_pkey PRIMARY KEY (user_id, sku_id);


--
-- Name: category category_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.category
    ADD CONSTRAINT category_pkey PRIMARY KEY (id);


--
-- Name: currency currency_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.currency
    ADD CONSTRAINT currency_pkey PRIMARY KEY (id);


--
-- Name: discount discount_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.discount
    ADD CONSTRAINT discount_pkey PRIMARY KEY (id, product_id);


--
-- Name: option name_unique; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.option
    ADD CONSTRAINT name_unique UNIQUE (name);


--
-- Name: option option_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.option
    ADD CONSTRAINT option_pkey PRIMARY KEY (id);


--
-- Name: option_value option_value_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.option_value
    ADD CONSTRAINT option_value_pkey PRIMARY KEY (id);


--
-- Name: order_item order_item_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.order_item
    ADD CONSTRAINT order_item_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: product product_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_pkey PRIMARY KEY (id);


--
-- Name: region region_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.region
    ADD CONSTRAINT region_pkey PRIMARY KEY (id);


--
-- Name: sku sku_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.sku
    ADD CONSTRAINT sku_pkey PRIMARY KEY (id);


--
-- Name: sku sku_sku_key; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.sku
    ADD CONSTRAINT sku_sku_key UNIQUE (sku);


--
-- Name: sku_value sku_value_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.sku_value
    ADD CONSTRAINT sku_value_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_phone_number_key; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_phone_number_key UNIQUE (phone_number);


--
-- Name: users users_photo_key; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_photo_key UNIQUE (photo);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_public_id_key; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_public_id_key UNIQUE (public_id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: cart cart_sku_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.cart
    ADD CONSTRAINT cart_sku_id_fkey FOREIGN KEY (sku_id) REFERENCES public.sku(id);


--
-- Name: cart cart_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.cart
    ADD CONSTRAINT cart_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: discount discount_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.discount
    ADD CONSTRAINT discount_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.product(id);


--
-- Name: discount discount_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.discount
    ADD CONSTRAINT discount_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: option option_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.option
    ADD CONSTRAINT option_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.category(id);


--
-- Name: option_value option_value_option_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.option_value
    ADD CONSTRAINT option_value_option_id_fkey FOREIGN KEY (option_id) REFERENCES public.option(id);


--
-- Name: order_item order_item_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.order_item
    ADD CONSTRAINT order_item_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: order_item order_item_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.order_item
    ADD CONSTRAINT order_item_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.product(id);


--
-- Name: orders orders_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: product product_brand_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES public.brand(id);


--
-- Name: product product_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.category(id);


--
-- Name: product product_region_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_region_id_fkey FOREIGN KEY (region_id) REFERENCES public.region(id);


--
-- Name: sku sku_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.sku
    ADD CONSTRAINT sku_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.product(id);


--
-- Name: sku_value sku_value_option_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.sku_value
    ADD CONSTRAINT sku_value_option_id_fkey FOREIGN KEY (option_id) REFERENCES public.option(id);


--
-- Name: sku_value sku_value_option_value_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.sku_value
    ADD CONSTRAINT sku_value_option_value_id_fkey FOREIGN KEY (option_value_id) REFERENCES public.option_value(id);


--
-- Name: sku_value sku_value_sku_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.sku_value
    ADD CONSTRAINT sku_value_sku_id_fkey FOREIGN KEY (sku_id) REFERENCES public.sku(id);


--
-- Name: users users_region_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: market
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_region_id_fkey FOREIGN KEY (region_id) REFERENCES public.region(id);


--
-- PostgreSQL database dump complete
--

