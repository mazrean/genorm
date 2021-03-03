package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mazrean/gopendb-generator/cmd/domain"
	"github.com/mazrean/gopendb-generator/cmd/usecases/code/mock_code"
	"github.com/mazrean/gopendb-generator/cmd/usecases/config"
	"github.com/mazrean/gopendb-generator/cmd/usecases/config/mock_config"
	"github.com/stretchr/testify/assert"
)

type mockProgressCounter struct {
	isStarted  bool
	isFinished bool
}

func (mpc *mockProgressCounter) SetTotal(total int) {}
func (mpc *mockProgressCounter) Set(progress int)   {}
func (mpc *mockProgressCounter) Start() error {
	mpc.isStarted = true
	return nil
}
func (mpc *mockProgressCounter) Finish() error {
	mpc.isFinished = true
	return nil
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	assertion := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mock_config.NewMockReader(ctrl)
	conf := mock_config.NewMockConfig(ctrl)
	table := mock_config.NewMockTable(ctrl)
	codeConfig := mock_code.NewMockConfig(ctrl)
	codeTable := mock_code.NewMockTable(ctrl)

	generate := NewGenerate(reader, conf, table, codeConfig, codeTable)

	type args struct {
		yamlPath string
		rootPath string
	}
	type mock struct {
		readErr           error
		config            *domain.Config
		tables            []*domain.Table
		primaryKeyNameMap map[string][]string
		columnsMap        map[string][]*domain.Column
		referenceMap      map[string][]*config.TableReference
		generateErr       error
	}
	type expect struct {
		isErr bool
		err   error
	}
	type test struct {
		description string
		args
		mock
		expect
	}
	testCases := []test{
		{
			description: "normal",
			args: args{
				yamlPath: "/yamlpath",
				rootPath: "/rootPath",
			},
			mock: mock{
				config: &domain.Config{
					DBMS:     domain.MySQL,
					Version:  "8.0",
					Database: "test",
				},
				tables: []*domain.Table{
					{
						ID: "test",
					},
				},
				primaryKeyNameMap: map[string][]string{
					"test": {},
				},
				columnsMap: map[string][]*domain.Column{
					"test": {},
				},
				referenceMap: map[string][]*config.TableReference{
					"test": {},
				},
			},
		},
		{
			description: "read yaml error",
			args: args{
				yamlPath: "/yamlpath",
				rootPath: "/rootPath",
			},
			mock: mock{
				readErr: errors.New("read yaml error"),
				config: &domain.Config{
					DBMS:     domain.MySQL,
					Version:  "8.0",
					Database: "test",
				},
				tables:            []*domain.Table{},
				primaryKeyNameMap: map[string][]string{},
				columnsMap:        map[string][]*domain.Column{},
				referenceMap:      map[string][]*config.TableReference{},
			},
			expect: expect{
				isErr: true,
			},
		},
		{
			description: "generate code error",
			args: args{
				yamlPath: "/yamlpath",
				rootPath: "/rootPath",
			},
			mock: mock{
				config: &domain.Config{
					DBMS:     domain.MySQL,
					Version:  "8.0",
					Database: "test",
				},
				tables: []*domain.Table{
					{
						ID: "test",
					},
				},
				primaryKeyNameMap: map[string][]string{
					"test": {},
				},
				columnsMap: map[string][]*domain.Column{
					"test": {},
				},
				referenceMap: map[string][]*config.TableReference{
					"test": {},
				},
				generateErr: errors.New("generate code error"),
			},
			expect: expect{
				isErr: true,
			},
		},
	}

	for _, testCase := range testCases {
		reader.EXPECT().ReadYAML(testCase.args.yamlPath).Return(testCase.mock.readErr)
		if testCase.mock.readErr == nil {
			conf.EXPECT().Get().Return(testCase.mock.config, nil)
			codeConfig.EXPECT().Set(testCase.mock.config).Return()
			table.EXPECT().GetAll().Return(testCase.mock.tables)
			for tableID, primaryKeyNames := range testCase.mock.primaryKeyNameMap {
				table.EXPECT().GetPrimaryKeyNames(tableID).Return(primaryKeyNames, nil)
			}
			for tableID, columns := range testCase.mock.columnsMap {
				table.EXPECT().GetColumns(tableID).Return(columns, nil)
			}
			for tableID, reference := range testCase.mock.referenceMap {
				table.EXPECT().GetReference(tableID).Return(reference, nil)
			}
			if testCase.mock.generateErr == nil {
				codeTable.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return(testCase.mock.generateErr).Times(len(testCase.mock.tables))
			} else {
				codeTable.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return(testCase.mock.generateErr).Times(1)
			}
		}

		pgc := mockProgressCounter{}
		err := generate.Service(context.Background(), testCase.yamlPath, testCase.rootPath, &pgc)

		if testCase.expect.isErr {
			if err == nil {
				t.Errorf("unexpected no error(%s)", testCase.description)
			}
			if testCase.expect.err != nil && !errors.Is(err, testCase.expect.err) {
				t.Errorf("invalid error(%s): %w", testCase.description, err)
			}
			continue
		} else {
			if err != nil {
				t.Error("unexpected error: %w", err)
			}
		}

		assertion.True(pgc.isStarted, testCase.description, "is started")
		assertion.True(pgc.isFinished, testCase.description, "is finished")
	}
}
