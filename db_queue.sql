--
-- PostgreSQL database dump
--

-- Dumped from database version 10.14 (Ubuntu 10.14-1.pgdg18.04+1)
-- Dumped by pg_dump version 12.4 (Ubuntu 12.4-1.pgdg18.04+1)

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

--
-- Name: func_consume_message(integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.func_consume_message(_id integer) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
	access_granted bool;
	del_tot int;
BEGIN
	WITH rows AS (
		delete from tbl_message where "id" = _id returning "id") 
	SELECT count(*) into del_tot FROM rows;
	if del_tot = 1
	then
		return true;
	end if;
	return false;
END;
$$;


ALTER FUNCTION public.func_consume_message(_id integer) OWNER TO postgres;

--
-- Name: func_create_queue(text, text, text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.func_create_queue(_name text, _username text, _password text) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
	access_granted bool;
	userid text;
BEGIN
	select true into access_granted from tbl_user where username = _username and "password" = _password;
	if access_granted
	then
		if (select true from tbl_queue where "name" = _name) is null
		then
			insert into tbl_queue ("name", "owner") values (_name, _username);
			return true;
		end if;
		return false;
	end if;
	return false;
END;
$$;


ALTER FUNCTION public.func_create_queue(_name text, _username text, _password text) OWNER TO postgres;

--
-- Name: func_delete_queue(text, text, text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.func_delete_queue(_name text, _username text, _password text) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
	access_granted bool;
	userid text;
BEGIN
	select true into access_granted from tbl_user where username = _username and "password" = _password;
	if access_granted
	then
		if (select true from tbl_queue where "owner" = _username and "name" = _name)
		then
			delete from tbl_queue where "owner" = _username and "name" = _name;
			return true;
		end if;
		return false;
	end if;
	return false;
END;
$$;


ALTER FUNCTION public.func_delete_queue(_name text, _username text, _password text) OWNER TO postgres;

--
-- Name: func_produce_message(text, text, text, text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.func_produce_message(_name text, _message text, _username text, _password text) RETURNS TABLE(_status boolean, _id integer)
    LANGUAGE plpgsql
    AS $$
DECLARE
	access_granted bool;
	idmsg int;
BEGIN
	select true into access_granted from tbl_user where username = _username and "password" = _password;
	if access_granted
	then
		insert into tbl_message (name_queue, "message") values (_name, _message) returning id into idmsg;
		return query
			select true, idmsg;
	else
		return query
			select false, 0;
	end if;
END;
$$;


ALTER FUNCTION public.func_produce_message(_name text, _message text, _username text, _password text) OWNER TO postgres;

SET default_tablespace = '';

--
-- Name: tbl_message; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tbl_message (
    name_queue text NOT NULL,
    message text,
    id integer NOT NULL
);


ALTER TABLE public.tbl_message OWNER TO postgres;

--
-- Name: tbl_message_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tbl_message_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tbl_message_id_seq OWNER TO postgres;

--
-- Name: tbl_message_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tbl_message_id_seq OWNED BY public.tbl_message.id;


--
-- Name: tbl_queue; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tbl_queue (
    name text NOT NULL,
    owner text NOT NULL
);


ALTER TABLE public.tbl_queue OWNER TO postgres;

--
-- Name: tbl_user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tbl_user (
    username text NOT NULL,
    password text NOT NULL
);


ALTER TABLE public.tbl_user OWNER TO postgres;

--
-- Name: tbl_message id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tbl_message ALTER COLUMN id SET DEFAULT nextval('public.tbl_message_id_seq'::regclass);


--
-- Data for Name: tbl_message; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tbl_message (name_queue, message, id) FROM stdin;
\.


--
-- Data for Name: tbl_queue; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tbl_queue (name, owner) FROM stdin;
q-1	admin
\.


--
-- Data for Name: tbl_user; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tbl_user (username, password) FROM stdin;
admin	5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8
\.


--
-- Name: tbl_message_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tbl_message_id_seq', 72, true);


--
-- Name: tbl_message tbl_message_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tbl_message
    ADD CONSTRAINT tbl_message_pk PRIMARY KEY (id);


--
-- Name: tbl_queue tbl_queue_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tbl_queue
    ADD CONSTRAINT tbl_queue_pk PRIMARY KEY (name);


--
-- Name: tbl_user tbl_user_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tbl_user
    ADD CONSTRAINT tbl_user_pk PRIMARY KEY (username);


--
-- Name: tbl_message_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX tbl_message_id_uindex ON public.tbl_message USING btree (id);


--
-- Name: tbl_queue_name_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX tbl_queue_name_uindex ON public.tbl_queue USING btree (name);


--
-- Name: tbl_user_username_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX tbl_user_username_uindex ON public.tbl_user USING btree (username);


--
-- Name: tbl_message tbl_message_tbl_queue_name_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tbl_message
    ADD CONSTRAINT tbl_message_tbl_queue_name_fk FOREIGN KEY (name_queue) REFERENCES public.tbl_queue(name) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: tbl_queue tbl_queue_tbl_user_username_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tbl_queue
    ADD CONSTRAINT tbl_queue_tbl_user_username_fk FOREIGN KEY (owner) REFERENCES public.tbl_user(username);


--
-- PostgreSQL database dump complete
--

