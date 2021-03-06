--
-- PostgreSQL database dump
--

-- Dumped from database version 12.5 (Ubuntu 12.5-0ubuntu0.20.04.1)
-- Dumped by pg_dump version 12.5 (Ubuntu 12.5-0ubuntu0.20.04.1)

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
-- Name: citext; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;


--
-- Name: EXTENSION citext; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';


--
-- Name: tablefunc; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS tablefunc WITH SCHEMA public;


--
-- Name: EXTENSION tablefunc; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION tablefunc IS 'functions that manipulate whole tables, including crosstab';

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: forums; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.forums (
    forumid integer NOT NULL,
    title character varying(100) NOT NULL,
    creator public.citext NOT NULL,
    forumname public.citext NOT NULL,
    posts_count integer DEFAULT 0 NOT NULL,
    threads_count integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.forums OWNER TO postgres;

--
-- Name: forums_forumid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.forums_forumid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.forums_forumid_seq OWNER TO postgres;

--
-- Name: forums_forumid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.forums_forumid_seq OWNED BY public.forums.forumid;


--
-- Name: posts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.posts (
    postid integer NOT NULL,
    parentid integer DEFAULT 0 NOT NULL,
    creator public.citext NOT NULL,
    message text,
    isedited boolean DEFAULT false NOT NULL,
    threadid integer NOT NULL,
    created timestamp(3) with time zone DEFAULT now() NOT NULL,
    path integer[] NOT NULL,
    forumname public.citext NOT NULL
);

ALTER TABLE public.posts OWNER TO postgres;

--
-- Name: posts_postid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.posts_postid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.posts_postid_seq OWNER TO postgres;

--
-- Name: posts_postid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.posts_postid_seq OWNED BY public.posts.postid;


--
-- Name: threads; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.threads (
    threadid integer NOT NULL,
    creator public.citext NOT NULL,
    title character varying(100) NOT NULL,
    forumname public.citext DEFAULT NULL::character varying,
    message text,
    created timestamp(3) with time zone DEFAULT now() NOT NULL,
    rating integer DEFAULT 0 NOT NULL,
    threadname public.citext
);


ALTER TABLE public.threads OWNER TO postgres;

--
-- Name: threads_threadid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.threads_threadid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.threads_threadid_seq OWNER TO postgres;

--
-- Name: threads_threadid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.threads_threadid_seq OWNED BY public.threads.threadid;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    userid integer NOT NULL,
    username public.citext NOT NULL,
    email public.citext NOT NULL,
    fullname character varying(100),
    description text
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_userid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_userid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_userid_seq OWNER TO postgres;

--
-- Name: users_userid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_userid_seq OWNED BY public.users.userid;


--
-- Name: votes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.votes (
    voteid integer NOT NULL,
    username public.citext NOT NULL,
    upvote boolean DEFAULT false NOT NULL,
    threadid integer NOT NULL
);


ALTER TABLE public.votes OWNER TO postgres;

--
-- Name: votes_voteid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.votes_voteid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.votes_voteid_seq OWNER TO postgres;

--
-- Name: votes_voteid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.votes_voteid_seq OWNED BY public.votes.voteid;


--
-- Name: forums forumid; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums ALTER COLUMN forumid SET DEFAULT nextval('public.forums_forumid_seq'::regclass);


--
-- Name: posts postid; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts ALTER COLUMN postid SET DEFAULT nextval('public.posts_postid_seq'::regclass);


--
-- Name: threads threadid; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads ALTER COLUMN threadid SET DEFAULT nextval('public.threads_threadid_seq'::regclass);


--
-- Name: users userid; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN userid SET DEFAULT nextval('public.users_userid_seq'::regclass);


--
-- Name: votes voteid; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes ALTER COLUMN voteid SET DEFAULT nextval('public.votes_voteid_seq'::regclass);


--
-- Name: forums forums_pk_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_pk_id PRIMARY KEY (forumid);


--
-- Name: forums forums_un_forumname; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_un_forumname UNIQUE (forumname);


--
-- Name: posts posts_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pk PRIMARY KEY (postid);


--
-- Name: threads threads_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_pk PRIMARY KEY (threadid);


--
-- Name: threads threads_un_threadname; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_un_threadname UNIQUE (threadname);


--
-- Name: users users_pk_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pk_id PRIMARY KEY (userid);


--
-- Name: users users_un_email; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_un_email UNIQUE (email);


--
-- Name: users users_un_username; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_un_username UNIQUE (username);


--
-- Name: votes votes_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_pk PRIMARY KEY (voteid);


--
-- Name: votes votes_un_username_threadID; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT "votes_un_username_threadID" UNIQUE (username, threadid);


--
-- Name: forums_creator_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX forums_creator_idx ON public.forums USING btree (creator);


--
-- Name: posts_forumname_creator_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX posts_forumname_creator_idx ON public.posts USING btree (forumname, creator);


--
-- Name: posts_path_start_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX posts_path_start_idx ON public.posts USING btree ((path[1]));


--
-- Name: posts_threadid_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX posts_threadid_idx ON public.posts USING btree (threadid);


--
-- Name: posts_tree_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX posts_tree_search_idx ON public.posts USING btree (threadid, path);


--
-- Name: threads_creator_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX threads_creator_idx ON public.threads USING btree (creator);


--
-- Name: threads_forumname_created_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX threads_forumname_created_idx ON public.threads USING btree (forumname, created);


--
-- Name: threads_forumname_creator_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX threads_forumname_creator_idx ON public.threads USING btree (forumname, creator);


--
-- Name: threads_forumname_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX threads_forumname_idx ON public.threads USING btree (forumname);


--
-- Name: forums forums_fk_creator; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.forums
    ADD CONSTRAINT forums_fk_creator FOREIGN KEY (creator) REFERENCES public.users(username) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: posts posts_fk_creator; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_fk_creator FOREIGN KEY (creator) REFERENCES public.users(username) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: posts posts_fk_forumname; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_fk_forumname FOREIGN KEY (forumname) REFERENCES public.forums(forumname) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: posts posts_fk_threadid; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_fk_threadid FOREIGN KEY (threadid) REFERENCES public.threads(threadid) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: threads threads_fk_creator; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_fk_creator FOREIGN KEY (creator) REFERENCES public.users(username) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: threads threads_fk_forumname; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_fk_forumname FOREIGN KEY (forumname) REFERENCES public.forums(forumname) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: votes votes_fk_threadid; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_fk_threadid FOREIGN KEY (threadid) REFERENCES public.threads(threadid) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: votes votes_fk_username; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT votes_fk_username FOREIGN KEY (username) REFERENCES public.users(username) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

