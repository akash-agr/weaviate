// Code generated by go-swagger; DO NOT EDIT.

package keys

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/creativesoftwarefdn/weaviate/models"
)

// WeaviateKeysGetReader is a Reader for the WeaviateKeysGet structure.
type WeaviateKeysGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WeaviateKeysGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewWeaviateKeysGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 401:
		result := NewWeaviateKeysGetUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 403:
		result := NewWeaviateKeysGetForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewWeaviateKeysGetNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 501:
		result := NewWeaviateKeysGetNotImplemented()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewWeaviateKeysGetOK creates a WeaviateKeysGetOK with default headers values
func NewWeaviateKeysGetOK() *WeaviateKeysGetOK {
	return &WeaviateKeysGetOK{}
}

/*WeaviateKeysGetOK handles this case with default header values.

Successful response.
*/
type WeaviateKeysGetOK struct {
	Payload *models.KeyGetResponse
}

func (o *WeaviateKeysGetOK) Error() string {
	return fmt.Sprintf("[GET /keys/{keyId}][%d] weaviateKeysGetOK  %+v", 200, o.Payload)
}

func (o *WeaviateKeysGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.KeyGetResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWeaviateKeysGetUnauthorized creates a WeaviateKeysGetUnauthorized with default headers values
func NewWeaviateKeysGetUnauthorized() *WeaviateKeysGetUnauthorized {
	return &WeaviateKeysGetUnauthorized{}
}

/*WeaviateKeysGetUnauthorized handles this case with default header values.

Unauthorized or invalid credentials.
*/
type WeaviateKeysGetUnauthorized struct {
}

func (o *WeaviateKeysGetUnauthorized) Error() string {
	return fmt.Sprintf("[GET /keys/{keyId}][%d] weaviateKeysGetUnauthorized ", 401)
}

func (o *WeaviateKeysGetUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewWeaviateKeysGetForbidden creates a WeaviateKeysGetForbidden with default headers values
func NewWeaviateKeysGetForbidden() *WeaviateKeysGetForbidden {
	return &WeaviateKeysGetForbidden{}
}

/*WeaviateKeysGetForbidden handles this case with default header values.

The used API-key has insufficient permissions.
*/
type WeaviateKeysGetForbidden struct {
}

func (o *WeaviateKeysGetForbidden) Error() string {
	return fmt.Sprintf("[GET /keys/{keyId}][%d] weaviateKeysGetForbidden ", 403)
}

func (o *WeaviateKeysGetForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewWeaviateKeysGetNotFound creates a WeaviateKeysGetNotFound with default headers values
func NewWeaviateKeysGetNotFound() *WeaviateKeysGetNotFound {
	return &WeaviateKeysGetNotFound{}
}

/*WeaviateKeysGetNotFound handles this case with default header values.

Successful query result but no resource was found.
*/
type WeaviateKeysGetNotFound struct {
}

func (o *WeaviateKeysGetNotFound) Error() string {
	return fmt.Sprintf("[GET /keys/{keyId}][%d] weaviateKeysGetNotFound ", 404)
}

func (o *WeaviateKeysGetNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewWeaviateKeysGetNotImplemented creates a WeaviateKeysGetNotImplemented with default headers values
func NewWeaviateKeysGetNotImplemented() *WeaviateKeysGetNotImplemented {
	return &WeaviateKeysGetNotImplemented{}
}

/*WeaviateKeysGetNotImplemented handles this case with default header values.

Not (yet) implemented.
*/
type WeaviateKeysGetNotImplemented struct {
}

func (o *WeaviateKeysGetNotImplemented) Error() string {
	return fmt.Sprintf("[GET /keys/{keyId}][%d] weaviateKeysGetNotImplemented ", 501)
}

func (o *WeaviateKeysGetNotImplemented) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
