// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package exif

import "github.com/vinymeuh/nifuda/pkg/tiff"

// TagDictionary contains Exif tags definitions (in order they appeared in chapter 4.6 of Exif 2.31)
var TagDictionary = tiff.TagDictionary{
	/*********************/
	/* Exif-specific IFD */
	/*********************/
	34665: {Name: "Exif IFD"},
	34853: {Name: "GPS IFD"},
	40965: {Name: "Interoperability IFD"},

	/****************************************************/
	/* TIFF Rev. 6.0 Attribute Information Used in Exif */
	/****************************************************/
	// A. Tags relating to image data structure
	256: {Name: "ImageWidth"},
	257: {Name: "ImageLength"},
	258: {Name: "BitsPerSample"},
	259: {Name: "Compression"},
	262: {Name: "PhotometricInterpretation"},
	274: {Name: "Orientation"},
	277: {Name: "SamplesPerPixel"},
	284: {Name: "PlanarConfiguration"},
	530: {Name: "YCbCrSubSampling"},
	531: {Name: "YCbCrPositioning"},
	282: {Name: "XResolution"},
	283: {Name: "YResolution"},
	296: {Name: "ResolutionUnit"},
	// B. Tags relating to recording offset
	273: {Name: "StripOffsets"},
	278: {Name: "RowsPerStrip"},
	279: {Name: "StripByteCounts"},
	513: {Name: "JPEGInterchangeFormat"},
	514: {Name: "JPEGInterchangeFormatLength"},
	// C. Tags relating to image data characteristics
	301: {Name: "TransferFunction"},
	318: {Name: "WhitePoint"},
	319: {Name: "PrimaryChromaticities"},
	529: {Name: "YCbCrCoefficients"},
	532: {Name: "ReferenceBlackWhite"},
	// D. Other tags
	306:   {Name: "DateTime"},
	270:   {Name: "ImageDescription"},
	271:   {Name: "Make"},
	272:   {Name: "Model"},
	305:   {Name: "Software"},
	315:   {Name: "Artist"},
	33432: {Name: "Copyright"},

	/**********************************/
	/* Exif IFD Attribute Information */
	/**********************************/
	// A. Tags Relating to Version
	36864: {Name: "ExifVersion"},
	40960: {Name: "FlashpixVersion"},
	// B. Tag Relating to Image Data Characteristics
	40961: {Name: "ColorSpace"},
	42240: {Name: "Gamma"},
	// C. Tags Relating to Image Configuration
	37121: {Name: "ComponentsConfiguration"},
	37122: {Name: "CompressedBitsPerPixel"},
	40962: {Name: "PixelXDimension"},
	40963: {Name: "PixelYDimension"},
	// D. Tags Relating to User Information
	37500: {Name: "MakerNote"},
	37510: {Name: "UserComment"},
	// E. Tag Relating to Related File Information
	40964: {Name: "RelatedSoundFile"},
	// F. Tags Relating to Date and Time
	36867: {Name: "DateTimeOriginal"},
	36868: {Name: "DateTimeDigitized"},
	36880: {Name: "OffsetTime"},
	36881: {Name: "OffsetTimeOriginal"},
	36882: {Name: "OffsetTimeDigitized"},
	37520: {Name: "SubSecTime"},
	37521: {Name: "SubSecTimeOriginal"},
	37522: {Name: "SubSecTimeDigitized"},
	// G. Tags Relating to Picture-Taking Conditions
	33434: {Name: "ExposureTime"},
	33437: {Name: "FNumber"},
	34850: {Name: "ExposureProgram"},
	34852: {Name: "SpectralSensitivity"},
	34855: {Name: "PhotographicSensitivity"},
	34856: {Name: "OECF"},
	34864: {Name: "SensitivityType"},
	34865: {Name: "StandardOutputSensitivity"},
	34866: {Name: "RecommendedExposureIndex"},
	34867: {Name: "ISOSpeed"},
	34868: {Name: "ISOSpeedLatitudeyyy"},
	34869: {Name: "ISOSpeedLatitudezzz"},
	37377: {Name: "ShutterSpeedValue"},
	37378: {Name: "ApertureValue"},
	37379: {Name: "BrightnessValue"},
	37380: {Name: "ExposureBiasValue"},
	37381: {Name: "MaxApertureValue"},
	37382: {Name: "SubjectDistance"},
	37383: {Name: "MeteringMode"},
	37384: {Name: "LightSource"},
	37385: {Name: "Flash"},
	37386: {Name: "FocalLength"},
	37396: {Name: "SubjectArea"},
	41483: {Name: "FlashEnergy"},
	41484: {Name: "SpatialFrequencyResponse"},
	41486: {Name: "FocalPlaneXResolution"},
	41487: {Name: "FocalPlaneYResolution"},
	41488: {Name: "FocalPlaneResolutionUnit"},
	41492: {Name: "SubjectLocation"},
	41493: {Name: "ExposureIndex"},
	41495: {Name: "SensingMethod"},
	41728: {Name: "FileSource"},
	41729: {Name: "SceneType"},
	41730: {Name: "CFAPattern"},
	41985: {Name: "CustomRendered"},
	41986: {Name: "ExposureMode"},
	41987: {Name: "WhiteBalance"},
	41988: {Name: "DigitalZoomRatio"},
	41989: {Name: "FocalLengthIn35mmFilm"},
	41990: {Name: "SceneCaptureType"},
	41991: {Name: "GainControl"},
	41992: {Name: "Contrast"},
	41993: {Name: "Saturation"},
	41994: {Name: "Sharpness"},
	41995: {Name: "DeviceSettingDescription"},
	41996: {Name: "SubjectDistanceRange"},
	// G2. Tags Relating to shooting situation
	37888: {Name: "Temperature"},
	37889: {Name: "Humidity"},
	37890: {Name: "Pressure"},
	37891: {Name: "WaterDepth"},
	37892: {Name: "Acceleration"},
	37893: {Name: "CameraElevationAngle"},
	// H. Other Tags
	42016: {Name: "ImageUniqueID"},
	42032: {Name: "CameraOwnerName"},
	42033: {Name: "BodySerialNumber"},
	42034: {Name: "LensSpecification"},
	42035: {Name: "LensMake"},
	42036: {Name: "LensModel"},
	42037: {Name: "LensSerialNumber"},

	/**********************/
	/* Other Private Tags */
	/**********************/
	// Ratings tag used by Windows
	18246: {Name: "Image.Rating"},
	18249: {Name: "Image.RatingPercent"},
}
