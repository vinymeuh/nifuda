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

TEMPLATE_JPEG="${SCRIPT_DIR}/data/jpeg.jpg"

# values for tag anonymisation 
ARTIST="Pink Panther"
SERIAL="123456789"

if [ $# -ne 1 ]; then
    echo "Usage: $(basename $0) IMAGE"
    exit 1
fi

IMAGE_SRC="$1"
IMAGE_TGT="${SCRIPT_DIR}/data/TEST_$(basename ${IMAGE_SRC})"

cp "${TEMPLATE_JPEG}" "${IMAGE_TGT}"

exiftool -TagsFromFile "${IMAGE_SRC}" "${IMAGE_TGT}" -overwrite_original \
    -Artist="${ARTIST}" -Copyright="${ARTIST}" -SerialNumber="${SERIAL}" \
    --ThumbnailImage -makernotes:all= -xmp:all=
