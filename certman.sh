#!/bin/bash

./certman/certman -c b00m.config 2>&1 | funnel -app=b00m-certman &
