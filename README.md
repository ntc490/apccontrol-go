# Go Based APC control application

This is a CLI based application to control an
APC networkable power strip via telnet.  It only
supports a small variety of models, and was ported
from a perl script I wrote many years ago.

Users can set aliases to port numbers to make
controlling devices easier.

The config file lives in ~/.config/apc/config.yml
by default.  The file is yaml format.  Users can
use the config file, or specify all information
on the command line.  The apc application needs
a hostname, user, and password to complete most
operations.


    user: apc
    password: password
    hostname: apc
