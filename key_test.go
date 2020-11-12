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
		{name: "testdata/SOFTWARE", args: args{filename: "testdata/SOFTWARE", path: "", n: -1}, want: []string{"7-Zip", "Classes", "Clients", "DefaultUserEnvironment", "Docker Inc.", "GitForWindows", "Google", "Intel", "Macromedia", "Microsoft", "MozillaPlugins", "ODBC", "OEM", "OpenSSH", "Oracle", "Partner", "Policies", "RegisteredApplications", "VMware, Inc.", "WOW6432Node"}},
		{name: "testdata/SOFTWARE", args: args{filename: "testdata/SOFTWARE", path: "", n: 1}, want: []string{"7-Zip"}},
		{name: "testdata/SOFTWARE", args: args{filename: "testdata/SOFTWARE", path: "", n: 5}, want: []string{"7-Zip", "Classes", "Clients", "DefaultUserEnvironment", "Docker Inc."}},
		{name: "testdata/SOFTWARE", args: args{filename: "testdata/SOFTWARE", path: "", n: 100}, want: []string{"7-Zip", "Classes", "Clients", "DefaultUserEnvironment", "Docker Inc.", "GitForWindows", "Google", "Intel", "Macromedia", "Microsoft", "MozillaPlugins", "ODBC", "OEM", "OpenSSH", "Oracle", "Partner", "Policies", "RegisteredApplications", "VMware, Inc.", "WOW6432Node"}},
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
		{name: "testdata/SOFTWARE 7-Zip", args: args{filename: "testdata/SOFTWARE", path: "7-Zip", n: -1}, want: []string{"Path", "Path64"}},
		{name: "testdata/SOFTWARE 7-Zip", args: args{filename: "testdata/SOFTWARE", path: "7-Zip", n: 1}, want: []string{"Path64"}},
		{name: "testdata/SOFTWARE 7-Zip", args: args{filename: "testdata/SOFTWARE", path: "7-Zip", n: 5}, want: []string{"Path", "Path64"}},
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
		{name: "testdata/SOFTWARE 7-Zip", args: args{filename: "testdata/SOFTWARE", path: "7-Zip", valuename: "Path"}, wantVal: `C:\Program Files\7-Zip\`, wantValtype: REG_SZ},
		{name: "testdata/SOFTWARE Puppet Labs\\Puppet", args: args{filename: "testdata/SOFTWARE", path: "WOW6432Node\\Puppet Labs\\Puppet", valuename: "RememberedInstallDir64"}, wantVal: `C:\Program Files\Puppet Labs\Puppet\`, wantValtype: REG_SZ},
		{name: "testdata/SOFTWARE WOW6432Node\\Runtime Software\\ShadowCopy", args: args{filename: "testdata/SOFTWARE", path: "WOW6432Node\\Runtime Software\\ShadowCopy", valuename: "Version"}, wantVal: `2.02.000`, wantValtype: REG_SZ},
		{name: `testdata/SOFTWARE WOW6432Node\Puppet Labs\PuppetInstaller RememberedPuppetAgentStartupMode`, args: args{filename: "testdata/SOFTWARE", path: `WOW6432Node\Puppet Labs\PuppetInstaller`, valuename: "RememberedPuppetAgentStartupMode"}, wantVal: "", wantValtype: REG_SZ},
		{name: `testdata/SOFTWARE WOW6432Node\ODBC\ODBCINST.INI\SQL Server Setup`, args: args{filename: "testdata/SOFTWARE", path: `WOW6432Node\ODBC\ODBCINST.INI\SQL Server`, valuename: "Setup"}, wantVal: `%WINDIR%\system32\sqlsrv32.dll`, wantValtype: REG_EXPAND_SZ},
		{name: `testdata/SOFTWARE Classes\*\shell\UpdateEncryptionSettingsWork AttributeMask`, args: args{filename: "testdata/SOFTWARE", path: `Classes\*\shell\UpdateEncryptionSettingsWork`, valuename: "AttributeMask"}, wantValtype: REG_DWORD_LITTLE_ENDIAN, wantErr: true, err: ErrUnexpectedType},
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
		{name: `testdata/SOFTWARE Classes\*\shell\UpdateEncryptionSettingsWork AttributeMask`, args: args{filename: "testdata/SOFTWARE", path: `Classes\*\shell\UpdateEncryptionSettingsWork`, valuename: "AttributeMask"}, wantVal: 8192, wantValtype: REG_DWORD_LITTLE_ENDIAN},
		{name: `testdata/SOFTWARE Classes\*\shell\UpdateEncryptionSettingsWork ImpliedSelectionModel`, args: args{filename: "testdata/SOFTWARE", path: `Classes\*\shell\UpdateEncryptionSettingsWork`, valuename: "ImpliedSelectionModel"}, wantVal: 0, wantValtype: REG_DWORD_LITTLE_ENDIAN},
		{name: `testdata/SOFTWARE Microsoft\.NETFramework\v2.0.50727\NGenService\StateP`, args: args{filename: "testdata/SOFTWARE", path: `Microsoft\.NETFramework\v2.0.50727\NGenService\State`, valuename: "LastSuccess"}, wantVal: 637207470905191361, wantValtype: REG_QWORD},

		{name: `testdata/SOFTWARE Classes\.3gp\OpenWithProgIds WMP11.AssocFile.3GP`, args: args{filename: "testdata/SOFTWARE", path: `Classes\.3gp\OpenWithProgIds`, valuename: "WMP11.AssocFile.3GP"}, wantErr: true, wantValtype: REG_NONE, err: ErrUnexpectedType},
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
		{name: `testdata/SOFTWARE WOW6432Node\Microsoft\WSDAPI\Reporting LastUploadTime`, args: args{filename: "testdata/SOFTWARE", path: `WOW6432Node\Microsoft\WSDAPI\Reporting`, valuename: "LastUploadTime"}, wantVal: make([]byte, 28), wantValtype: REG_BINARY, wantN: 28},
		{name: `testdata/SOFTWARE Classes\*\shell\UpdateEncryptionSettingsWork AttributeMask`, args: args{filename: "testdata/SOFTWARE", path: `Classes\*\shell\UpdateEncryptionSettingsWork`, valuename: "AttributeMask"}, wantVal: []byte{0, 32, 0, 0}, wantValtype: REG_DWORD_LITTLE_ENDIAN, wantN: dataSizeFromType(REG_DWORD_LITTLE_ENDIAN)},
		{name: `testdata/SOFTWARE Microsoft\Cryptography\OID\EncodingType 0\CryptsvcDllCtrl\DEFAULT`, args: args{filename: "testdata/SOFTWARE", path: `Microsoft\Cryptography\OID\EncodingType 0\CryptsvcDllCtrl\DEFAULT`, valuename: "Dll"}, wantVal: []byte{'C', ':', '\\', 'W', 'i', 'n', 'd', 'o', 'w', 's', '\\', 'S', 'y', 's', 't', 'e', 'm', '3', '2', '\\', 'c', 'r', 'y', 'p', 't', 't', 'p', 'm', 'e', 'k', 's', 'v', 'c', '.', 'd', 'l', 'l', 0, 'C', ':', '\\', 'W', 'i', 'n', 'd', 'o', 'w', 's', '\\', 'S', 'y', 's', 't', 'e', 'm', '3', '2', '\\', 'c', 'r', 'y', 'p', 't', 'c', 'a', 't', 's', 'v', 'c', '.', 'd', 'l', 'l', 0, 'C', ':', '\\', 'W', 'i', 'n', 'd', 'o', 'w', 's', '\\', 'S', 'y', 's', 't', 'e', 'm', '3', '2', '\\', 'w', 'e', 'b', 'a', 'u', 't', 'h', 'n', '.', 'd', 'l', 'l', 0}, wantValtype: REG_MULTI_SZ, wantN: 107},
		{name: `testdata/SOFTWARE Microsoft\.NETFramework\v2.0.50727\NGenService\StateP`, args: args{filename: "testdata/SOFTWARE", path: `Microsoft\.NETFramework\v2.0.50727\NGenService\State`, valuename: "LastSuccess"}, wantVal: []byte{193, 107, 135, 151, 209, 208, 215, 8}, wantValtype: REG_QWORD, wantN: dataSizeFromType(REG_QWORD)},
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
			name:        `testdata/SOFTWARE WOW6432Node\Microsoft\WSDAPI\Reporting LastUploadTime`,
			args:        args{filename: "testdata/SOFTWARE", path: `WOW6432Node\Microsoft\WSDAPI\Reporting`, valuename: "LastUploadTime"},
			wantValtype: REG_BINARY,
			wantN:       28,
		},
		{
			name:        `testdata/SOFTWARE Classes\*\shell\UpdateEncryptionSettingsWork AttributeMask`,
			args:        args{filename: "testdata/SOFTWARE", path: `Classes\*\shell\UpdateEncryptionSettingsWork`, valuename: "AttributeMask", buf: make([]byte, dataSizeFromType(REG_DWORD_LITTLE_ENDIAN))},
			wantVal:     []byte{0, 32, 0, 0},
			wantValtype: REG_DWORD_LITTLE_ENDIAN,
			wantN:       dataSizeFromType(REG_DWORD_LITTLE_ENDIAN),
		},
		{
			name:        `testdata/SOFTWARE Microsoft\Cryptography\OID\EncodingType 0\CryptsvcDllCtrl\DEFAULT no buf`,
			args:        args{filename: "testdata/SOFTWARE", path: `Microsoft\Cryptography\OID\EncodingType 0\CryptsvcDllCtrl\DEFAULT`, valuename: "Dll"},
			wantValtype: REG_MULTI_SZ,
			wantN:       107,
		},
		{
			name:        `testdata/SOFTWARE Microsoft\Cryptography\OID\EncodingType 0\CryptsvcDllCtrl\DEFAULT buf`,
			args:        args{filename: "testdata/SOFTWARE", path: `Microsoft\Cryptography\OID\EncodingType 0\CryptsvcDllCtrl\DEFAULT`, valuename: "Dll", buf: make([]byte, 107)},
			wantVal:     []byte{'C', ':', '\\', 'W', 'i', 'n', 'd', 'o', 'w', 's', '\\', 'S', 'y', 's', 't', 'e', 'm', '3', '2', '\\', 'c', 'r', 'y', 'p', 't', 't', 'p', 'm', 'e', 'k', 's', 'v', 'c', '.', 'd', 'l', 'l', 0, 'C', ':', '\\', 'W', 'i', 'n', 'd', 'o', 'w', 's', '\\', 'S', 'y', 's', 't', 'e', 'm', '3', '2', '\\', 'c', 'r', 'y', 'p', 't', 'c', 'a', 't', 's', 'v', 'c', '.', 'd', 'l', 'l', 0, 'C', ':', '\\', 'W', 'i', 'n', 'd', 'o', 'w', 's', '\\', 'S', 'y', 's', 't', 'e', 'm', '3', '2', '\\', 'w', 'e', 'b', 'a', 'u', 't', 'h', 'n', '.', 'd', 'l', 'l', 0},
			wantValtype: REG_MULTI_SZ,
			wantN:       107,
		},
		{
			name:        `testdata/SOFTWARE Microsoft\.NETFramework\v2.0.50727\NGenService\StateP`,
			args:        args{filename: "testdata/SOFTWARE", path: `Microsoft\.NETFramework\v2.0.50727\NGenService\State`, valuename: "LastSuccess", buf: make([]byte, dataSizeFromType(REG_QWORD))},
			wantVal:     []byte{193, 107, 135, 151, 209, 208, 215, 8},
			wantValtype: REG_QWORD,
			wantN:       dataSizeFromType(REG_QWORD),
		},
		{
			name:        `testdata/SOFTWARE Microsoft\.NETFramework\v2.0.50727\NGenService\StateP`,
			args:        args{filename: "testdata/SOFTWARE", path: `Microsoft\.NETFramework\v2.0.50727\NGenService\State`, valuename: "LastSuccess", buf: make([]byte, 2)},
			wantVal:     []byte{193, 107, 135, 151, 209, 208, 215, 8},
			wantValtype: REG_QWORD,
			wantN:       dataSizeFromType(REG_QWORD),
			wantErr:     true,
			err:         ErrShortBuffer,
		},
		{
			name:        `testdata/SOFTWARE Classes\WOW6432Node\CLSID\{0BFCF7B7-E7B6-433a-B205-2904FCF040DD}\InProcServer32`,
			args:        args{filename: "testdata/SOFTWARE", path: `Classes\WOW6432Node\CLSID\{0BFCF7B7-E7B6-433a-B205-2904FCF040DD}\InProcServer32`, valuename: "(default)", buf: make([]byte, 32)},
			wantVal:     []byte{'%', 'S', 'y', 's', 't', 'e', 'm', 'R', 'o', 'o', 't', '%', '\\', 'S', 'y', 's', 't', 'e', 'm', '3', '2', '\\', 'a', 'p', 'p', 'w', 'i', 'z', '.', 'c', 'p', 'l'},
			wantValtype: REG_EXPAND_SZ,
			wantN:       32,
		},
		{
			name:        `testdata/SOFTWARE Classes\AudioCD EditFlags`,
			args:        args{filename: "testdata/SOFTWARE", path: `Classes\AudioCD`, valuename: "EditFlags", buf: nil},
			wantVal:     nil,
			wantN:       4,
			wantValtype: REG_BINARY,
		},
		{
			name:        `testdata/SOFTWARE Classes\AudioCD EditFlags`,
			args:        args{filename: "testdata/SOFTWARE", path: `Classes\AudioCD`, valuename: "EditFlags", buf: make([]byte, 4)},
			wantVal:     []byte{2, 0, 16, 0},
			wantN:       4,
			wantValtype: REG_BINARY,
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
		{name: `testdata/SOFTWARE Microsoft\Cryptography\OID\EncodingType 0\CryptsvcDllCtrl\DEFAULT`, args: args{filename: "testdata/SOFTWARE", path: `Microsoft\Cryptography\OID\EncodingType 0\CryptsvcDllCtrl\DEFAULT`, valuename: "Dll"}, wantVal: []string{`C:\Windows\System32\crypttpmeksvc.dll`, `C:\Windows\System32\cryptcatsvc.dll`, `C:\Windows\System32\webauthn.dll`}, wantValtype: REG_MULTI_SZ},
		{name: "testdata/SOFTWARE 7-Zip", args: args{filename: "testdata/SOFTWARE", path: "7-Zip", valuename: "Path"}, wantValtype: REG_SZ, wantErr: true, err: ErrUnexpectedType},
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
		{name: `testdata/SOFTWARE WOW6432Node\Microsoft\WSDAPI\Reporting LastUploadTime`, args: args{filename: "testdata/SOFTWARE", path: `WOW6432Node\Microsoft\WSDAPI\Reporting`, valuename: "LastUploadTime"}, wantVal: make([]byte, 28), wantValtype: REG_BINARY},
		{name: "testdata/SOFTWARE 7-Zip", args: args{filename: "testdata/SOFTWARE", path: "7-Zip", valuename: "Path"}, wantValtype: REG_SZ, wantErr: true, err: ErrUnexpectedType},
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
