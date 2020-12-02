#!/bin/bash

inputDIR="/sample/path"
inputTS="sample-2020-10-21.ts"
inputVDR="./sample.vdr"

cnt=0
times=()
IFS=$'\n'
# parse input file(skip line 1)
for line in $(tail -n +2 $inputVDR | cut -d' ' -f 1); do
    times+=("${line}")
done

# read times[0] - times[n-2]
for i in $(seq 0 $(expr ${#times[*]} - 2)); do
    # get length(sec) / sec = (time2 - time1)[cast -> sec]
    if [ $(expr $i % 2) -eq 0 ]; then
        tmpLabel=$(expr $i / 2)
        sec=$(echo "scale=2; $(date '+%s.%2N' -d ${times[i + 1]}) - $(date '+%s.%2N' -d ${times[i]})" | bc)
        # entrypoint : length(seconds)
        echo "[$tmpLabel]${times[i]} : $sec"
        ffmpeg -i $inputDIR/$inputTS -ss ${times[i]} -t $sec -c copy $inputDIR/output_$tmpLabel.ts
        echo "file '$inputDIR/output_$tmpLabel.ts'" >>$inputDIR/outputs.txt
    fi
done

ffmpeg -f concat -i $inputDIR/outputs.txt -c copy "CMCUT_$inputTS"
