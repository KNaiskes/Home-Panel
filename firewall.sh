#!/bin/bash


if [ "$EUID" -ne 0 ] 
then
	echo "You must me root, to run this script"
	exit
fi

sudo iptables -I INPUT -p tcp --dport 8080 -j ACCEPT
sudo iptables -I INPUT -p tcp --dport 1883 -j ACCEPT

sudo ip6tables -I INPUT -p tcp --dport 1883 -j ACCEPT
