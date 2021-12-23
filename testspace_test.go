package testspace

import (
	"errors"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShellSpace(t *testing.T) {
	currentPath, _ := os.Getwd()

	type args struct {
		options []CreateOption
	}
	tests := []struct {
		name    string
		args    args
		want    Space
		wantErr bool
	}{
		{
			name: "only_path_parameter",
			args: args{
				options: []CreateOption{
					WithPathOption("tmp"),
				},
			},
			want: &workSpace{
				path: path.Join(currentPath, "tmp"),
				env: []string{
					fmt.Sprintf("HOME=%s/tmp", currentPath),
					"GIT_AUTHOR_EMAIL=author@example.com",
					"GIT_AUTHOR_NAME='A U Thor'",
					"GIT_COMMITTER_EMAIL=committer@example.com",
					"GIT_COMMITTER_NAME='C O Mitter'",
				},
				template: `
test_tick () {
        if test -z "${test_tick+set}"
        then
                test_tick=1112911993
        else
                test_tick=$(($test_tick + 60))
        fi
        GIT_COMMITTER_DATE="$test_tick -0700"
        GIT_AUTHOR_DATE="$test_tick -0700"
        export GIT_COMMITTER_DATE GIT_AUTHOR_DATE
}
`,
				customShell: "",
				output:      "",
			},
			wantErr: false,
		},
		{
			name: "path_and_env_parameter",
			args: args{
				options: []CreateOption{
					WithPathOption("tmp"),
					WithEnvironmentsOption("goshelltest1=111", "goshelltest2=222"),
				},
			},
			want: &workSpace{
				path: path.Join(currentPath, "tmp"),
				env: []string{
					fmt.Sprintf("HOME=%s/tmp", currentPath),
					"GIT_AUTHOR_EMAIL=author@example.com",
					"GIT_AUTHOR_NAME='A U Thor'",
					"GIT_COMMITTER_EMAIL=committer@example.com",
					"GIT_COMMITTER_NAME='C O Mitter'",
					"goshelltest1=111",
					"goshelltest2=222",
				},
				template: `
test_tick () {
        if test -z "${test_tick+set}"
        then
                test_tick=1112911993
        else
                test_tick=$(($test_tick + 60))
        fi
        GIT_COMMITTER_DATE="$test_tick -0700"
        GIT_AUTHOR_DATE="$test_tick -0700"
        export GIT_COMMITTER_DATE GIT_AUTHOR_DATE
}
`,
				customShell: "",
				output:      "",
			},
			wantErr: false,
		},
		{
			name: "path_and_env_and_template_parameter",
			args: args{
				options: []CreateOption{
					WithPathOption("tmp"),
					WithEnvironmentsOption("goshelltest1=111", "goshelltest2=222"),
					WithTemplateOption(`
test(){
	echo hello
}
`),
				},
			},
			want: &workSpace{
				path: path.Join(currentPath, "tmp"),
				env: []string{
					fmt.Sprintf("HOME=%s/tmp", currentPath),
					"GIT_AUTHOR_EMAIL=author@example.com",
					"GIT_AUTHOR_NAME='A U Thor'",
					"GIT_COMMITTER_EMAIL=committer@example.com",
					"GIT_COMMITTER_NAME='C O Mitter'",
					"goshelltest1=111",
					"goshelltest2=222",
				},
				template: `
test_tick () {
        if test -z "${test_tick+set}"
        then
                test_tick=1112911993
        else
                test_tick=$(($test_tick + 60))
        fi
        GIT_COMMITTER_DATE="$test_tick -0700"
        GIT_AUTHOR_DATE="$test_tick -0700"
        export GIT_COMMITTER_DATE GIT_AUTHOR_DATE
}


test(){
	echo hello
}
`,
				customShell: "",
				output:      "",
			},
			wantErr: false,
		},
		{
			name: "path_and_env_and_template_and_customshell_parameter",
			args: args{
				options: []CreateOption{
					WithPathOption("tmp"),
					WithEnvironmentsOption("goshelltest1=111", "goshelltest2=222"),
					WithTemplateOption(`
test(){
	echo hello
}
`),
					WithShellOption("test"),
				},
			},
			want: &workSpace{
				path: path.Join(currentPath, "tmp"),
				env: []string{
					fmt.Sprintf("HOME=%s/tmp", currentPath),
					"GIT_AUTHOR_EMAIL=author@example.com",
					"GIT_AUTHOR_NAME='A U Thor'",
					"GIT_COMMITTER_EMAIL=committer@example.com",
					"GIT_COMMITTER_NAME='C O Mitter'",
					"goshelltest1=111",
					"goshelltest2=222",
				},
				template: `
test_tick () {
        if test -z "${test_tick+set}"
        then
                test_tick=1112911993
        else
                test_tick=$(($test_tick + 60))
        fi
        GIT_COMMITTER_DATE="$test_tick -0700"
        GIT_AUTHOR_DATE="$test_tick -0700"
        export GIT_COMMITTER_DATE GIT_AUTHOR_DATE
}


test(){
	echo hello
}
`,
				customShell: "test",
				output:      "hello\n",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Create(tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer actual.Cleanup()

			assert.Equal(t, tt.want.GetPath(""), actual.GetPath(""))
			assert.Equal(t, tt.want.GetTemplateStr(), actual.GetTemplateStr())
			envs1 := tt.want.GetEnvStr()
			envs2 := actual.GetEnvStr()
			if assert.True(t,
				len(envs1) <= len(envs2),
				"actual envs less than expect envs") {
				for i := range envs1 {
					assert.Equal(t, envs1[i], envs2[i])
				}

			}
			_, err = os.Stat(actual.GetPath(".git"))
			if err != nil {
				t.Error("git init failed")
			}
		})
	}
}

func Test_workSpace_RegistrationCustomCleaner(t *testing.T) {
	type fields struct {
		path        string
		env         []string
		template    string
		customShell string
		output      string
		outErr      string
		cleaners    []CustomCleaner
	}
	type args struct {
		cleaners []CustomCleaner
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		expectCount int
	}{
		{
			name: "add_nil",
			fields: fields{
				cleaners: nil,
			},
			args: args{
				cleaners: nil,
			},
			expectCount: 0,
		}, {
			name:   "add_one",
			fields: fields{},
			args: args{
				cleaners: []CustomCleaner{
					func() error {
						return nil
					},
				},
			},
			expectCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &workSpace{
				path:        tt.fields.path,
				env:         tt.fields.env,
				template:    tt.fields.template,
				customShell: tt.fields.customShell,
				output:      tt.fields.output,
				outErr:      tt.fields.outErr,
				cleaners:    tt.fields.cleaners,
			}
			w.RegistrationCustomCleaner(tt.args.cleaners...)

			if len(w.cleaners) != tt.expectCount {
				t.Errorf("the cleanner count invalid, expected: %d, acture: %d", tt.expectCount, len(w.cleaners))
			}
		})
	}
}

func Test_workSpace_Cleanup(t *testing.T) {
	tmpPath := func() string {
		s, err := os.MkdirTemp("", "test*")
		if err != nil {
			panic(err)
		}

		return s
	}()

	// Just used for cleaner
	var tmpNum int

	type fields struct {
		path        string
		env         []string
		template    string
		customShell string
		output      string
		outErr      string
		cleaners    []CustomCleaner
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		expectCheck func() bool
	}{
		{
			name: "no_customer_cleaner",
			fields: fields{
				path: tmpPath,
			},
			wantErr: false,
		},
		{
			name:    "no_path",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "with_cleaner",
			fields: fields{
				path: tmpPath,
				cleaners: []CustomCleaner{
					func() error {
						tmpNum = 100
						return nil
					},
				},
			},
			wantErr: false,
			expectCheck: func() bool {
				return tmpNum == 100
			},
		},
		{
			name: "with_multiple_cleaner",
			fields: fields{
				path: tmpPath,
				cleaners: []CustomCleaner{
					func() error {
						tmpNum = 100
						return nil
					},
					func() error {
						tmpNum += 100
						return nil
					},
				},
			},
			wantErr: false,
			expectCheck: func() bool {
				return tmpNum == 200
			},
		},
		{
			name: "one_cleaner_error",
			fields: fields{
				path: tmpPath,
				cleaners: []CustomCleaner{
					func() error {
						tmpNum = 100
						return nil
					},
					func() error {
						return errors.New("something went wrong")
					},
				},
			},
			wantErr:     true,
			expectCheck: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &workSpace{
				path:        tt.fields.path,
				env:         tt.fields.env,
				template:    tt.fields.template,
				customShell: tt.fields.customShell,
				output:      tt.fields.output,
				outErr:      tt.fields.outErr,
				cleaners:    tt.fields.cleaners,
			}
			if err := w.Cleanup(); err != nil != tt.wantErr {
				t.Errorf("got error: %v", err)
			}

			if tt.expectCheck != nil {
				if !tt.expectCheck() {
					t.Error("the cleaner invalid")
				}
			}
		})
	}
}
