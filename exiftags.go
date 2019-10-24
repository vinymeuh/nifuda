// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import "encoding/binary"

// ExifTags contains Exif tags.
// Fields are defined in order they appeared in chapter 4.6 of Exif 2.31
type ExifTags struct {
	/****************************************************/
	/* TIFF Rev. 6.0 Attribute Information Used in Exif */
	/****************************************************/
	// A. Tags relating to image data structure
	ExifIFD             uint32
	GpsIFD              uint32
	InteroperabilityIFD uint32
	//ImageWidth SHORT or LONG :(
	//ImageLength SHORT or LONG :(
	BitsPerSample             uint16
	Compression               string
	PhotometricInterpretation string
	SamplesPerPixel           uint16
	PlanarConfiguration       string
	YCbCrPositioning          string
	ResolutionUnit            string
	// B. Tags relating to recording offset
	// C. Tags relating to image data characteristics
	// D. Other tags
	DateTime         string
	ImageDescription string
	Make             string
	Model            string
	Software         string
	Artist           string
	Copyright        string
	/**********************************/
	/* Exif IFD Attribute Information */
	/**********************************/
	// A. Tags Relating to Version
	ExifVersion     string
	FlashpixVersion string
	// B. Tag Relating to Image Data Characteristics
	// C. Tags Relating to Image Configuration
	// D. Tags Relating to User Information
	// E. Tag Relating to Related File Information
	// F. Tags Relating to Date and Time
	DateTimeOriginal    string
	DateTimeDigitized   string
	OffsetTime          string
	OffsetTimeOriginal  string
	OffsetTimeDigitized string
	SubSecTime          string
	SubSecTimeOriginal  string
	SubSecTimeDigitized string
	// G. Tags Relating to Picture-Taking Conditions
	ExposureProgram     string
	SpectralSensitivity string
	MeteringMode        string
}

func parseIFDTagsAsExifTag(ifd *tiffIFD, bo binary.ByteOrder) ExifTags {
	var x ExifTags
	for _, ifdtag := range ifd.tags {
		switch ifdtag.id {
		case 34665: // Exif IFD
			x.ExifIFD = ifdtag.longToUint32(bo)[0]
		case 34853: // GPS IFD
			x.GpsIFD = ifdtag.longToUint32(bo)[0]
		case 258: // BitsPerSample
			x.BitsPerSample = ifdtag.shortToUint16(bo)[0]
		case 259: // Compression
			switch ifdtag.shortToUint16(bo)[0] {
			case 1:
				x.Compression = "uncompressed"
			case 6:
				x.Compression = "JPEG compression"
			}
		case 262: // PhotometricInterpretation
			switch ifdtag.shortToUint16(bo)[0] {
			case 2:
				x.PhotometricInterpretation = "RGB"
			case 6:
				x.PhotometricInterpretation = "YCbCr"
			}
		case 274: // Orientation
		case 277: // SamplesPerPixel
			x.SamplesPerPixel = ifdtag.shortToUint16(bo)[0]
		case 284: // PlanarConfiguration
			switch ifdtag.shortToUint16(bo)[0] {
			case 1:
				x.PlanarConfiguration = "chunky format"
			case 2:
				x.PlanarConfiguration = "planar format"
			}
		case 531: // YCbCrPositioning
			switch ifdtag.shortToUint16(bo)[0] {
			case 1:
				x.YCbCrPositioning = "centered"
			case 2:
				x.YCbCrPositioning = "co-sited"
			}
		case 296: // ResolutionUnit
			switch ifdtag.shortToUint16(bo)[0] {
			case 2:
				x.ResolutionUnit = "inches"
			case 3:
				x.ResolutionUnit = "centimeters"
			}
		case 306: // DateTime
			x.DateTime = ifdtag.asciiToString()
		case 270: // ImageDescription
			x.ImageDescription = ifdtag.asciiToString()
		case 271: // Make
			x.Make = ifdtag.asciiToString()
		case 272: // Model
			x.Model = ifdtag.asciiToString()
		case 305: // Software
			x.Software = ifdtag.asciiToString()
		case 315: // Artist
			x.Artist = ifdtag.asciiToString()
		case 33432: // Copyright
			x.Copyright = ifdtag.asciiToString()
		case 36864: // ExifVersion
			x.ExifVersion = ifdtag.undefinedToString()
		case 40960: // FlashpixVersion
			x.FlashpixVersion = ifdtag.undefinedToString()
		case 36867: // DateTimeOriginal
			x.DateTimeOriginal = ifdtag.asciiToString()
		case 36868: // DateTimeDigitized
			x.DateTimeDigitized = ifdtag.asciiToString()
		case 36880: // OffsetTime
			x.OffsetTime = ifdtag.asciiToString()
		case 36881: // OffsetTimeOriginal
			x.OffsetTimeOriginal = ifdtag.asciiToString()
		case 36882: // OffsetTimeDigitized
			x.OffsetTimeDigitized = ifdtag.asciiToString()
		case 37520: // SubsecTime
			x.SubSecTime = ifdtag.asciiToString()
		case 37521: // SubsecTimeOriginal
			x.SubSecTimeOriginal = ifdtag.asciiToString()
		case 37522: // SubSecTimeDigitized
			x.SubSecTimeDigitized = ifdtag.asciiToString()
		case 34850: // ExposureProgram
			switch ifdtag.shortToUint16(bo)[0] {
			case 0:
				x.ExposureProgram = "not defined"
			case 1:
				x.ExposureProgram = "manual"
			case 2:
				x.ExposureProgram = "normal program"
			case 3:
				x.ExposureProgram = "aperture priority"
			case 4:
				x.ExposureProgram = "shutter priority"
			case 5:
				x.ExposureProgram = "creative program"
			case 6:
				x.ExposureProgram = "action program"
			case 7:
				x.ExposureProgram = "portrait mode"
			case 8:
				x.ExposureProgram = "landscape mode"
			}
		case 34852: // SpectralSensitivity
			x.SpectralSensitivity = ifdtag.asciiToString()
		case 34855: // PhotographicSensitivity (ISO 12232)
		case 34856: // OECF (ISO 14524)
		case 34864: // SensitivityType (ISO 12232)
		case 34865: // StandardOutputSensitivity (ISO 12232)
		case 37383: // MeteringMode
			switch ifdtag.shortToUint16(bo)[0] {
			case 0:
				x.MeteringMode = "unknown"
			case 1:
				x.MeteringMode = "average"
			case 2:
				x.MeteringMode = "center-weighted average"
			case 3:
				x.MeteringMode = "spot"
			case 4:
				x.MeteringMode = "multispot"
			case 5:
				x.MeteringMode = "pattern"
			case 6:
				x.MeteringMode = "partial"
			case 255:
				x.MeteringMode = "other"
			}
		}
	}
	return x
}

