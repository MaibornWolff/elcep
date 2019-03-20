#!/bin/sh

get_random () {
    tr -cd 0-${1} </dev/urandom | head -c ${2}
}

sleep_random_time() {
    # wait longer, sometimes (makes the logs more interesting)
    if [[ "$(get_random 9 2)" -lt "70" ]]; then
        sleep $(get_random 9 1)
    fi
    sleep 0.$(get_random 9 2)
}

generate_log() {
    echo "{\
        \"@timestamp\": \"$1\",\
        \"log\": \"Exception this is a sample log message\",\
        \"bucket\": true,\
        \"key1\": $(get_random 3 1),\
        \"key2\": $(get_random 3 1),\
        \"somefield\": \"somevalue\"\
    }"
}

push_log() {
    # 2019-01-14T10:15:43.000Z
    date=$(date +'%FT%T.')$(adjtimex | awk '/(time.tv_usec):/ { printf("%06d\n", $2) }' | head -c3)Z
    log=$(generate_log ${date})
    curl -s -X POST "http://elasticsearch:9200/myindex/mydoc" -H 'Content-Type: application/json' -d"$log" || sleep 5
}

# uncomment to delete all logs on elastic
# curl -XDELETE "http://elasticsearch:9200/myindex"

while true
do
    sleep_random_time
    push_log
done
