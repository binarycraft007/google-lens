package lens

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var ErrUnsupportImage = errors.New("unsupported image")
var ErrAFCallbackNotFound = errors.New("AF_initDataCallback not found")

type ChromeVerion struct {
	Major  string
	Suffix string
}

type Lens struct {
	config *Config
	client *http.Client
}

type Config struct {
	verion    ChromeVerion
	userAgent string
}

type Option func(*Config)

func WithUserAgent(ua string) Option {
	return func(c *Config) {
		c.userAgent = ua
	}
}

func WithChromeVerson(ver ChromeVerion) Option {
	return func(c *Config) {
		c.verion = ver
	}
}

func NewLens(opts ...Option) *Lens {
	userAgent := "Mozilla/5.0 " +
		"(Windows NT 10.0; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/122.0.0.0 Safari/537.36"
	config := &Config{
		verion: ChromeVerion{
			Major:  "122",
			Suffix: ".0.6261.112",
		},
		userAgent: userAgent,
	}
	for _, opt := range opts {
		opt(config)
	}
	return &Lens{
		config: config,
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (l *Lens) Scan(r io.ReadSeeker) (*[]byte, error) {
	br := base64.NewDecoder(base64.StdEncoding, r)
	image, ext, err := image.Decode(br)
	if err != nil {
		return nil, err
	}

	mime, ok := SupportedFormats[ext]
	if !ok {
		return nil, ErrUnsupportImage
	}
	r.Seek(0, io.SeekStart)
	b, err := io.ReadAll(br)
	if err != nil {
		return nil, err
	}
	fetchOptions, err := imageToFormData(b, image.Bounds(), mime)
	if err != nil {
		return nil, err
	}
	h := l.GenerateHeader()
	_ = l.Fetch(h, fetchOptions)
	return &b, nil
}

type FetchOptions struct {
	formData   *bytes.Buffer
	contenType *string
}

func (l *Lens) Fetch(h *http.Header, opt *FetchOptions) error {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	params := url.Values{}
	params.Add("hl", "en-US")
	params.Add("re", "df")
	params.Add("stcs", timestamp)
	params.Add("vpw", "871")
	params.Add("vph", "919")
	params.Add("ep", "subb")
	u, _ := url.Parse(Endpoint)
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("POST", u.String(), opt.formData)
	if err != nil {
		return err
	}

	req.Header = *h
	req.Header.Set("Content-Type", *opt.contenType)
	resp, err := l.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var text string
	if resp.Header.Get("Content-Encoding") == "gzip" {
		// Create a Gzip reader to decompress the response
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		defer reader.Close()

		// Read the decompressed data
		decompressedData, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		// Process the decompressed data (e.g., print or parse)
		text = string(decompressedData)
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		text = string(body)
	}
	fmt.Println(text)

	//re := regexp.MustCompile(`AF_initDataCallback\((\{.*?\})\)`)
	//callbacks := re.FindAllString(text, -1)

	//var lensCallback string
	//for _, c := range callbacks {
	//	if containsDetectedObject(c) {
	//		lensCallback = c
	//		break
	//	}
	//}

	//if lensCallback == "" {
	//	fmt.Println(callbacks)
	//	return ErrAFCallbackNotFound
	//}

	//re = regexp.MustCompile(`AF_initDataCallback\((\{.*?\})\)`)
	//match := re.FindString(lensCallback)

	//fmt.Println(match[1])

	return nil
}

func containsDetectedObject(s string) bool {
	return strings.Contains(s, "DetectedObject")
}

func imageToFormData(b []byte, bounds image.Rectangle, mime string) (*FetchOptions, error) {
	formData := &bytes.Buffer{}
	writer := multipart.NewWriter(formData)
	// Close writer before use it in post request
	defer writer.Close()

	// Create a new form field for the file.
	image, err := CreateFormFile(writer, "encoded_image", "test.png", mime)
	if err != nil {
		return nil, err
	}
	image.Write(b)
	writer.WriteField("original_width", strconv.FormatInt(int64(bounds.Dx()), 10))
	writer.WriteField("original_height", strconv.FormatInt(int64(bounds.Dy()), 10))
	writer.WriteField("processed_image_dimensions", fmt.Sprintf(
		"%d,%d",
		bounds.Dx(),
		bounds.Dy(),
	))
	contenType := writer.FormDataContentType()
	return &FetchOptions{
		formData:   formData,
		contenType: &contenType,
	}, nil
}

// CreateFormFile is a convenience wrapper around CreatePart. It creates
// a new form-data header with the provided field name and file name.
func CreateFormFile(w *multipart.Writer, fieldname, filename, mime string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	h.Set("Content-Type", mime)
	return w.CreatePart(h)
}

func (l *Lens) GenerateHeader() *http.Header {
	h := &HttpHeader
	h.Add(
		"Sec-Ch-Ua",
		fmt.Sprintf(`"Google Chrome";v="%s"`, l.config.verion.Major),
	)
	h.Add(
		"Sec-Ch-Ua",
		fmt.Sprintf(`"Chromium";v="%s"`, l.config.verion.Major),
	)
	h.Add(
		"Sec-Ch-Ua-Full-Version-List",
		fmt.Sprintf(`"Google Chrome";v="%s"`, l.config.verion.Major),
	)
	h.Add(
		"Sec-Ch-Ua-Full-Version-List",
		fmt.Sprintf(`"Chromium";v="%s"`, l.config.verion.Major),
	)
	h.Add(
		"Sec-Ch-Ua-Full-Version",
		fmt.Sprintf(`"%s"`, l.config.verion.Major+l.config.verion.Suffix),
	)
	h.Add(
		"User-Agent",
		fmt.Sprintf(`%s`, l.config.userAgent),
	)
	return h
}
