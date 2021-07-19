package testspace

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"
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
			want: &WorkSpace{
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
			want: &WorkSpace{
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
			want: &WorkSpace{
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
			want: &WorkSpace{
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
			got, err := Create(tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer got.Cleanup()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}

			_, err = os.Stat("tmp/.git")
			if err != nil {
				t.Error("git init failed")
			}
		})
	}
}
