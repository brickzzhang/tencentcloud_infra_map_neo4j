package register

import (
	"tencentcloud/product/base"
	"tencentcloud/product/cvm"
	"tencentcloud/product/image"
	"tencentcloud/product/region"
)

var ServiceMap = make(map[string]base.CloudService)

func init() {
	ServiceMap["CVM"] = new(cvm.CvmService)
	ServiceMap["Region"] = new(region.RegionService)
	ServiceMap["Image"] = new(image.ImageService)
}
