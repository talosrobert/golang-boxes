# golang-boxes

## create boxes database

1. build the container image for the PostgreSQL database using the Containerfile
2. run the container using that image and publish all exposed ports to random ports. This means that the port tcp/5432 should be exposed at some random ephemeral port on your host.
3. connect to the database but the port number might change. Make sure you're trying to connect to the correct one. Check with `ss -tlnp` for the listenting tcp ports. Or use `podman container inspect` to see the exporter ports.

~~~bash
$ podman image build -t boxesdb -f Containerfile .
$ podman container run -d -P --name boxesdb boxesdb
$ psql -U postgres -d boxes -h localhost -p 43089 -W
~~~

## configure boxes database

Once connected to the running PostgreSQL container instance, execute the following SQL commands. They'll create two tables in the public schema. One for holding our _boxes_, which will be filled with a few smaple entries. And another table for holding the client session data.

~~~sql
CREATE TABLE boxes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title varchar(100) NOT NULL,
    content text NOT NULL,
    created timestamp NOT NULL DEFAULT now(),
    expires timestamp NOT NULL
);

CREATE INDEX boxes_idx ON boxes(created);

INSERT INTO boxes (title, content, expires) VALUES 
    ('box1', 'aaaaaaa', now() + interval '365 days'),
    ('box2', 'bbbbbbb', now() + interval '365 days'),
    ('box3', 'ccccccc', now() + interval '365 days');

CREATE TABLE sessions (
	token text PRIMARY KEY,
	data bytea NOT NULL,
	expiry timestamp NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

CREATE EXTENSION pgcrypto;

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(100) NOT NULL,
    email text NOT NULL UNIQUE,
    pswhash text NOT NULL, 
    created timestamp NOT NULL DEFAULT now()
);
~~~

Create a service user with read-only permissions.

~~~sql
CREATE USER web WITH PASSWORD 'webesjelszo';
GRANT pg_read_all_data TO web;
~~~
