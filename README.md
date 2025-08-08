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

The PostgreSQL container instance will be configured via an initialization scripts. It'll create three tables in the public schema. One for holding our _boxes_, which will be filled with a few smaple entries. One table for holding the client session data and one for user data.

