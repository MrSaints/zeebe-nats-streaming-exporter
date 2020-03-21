--
-- PostgreSQL database dump
--

-- Dumped from database version 11.1
-- Dumped by pg_dump version 12.2

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

SET default_tablespace = '';

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: zeebe-dev-user
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO "zeebe-dev-user";

--
-- Name: workflow_instance_elements; Type: TABLE; Schema: public; Owner: zeebe-dev-user
--

CREATE TABLE public.workflow_instance_elements (
    id character varying(255) NOT NULL,
    key bigint NOT NULL,
    partition_id integer NOT NULL,
    "position" bigint NOT NULL,
    workflow_key bigint NOT NULL,
    workflow_instance_key bigint NOT NULL,
    intent character varying(255) NOT NULL,
    element_id character varying(255) NOT NULL,
    element_type character varying(255) NOT NULL,
    flow_scope_key bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    indexed_at timestamp without time zone NOT NULL
);


ALTER TABLE public.workflow_instance_elements OWNER TO "zeebe-dev-user";

--
-- Name: workflow_instances; Type: TABLE; Schema: public; Owner: zeebe-dev-user
--

CREATE TABLE public.workflow_instances (
    id character varying(255) NOT NULL,
    key bigint NOT NULL,
    partition_id integer NOT NULL,
    workflow_key bigint NOT NULL,
    process_id character varying(255) NOT NULL,
    version integer NOT NULL,
    intent character varying(255) NOT NULL,
    parent_workflow_instance_key bigint NOT NULL,
    parent_element_instance_key bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    indexed_at timestamp without time zone NOT NULL
);


ALTER TABLE public.workflow_instances OWNER TO "zeebe-dev-user";

--
-- Name: workflows; Type: TABLE; Schema: public; Owner: zeebe-dev-user
--

CREATE TABLE public.workflows (
    id character varying(255) NOT NULL,
    key bigint NOT NULL,
    process_id character varying(255) NOT NULL,
    version integer NOT NULL,
    resource text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    indexed_at timestamp without time zone NOT NULL
);


ALTER TABLE public.workflows OWNER TO "zeebe-dev-user";

--
-- Name: workflow_instance_elements workflow_instance_elements_pkey; Type: CONSTRAINT; Schema: public; Owner: zeebe-dev-user
--

ALTER TABLE ONLY public.workflow_instance_elements
    ADD CONSTRAINT workflow_instance_elements_pkey PRIMARY KEY (id);


--
-- Name: workflow_instances workflow_instances_pkey; Type: CONSTRAINT; Schema: public; Owner: zeebe-dev-user
--

ALTER TABLE ONLY public.workflow_instances
    ADD CONSTRAINT workflow_instances_pkey PRIMARY KEY (id);


--
-- Name: workflows workflows_pkey; Type: CONSTRAINT; Schema: public; Owner: zeebe-dev-user
--

ALTER TABLE ONLY public.workflows
    ADD CONSTRAINT workflows_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: workflow_instance_elements_created_at_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflow_instance_elements_created_at_idx ON public.workflow_instance_elements USING btree (created_at);


--
-- Name: workflow_instance_elements_intent_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflow_instance_elements_intent_idx ON public.workflow_instance_elements USING btree (intent);


--
-- Name: workflow_instance_elements_key_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflow_instance_elements_key_idx ON public.workflow_instance_elements USING btree (key);


--
-- Name: workflow_instance_elements_partition_id_position_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflow_instance_elements_partition_id_position_idx ON public.workflow_instance_elements USING btree (partition_id, "position");


--
-- Name: workflow_instance_elements_workflow_instance_key_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflow_instance_elements_workflow_instance_key_idx ON public.workflow_instance_elements USING btree (workflow_instance_key);


--
-- Name: workflow_instance_elements_workflow_key_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflow_instance_elements_workflow_key_idx ON public.workflow_instance_elements USING btree (workflow_key);


--
-- Name: workflow_instances_created_at_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflow_instances_created_at_idx ON public.workflow_instances USING btree (created_at);


--
-- Name: workflow_instances_intent_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflow_instances_intent_idx ON public.workflow_instances USING btree (intent);


--
-- Name: workflow_instances_key_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflow_instances_key_idx ON public.workflow_instances USING btree (key);


--
-- Name: workflow_instances_process_id_version_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflow_instances_process_id_version_idx ON public.workflow_instances USING btree (process_id, version);


--
-- Name: workflow_instances_workflow_key_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflow_instances_workflow_key_idx ON public.workflow_instances USING btree (workflow_key);


--
-- Name: workflows_created_at_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflows_created_at_idx ON public.workflows USING btree (created_at);


--
-- Name: workflows_key_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE INDEX workflows_key_idx ON public.workflows USING btree (key);


--
-- Name: workflows_process_id_version_idx; Type: INDEX; Schema: public; Owner: zeebe-dev-user
--

CREATE UNIQUE INDEX workflows_process_id_version_idx ON public.workflows USING btree (process_id, version);


--
-- PostgreSQL database dump complete
--

