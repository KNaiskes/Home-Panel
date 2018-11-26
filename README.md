# Home-Panel

### A Web application to control all of my IoT devices from [Home-IoT](https://github.com/KNaiskes/Home-IoT) repository.

### Building
```
$ git clone https://github.com/KNaiskes/Home-Panel
$ cd Home-Panel
$ go get -d ./...
$ go install
```

## Install Mosquitto

#### Debian based distros
```
$ sudo apt-get install mosquitto mosquitto-clients
```

#### Arch
```
$ sudo pacman -S mosquitto mosquitto-clients
```

## Configure Mosquitto

Append /etc/mosquitto/mosquitto.conf

```
sudo vim /etc/mosquitto/mosquitto.conf
allow_anonymous false
password_file /etc/mosquitto/pwfile
listener 1883
```

#### Generate a new username and password

```
sudo mosquitto_passwd -c /etc/mosquitto/pwfile username
sudo systemctl restart mosquitto
```

#### Add username and password to config file of the project

Open config.json and add your username and password in the appropriate fields

``` json
"mqtt_username": "myUsername"
"mqtt_password": "myPassword"

```
