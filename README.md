# golang-boxes

## create boxes database

1. build the container image for the PostgreSQL database using the Containerfile
2. run the container using that image and publish all exposed ports to random ports. This means that the port tcp/5432 should be exposed at some random ephemeral port on your host.
3. connect to the database but the port number might change. Make sure you're trying to connect to the correct one. Check with `ss -tlnp` for the listenting tcp ports.

~~~bash
$ podman image build -t boxesdb -f Containerfile .
$ podman container run -d -P --name boxesdb boxesdb
$ psql -U postgres -d boxes -h localhost -p 43089 -W
~~~

## configure boxes database

Once connected to the running PostgreSQL container instance, execute the following SQL commands.

~~~sql
CREATE TABLE boxes (
    id serial PRIMARY KEY,
    title varchar(100) NOT NULL,
    content text NOT NULL,
    created timestamp NOT NULL,
    expires timestamp NOT NULL
);

CREATE INDEX boxes_idx ON boxes(created);

INSERT INTO boxes (title, content, created, expires) VALUES 
    ('box1', 'aaaaaaa', now(), now() + interval '365 days'),
    ('box2', 'bbbbbbb', now(), now() + interval '365 days'),
    ('box3', 'ccccccc', now(), now() + interval '365 days');
~~~

Create a service user with read-only permissions.

~~~sql
CREATE USER web WITH PASSWORD 'webesjelszo';
GRANT pg_read_all_data TO web;
~~~