// const (
// 	tagExifIfd             = 34665
// 	tagGpsIfd              = 34853
// 	tagInteroperabilityIfd = 40965
// )

// exifTags contains Exif tags definitions (in order they appeared in chapter 4.6 of Exif 2.31)
// var dictExif = tagDictionary{
// 	/*********************/
// 	/* Exif-specific IFD */
// 	/*********************/
// 	tagExifIfd:             {Name: "Exif IFD"},
// 	tagGpsIfd:              {Name: "GPS IFD"},
// 	tagInteroperabilityIfd: {Name: "Interoperability IFD"},

// 	/****************************************************/
// 	/* TIFF Rev. 6.0 Attribute Information Used in Exif */
// 	/****************************************************/
// 	// A. Tags relating to image data structure
// 	256: {Name: "ImageWidth"},
// 	257: {Name: "ImageLength"},
// 	258: {Name: "BitsPerSample"},
// 	259: {Name: "Compression"},
// 	262: {Name: "PhotometricInterpretation"},
// 	274: {Name: "Orientation"},
// 	277: {Name: "SamplesPerPixel"},
// 	284: {Name: "PlanarConfiguration"},
// 	530: {Name: "YCbCrSubSampling"},
// 	531: {Name: "YCbCrPositioning"},
// 	282: {Name: "XResolution"},
// 	283: {Name: "YResolution"},
// 	296: {Name: "ResolutionUnit"},
// 	// B. Tags relating to recording offset
// 	273: {Name: "StripOffsets"},
// 	278: {Name: "RowsPerStrip"},
// 	279: {Name: "StripByteCounts"},
// 	513: {Name: "JPEGInterchangeFormat"},
// 	514: {Name: "JPEGInterchangeFormatLength"},
// 	// C. Tags relating to image data characteristics
// 	301: {Name: "TransferFunction"},
// 	318: {Name: "WhitePoint"},
// 	319: {Name: "PrimaryChromaticities"},
// 	529: {Name: "YCbCrCoefficients"},
// 	532: {Name: "ReferenceBlackWhite"},
// 	// D. Other tags
// 	306:   {Name: "DateTime"},
// 	270:   {Name: "ImageDescription"},
// 	271:   {Name: "Make"},
// 	272:   {Name: "Model"},
// 	305:   {Name: "Software"},
// 	315:   {Name: "Artist"},
// 	33432: {Name: "Copyright"},

