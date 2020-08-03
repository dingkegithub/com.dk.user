#!/bin/bash

LOG_DIR=/tmp
LOG="${LOG_DIR}/disk.log"
LOG_SIZE=52428800

MANAGE_DIR=(
"/data/user_behaviors"
"/data/user_behaviors_test"
)

mkdir -p ${LOG_DIR}

log() {
    if [ ! -f ${LOG} ]; then
        touch ${LOG}
    else
        # shellcheck disable=SC2012
        log_size=$(ls -l ${LOG} | cut -d ' ' -f5)
        too_bigger=$(echo "${log_size} >= ${LOG_SIZE}" | bc)

        if [ "${too_bigger}" -eq 1 ]; then
            rm ${LOG}
            touch ${LOG}
        fi
    fi

    echo "$(date +%Y-%m-%dT%H:%M:%S) $1" >> ${LOG}
}

remove() {
    cur_dir=$(pwd)
    dir_name=$(dirname $1)
    base_name=$(basename $1)
    # shellcheck disable=SC2164
    cd "${dir_name}"
    log "INFO remove ${dir_name} ${base_name}"
    sep=$(echo "${base_name}" | grep '/')
    if [ -n "${sep}" ]; then
        log "ERROR can't rm abs path ${dir_name}, ${base_name}"
    else
        rm -rf ${base_name}
    fi
    # shellcheck disable=SC2164
    cd "${cur_dir}"
}

clear_dir() {
    disk_dir=$1
    # shellcheck disable=SC2045
    for date_name in $(ls "${disk_dir}")
    do
      date_record_dir="${disk_dir}/${date_name}"
      if [ ! -d "${date_record_dir}" ]; then
        # shellcheck disable=SC2154
        log "ERROR record file not exist ${record_dir}"
        continue
      fi

      last_modify_date=$(stat -c %Y "${date_record_dir}")
      last_date_format=$(date +'%Y%m%d' -d @"${last_modify_date}")
      cur_date_format=$(date +"%Y%m%d")

      timeout=$(echo $(( ($(date -d "${cur_date_format}" +%s) - $(date -d "${last_date_format}" +%s))/(24*60*60) )))

      log "INFO ${date_name} timeout: ${cur_date_format} - ${last_date_format} = ${timeout}"
      if [ ${timeout} -ge 2 ]; then
        log "INFO timeout will remove ${date_record_dir}"
        #remove "${date_record_dir}"
      fi
            continue
    done
}

function main() {
    for dir in "${MANAGE_DIR[@]}"
    do
      clear_dir "${dir}"
    done
}

log "INFO start check disk"
main
log "INFO end   check disk"