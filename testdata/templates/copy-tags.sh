#!/usr/bin/env bash
#
# Create a test image copying tags from a existing one using exiftool.
#
# Anonymize:
#  * Author & Copyright
#  * Serial Number
#
# Remove:
#  * Thumbnail embbeded image
#  * Makernotes
#  * XMP tags

SCRIPT_DIR=$(cd $(dirname $0) && pwd)

# values for tag anonymisation 
ARTIST="Pink Panther"
SERIAL="123456789"

if [ $# -ne 2 ]; then
    echo "Usage: $(basename $0) EMPTY_IMAGE EXIF_IMAGE_SOURCE"
    exit 1
fi

EMPTY_IMAGE=$1
EXIF_IMAGE_SOURCE=$2

IMAGE_TGT="$(dirname ${EXIF_IMAGE_SOURCE})/TEST_$(basename ${EXIF_IMAGE_SOURCE})"

cp "${EMPTY_IMAGE}" "${IMAGE_TGT}"

exiftool -TagsFromFile "${EXIF_IMAGE_SOURCE}" "${IMAGE_TGT}" -overwrite_original \
    -Artist="${ARTIST}" -Copyright="${ARTIST}" -SerialNumber="${SERIAL}" \
    --ThumbnailImage -makernotes:all= -xmp:all=
