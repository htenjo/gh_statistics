CREATE TABLE public.gh_user
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(256) NOT NULL,
    email VARCHAR(256) NOT NULL,
    session_id VARCHAR(50) NOT NULL,
    access_token VARCHAR(256) NOT NULL,
    repos text
);

ALTER TABLE public.gh_user OWNER to postgres;