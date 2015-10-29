#!/bin/sh -e

echo "Running smoke test on ELB: http://godspeed-hack-1631427837.ap-southeast-2.elb.amazonaws.com"
response=$(curl -s -o /dev/null -w %{http_code} http://godspeed-hack-1631427837.ap-southeast-2.elb.amazonaws.com)
echo "Received response code: $response"

if [ "$response" -ne "200" ]; then
	echo "Received non-200 exit code: $response"
	exit 1
fi

echo "Smoke test passed!"

exit 0
