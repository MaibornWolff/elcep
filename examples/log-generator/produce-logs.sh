#!/bin/sh

get_random_number () {
    tr -cd 0-9 </dev/urandom | head -c 2
}

sleep_random_time() {
    # wait longer, sometimes (makes the logs more interesting)
    if [ "$(get_random_number)" -lt "70" ]; then
        sleep $(get_random_number | head -c 1)
    fi
    sleep 0.$(get_random_number)
}

generate_log() {
    echo "{\
        \"@timestamp\": \"$1\",\
        \"log\": \"Exception this is a sample log message\",\
        \"somefield\": \"somevalue\"\
    }"
}

push_log() {
    # 2019-01-14T10:15:43.000Z
    date=$(date +'%FT%T.')$(adjtimex | awk '/(time.tv_usec):/ { printf("%06d\n", $2) }' | head -c3)Z
    log=$(generate_log $date)
    curl -s -X POST "http://elasticsearch:9200/myindex/mydoc" -H 'Content-Type: application/json' -d"$log" || sleep 5
}

while true
do
    sleep_random_time
    push_log
done
