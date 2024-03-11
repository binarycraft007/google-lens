package lens

import "net/http"

const (
	Endpoint = "https://lens.google.com/v3/upload"
)

var SupportedFormats = map[string]string{
	//"ico":  "image/x-icon",
	"bmp":  "image/bmp",
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"tiff": "image/tiff",
	"webp": "image/webp",
	//"heic": "image/heic",
}

var HttpHeader = http.Header{
	"Accept": {
		"text/html",
		"application/xhtml+xml",
		"application/xml;q=0.9",
		"image/avif",
		"image/webp",
		"image/apng",
		"*/*;q=0.8",
		"application/signed-exchange;v=b3;q=0.7",
	},
	"Accept-Encoding":             {"gzip", "deflate", "br"},
	"Accept-Language":             {"en-US", "en;q=0.9"},
	"Cache-Control":               {"max-age=0"},
	"Origin":                      {"https://lens.google.com"},
	"Referer":                     {"https://lens.google.com\\"},
	"Sec-Ch-Ua":                   {`"Not A(Brand";"v="99"`},
	"Sec-Ch-Ua-Arch":              {`"x86"`},
	"Sec-Ch-Ua-Bitness":           {`"64"`},
	"Sec-Ch-Ua-Full-Version-List": {`"Not A(Brand";v="99.0.0.0"`},
	"Sec-Ch-Ua-Mobile":            {"?0"},
	"Sec-Ch-Ua-Model":             {`""`},
	"Sec-Ch-Ua-Platform":          {`"Windows"`},
	"Sec-Ch-Ua-Platform-Version":  {`"15.0.0"`},
	"Sec-Ch-Ua-Wow64":             {"?0"},
	"Sec-Fetch-Dest":              {"document"},
	"Sec-Fetch-Mode":              {"navigate"},
	"Sec-Fetch-Site":              {"same-origin"},
	"Sec-Fetch-User":              {"?1"},
	"Upgrade-Insecure-Requests":   {"1"},
	"X-Client-Data":               {"CIW2yQEIorbJAQipncoBCIH+ygEIlqHLAQjtocsBCPWYzQEIhaDNAQjd7M0BCMruzQEIg/DNAQjW8c0BCIDyzQEIz/TNAQiQ9c0BCLb3zQEYp+rNARib+M0BGMr4zQE="},
	//X-Client-Data: CIW2yQEIorbJAQipncoBCIH+ygEIlqHLAQjtocsBCPWYzQEIhaDNAQjd7M0BCMruzQEIg/DNAQjW8c0BCIDyzQEIz/TNAQiQ9c0BCLb3zQEYp+rNARib+M0BGMr4zQE=
	//Decoded:
	//message ClientVariations {
	//    // Active client experiment variation IDs.
	//    repeated int32 variation_id = [3300101, 3300130, 3313321, 3325697, 3330198, 3330285, 3361909, 3362821, 3372637, 3372874, 3373059, 3373270, 3373312, 3373647, 3373712, 3374006];
	//    // Active client experiment variation IDs that trigger server-side behavior.
	//    repeated int32 trigger_variation_id = [3372327, 3374107, 3374154];
	//}
}
