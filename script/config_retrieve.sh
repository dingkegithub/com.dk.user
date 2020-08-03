#!/bin/bash

key='config_retrieve'
year=$(date +%Y -d "1 days ago")
month=$(date +%m -d "1 days ago")
day=$(date +%d -d "1 days ago")
mac=$(/sbin/ifconfig | grep eth0 | awk -F ' ' '{print $5}' | sed 's/://g' | tr 'a-z' 'A-Z')
logFile='/tmp/config_retrieve.log'

oss_dir="oss://gz-log/${year}-${month}-${day}/${key}"
cos_dir="preprocess/gou_zheng/config_retrieve/business_dt=${year}-${month}-${day}"
source /root/.pyenv/versions/gzsrv/bin/activate

if [ -n "$1" ]; then
  if [ -z "$2" ] || [ -z "$3" ]; then
    echo "need year, month, day"
    exit 0
  fi
  year=$1
  month=$2
  day=$3
fi

cos_path='preprocess/gou_zheng'
oss_path='oss://gz-log'
config_dir="/data/user_behaviors/${year}-${month}-${day}/${mac}/${key}/v1.1"

if [ ! -f "${logFile}" ]; then
  touch "${logFile}"
else
  l=$(cat "${logFile}" | wc -l)
  if [ "${l}" -gt 1000 ]; then
    rm ${logFile}
    touch "${logFile}"
  fi
fi

function log() {
    echo "$(date +%Y-%m-%d' '%H:%M:%S) $1" >> ${logFile}
}

if [ ! -d "${config_dir}" ]; then
  log "not exist: ${config_dir}"
fi


for fileName in $(ls "${config_dir}")
do
  zipFile=""
  size=${#fileName}
  suffix=${fileName##*.}

  zipFileName=""
  dstOssFileName=""
  dstCosFileName=""

  if [ "${suffix}" != "log" ]; then
    if [ "${suffix}" != "bz2" ] || [ ${size} -ne 25 ]; then
      continue
    fi

    zipFileName="${fileName}"
    dstOssFileName="${mac}-${zipFileName}"
    dstCosFileName="${mac}-v1.1-${zipFileName}"
  else
    bzip2 "${config_dir}/${fileName}"

    if [ ! -f "${config_dir}/${fileName}.bz2" ]; then
      log "bzip2 compress failed, ${config_dir}/${fileName}.bz2"
      continue
    fi

    zipFileName="${fileName}.bz2"
    dstOssFileName="${mac}-${zipFileName}"
    dstCosFileName="${mac}-v1.1-${zipFileName}"
  fi

  zipFile="${config_dir}/${zipFileName}"
  if [ ! -f "${zipFile}" ]; then
    log "compress file not exist: ${zipFile}"
    continue
  else
    log "/bin/ossutil64  --config-file /root/.ossutilconfig cp ${zipFile} ${oss_dir}/${dstOssFileName}"
    /bin/ossutil64  --config-file /root/.ossutilconfig cp "${zipFile}" "${oss_dir}/${dstOssFileName}" -f

    log "coscmd upload ${zipFile} ${cos_dir}/${dstCosFileName}"
    /root/.pyenv/versions/gzsrv/bin/coscmd upload "${zipFile}" "${cos_dir}/${dstCosFileName}"
  fi
done