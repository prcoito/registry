package registry

import (
	"encoding/hex"
	"reflect"
	"strings"
	"testing"
)

func TestKey_ReadSubKeyNames(t *testing.T) {
	type args struct {
		filename string
		path     string
		n        int
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "testdata/NTUSER.DAT", args: args{filename: "testdata/NTUSER.DAT", path: "SOFTWARE", n: -1}, want: []string{"Google", "Microsoft", "Policies"}},
		{name: "testdata/NTUSER.DAT", args: args{filename: "testdata/NTUSER.DAT", path: "SOFTWARE", n: 1}, want: []string{"Google"}},
		{name: "testdata/NTUSER.DAT", args: args{filename: "testdata/NTUSER.DAT", path: "SOFTWARE", n: 100}, want: []string{"Google", "Microsoft", "Policies"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := OpenKey(tt.args.filename, tt.args.path)
			if err != nil {
				t.Errorf("OpenKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer k.Close()

			got, err := k.ReadSubKeyNames(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.ReadSubKeyNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Key.ReadSubKeyNames()\n[%v]\nwant\n[%v]", strings.Join(got, ", "), strings.Join(tt.want, ", "))
			}

		})
	}
}

func TestKey_ReadValueNames(t *testing.T) {
	type args struct {
		filename string
		path     string
		n        int
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: `testdata/NTUSER.DAT SOFTWARE\Microsoft\CTF\Assemblies\0x00000816\{34745C63-B2F0-4784-8B67-5E12C8701A31}`, args: args{filename: "testdata/NTUSER.DAT", path: `SOFTWARE\Microsoft\CTF\Assemblies\0x00000816\{34745C63-B2F0-4784-8B67-5E12C8701A31}`, n: -1}, want: []string{"Default", "KeyboardLayout", "Profile"}},
		{name: `testdata/NTUSER.DAT SOFTWARE\Microsoft\CTF\Assemblies\0x00000816\{34745C63-B2F0-4784-8B67-5E12C8701A31}`, args: args{filename: "testdata/NTUSER.DAT", path: `SOFTWARE\Microsoft\CTF\Assemblies\0x00000816\{34745C63-B2F0-4784-8B67-5E12C8701A31}`, n: 1}, want: []string{"Default"}},
		{name: `testdata/NTUSER.DAT SOFTWARE\Microsoft\CTF\Assemblies\0x00000816\{34745C63-B2F0-4784-8B67-5E12C8701A31}`, args: args{filename: "testdata/NTUSER.DAT", path: `SOFTWARE\Microsoft\CTF\Assemblies\0x00000816\{34745C63-B2F0-4784-8B67-5E12C8701A31}`, n: 5}, want: []string{"Default", "KeyboardLayout", "Profile"}},
		{name: `testdata/NTUSER.DAT`, args: args{filename: "testdata/NTUSER.DAT", path: ``, n: 5}, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := OpenKey(tt.args.filename, tt.args.path)
			if err != nil {
				t.Errorf("OpenKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer k.Close()

			got, err := k.ReadValueNames(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.ReadValueNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Key.ReadValueNames() = [%v], want [%v]", strings.Join(got, ", "), strings.Join(tt.want, ", "))
			}
		})
	}
}

func TestKey_GetStringValue(t *testing.T) {
	type args struct {
		filename  string
		path      string
		valuename string
	}
	tests := []struct {
		name        string
		args        args
		wantVal     string
		wantValtype uint32
		wantErr     bool
		err         error
	}{
		{name: `testdata/NTUSER.DAT SOFTWARE\Google\Chrome\NativeMessagingHosts\com.microsoft.browsercore`, args: args{filename: "testdata/NTUSER.DAT", path: `SOFTWARE\Google\Chrome\NativeMessagingHosts\com.microsoft.browsercore`, valuename: "(default)"}, wantVal: `C:\Program Files\Windows Security\BrowserCore\manifest.json`, wantValtype: REG_SZ},
		{name: `testdata/NTUSER.DAT Environment`, args: args{filename: "testdata/NTUSER.DAT", path: `Environment`, valuename: "Path"}, wantVal: `%USERPROFILE%\AppData\Local\Microsoft\WindowsApps;`, wantValtype: REG_EXPAND_SZ},
		{name: `testdata/NTUSER.DAT Control Panel\PowerCfg\PowerPolicies\5 Policies`, args: args{filename: "testdata/NTUSER.DAT", path: `Control Panel\PowerCfg\PowerPolicies\5`, valuename: "Policies"}, wantValtype: REG_BINARY, wantErr: true, err: ErrUnexpectedType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := OpenKey(tt.args.filename, tt.args.path)
			if err != nil {
				t.Errorf("OpenKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer k.Close()

			gotVal, gotValtype, err := k.GetStringValue(tt.args.valuename)
			if (err != nil) != tt.wantErr || err != tt.err {
				t.Errorf("Key.GetStringValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Key.GetStringValue() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotValtype != tt.wantValtype {
				t.Errorf("Key.GetStringValue() gotValtype = `%v`, want `%v`", gotValtype, tt.wantValtype)
			}
		})
	}
}

func TestKey_GetIntegerValue(t *testing.T) {
	type args struct {
		filename  string
		path      string
		valuename string
	}
	tests := []struct {
		name        string
		args        args
		wantVal     uint64
		wantValtype uint32
		wantErr     bool
		err         error
	}{
		{name: `testdata/NTUSER.DAT SOFTWARE\Microsoft\InputPersonalization`, args: args{filename: "testdata/NTUSER.DAT", path: `SOFTWARE\Microsoft\InputPersonalization`, valuename: "RestrictImplicitInkCollection"}, wantVal: 0, wantValtype: REG_DWORD_LITTLE_ENDIAN},
		{name: `testdata/NTUSER.DAT Classes\*\shell\UpdateEncryptionSettingsWork ImpliedSelectionModel`, args: args{filename: "testdata/NTUSER.DAT", path: `SOFTWARE\Microsoft\InputPersonalization\TrainedDataStore`, valuename: "HarvestContacts"}, wantVal: 1, wantValtype: REG_DWORD_LITTLE_ENDIAN},

		{name: `testdata/NTUSER.DAT Control Panel\PowerCfg\PowerPolicies\5 Policies`, args: args{filename: "testdata/NTUSER.DAT", path: `Control Panel\PowerCfg\PowerPolicies\5`, valuename: "Policies"}, wantValtype: REG_BINARY, wantErr: true, err: ErrUnexpectedType},
		// TODO: REG_DWORD_BIG_ENDIAN support
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := OpenKey(tt.args.filename, tt.args.path)
			if err != nil {
				t.Errorf("OpenKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer k.Close()

			gotVal, gotValtype, err := k.GetIntegerValue(tt.args.valuename)
			if (err != nil) != tt.wantErr || err != tt.err {
				t.Errorf("Key.GetIntegerValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Key.GetIntegerValue() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotValtype != tt.wantValtype {
				t.Errorf("Key.GetIntegerValue() gotValtype = %v, want %v", gotValtype, tt.wantValtype)
			}
		})
	}
}

func TestKey_GetValue(t *testing.T) {
	type args struct {
		filename  string
		path      string
		valuename string
	}
	tests := []struct {
		name        string
		args        args
		wantN       int
		wantVal     []byte
		wantValtype uint32
		wantErr     bool
	}{
		{name: `testdata/NTUSER.DAT Control Panel\Input Method\Hot Keys\00000010`, args: args{filename: "testdata/NTUSER.DAT", path: `Control Panel\Input Method\Hot Keys\00000010`, valuename: "Key Modifiers"}, wantVal: []byte{'\x02', '\xc0', '\x00', '\x00'}, wantValtype: REG_BINARY, wantN: 4},
		{name: `testdata/NTUSER.DAT SOFTWARE\Google\Chrome\NativeMessagingHosts\com.microsoft.browsercore`, args: args{filename: "testdata/NTUSER.DAT", path: `SOFTWARE\Google\Chrome\NativeMessagingHosts\com.microsoft.browsercore`, valuename: "(default)"}, wantVal: []byte(`C:\Program Files\Windows Security\BrowserCore\manifest.json`), wantValtype: REG_SZ, wantN: 59},
		{name: `testdata/NTUSER.DAT SOFTWARE\Microsoft\InputPersonalization`, args: args{filename: "testdata/NTUSER.DAT", path: `SOFTWARE\Microsoft\InputPersonalization`, valuename: "RestrictImplicitInkCollection"}, wantVal: []byte{0, 0, 0, 0}, wantValtype: REG_DWORD_LITTLE_ENDIAN, wantN: dataSizeFromType(REG_DWORD_LITTLE_ENDIAN)},
		{name: `testdata/NTUSER.DAT Control Panel\International\User Profile`, args: args{filename: "testdata/NTUSER.DAT", path: `Control Panel\International\User Profile`, valuename: "Languages"}, wantVal: append([]byte(`pt-PT`), 0), wantValtype: REG_MULTI_SZ, wantN: 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := OpenKey(tt.args.filename, tt.args.path)
			if err != nil {
				t.Errorf("OpenKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer k.Close()

			gotN, gotValtype, err := k.GetValue(tt.args.valuename, nil)
			buf := make([]byte, gotN)

			gotN, gotValtype, err = k.GetValue(tt.args.valuename, buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Key.GetValue() gotN = %v, want %v", gotN, tt.wantN)
			}
			if gotValtype != tt.wantValtype {
				t.Errorf("Key.GetValue() gotValtype = %v, want %v", gotValtype, tt.wantValtype)
			}

			if !reflect.DeepEqual(tt.wantVal, buf) {
				t.Errorf("Key.GetValue()\nval  = %v\nwant = %v\n%v%v", buf, tt.wantVal, hex.Dump(buf), hex.Dump(tt.wantVal))
			}
		})
	}
}

func TestKey_GetValue2(t *testing.T) {
	type args struct {
		filename  string // registry path
		path      string // Key path
		valuename string // Key value name
		buf       []byte
	}
	tests := []struct {
		name        string
		args        args
		wantN       int
		wantVal     []byte
		wantValtype uint32
		wantErr     bool
		err         error
	}{
		{
			name:        `testdata/NTUSER.DAT Control Panel\Input Method\Hot Keys\00000010`,
			args:        args{filename: "testdata/NTUSER.DAT", path: `Control Panel\Input Method\Hot Keys\00000010`, valuename: "Key Modifiers"},
			wantValtype: REG_BINARY,
			wantN:       4,
		},
		{
			name:        `testdata/NTUSER.DAT Control Panel\Input Method\Hot Keys\00000010`,
			args:        args{filename: "testdata/NTUSER.DAT", path: `Control Panel\Input Method\Hot Keys\00000010`, valuename: "Key Modifiers", buf: make([]byte, 4)},
			wantVal:     []byte{'\x02', '\xc0', '\x00', '\x00'},
			wantValtype: REG_BINARY,
			wantN:       4,
		},
		{
			name:        `testdata/NTUSER.DAT SOFTWARE\Google\Chrome\NativeMessagingHosts\com.microsoft.browsercore`,
			args:        args{filename: "testdata/NTUSER.DAT", path: `SOFTWARE\Google\Chrome\NativeMessagingHosts\com.microsoft.browsercore`, valuename: "(default)"},
			wantValtype: REG_SZ,
			wantN:       59,
		},
		{
			name:        `testdata/NTUSER.DAT excess memory`,
			args:        args{filename: "testdata/NTUSER.DAT", path: `SOFTWARE\Google\Chrome\NativeMessagingHosts\com.microsoft.browsercore`, valuename: "(default)", buf: make([]byte, 62)},
			wantVal:     append([]byte(`C:\Program Files\Windows Security\BrowserCore\manifest.json`), 0, 0, 0),
			wantValtype: REG_SZ,
			wantN:       59,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := OpenKey(tt.args.filename, tt.args.path)
			if err != nil {
				t.Errorf("OpenKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer k.Close()

			gotN, gotValtype, err := k.GetValue(tt.args.valuename, tt.args.buf)
			if (err != nil) != tt.wantErr || err != tt.err {
				t.Errorf("Key.GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Key.GetValue() gotN = %v, want %v", gotN, tt.wantN)
			}
			if gotValtype != tt.wantValtype {
				t.Errorf("Key.GetValue() gotValtype = %v, want %v", gotValtype, tt.wantValtype)
			}

			if err == nil && !reflect.DeepEqual(tt.wantVal, tt.args.buf) {
				t.Errorf("Key.GetValue()\nval  = %v\nwant = %v\n%v%v", tt.args.buf, tt.wantVal, hex.Dump(tt.args.buf), hex.Dump(tt.wantVal))
			}
		})
	}
}

func TestKey_GetStringsValue(t *testing.T) {
	type args struct {
		filename  string
		path      string
		valuename string
	}
	tests := []struct {
		name        string
		args        args
		wantVal     []string
		wantValtype uint32
		wantErr     bool
		err         error
	}{
		{name: `testdata/NTUSER.DAT Control Panel\International\User Profile`, args: args{filename: "testdata/NTUSER.DAT", path: `Control Panel\International\User Profile`, valuename: "Languages"}, wantVal: []string{`pt-PT`}, wantValtype: REG_MULTI_SZ},
		{name: `testdata/NTUSER.DAT Control Panel\PowerCfg\PowerPolicies\5 Policies`, args: args{filename: "testdata/NTUSER.DAT", path: `Control Panel\PowerCfg\PowerPolicies\5`, valuename: "Policies"}, wantValtype: REG_BINARY, wantErr: true, err: ErrUnexpectedType}, {name: `testdata/NTUSER.DAT Control Panel\PowerCfg\PowerPolicies\5 Policies`, args: args{filename: "testdata/NTUSER.DAT", path: `Control Panel\PowerCfg\PowerPolicies\5`, valuename: "Policies"}, wantValtype: REG_BINARY, wantErr: true, err: ErrUnexpectedType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := OpenKey(tt.args.filename, tt.args.path)
			if err != nil {
				t.Errorf("OpenKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer k.Close()

			gotVal, gotValtype, err := k.GetStringsValue(tt.args.valuename)
			if (err != nil) != tt.wantErr || err != tt.err {
				t.Errorf("Key.GetStringsValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Key.GetStringsValue() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotValtype != tt.wantValtype {
				t.Errorf("Key.GetStringsValue() gotValtype = %v, want %v", gotValtype, tt.wantValtype)
			}
		})
	}
}

func TestKey_GetBinaryValue(t *testing.T) {
	type args struct {
		filename  string
		path      string
		valuename string
	}
	tests := []struct {
		name        string
		args        args
		wantVal     []byte
		wantValtype uint32
		wantErr     bool
		err         error
	}{
		{name: `testdata/NTUSER.DAT Control Panel\Input Method\Hot Keys\00000010`, args: args{filename: "testdata/NTUSER.DAT", path: `Control Panel\Input Method\Hot Keys\00000010`, valuename: "Key Modifiers"}, wantVal: []byte{'\x02', '\xc0', '\x00', '\x00'}, wantValtype: REG_BINARY},
		{name: `testdata/NTUSER.DAT Classes\*\shell\UpdateEncryptionSettingsWork ImpliedSelectionModel`, args: args{filename: "testdata/NTUSER.DAT", path: `SOFTWARE\Microsoft\InputPersonalization\TrainedDataStore`, valuename: "HarvestContacts"}, wantValtype: REG_DWORD_LITTLE_ENDIAN, wantErr: true, err: ErrUnexpectedType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := OpenKey(tt.args.filename, tt.args.path)
			if err != nil {
				t.Errorf("OpenKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer k.Close()

			gotVal, gotValtype, err := k.GetBinaryValue(tt.args.valuename)
			if (err != nil) != tt.wantErr || err != tt.err {
				t.Errorf("Key.GetBinaryValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Key.GetBinaryValue() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotValtype != tt.wantValtype {
				t.Errorf("Key.GetBinaryValue() gotValtype = %v, want %v", gotValtype, tt.wantValtype)
			}
		})
	}
}
