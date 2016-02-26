package postal

/*
#cgo pkg-config: libpostal
#include <libpostal/libpostal.h>
#include <stdlib.h>
*/
import "C"
import "log"
import "unsafe"

func init() {
    if (!bool(C.libpostal_setup()) || !bool(C.libpostal_setup_parser())) {
        log.Fatal("Could not load libpostal")
    }
}

type ParserOptions struct {
    Language string
    Country string
}


func getDefaultParserOptions() ParserOptions {
    return ParserOptions {
        Language: "",
        Country: "",
    }
}

var parserDefaultOptions = getDefaultParserOptions()

type ParsedComponent struct {
    Value string
    Label string
}

func ParseAddressOptions(address string, options ParserOptions) []ParsedComponent {
    cAddress := C.CString(address)
    defer C.free(unsafe.Pointer(cAddress))

    //var b *C.char
    //ptr_size := unsafe.Sizeof(b)

    cOptions := C.get_libpostal_address_parser_default_options()
    if options.Language != "" {
        cLanguage := C.CString(options.Language)
        defer C.free(unsafe.Pointer(cLanguage))

        cOptions.language = cLanguage
    }

    if options.Country != "" {
        cCountry := C.CString(options.Country)
        defer C.free(unsafe.Pointer(cCountry))

        cOptions.country = cCountry
    }

    cAddressParserResponsePtr := C.parse_address(cAddress, cOptions)

    if cAddressParserResponsePtr == nil {
        return nil
    }

    cAddressParserResponse := *cAddressParserResponsePtr

    cNumComponents := cAddressParserResponse.num_components
    cComponents := cAddressParserResponse.components
    cLabels := cAddressParserResponse.labels

    numComponents := uint64(cNumComponents)

    parsedComponents := make([]ParsedComponent, numComponents)

    // Accessing a C array
    cComponentsPtr := (*[1<<30](*C.char))(unsafe.Pointer(cComponents))
    cLabelsPtr := (*[1<<30](*C.char))(unsafe.Pointer(cLabels))

    var i uint64
    for i = 0; i < numComponents; i++ {
        parsedComponents[i] = ParsedComponent{
            Value: C.GoString(cComponentsPtr[i]),
            Label: C.GoString(cLabelsPtr[i]),
        }
    }

    C.address_parser_response_destroy(cAddressParserResponsePtr)

    return parsedComponents
}

func ParseAddress(address string) []ParsedComponent {
    parsedComponents := ParseAddressOptions(address, parserDefaultOptions);
    return parsedComponents
}
