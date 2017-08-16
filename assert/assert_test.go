package assert

import (
	"errors"
	"testing"
)

func TestTrue(t *testing.T) {
	mockT := new(testing.T)

	succes := True(mockT, true)
	if !succes {
		t.Fail()
	}
	succes = True(mockT, false)
	if succes {
		t.Fail()
	}
}

func TestFalse(t *testing.T) {
	mockT := new(testing.T)

	succes := False(mockT, false)
	if !succes {
		t.Fail()
	}
	succes = False(mockT, true)
	if succes {
		t.Fail()
	}
}

func TestError(t *testing.T) {
	mockT := new(testing.T)
	mockError := errors.New("some mock error")
	var mockNilError error

	succes := Error(mockT, mockError)
	if !succes {
		t.Fail()
	}
	succes = Error(mockT, mockNilError)
	if succes {
		t.Fail()
	}
	succes = Error(mockT, nil)
	if succes {
		t.Fail()
	}
}

func TestNoError(t *testing.T) {
	mockT := new(testing.T)
	mockError := errors.New("some mock error")
	var mockNilError error

	succes := NoError(mockT, mockNilError)
	if !succes {
		t.Fail()
	}
	succes = NoError(mockT, nil)
	if !succes {
		t.Fail()
	}
	succes = NoError(mockT, mockError)
	if succes {
		t.Fail()
	}
}

func TestEqual(t *testing.T) {
	mockT := new(testing.T)

	// nils
	succes := Equal(mockT, nil, nil)
	if !succes {
		t.Fail()
	}
	succes = Equal(mockT, nil, 0)
	if succes {
		t.Fail()
	}

	// strings
	str1 := "string1"
	str2 := "string2"

	succes = Equal(mockT, str1, str1)
	if !succes {
		t.Fail()
	}
	succes = Equal(mockT, str1, str2)
	if succes {
		t.Fail()
	}

	// ints
	succes = Equal(mockT, 1, 1)
	if !succes {
		t.Fail()
	}
	succes = Equal(mockT, 1, 2)
	if succes {
		t.Fail()
	}

	// byte slice
	bs1 := []byte{
		1, 2, 3, 4,
	}
	bs2 := []byte{
		5, 6, 7, 8,
	}

	succes = Equal(mockT, bs1, bs1)
	if !succes {
		t.Fail()
	}
	succes = Equal(mockT, bs1, bs2)
	if succes {
		t.Fail()
	}
	succes = Equal(mockT, bs1, str1)
	if succes {
		t.Fail()
	}

	var bs3 []byte
	succes = Equal(mockT, bs3, bs3)
	if !succes {
		t.Fail()
	}
}

func TestNotEqual(t *testing.T) {
	mockT := new(testing.T)

	// nils
	succes := NotEqual(mockT, nil, 0)
	if !succes {
		t.Fail()
	}
	succes = NotEqual(mockT, nil, nil)
	if succes {
		t.Fail()
	}

	// strings
	str1 := "string1"
	str2 := "string2"

	succes = NotEqual(mockT, str1, str2)
	if !succes {
		t.Fail()
	}
	succes = NotEqual(mockT, str1, str1)
	if succes {
		t.Fail()
	}

	// ints
	succes = NotEqual(mockT, 1, 2)
	if !succes {
		t.Fail()
	}
	succes = NotEqual(mockT, 1, 1)
	if succes {
		t.Fail()
	}

	// byte slice
	bs1 := []byte{
		1, 2, 3, 4,
	}
	bs2 := []byte{
		5, 6, 7, 8,
	}

	succes = NotEqual(mockT, bs1, bs2)
	if !succes {
		t.Fail()
	}
	succes = NotEqual(mockT, bs1, bs1)
	if succes {
		t.Fail()
	}
}
