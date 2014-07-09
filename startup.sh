#!/bin/bash
cd /usr/local/IBJts
java -cp jts.jar:total.2013.jar -Dsun.java2d.noddraw=true -Xmx512M ibgateway.GWClient .
