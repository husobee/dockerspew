package content

import (
	"encoding/json"
	"encoding/xml"
	"gopkg.in/unrolled/render.v1"
	"io"
	"log"
)

// DecoderInterface - an interface that encoding/* can implement
type DecoderInterface interface {
	Decode(interface{}) error
}

// GetContentTypeDecoder - given a contentType and a request body io.Reader, return a decoder for what is
// specified on the content-type header
func GetContentTypeDecoder(contentType interface{}, reader io.Reader) DecoderInterface {
	log.Print("[DEBUG] GetContentTypeDecoder - contentType is ", contentType)
	switch contentType {
	case render.ContentJSON:
		log.Print("[DEBUG] GetContentTypeDecoder - returning a json decoder")
		return json.NewDecoder(reader)
	case render.ContentXML:
		log.Print("[DEBUG] GetContentTypeDecoder - returning an xml decoder")
		return xml.NewDecoder(reader)
	default:
		log.Print("[DEBUG] GetContentTypeDecoder - returning a json decoder")
		return json.NewDecoder(reader)
	}
}
