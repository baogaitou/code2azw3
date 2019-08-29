#!/bin/sh


function buildAzw3() {
    echo -e "\033[36m -- Building ${1}.azw2 \033[0m"
    ~/Downloads/Dev/Kindle/KindleGen_Mac_i386_v2_9/kindlegen -dont_append_source -c2 -o ${1}.azw3 ./${1}/${1}.opf
}


if [ x$1 == x ]; then
    echo "Useage: ./build.sh -n xxxxxx"
    echo ""
fi

while getopts "n:" opt; do
    case $opt in
    n)
        echo -e "\033[36m -- Fetching certs for: ${OPTARG} \033[0m"
        buildAzw3 ${OPTARG}
    ;;
    \?)
        echo ""
        echo -e "\033[31m -- [X] Invalid option: -${OPTARG} \033[0m"
    ;;
    esac
done
