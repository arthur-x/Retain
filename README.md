# Retain

## Usage
1. Run your server using this:
```shell
bin/retainServe -s <service> -p <port> -l -d (BlockStoreAddr*)
```
Here, `service` should be one of three values: meta, block, or both. This is used to specify the service provided by the server. `port` defines the port number that the server listens to (default=8080). `-l` configures the server to only listen on localhost. `-d` configures the server to output log statements. Lastly, (BlockStoreAddr\*) is the BlockStore address that the server is configured with. If `service=both` then the BlockStoreAddr should be the `ip:port` of this server.

2. Run your client using this:
```shell
bin/retainSync -d <meta_addr:port> <base_dir> <block_size>
```

### Examples
```shell
bin/retainServe -s both -p 8081 -l localhost:8081
```
This starts a server that listens only to localhost on port 8081 and services both the BlockStore and MetaStore interface.

```shell
Run the commands below on separate terminals (or nodes)
> bin/retainServe -s block -p 8081 -l
> bin/retainServe -s meta -l localhost:8081
```
The first line starts a server that services only the BlockStore interface and listens only to localhost on port 8081. The second line starts a server that services only the MetaStore interface, listens only to localhost on port 8080, and references the BlockStore we created as the underlying BlockStore. (Note: if these are on separate nodes, then you should use the public ip address and remove `-l`)

From a new terminal (or a new node), choose a base directory with some files in it.
```shell
> mkdir dataA
> cp ~/pic.jpg dataA/ 
> bin/retainSync server_addr:port dataA 4096
```
This would sync pic.jpg to the server hosted on `server_addr:port`, using `dataA` as the base directory, with a block size of 4096 bytes.

From another terminal (or a new node), run the client to sync with the server.
```shell
> mkdir dataB
> bin/retainSync server_addr:port dataB 4096
> ls dataB/
pic.jpg index.txt
```
We observe that pic.jpg has been synced to this client.
