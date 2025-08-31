CREATE TABLE public.push_subscriptions (
	id bigserial NOT NULL,
	user_id text NOT NULL,
	endpoint text NOT NULL,
	p256dh text NOT NULL,
	auth text NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	CONSTRAINT push_subscriptions_endpoint_key UNIQUE (endpoint),
	CONSTRAINT push_subscriptions_pkey PRIMARY KEY (id)
);


-- public.push_subscriptions foreign keys

ALTER TABLE public.push_subscriptions ADD CONSTRAINT push_subscriptions_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id);