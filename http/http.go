package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

var globalHeader *http.Header
var globalCookies []*http.Cookie
var defaultHTTPClient = http.DefaultClient

// NewOption 创建一个 NewOption
func NewOption() *Option {
	return &Option{
		Header:  &http.Header{},
		Cookies: []*http.Cookie{},
	}
}

// SetClient 设置默认的 Client
func SetClient(client *http.Client) {
	defaultHTTPClient = client
}

// GetClient 获取默认的 Client
func GetClient() *http.Client {
	return defaultHTTPClient
}

// SetHeader 设置全局的 header
func SetHeader(key, value string) {
	if globalHeader == nil {
		globalHeader = &http.Header{}
	}
	globalHeader.Set(key, value)
}

// SetCookie 设置全局的 Cookie
func SetCookie(cookie *http.Cookie) {
	if globalHeader == nil {
		globalCookies = []*http.Cookie{}
	}
	globalCookies = append(globalCookies, cookie)
}

// ClearHeader 清空全局 header
func ClearHeader() {
	globalHeader = nil
}

// ClearCookies 清空全局 Cookies
func ClearCookies() {
	globalCookies = nil
}

// Option options
type Option struct {
	Header  *http.Header
	Cookies []*http.Cookie
}

//Get get 请求
func Get(uri string, options ...*Option) ([]byte, error) {
	reqest, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	loadOption(reqest, options...)

	response, err := defaultHTTPClient.Do(reqest)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

//GetJSON 请求JSON
func GetJSON(uri string, dest interface{}, options ...*Option) error {
	var err error

	data, err := Get(uri, options...)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		return fmt.Errorf("json unmarshal http response body error : uri=%v , err=%v", uri, err)
	}

	return nil
}

//GetXML 请求XML
func GetXML(uri string, dest interface{}, options ...*Option) error {
	var err error

	data, err := Get(uri, options...)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, dest)
	if err != nil {
		return fmt.Errorf("xml unmarshal http response body error : uri=%v , err=%v", uri, err)
	}

	return nil
}

//MultipartFormField 保存文件或其他字段信息
type MultipartFormField struct {
	IsFile    bool
	Fieldname string
	Value     []byte
	Filename  string
}

//PostFile 上传文件
func PostFile(uri, fieldname, filename string, options ...*Option) ([]byte, error) {
	fields := []MultipartFormField{
		{
			IsFile:    true,
			Fieldname: fieldname,
			Filename:  filename,
		},
	}
	return PostMultipartForm(uri, fields)
}

//PostXML perform a HTTP/POST request with XML body
func PostXML(uri string, obj interface{}, options ...*Option) ([]byte, error) {
	xmlData, err := xml.Marshal(obj)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(xmlData)
	response, err := defaultHTTPClient.Post(uri, "application/xml;charset=utf-8", body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

//Post post 数据请求
func Post(uri string, contentType string, body io.Reader, options ...*Option) ([]byte, error) {
	reqest, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	loadOption(reqest, options...)

	reqest.Header.Set("Content-Type", contentType)
	response, err := defaultHTTPClient.Do(reqest)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}

	return ioutil.ReadAll(response.Body)
}

//PostJSON post json 数据请求
func PostJSON(uri string, obj interface{}, options ...*Option) ([]byte, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)

	body := bytes.NewBuffer(jsonData)

	return Post(uri, "application/json;charset=utf-8", body, options...)
}

//PostMultipartForm 上传文件或其他多个字段
func PostMultipartForm(uri string, fields []MultipartFormField, options ...*Option) (respBody []byte, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for _, field := range fields {
		if field.IsFile {
			fileWriter, e := bodyWriter.CreateFormFile(field.Fieldname, field.Filename)
			if e != nil {
				err = fmt.Errorf("error writing to buffer , err=%v", e)
				return
			}

			fh, e := os.Open(field.Filename)
			if e != nil {
				err = fmt.Errorf("error opening file , err=%v", e)
				return
			}
			defer fh.Close()

			if _, err = io.Copy(fileWriter, fh); err != nil {
				return
			}
		} else {
			partWriter, e := bodyWriter.CreateFormField(field.Fieldname)
			if e != nil {
				err = e
				return
			}
			valueReader := bytes.NewReader(field.Value)
			if _, err = io.Copy(partWriter, valueReader); err != nil {
				return
			}
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	return Post(uri, contentType, bodyBuf, options...)
}

func loadOption(reqest *http.Request, options ...*Option) {

	// 拼接全局的 option
	integralOptions := append(options, &Option{
		Header:  globalHeader,
		Cookies: globalCookies,
	})

	if integralOptions == nil || reqest == nil {
		return
	}

	for _, option := range integralOptions {
		if option.Header != nil {
			for key, values := range *option.Header {
				if values != nil {
					for _, value := range values {
						reqest.Header.Add(key, value)
					}
				}
			}
		}

		if option.Cookies != nil {
			for _, cookie := range option.Cookies {
				if cookie != nil {
					reqest.AddCookie(cookie)
				}
			}
		}
	}
}
