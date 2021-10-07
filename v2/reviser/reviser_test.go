package reviser

import (
	"io/ioutil"
	"testing"

	_ "github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	type args struct {
		projectName string
		filePath    string
		fileContent string
	}

	tests := []struct {
		name       string
		args       args
		want       string
		wantChange bool
		wantErr    bool
	}{
		{
			name: "success with comments",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata

import (
	"log"

	"github.com/psawicki5/goimports-reviser/testdata/innderpkg"

	"bytes"

	"github.com/pkg/errors"
)

// nolint:gomnd
`,
			},
			want: `package testdata

import (
	"bytes"
	"log"

	"github.com/psawicki5/goimports-reviser/testdata/innderpkg"

	"github.com/pkg/errors"
)

// nolint:gomnd
`,
			wantChange: true,
			wantErr:    false,
		},

		{
			name: "success with std & project deps",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata

import (
	"log"

	"github.com/psawicki5/goimports-reviser/testdata/innderpkg"

	"bytes"


)

// nolint:gomnd
`,
			},
			want: `package testdata

import (
	"bytes"
	"log"

	"github.com/psawicki5/goimports-reviser/testdata/innderpkg"
)

// nolint:gomnd
`,
			wantChange: true,
			wantErr:    false,
		},

		{
			name: "success with std & third-party deps",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata
		
import (
"log"

"bytes"

"github.com/pkg/errors"
)

// nolint:gomnd
`,
			},
			want: `package testdata

import (
	"bytes"
	"log"

	"github.com/pkg/errors"
)

// nolint:gomnd
`,
			wantChange: true,
			wantErr:    false,
		},

		{
			name: "success with std deps only",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata
		
import (
"log"

"bytes"
)

// nolint:gomnd
`,
			},
			want: `package testdata

import (
	"bytes"
	"log"
)

// nolint:gomnd
`,
			wantChange: true,
			wantErr:    false,
		},

		{
			name: "success with third-party deps only",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata

import (

	"github.com/pkg/errors"
)

// nolint:gomnd
`,
			},
			want: `package testdata

import (
	"github.com/pkg/errors"
)

// nolint:gomnd
`,
			wantChange: true,
			wantErr:    false,
		},

		{
			name: "success with project deps only",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata

import (

	"github.com/psawicki5/goimports-reviser/testdata/innderpkg"

)

// nolint:gomnd
`,
			},
			want: `package testdata

import (
	"github.com/psawicki5/goimports-reviser/testdata/innderpkg"
)

// nolint:gomnd
`,
			wantChange: true,
			wantErr:    false,
		},

		{
			name: "success with clear doc for import",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata

import (
	"fmt"


	// test
	"github.com/psawicki5/goimports-reviser/testdata/innderpkg"
)

// nolint:gomnd
`,
			},
			want: `package testdata

import (
	"fmt"

	"github.com/psawicki5/goimports-reviser/testdata/innderpkg"
)

// nolint:gomnd
`,
			wantChange: true,
			wantErr:    false,
		},

		{
			name: "success with comment for import",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata

import (
	"github.com/psawicki5/goimports-reviser/testdata/innderpkg" // test1
	
	"fmt" //test2
	// this should be skipped
)

// nolint:gomnd
`,
			},
			want: `package testdata

import (
	"fmt" // test2

	"github.com/psawicki5/goimports-reviser/testdata/innderpkg" // test1
)

// nolint:gomnd
`,
			wantChange: true,
			wantErr:    false,
		},

		{
			name: "success with no changes",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata

import (
	"github.com/psawicki5/goimports-reviser/testdata/innderpkg"
)

// nolint:gomnd
`,
			},
			want: `package testdata

import (
	"github.com/psawicki5/goimports-reviser/testdata/innderpkg"
)

// nolint:gomnd
`,
			wantChange: false,
			wantErr:    false,
		},
		{
			name: "success no changes by imports and comments",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/psawicki5/goimports-reviser/pkg/somepkg"

	_ "github.com/lib/pq" // configure database/sql with postgres driver
	"github.com/pkg/errors"
	"go.uber.org/fx"
)
`,
			},
			want: `package testdata

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/psawicki5/goimports-reviser/pkg/somepkg"

	_ "github.com/lib/pq" // configure database/sql with postgres driver
	"github.com/pkg/errors"
	"go.uber.org/fx"
)
`,
			wantChange: false,
			wantErr:    false,
		},
		{
			name: "success with multiple import statements",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata

	import "sync" //test comment
	import "testing"

	// yolo
	import "fmt"


	// not sure why this is here but we shall find out soon enough
	import "io"
`,
			},
			want: `package testdata

import (
	"fmt"
	"io"
	"sync" // test comment
	"testing"
)
`,
			wantChange: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		if err := ioutil.WriteFile(tt.args.filePath, []byte(tt.args.fileContent), 0644); err != nil {
			t.Errorf("write test file failed: %s", err)
		}

		t.Run(tt.name, func(t *testing.T) {
			got, hasChange, err := Execute(tt.args.projectName, tt.args.filePath, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantChange, hasChange)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestExecute_WithRemoveUnusedImports(t *testing.T) {
	type args struct {
		projectName string
		filePath    string
		fileContent string
	}

	tests := []struct {
		name       string
		args       args
		want       string
		wantChange bool
		wantErr    bool
	}{
		{
			name: "remove unused import",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata

import (
	"fmt" //fmt package
	"github.com/pkg/errors" //custom package
)

// nolint:gomnd
func main(){
  _ = fmt.Println("test")
}
`,
			},
			want: `package testdata

import (
	"fmt" // fmt package
)

// nolint:gomnd
func main() {
	_ = fmt.Println("test")
}
`,
			wantChange: true,
			wantErr:    false,
		},

		{
			name: "remove unused import with alias",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata

import (
	"fmt" //fmt package
	p "github.com/pkg/errors" //p package
)

// nolint:gomnd
func main(){
  _ = fmt.Println("test")
}
`,
			},
			want: `package testdata

import (
	"fmt" // fmt package
)

// nolint:gomnd
func main() {
	_ = fmt.Println("test")
}
`,
			wantChange: true,
			wantErr:    false,
		},

		{
			name: "use loaded import but not used",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata

import (
	"fmt" //fmt package
	_ "github.com/pkg/errors" //custom package
)

// nolint:gomnd
func main(){
  _ = fmt.Println("test")
}
`,
			},
			want: `package testdata

import (
	"fmt" // fmt package

	_ "github.com/pkg/errors" // custom package
)

// nolint:gomnd
func main() {
	_ = fmt.Println("test")
}
`,
			wantChange: true,
			wantErr:    false,
		},
		{
			name: "success with comments before imports",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `// Some comments are here
package testdata

// test
import (
	"fmt" //fmt package
	_ "github.com/pkg/errors" //custom package
)

// nolint:gomnd
func main(){
  _ = fmt.Println("test")
}
`,
			},
			want: `// Some comments are here
package testdata

// test
import (
	"fmt" // fmt package

	_ "github.com/pkg/errors" // custom package
)

// nolint:gomnd
func main() {
	_ = fmt.Println("test")
}
`,
			wantChange: true,
			wantErr:    false,
		},
		{
			name: "success without imports",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `// Some comments are here
package main

// OutputDir the output directory where the built version of Authelia is located.
var OutputDir = "dist"

// DockerImageName the official name of Authelia docker image.
var DockerImageName = "authelia/authelia"

// IntermediateDockerImageName local name of the docker image.
var IntermediateDockerImageName = "authelia:dist"

const masterTag = "master"
const stringFalse = "false"
const stringTrue = "true"
const suitePathPrefix = "PathPrefix"
const webDirectory = "web"
`,
			},
			want: `// Some comments are here
package main

// OutputDir the output directory where the built version of Authelia is located.
var OutputDir = "dist"

// DockerImageName the official name of Authelia docker image.
var DockerImageName = "authelia/authelia"

// IntermediateDockerImageName local name of the docker image.
var IntermediateDockerImageName = "authelia:dist"

const masterTag = "master"
const stringFalse = "false"
const stringTrue = "true"
const suitePathPrefix = "PathPrefix"
const webDirectory = "web"
`,
			wantChange: false,
			wantErr:    false,
		},
		{
			name: "cleanup empty import block",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `// Some comments are here
package testdata

import (
	"fmt"
)

// nolint:gomnd
func main(){
}
`,
			},
			want: `// Some comments are here
package testdata

// nolint:gomnd
func main() {
}
`,
			wantChange: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		if err := ioutil.WriteFile(tt.args.filePath, []byte(tt.args.fileContent), 0644); err != nil {
			t.Errorf("write test file failed: %s", err)
		}

		t.Run(tt.name, func(t *testing.T) {
			got, hasChange, err := Execute(tt.args.projectName, tt.args.filePath, "", OptionRemoveUnusedImports)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantChange, hasChange)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestExecute_WithAliasForVersionSuffix(t *testing.T) {
	type args struct {
		projectName string
		filePath    string
		fileContent string
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantChange bool
		wantErr    bool
	}{
		{
			name: "success with set alias",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package main
import(
	"fmt"
	"github.com/go-pg/pg/v9"
	"strconv"
)

func main(){
	_ = strconv.Itoa(1)
	fmt.Println(pg.In([]string{"test"}))
}`,
			},
			want: `package main

import (
	"fmt"
	"strconv"

	pg "github.com/go-pg/pg/v9"
)

func main() {
	_ = strconv.Itoa(1)
	fmt.Println(pg.In([]string{"test"}))
}
`,
			wantChange: true,
			wantErr:    false,
		},
		{
			name: "success with github.com/pkg/errors",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package main
import(
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

func main(){
	_ = strconv.Itoa(1)
	fmt.Println(pg.In([]string{"test"}))
}`,
			},
			want: `package main

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

func main() {
	_ = strconv.Itoa(1)
	fmt.Println(pg.In([]string{"test"}))
}
`,
			wantChange: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		if err := ioutil.WriteFile(tt.args.filePath, []byte(tt.args.fileContent), 0644); err != nil {
			t.Errorf("write test file failed: %s", err)
		}

		t.Run(tt.name, func(t *testing.T) {
			got, hasChange, err := Execute(tt.args.projectName, tt.args.filePath, "", OptionUseAliasForVersionSuffix)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantChange, hasChange)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestExecute_WithLocalPackagePrefixes(t *testing.T) {
	type args struct {
		projectName      string
		filePath         string
		fileContent      string
		localPkgPrefixes string
	}

	tests := []struct {
		name       string
		args       args
		want       string
		wantChange bool
		wantErr    bool
	}{
		{
			name: "group local packages",
			args: args{
				projectName:      "github.com/szwagier-company",
				localPkgPrefixes: "github.com/szwagier-company/srv-serwus",
				filePath:         "./testdata/example.go",
				fileContent: `package testdata

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/szwagier-company/mirek"

	"github.com/szwagier-company/srv-serwus/pkg/book/domain"
	"github.com/szwagier-company/srv-serwus/pkg/platform/database"
)

// nolint:gomnd
func main(){
  _ = fmt.Println("test")
}
`,
			},
			want: `package testdata

import (
	"context"

	"github.com/szwagier-company/srv-serwus/pkg/book/domain"
	"github.com/szwagier-company/srv-serwus/pkg/platform/database"

	"github.com/szwagier-company/mirek"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// nolint:gomnd
func main() {
	_ = fmt.Println("test")
}
`,
			wantChange: true,
			wantErr:    false,
		},
		{
			name: "group local packages",
			args: args{
				projectName:      "github.com/psawicki5/goimports-reviser",
				localPkgPrefixes: "goimports-reviser",
				filePath:         "./testdata/example.go",
				fileContent: `package testdata

import (
	"fmt" //fmt package
	"github.com/pkg/errors" //custom package
	"github.com/psawicki5/goimports-reviser/pkg"
	"goimports-reviser/pkg"
)

// nolint:gomnd
func main(){
  _ = fmt.Println("test")
}
`,
			},
			want: `package testdata

import (
	"fmt" // fmt package

	"goimports-reviser/pkg"

	"github.com/psawicki5/goimports-reviser/pkg"

	"github.com/pkg/errors" // custom package
)

// nolint:gomnd
func main() {
	_ = fmt.Println("test")
}
`,
			wantChange: true,
			wantErr:    false,
		},
		{
			name: "group local packages",
			args: args{
				projectName:      "goimports-reviser",
				localPkgPrefixes: "github.com/psawicki5/goimports-reviser",
				filePath:         "./testdata/example.go",
				fileContent: `package testdata

import (
	"fmt" //fmt package
	"github.com/pkg/errors" //custom package
	"github.com/psawicki5/goimports-reviser/pkg"
	"goimports-reviser/pkg"
)
// nolint:gomnd
func main(){
  _ = fmt.Println("test")
}
`,
			},
			want: `package testdata

import (
	"fmt" // fmt package

	"github.com/psawicki5/goimports-reviser/pkg"

	"goimports-reviser/pkg"

	"github.com/pkg/errors" // custom package
)

// nolint:gomnd
func main() {
	_ = fmt.Println("test")
}
`,
			wantChange: true,
			wantErr:    false,
		},
		{
			name: "group local packages separately from project files",
			args: args{
				projectName:      "github.com/psawicki5/goimports-reviser/code",
				localPkgPrefixes: "github.com/psawicki5/goimports-reviser/code/thispkg",
				filePath:         "./testdata/example.go",
				fileContent: `package testdata

import (
	"fmt"
	"github.com/3rdparty/pkg"
	"github.com/psawicki5/goimports-reviser/code/foopkg"
	"github.com/psawicki5/goimports-reviser/code/otherpkg"
	"github.com/psawicki5/goimports-reviser/code/thispkg/stuff"
	"github.com/psawicki5/goimports-reviser/code/thispkg/morestuff"
)

// nolint:gomnd
func main(){
  _ = fmt.Println("test")
}
`,
			},
			want: `package testdata

import (
	"fmt"

	"github.com/psawicki5/goimports-reviser/code/thispkg/morestuff"
	"github.com/psawicki5/goimports-reviser/code/thispkg/stuff"

	"github.com/psawicki5/goimports-reviser/code/foopkg"
	"github.com/psawicki5/goimports-reviser/code/otherpkg"

	"github.com/3rdparty/pkg"
)

// nolint:gomnd
func main() {
	_ = fmt.Println("test")
}
`,
			wantChange: true,
			wantErr:    false,
		},
		{
			name: "check without local packages",
			args: args{
				projectName:      "github.com/psawicki5/goimports-reviser/code/thispkg",
				localPkgPrefixes: "",
				filePath:         "./testdata/example.go",
				fileContent: `package testdata

import (
	"fmt"
	"github.com/3rdparty/pkg"
	"github.com/psawicki5/goimports-reviser/code/foopkg"
	"github.com/psawicki5/goimports-reviser/code/otherpkg"
	"github.com/psawicki5/goimports-reviser/code/thispkg/stuff"
	"github.com/psawicki5/goimports-reviser/code/thispkg/morestuff"
)

// nolint:gomnd
func main(){
  _ = fmt.Println("test")
}
`,
			},
			want: `package testdata

import (
	"fmt"

	"github.com/psawicki5/goimports-reviser/code/thispkg/morestuff"
	"github.com/psawicki5/goimports-reviser/code/thispkg/stuff"

	"github.com/3rdparty/pkg"
	"github.com/psawicki5/goimports-reviser/code/foopkg"
	"github.com/psawicki5/goimports-reviser/code/otherpkg"
)

// nolint:gomnd
func main() {
	_ = fmt.Println("test")
}
`,
			wantChange: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		if err := ioutil.WriteFile(tt.args.filePath, []byte(tt.args.fileContent), 0644); err != nil {
			t.Errorf("write test file failed: %s", err)
		}

		t.Run(tt.name, func(t *testing.T) {
			got, hasChange, err := Execute(tt.args.projectName, tt.args.filePath, tt.args.localPkgPrefixes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantChange, hasChange)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestExecute_WithFormat(t *testing.T) {
	type args struct {
		projectName string
		filePath    string
		fileContent string
	}

	tests := []struct {
		name       string
		args       args
		want       string
		wantChange bool
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata
type SomeStruct struct{}
type SomeStruct1 struct{}
// SomeStruct2 comments
type SomeStruct2 struct{}
func (s *SomeStruct2) test() {}
func test(){}
func test1(){}
`,
			},
			want: `package testdata

type SomeStruct struct{}

type SomeStruct1 struct{}

// SomeStruct2 comments
type SomeStruct2 struct{}

func (s *SomeStruct2) test() {}

func test() {}

func test1() {}
`,
			wantChange: true,
			wantErr:    false,
		},
		{
			name: "success with comments",
			args: args{
				projectName: "github.com/psawicki5/goimports-reviser",
				filePath:    "./testdata/example.go",
				fileContent: `package testdata
// test -  test comment
func test(){}
// test1 -  test comment
func test1(){}
`,
			},
			want: `package testdata

// test -  test comment
func test() {}

// test1 -  test comment
func test1() {}
`,
			wantChange: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		if err := ioutil.WriteFile(tt.args.filePath, []byte(tt.args.fileContent), 0644); err != nil {
			t.Errorf("write test file failed: %s", err)
		}

		t.Run(tt.name, func(t *testing.T) {
			got, hasChange, err := Execute(tt.args.projectName, tt.args.filePath, "", OptionFormat)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantChange, hasChange)
			assert.Equal(t, tt.want, string(got))
		})
	}
}
