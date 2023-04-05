# 1. Build a simple http server on Ubuntu

1. Install apache2

    `sudo apt-get update`

    `sudo apt-get install apache2`

2. After apache2 is successfully installed, we can see an **index.html** in the **/var/www/** directory. We just need to restart the apache2 service, and then the devices in the same LAN can visit this server according to the specific IP address.

    `sudo /etc/init.d/apache2 restart`

    And command `ifconfig` cna be used to obtain the IP address.

3. Manage files

    First of all, we need delete the default **index.html**.

    Then we can add, delete and modify the files in the folder  **/var/www/html**.
