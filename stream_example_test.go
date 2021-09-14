package molecule

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func ExampleNew() {
	/* Encoding the following:
	 *
	 * message SearchRequest {
	 *   string query = 1;
	 *   int32 page_number = 2;
	 *   int32 result_per_page = 3;
	 * }
	 */
	var err error
	output := bytes.NewBuffer([]byte{})
	ps := NewProtoStream(output)

	// values copied from the .proto file
	const fieldQuery int = 1
	const fieldPageNumber int = 2
	const fieldResultPerPage int = 3

	err = ps.String(fieldQuery, "q=streaming+protobufs")
	if err != nil {
		panic(err)
	}

	err = ps.Int32(fieldPageNumber, 2)
	if err != nil {
		panic(err)
	}

	err = ps.Int32(fieldResultPerPage, 100)
	if err != nil {
		panic(err)
	}

	// The encoded result is in `output.Bytes()`.
}

func ExampleProtoStream_Embedded() {
	/* Encoding the following:
	     *
	     * message MultiSearch {
	     *   string api_key = 10;
	     *   repeated SearchRequest request = 11;
	     * }
		 *
		 * message SearchRequest {
		 *   string query = 1;
		 *   int32 page_number = 2;
		 *   int32 result_per_page = 3;
		 * }
	*/
	var err error
	output := bytes.NewBuffer([]byte{})
	ps := NewProtoStream(output)

	// values copied from the .proto file
	const fieldAPIKey = 10
	const fieldRequest = 11
	const fieldQuery int = 1
	const fieldPageNumber int = 2
	const fieldResultPerPage int = 3

	err = ps.String(fieldAPIKey, "abc-123")
	if err != nil {
		panic(err)
	}

	err = ps.Embedded(fieldRequest, func(ps *ProtoStream) error {
		err = ps.String(fieldQuery, "author=octavia+butler")
		if err != nil {
			return err
		}

		err = ps.Int32(fieldPageNumber, 2)
		if err != nil {
			return err
		}

		err = ps.Int32(fieldResultPerPage, 100)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	err = ps.Embedded(fieldRequest, func(ps *ProtoStream) error {
		err = ps.String(fieldQuery, "author=margaret+atwood")
		if err != nil {
			return err
		}

		err = ps.Int32(fieldPageNumber, 0)
		if err != nil {
			return err
		}

		err = ps.Int32(fieldResultPerPage, 10)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	// The encoded result is in `output.Bytes()`.
}

func ExampleProtoStream_Sint32Packed() {
	/* Encoding the following:
	 *
	 * message Numbers {
	 *   repeated int32 number = 22;
	 * }
	 */
	var err error
	output := bytes.NewBuffer([]byte{})
	ps := NewProtoStream(output)

	const fieldNumber = 22

	numbers := []int32{20, -30, -31, 1999}

	err = ps.Sint32Packed(fieldNumber, numbers)
	if err != nil {
		panic(err)
	}

	res := bytes.NewReader(output.Bytes())
	key, _ := binary.ReadUvarint(res)
	fmt.Printf("key: 0x%x = 22<<3 + 2\n", key)
	leng, _ := binary.ReadUvarint(res)
	fmt.Printf("length: 0x%x\n", leng)
	v, _ := binary.ReadUvarint(res)
	fmt.Printf("v[0]: 0x%x\n", v)
	v, _ = binary.ReadUvarint(res)
	fmt.Printf("v[1]: 0x%x\n", v)
	v, _ = binary.ReadUvarint(res)
	fmt.Printf("v[2]: 0x%x\n", v)
	v, _ = binary.ReadUvarint(res)
	fmt.Printf("v[3]: 0x%x\n", v)
	// Output:
	// key: 0xb2 = 22<<3 + 2
	// length: 0x5
	// v[0]: 0x28
	// v[1]: 0x3b
	// v[2]: 0x3d
	// v[3]: 0xf9e
}
