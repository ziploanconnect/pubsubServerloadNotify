# Notification for pubsub server on server load

This service works only when you have nats server running and a service that emits `stats.loadaverage` and gives the load of server. 

### Mandatory 

Install nats server and make sure `gnatsd` is available as a command in the system.



Here is how to create the PUBLISH service for server load.

- `cd ~ && wget https://github.com/yuce/catnats/raw/0.1.2/catnats.py`

- `chmod +x catnats.py`

- `sudo mv catnats.py /srv/nats/bin/catnats`

- Test it with `printf "PING\r\n" | /srv/nats/bin/catnats --addr 127.0.0.1:4222` (Only works when nats server is running on port 4222)

- Schedule it with crontab. `crontab -e`

- `*/1 * * * * bash /home/ubuntu/publish_load_average.sh`


## Install this subscription service 

- Make sure you have golang
-  `go get github.com/ziploanconnect/pubsubServerloadNotify` and run `$HOME/go/bin/pubsubServerloadNotify` and it should work
- I have created a natsnoty.service for ubuntu see if that works.

## Testing 

- Install stress to give system some stress testing`sudo apt-get install -y stress`
- Run `stress --cpu $(getconf _NPROCESSORS_ONLN)` for fews to test. You should start receiving emails