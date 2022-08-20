#!/bin/bash

DRIVERS_FILE=pkg/drivers/drivers.go
DRIVERS_FILE_SLIM=pkg/drivers/drivers_slim.go

listDrivers() {
    grep 'DriverName =' $DRIVERS_FILE | awk -F '"' '{print $2}'
}

varnameForDriver() {
    grep 'DriverName =' $DRIVERS_FILE | \
    grep "\"$1\"" | \
    awk -F ' ' '{print $1}'
}

containsElement() {
    local e
    for e in "${@:2}"; do [[ "$e" == "$1" ]] && return 0; done
    return 1
}

checkDriversEnabled() {
    ENABLED_DRIVERS=($(listDrivers))
    DESIRED_DRIVERS=($@)
    for driver in "${DESIRED_DRIVERS[@]}"; do
        if ! containsElement "$driver" "${ENABLED_DRIVERS[@]}"; then
            echo "Driver $driver is not enabled"
            exit 1
        fi
    done
}

disabledDrivers() {
    local DISABLED=()
    for driver in "${ENABLED_DRIVERS[@]}"; do
        sdriver=$(echo $driver | awk -F'-' '{print $1}')
        if ! containsElement "$driver" "${DESIRED_DRIVERS[@]}"; then
            DISABLED+=("$driver")
        fi
    done
    echo "${DISABLED[@]}"
}

disabledImports() {
    local DISABLED=()
    for driver in $(disabledDrivers); do
        sdriver=$(echo $driver | awk -F'-' '{print $1}')
        found="false"
        for edriver in ${DESIRED_DRIVERS[@]}; do
            sedriver=$(echo $edriver | awk -F'-' '{print $1}')
            if [ "$sdriver" == "$sedriver" ]; then
                found="true"
                break
            fi
        done
        if [ "$found" == "false" ]; then
            DISABLED+=("$sdriver")
        fi
    done
    echo "${DISABLED[@]}" | tr ' ' '\n' | sort | uniq | xargs
}

grepOutDrivers() {
    local rmdrivers=($@)
    TEMP_FILE=$(mktemp)
    cat $DRIVERS_FILE > $TEMP_FILE
    for driver in "${rmdrivers[@]}"; do
        sdriver=`echo $driver | awk -F'-' '{print $1}'`
        grep -v "DriverName = \"$driver\"" $TEMP_FILE > ${TEMP_FILE}.new && mv ${TEMP_FILE}.new $TEMP_FILE
        sed -e "/case $(varnameForDriver $driver):/,+1d" $TEMP_FILE > ${TEMP_FILE}.new && mv ${TEMP_FILE}.new $TEMP_FILE
    done
    for driver in `disabledImports`; do
        sdriver=`echo $driver | awk -F'-' '{print $1}'`
        grep -v "pushx/drivers/$sdriver" $TEMP_FILE > ${TEMP_FILE}.new && mv ${TEMP_FILE}.new $TEMP_FILE
    done
    
    mv $DRIVERS_FILE $DRIVERS_FILE.bak
    mv $TEMP_FILE $DRIVERS_FILE_SLIM
}

list() {
    listDrivers
}

build() {
    DRIVERS=($@)
    if [ ${#DRIVERS[@]} -eq 0 ]; then
        DRIVERS=($(listDrivers))
    fi
    checkDriversEnabled "${DRIVERS[@]}"
    grepOutDrivers "$(disabledDrivers)"
    make
    restore
    rm -f $DRIVERS_FILE_SLIM
}

restore() {
    mv $DRIVERS_FILE.bak $DRIVERS_FILE
}

main() {
case "$1" in
    list)
        list
        ;;
    build)
        build "${@:2}"
        ;;
    *)
        echo "Usage: $0 {list|build [drivers]}"
        exit 1
        ;;
esac
}
main "$@"