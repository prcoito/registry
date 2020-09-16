package registry

import (
	"os"
	"testing"
)

func Test_valueKey_Read(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    valueKey
		wantErr bool
	}{
		{
			name: "VK big int", file: "testdata/unit/vk_big_int",
			want:    valueKey{dataSize: 2, data: uint32(7), nameSize: 7, name: "BIG_INT", dataOffset: 117440512, flags: 0, signature: "vk", valueOffset: 0, binOffset: 0, dataType: REG_DWORD_BIG_ENDIAN},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp, err := os.Open(tt.file)
			if err != nil {
				t.Fatal(err)
			}
			defer fp.Close()

			vk := newValueKey(fp, 0, 0)
			if err := vk.Read(); (err != nil) != tt.wantErr {
				t.Errorf("valueKey.Read() error = %v, wantErr %v", err, tt.wantErr)
			}

			vk.rws = nil
			if *vk != tt.want {
				t.Errorf("Read error:\nvk      = %+v;\ntt.want = %+v", *vk, tt.want)
			}
		})
	}
}
