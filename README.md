# Battleship

first change the IP address at the line 37 from the index.html file by the one from your computer:
socket = new WebSocket("ws://*your_ip_adress*/ws");

then run *go run main.go* in your terminal and now you can connect to the server from any device connected to the same network by going to *http://your_ip_address:8080/*
Also make sure that you're allowing connexion from other device on the port 8080 of your device