// 	/**********************************/
// 	/* Exif IFD Attribute Information */
// 	/**********************************/
// 	// A. Tags Relating to Version
// 	36864: {Name: "ExifVersion"},
// 	40960: {Name: "FlashpixVersion"},
// 	// B. Tag Relating to Image Data Characteristics
// 	40961: {Name: "ColorSpace"},
// 	42240: {Name: "Gamma"},
// 	// C. Tags Relating to Image Configuration
// 	37121: {Name: "ComponentsConfiguration"},
// 	37122: {Name: "CompressedBitsPerPixel"},
// 	40962: {Name: "PixelXDimension"},
// 	40963: {Name: "PixelYDimension"},
// 	// D. Tags Relating to User Information
// 	37500: {Name: "MakerNote"},
// 	37510: {Name: "UserComment"},
// 	// E. Tag Relating to Related File Information
// 	40964: {Name: "RelatedSoundFile"},
// 	// F. Tags Relating to Date and Time
// 	36867: {Name: "DateTimeOriginal"},
// 	36868: {Name: "DateTimeDigitized"},
// 	36880: {Name: "OffsetTime"},
// 	36881: {Name: "OffsetTimeOriginal"},
// 	36882: {Name: "OffsetTimeDigitized"},
// 	37520: {Name: "SubSecTime"},
// 	37521: {Name: "SubSecTimeOriginal"},
// 	37522: {Name: "SubSecTimeDigitized"},
// 	// G. Tags Relating to Picture-Taking Conditions
// 	33434: {Name: "ExposureTime"},
// 	33437: {Name: "FNumber"},
// 	34850: {Name: "ExposureProgram"},
// 	34852: {Name: "SpectralSensitivity"},
// 	34855: {Name: "PhotographicSensitivity"},
// 	34856: {Name: "OECF"},
// 	34864: {Name: "SensitivityType"},
// 	34865: {Name: "StandardOutputSensitivity"},
// 	34866: {Name: "RecommendedExposureIndex"},
// 	34867: {Name: "ISOSpeed"},
// 	34868: {Name: "ISOSpeedLatitudeyyy"},
// 	34869: {Name: "ISOSpeedLatitudezzz"},
// 	37377: {Name: "ShutterSpeedValue"},
// 	37378: {Name: "ApertureValue"},
// 	37379: {Name: "BrightnessValue"},
// 	37380: {Name: "ExposureBiasValue"},
// 	37381: {Name: "MaxApertureValue"},
// 	37382: {Name: "SubjectDistance"},
// 	37383: {Name: "MeteringMode"},
// 	37384: {Name: "LightSource"},
// 	37385: {Name: "Flash"},
// 	37386: {Name: "FocalLength"},
// 	37396: {Name: "SubjectArea"},
// 	41483: {Name: "FlashEnergy"},
// 	41484: {Name: "SpatialFrequencyResponse"},
// 	41486: {Name: "FocalPlaneXResolution"},
// 	41487: {Name: "FocalPlaneYResolution"},
// 	41488: {Name: "FocalPlaneResolutionUnit"},
// 	41492: {Name: "SubjectLocation"},
// 	41493: {Name: "ExposureIndex"},
// 	41495: {Name: "SensingMethod"},
// 	41728: {Name: "FileSource"},
// 	41729: {Name: "SceneType"},
// 	41730: {Name: "CFAPattern"},
// 	41985: {Name: "CustomRendered"},
// 	41986: {Name: "ExposureMode"},
// 	41987: {Name: "WhiteBalance"},
// 	41988: {Name: "DigitalZoomRatio"},
// 	41989: {Name: "FocalLengthIn35mmFilm"},
// 	41990: {Name: "SceneCaptureType"},
// 	41991: {Name: "GainControl"},
// 	41992: {Name: "Contrast"},
// 	41993: {Name: "Saturation"},
// 	41994: {Name: "Sharpness"},
// 	41995: {Name: "DeviceSettingDescription"},
// 	41996: {Name: "SubjectDistanceRange"},
// 	// G2. Tags Relating to shooting situation
// 	37888: {Name: "Temperature"},
// 	37889: {Name: "Humidity"},
// 	37890: {Name: "Pressure"},
// 	37891: {Name: "WaterDepth"},
// 	37892: {Name: "Acceleration"},
// 	37893: {Name: "CameraElevationAngle"},
// 	// H. Other Tags
// 	42016: {Name: "ImageUniqueID"},
// 	42032: {Name: "CameraOwnerName"},
// 	42033: {Name: "BodySerialNumber"},
// 	42034: {Name: "LensSpecification"},
// 	42035: {Name: "LensMake"},
// 	42036: {Name: "LensModel"},
// 	42037: {Name: "LensSerialNumber"},

// 	/**********************/
// 	/* Other Private Tags */
// 	/**********************/
// 	// Ratings tag used by Windows
// 	18246: {Name: "Image.Rating"},
// 	18249: {Name: "Image.RatingPercent"},
// }
