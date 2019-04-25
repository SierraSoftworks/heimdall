package file

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTailReader(t *testing.T) {
	Convey("TailReader", t, func(c C) {
		os.Remove("tail_reader_test.dat")
		defer os.Remove("tail_reader_test.dat")

		f, err := os.OpenFile("tail_reader_test.dat", os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0664)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		defer f.Close()

		fr, err := os.OpenFile("tail_reader_test.dat", os.O_RDONLY|os.O_SYNC, 0664)
		So(err, ShouldBeNil)
		So(fr, ShouldNotBeNil)
		defer fr.Close()

		tr := NewTailReader(fr, 10*time.Millisecond)
		So(tr, ShouldNotBeNil)

		written := false
		test_data := []byte("test")
		go func() {
			time.Sleep(30 * time.Millisecond)
			_, err := f.Write(test_data)
			written = true
			c.So(err, ShouldBeNil)
		}()

		read_data := make([]byte, 10)
		n, err := tr.Read(read_data)
		So(err, ShouldBeNil)
		So(n, ShouldEqual, 4)
		So(read_data[:n], ShouldResemble, test_data)
		So(written, ShouldBeTrue)
	})
}
